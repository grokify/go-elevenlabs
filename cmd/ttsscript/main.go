// Command ttsscript generates TTS audio from a JSON script file using ElevenLabs.
//
// Usage:
//
//	ttsscript [flags] <script.json>
//
// Flags:
//
//	-lang string      Language code to generate (default "en")
//	-output string    Output directory (default "./output")
//	-per-slide        Concatenate segments into per-slide audio files (requires ffmpeg)
//	-manifest         Generate manifest JSON file (default true)
//	-dry-run          Show what would be generated without calling API
//	-model string     ElevenLabs model ID (default "eleven_multilingual_v2")
//
// Environment:
//
//	ELEVENLABS_API_KEY    Required API key for ElevenLabs
package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"

	elevenlabs "github.com/agentplexus/go-elevenlabs"
	"github.com/agentplexus/go-elevenlabs/ttsscript"
)

func main() {
	// Parse flags
	lang := flag.String("lang", "en", "Language code to generate")
	outputDir := flag.String("output", "./output", "Output directory")
	perSlide := flag.Bool("per-slide", false, "Concatenate segments into per-slide audio files (requires ffmpeg)")
	manifest := flag.Bool("manifest", true, "Generate manifest JSON file")
	dryRun := flag.Bool("dry-run", false, "Show what would be generated without calling API")
	modelID := flag.String("model", "eleven_multilingual_v2", "ElevenLabs model ID")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [flags] <script.json>\n\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "Generate TTS audio from a JSON script file using ElevenLabs.\n\n")
		fmt.Fprintf(os.Stderr, "Flags:\n")
		flag.PrintDefaults()
		fmt.Fprintf(os.Stderr, "\nEnvironment:\n")
		fmt.Fprintf(os.Stderr, "  ELEVENLABS_API_KEY    Required API key for ElevenLabs\n")
	}

	flag.Parse()

	if flag.NArg() < 1 {
		flag.Usage()
		os.Exit(1)
	}

	scriptPath := flag.Arg(0)

	// Check for API key (unless dry run)
	if !*dryRun && os.Getenv("ELEVENLABS_API_KEY") == "" {
		log.Fatal("ELEVENLABS_API_KEY environment variable is required")
	}

	// Check for ffmpeg if per-slide mode
	if *perSlide {
		if _, err := exec.LookPath("ffmpeg"); err != nil {
			log.Fatal("ffmpeg is required for --per-slide mode but was not found in PATH")
		}
	}

	// Load script
	script, err := ttsscript.LoadScript(scriptPath)
	if err != nil {
		log.Fatalf("Failed to load script: %v", err)
	}

	// Validate script
	if issues := script.Validate(); len(issues) > 0 {
		log.Fatalf("Script validation failed:\n  - %s", strings.Join(issues, "\n  - "))
	}

	fmt.Printf("Script: %s\n", script.Title)
	fmt.Printf("Language: %s\n", *lang)
	fmt.Printf("Slides: %d, Segments: %d\n", script.SlideCount(), script.SegmentCount())

	// Compile script
	compiler := ttsscript.NewCompiler()
	segments, err := compiler.Compile(script, *lang)
	if err != nil {
		log.Fatalf("Failed to compile script: %v", err)
	}

	// Format for ElevenLabs
	formatter := ttsscript.NewElevenLabsFormatter()
	jobs := formatter.Format(segments)

	fmt.Printf("Generated %d TTS jobs\n\n", len(jobs))

	// Create output directory
	if err := os.MkdirAll(*outputDir, 0750); err != nil {
		log.Fatalf("Failed to create output directory: %v", err)
	}

	// Generate batch config
	config := ttsscript.NewBatchConfig(*outputDir)
	config.IncludeLanguageInFilename = true

	// Generate manifest
	manifestEntries := ttsscript.GenerateManifest(jobs, config, *lang)

	if *dryRun {
		fmt.Println("Dry run - would generate:")
		for _, entry := range manifestEntries {
			segType := "segment"
			if entry.IsTitleSegment {
				segType = "title"
			}
			fmt.Printf("  [%s] %s\n", segType, entry.OutputFile)
			fmt.Printf("    Text: %s\n", truncate(entry.Text, 60))
			fmt.Printf("    Voice: %s\n", entry.VoiceID)
		}

		if *perSlide {
			fmt.Println("\nPer-slide output:")
			slideFiles := getSlideOutputFiles(manifestEntries, config, *lang)
			for slide, file := range slideFiles {
				fmt.Printf("  Slide %d: %s\n", slide+1, file)
			}
		}
		return
	}

	// Create ElevenLabs client
	client, err := elevenlabs.NewClient()
	if err != nil {
		log.Fatalf("Failed to create ElevenLabs client: %v", err)
	}

	ctx := context.Background()

	// Generate audio for each segment
	generatedFiles := make([]string, 0, len(jobs))
	for i, job := range jobs {
		if job.VoiceID == "" {
			log.Printf("Skipping segment %d: no voice ID configured", i+1)
			continue
		}

		outputFile := config.GenerateFilename(job, *lang)

		segType := "segment"
		if job.IsTitleSegment {
			segType = "title"
		}

		fmt.Printf("[%d/%d] Generating %s: %s\n", i+1, len(jobs), segType, truncate(job.Text, 50))

		resp, err := client.TextToSpeech().Generate(ctx, &elevenlabs.TTSRequest{
			VoiceID:       job.VoiceID,
			Text:          job.Text,
			ModelID:       *modelID,
			VoiceSettings: elevenlabs.DefaultVoiceSettings(),
		})
		if err != nil {
			log.Printf("  ERROR: %v", err)
			continue
		}
		audio := resp.Audio

		f, err := os.Create(outputFile)
		if err != nil {
			log.Printf("  ERROR creating file: %v", err)
			continue
		}

		_, err = io.Copy(f, audio)
		f.Close()
		if err != nil {
			log.Printf("  ERROR writing file: %v", err)
			continue
		}

		fmt.Printf("  Saved: %s\n", outputFile)
		generatedFiles = append(generatedFiles, outputFile)
	}

	// Write manifest
	if *manifest {
		manifestPath := filepath.Join(*outputDir, fmt.Sprintf("manifest_%s.json", *lang))
		manifestData, err := json.MarshalIndent(manifestEntries, "", "  ")
		if err != nil {
			log.Printf("Failed to marshal manifest: %v", err)
		} else if err := os.WriteFile(manifestPath, manifestData, 0600); err != nil {
			log.Printf("Failed to write manifest: %v", err)
		} else {
			fmt.Printf("\nManifest saved: %s\n", manifestPath)
		}
	}

	// Concatenate per-slide if requested
	if *perSlide {
		fmt.Println("\nConcatenating per-slide audio...")
		concatenatePerSlide(manifestEntries, *lang, *outputDir)
	}

	fmt.Printf("\nDone! Generated %d audio files.\n", len(generatedFiles))
}

// concatenatePerSlide uses ffmpeg to concatenate segment audio files into per-slide files.
func concatenatePerSlide(entries []ttsscript.ManifestEntry, language, outputDir string) {
	// Group entries by slide
	slideSegments := make(map[int][]ttsscript.ManifestEntry)
	for _, entry := range entries {
		slideSegments[entry.SlideIndex] = append(slideSegments[entry.SlideIndex], entry)
	}

	// Get sorted slide indices
	slideIndices := make([]int, 0, len(slideSegments))
	for idx := range slideSegments {
		slideIndices = append(slideIndices, idx)
	}
	sort.Ints(slideIndices)

	for _, slideIdx := range slideIndices {
		segments := slideSegments[slideIdx]

		// Sort segments: title first (SegmentIndex -1), then by segment index
		sort.Slice(segments, func(i, j int) bool {
			return segments[i].SegmentIndex < segments[j].SegmentIndex
		})

		// Skip if only one segment (no need to concatenate)
		if len(segments) == 1 {
			// Just copy/rename to slide output
			slideOutput := filepath.Join(outputDir, fmt.Sprintf("slide%02d_%s.mp3", slideIdx+1, language))
			if err := copyFile(segments[0].OutputFile, slideOutput); err != nil {
				log.Printf("  Slide %d: failed to copy: %v", slideIdx+1, err)
				continue
			}
			fmt.Printf("  Slide %d: %s (1 segment)\n", slideIdx+1, slideOutput)
			continue
		}

		// Create concat list file for ffmpeg
		listFile := filepath.Join(outputDir, fmt.Sprintf(".concat_slide%02d.txt", slideIdx+1))
		var listContent strings.Builder

		for i, seg := range segments {
			// Add pause before (as silence) if needed
			if seg.PauseBeforeMs > 0 && i > 0 {
				silenceFile, err := generateSilence(outputDir, seg.PauseBeforeMs, slideIdx, i, "before")
				if err != nil {
					log.Printf("  Warning: failed to generate silence: %v", err)
				} else {
					listContent.WriteString(fmt.Sprintf("file '%s'\n", filepath.Base(silenceFile)))
				}
			}

			// Add the audio file
			listContent.WriteString(fmt.Sprintf("file '%s'\n", filepath.Base(seg.OutputFile)))

			// Add pause after (as silence) if needed
			if seg.PauseAfterMs > 0 {
				silenceFile, err := generateSilence(outputDir, seg.PauseAfterMs, slideIdx, i, "after")
				if err != nil {
					log.Printf("  Warning: failed to generate silence: %v", err)
				} else {
					listContent.WriteString(fmt.Sprintf("file '%s'\n", filepath.Base(silenceFile)))
				}
			}
		}

		if err := os.WriteFile(listFile, []byte(listContent.String()), 0600); err != nil {
			log.Printf("  Slide %d: failed to write concat list: %v", slideIdx+1, err)
			continue
		}

		// Run ffmpeg to concatenate
		slideOutput := filepath.Join(outputDir, fmt.Sprintf("slide%02d_%s.mp3", slideIdx+1, language))
		cmd := exec.Command("ffmpeg", "-y", "-f", "concat", "-safe", "0", "-i", listFile, "-c", "copy", slideOutput)
		cmd.Dir = outputDir
		if output, err := cmd.CombinedOutput(); err != nil {
			log.Printf("  Slide %d: ffmpeg failed: %v\n%s", slideIdx+1, err, string(output))
			continue
		}

		// Clean up temp files
		os.Remove(listFile)
		cleanupSilenceFiles(outputDir, slideIdx)

		fmt.Printf("  Slide %d: %s (%d segments)\n", slideIdx+1, slideOutput, len(segments))
	}
}

// generateSilence creates a silent audio file of the specified duration.
func generateSilence(outputDir string, durationMs, slideIdx, segIdx int, position string) (string, error) {
	filename := filepath.Join(outputDir, fmt.Sprintf(".silence_s%02d_%02d_%s.mp3", slideIdx, segIdx, position))
	duration := float64(durationMs) / 1000.0

	// #nosec G204 -- filename is constructed from user-controlled outputDir flag, which is intentional for CLI tools
	cmd := exec.Command("ffmpeg", "-y", "-f", "lavfi", "-i",
		fmt.Sprintf("anullsrc=r=44100:cl=mono:d=%.3f", duration),
		"-c:a", "libmp3lame", "-q:a", "9", filename)

	if output, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("ffmpeg silence generation failed: %v\n%s", err, string(output))
	}

	return filename, nil
}

// cleanupSilenceFiles removes temporary silence files for a slide.
func cleanupSilenceFiles(outputDir string, slideIdx int) {
	pattern := filepath.Join(outputDir, fmt.Sprintf(".silence_s%02d_*.mp3", slideIdx))
	files, _ := filepath.Glob(pattern)
	for _, f := range files {
		os.Remove(f)
	}
}

// copyFile copies a file from src to dst.
func copyFile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, in)
	return err
}

// getSlideOutputFiles returns a map of slide index to output file path.
func getSlideOutputFiles(entries []ttsscript.ManifestEntry, config *ttsscript.BatchConfig, language string) map[int]string {
	slides := make(map[int]string)
	for _, entry := range entries {
		if _, exists := slides[entry.SlideIndex]; !exists {
			slides[entry.SlideIndex] = filepath.Join(config.OutputDir, fmt.Sprintf("slide%02d_%s.mp3", entry.SlideIndex+1, language))
		}
	}
	return slides
}

func truncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen-3] + "..."
}
