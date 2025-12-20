// Package voices provides reference information for ElevenLabs voices.
//
// This package contains constants and metadata for ElevenLabs' pre-made voices,
// making it easier to reference voices by name rather than ID.
//
// Note: Voice IDs and availability may change. Use client.Voices().List() for
// the authoritative list of available voices for your account.
package voices

// Pre-made voice IDs from ElevenLabs.
// These are the default voices available to all users.
const (
	// Rachel - Calm, young, American female. Great for narration and audiobooks.
	Rachel = "21m00Tcm4TlvDq8ikWAM"

	// Domi - Strong, young, American female. Confident and clear.
	Domi = "AZnzlk1XvdvUeBnXmlld"

	// Bella - Soft, young, American female. Warm and friendly.
	Bella = "EXAVITQu4vr4xnSDxMaL"

	// Antoni - Well-rounded, young, American male. Professional and warm.
	Antoni = "ErXwobaYiN019PkySvjV"

	// Elli - Emotional, young, American female. Expressive range.
	Elli = "MF3mGyEYCl7XYWbV9V6O"

	// Josh - Deep, young, American male. Authoritative and clear.
	Josh = "TxGEqnHWrfWFTfGW9XjX"

	// Arnold - Crisp, middle-aged, American male. Confident narrator.
	Arnold = "VR6AewLTigWG4xSOukaG"

	// Adam - Deep, middle-aged, American male. Rich and warm.
	Adam = "pNInz6obpgDQGcFmaJgB"

	// Sam - Raspy, young, American male. Casual and friendly.
	Sam = "yoZ06aMxZJJ28mfd3POQ"

	// Nicole - Soft, young, American female. Whispery and intimate.
	Nicole = "piTKgcLEGmPE4e6mEKli"

	// Glinda - Witch-like, middle-aged female. Theatrical and dramatic.
	Glinda = "z9fAnlkpzviPz146aGWa"

	// Clyde - War veteran, middle-aged, American male. Gruff and experienced.
	Clyde = "2EiwWnXFnvU5JabPnv8n"

	// Dave - Conversational, young, British-Essex male. Casual and natural.
	Dave = "CYw3kZ02Hs0563khs1Fj"

	// Fin - Sailor, old, Irish male. Weathered and characterful.
	Fin = "D38z5RcWu1voky8WS1ja"

	// Sarah - Soft, young, American female. News presenter style.
	Sarah = "EXAVITQu4vr4xnSDxMaL"

	// Charlotte - Seductive, middle-aged, Swedish female. Sophisticated.
	Charlotte = "XB0fDUnXU5powFXDhCwa"

	// Callum - Intense, middle-aged, Transatlantic male. Dramatic.
	Callum = "N2lVS1w4EtoT3dr4eOWO"

	// Matilda - Warm, middle-aged, American female. Friendly and approachable.
	Matilda = "XrExE9yKIg1WjnnlVkGX"

	// Grace - Southern, young, American female. Sweet and melodic.
	Grace = "oWAxZDx7w5VEj9dCyTzz"

	// Lily - Raspy, middle-aged, British female. Expressive and characterful.
	Lily = "pFZP5JQG7iQjIQuC4Bku"

	// Serena - Pleasant, middle-aged, American female. Calm and professional.
	Serena = "pMsXgVXv3BLzUgSXRplE"

	// Michael - Old, American male. Wise and grandfatherly.
	Michael = "flq6f7yk4E4fJM5XTYuZ"

	// Emily - Calm, young, American female. Clear and professional.
	Emily = "LcfcDJNUP1GQjkzn1xUU"

	// Ethan - Young, American male. Energetic and youthful.
	Ethan = "g5CIjZEefAph4nQFvHAz"

	// Brian - Deep, middle-aged, American male. Narrator quality.
	Brian = "nPczCjzI2devNBz1zQrb"

	// George - Warm, middle-aged, British male. Refined and articulate.
	George = "JBFqnCBsd6RMkjVDRZzb"

	// Gigi - Childlike, young, American female. Playful and animated.
	Gigi = "jBpfuIE2acCO8z3wKNLl"

	// Freya - Young, American female. Expressive and clear.
	Freya = "jsCqWAovK2LkecY7zXl4"

	// Harry - Anxious, young, American male. Nervous energy.
	Harry = "SOYHLrjzK2X1ezoPC6cr"

	// Jeremy - Young, American male. Conversational and natural.
	Jeremy = "bVMeCyTHy58xNoL34h3p"

	// Joseph - Middle-aged, British male. Authoritative narrator.
	Joseph = "Zlb1dXrM653N07WRdFW3"

	// Jessie - Raspy, old, American male. Weathered and experienced.
	Jessie = "t0jbNlBVZ17f02VDIeMI"

	// Drew - Well-rounded, middle-aged, American male. Versatile.
	Drew = "29vD33N1CtxCmqQRPOHJ"

	// Paul - Ground reporter, middle-aged, American male. Professional.
	Paul = "5Q0t7uMcjvnagumLfvZi"

	// River - Young, American non-binary. Modern and inclusive.
	River = "SAz9YHcvj6GT2YYXdXww"

	// Dorothy - Pleasant, young, British female. Refined and clear.
	Dorothy = "ThT5KcBeYPX3keUQqHPh"

	// Chris - Casual, middle-aged, American male. Relaxed and natural.
	Chris = "iP95p4xoKVk53GoZ742B"

	// Liam - Young, American male. Articulate and clear.
	Liam = "TX3LPaxmHKxFdv7VOQHJ"

	// James - Old, Australian male. Warm and experienced.
	James = "ZQe5CZNOzWyzPSCn5a3c"
)

// Voice represents metadata about an ElevenLabs voice.
type Voice struct {
	// ID is the unique voice identifier.
	ID string `json:"id"`

	// Name is the display name.
	Name string `json:"name"`

	// Description describes the voice characteristics.
	Description string `json:"description"`

	// Gender is the voice gender (male, female, non-binary).
	Gender string `json:"gender"`

	// Age is the approximate age category (young, middle-aged, old).
	Age string `json:"age"`

	// Accent is the primary accent/nationality.
	Accent string `json:"accent"`

	// UseCase suggests ideal use cases for this voice.
	UseCase string `json:"use_case"`

	// Category is the voice category (premade, cloned, designed).
	Category string `json:"category"`
}

// PremadeVoices returns metadata for all pre-made voices.
func PremadeVoices() []Voice {
	return []Voice{
		{ID: Rachel, Name: "Rachel", Description: "Calm and composed", Gender: "female", Age: "young", Accent: "American", UseCase: "Narration, audiobooks", Category: "premade"},
		{ID: Domi, Name: "Domi", Description: "Strong and confident", Gender: "female", Age: "young", Accent: "American", UseCase: "Presentations, announcements", Category: "premade"},
		{ID: Bella, Name: "Bella", Description: "Soft and warm", Gender: "female", Age: "young", Accent: "American", UseCase: "Podcasts, friendly content", Category: "premade"},
		{ID: Antoni, Name: "Antoni", Description: "Well-rounded and professional", Gender: "male", Age: "young", Accent: "American", UseCase: "Business, education", Category: "premade"},
		{ID: Elli, Name: "Elli", Description: "Emotional and expressive", Gender: "female", Age: "young", Accent: "American", UseCase: "Storytelling, drama", Category: "premade"},
		{ID: Josh, Name: "Josh", Description: "Deep and authoritative", Gender: "male", Age: "young", Accent: "American", UseCase: "Documentaries, news", Category: "premade"},
		{ID: Arnold, Name: "Arnold", Description: "Crisp and confident", Gender: "male", Age: "middle-aged", Accent: "American", UseCase: "Narration, commercials", Category: "premade"},
		{ID: Adam, Name: "Adam", Description: "Deep and warm", Gender: "male", Age: "middle-aged", Accent: "American", UseCase: "Audiobooks, meditation", Category: "premade"},
		{ID: Sam, Name: "Sam", Description: "Raspy and casual", Gender: "male", Age: "young", Accent: "American", UseCase: "Casual content, vlogs", Category: "premade"},
		{ID: Nicole, Name: "Nicole", Description: "Soft and whispery", Gender: "female", Age: "young", Accent: "American", UseCase: "ASMR, intimate content", Category: "premade"},
		{ID: Clyde, Name: "Clyde", Description: "Gruff war veteran", Gender: "male", Age: "middle-aged", Accent: "American", UseCase: "Character voices, gaming", Category: "premade"},
		{ID: Dave, Name: "Dave", Description: "Conversational British-Essex", Gender: "male", Age: "young", Accent: "British", UseCase: "Casual content, tutorials", Category: "premade"},
		{ID: Fin, Name: "Fin", Description: "Weathered Irish sailor", Gender: "male", Age: "old", Accent: "Irish", UseCase: "Character voices, storytelling", Category: "premade"},
		{ID: Charlotte, Name: "Charlotte", Description: "Seductive and sophisticated", Gender: "female", Age: "middle-aged", Accent: "Swedish", UseCase: "Luxury brands, dramatic content", Category: "premade"},
		{ID: Callum, Name: "Callum", Description: "Intense and dramatic", Gender: "male", Age: "middle-aged", Accent: "Transatlantic", UseCase: "Trailers, dramatic narration", Category: "premade"},
		{ID: Matilda, Name: "Matilda", Description: "Warm and friendly", Gender: "female", Age: "middle-aged", Accent: "American", UseCase: "Customer service, education", Category: "premade"},
		{ID: Grace, Name: "Grace", Description: "Southern and sweet", Gender: "female", Age: "young", Accent: "American Southern", UseCase: "Friendly content, hospitality", Category: "premade"},
		{ID: Lily, Name: "Lily", Description: "Raspy British", Gender: "female", Age: "middle-aged", Accent: "British", UseCase: "Character voices, audiobooks", Category: "premade"},
		{ID: Serena, Name: "Serena", Description: "Pleasant and calm", Gender: "female", Age: "middle-aged", Accent: "American", UseCase: "Corporate, meditation", Category: "premade"},
		{ID: Michael, Name: "Michael", Description: "Wise and grandfatherly", Gender: "male", Age: "old", Accent: "American", UseCase: "Storytelling, wisdom content", Category: "premade"},
		{ID: Emily, Name: "Emily", Description: "Calm and professional", Gender: "female", Age: "young", Accent: "American", UseCase: "News, professional content", Category: "premade"},
		{ID: Ethan, Name: "Ethan", Description: "Energetic and youthful", Gender: "male", Age: "young", Accent: "American", UseCase: "Gaming, youth content", Category: "premade"},
		{ID: Brian, Name: "Brian", Description: "Deep narrator quality", Gender: "male", Age: "middle-aged", Accent: "American", UseCase: "Documentaries, audiobooks", Category: "premade"},
		{ID: George, Name: "George", Description: "Warm and refined British", Gender: "male", Age: "middle-aged", Accent: "British", UseCase: "Narration, sophisticated content", Category: "premade"},
		{ID: Gigi, Name: "Gigi", Description: "Childlike and playful", Gender: "female", Age: "young", Accent: "American", UseCase: "Children's content, animation", Category: "premade"},
		{ID: Freya, Name: "Freya", Description: "Expressive and clear", Gender: "female", Age: "young", Accent: "American", UseCase: "Storytelling, presentations", Category: "premade"},
		{ID: Harry, Name: "Harry", Description: "Anxious energy", Gender: "male", Age: "young", Accent: "American", UseCase: "Character voices, comedy", Category: "premade"},
		{ID: Jeremy, Name: "Jeremy", Description: "Conversational and natural", Gender: "male", Age: "young", Accent: "American", UseCase: "Podcasts, casual content", Category: "premade"},
		{ID: Joseph, Name: "Joseph", Description: "Authoritative British", Gender: "male", Age: "middle-aged", Accent: "British", UseCase: "Documentaries, formal content", Category: "premade"},
		{ID: Jessie, Name: "Jessie", Description: "Raspy and weathered", Gender: "male", Age: "old", Accent: "American", UseCase: "Character voices, westerns", Category: "premade"},
		{ID: Drew, Name: "Drew", Description: "Well-rounded and versatile", Gender: "male", Age: "middle-aged", Accent: "American", UseCase: "General purpose, narration", Category: "premade"},
		{ID: Paul, Name: "Paul", Description: "Professional reporter style", Gender: "male", Age: "middle-aged", Accent: "American", UseCase: "News, journalism", Category: "premade"},
		{ID: River, Name: "River", Description: "Modern and inclusive", Gender: "non-binary", Age: "young", Accent: "American", UseCase: "Modern content, inclusive brands", Category: "premade"},
		{ID: Dorothy, Name: "Dorothy", Description: "Pleasant and refined British", Gender: "female", Age: "young", Accent: "British", UseCase: "Narration, elegant content", Category: "premade"},
		{ID: Chris, Name: "Chris", Description: "Casual and relaxed", Gender: "male", Age: "middle-aged", Accent: "American", UseCase: "Casual content, tutorials", Category: "premade"},
		{ID: Liam, Name: "Liam", Description: "Articulate and clear", Gender: "male", Age: "young", Accent: "American", UseCase: "Education, presentations", Category: "premade"},
		{ID: James, Name: "James", Description: "Warm Australian", Gender: "male", Age: "old", Accent: "Australian", UseCase: "Narration, travel content", Category: "premade"},
	}
}

// GetVoice returns voice metadata by ID.
func GetVoice(id string) *Voice {
	for _, v := range PremadeVoices() {
		if v.ID == id {
			return &v
		}
	}
	return nil
}

// GetVoiceByName returns voice metadata by name (case-insensitive).
func GetVoiceByName(name string) *Voice {
	for _, v := range PremadeVoices() {
		if equalFold(v.Name, name) {
			return &v
		}
	}
	return nil
}

// FilterByGender returns voices matching the specified gender.
func FilterByGender(gender string) []Voice {
	var result []Voice
	for _, v := range PremadeVoices() {
		if equalFold(v.Gender, gender) {
			result = append(result, v)
		}
	}
	return result
}

// FilterByAccent returns voices matching the specified accent.
func FilterByAccent(accent string) []Voice {
	var result []Voice
	for _, v := range PremadeVoices() {
		if containsFold(v.Accent, accent) {
			result = append(result, v)
		}
	}
	return result
}

// FilterByAge returns voices matching the specified age category.
func FilterByAge(age string) []Voice {
	var result []Voice
	for _, v := range PremadeVoices() {
		if equalFold(v.Age, age) {
			result = append(result, v)
		}
	}
	return result
}

// equalFold is a simple case-insensitive string comparison.
func equalFold(a, b string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := 0; i < len(a); i++ {
		ca, cb := a[i], b[i]
		if ca >= 'A' && ca <= 'Z' {
			ca += 'a' - 'A'
		}
		if cb >= 'A' && cb <= 'Z' {
			cb += 'a' - 'A'
		}
		if ca != cb {
			return false
		}
	}
	return true
}

// containsFold checks if b is contained in a (case-insensitive).
func containsFold(a, b string) bool {
	if len(b) > len(a) {
		return false
	}
	for i := 0; i <= len(a)-len(b); i++ {
		if equalFold(a[i:i+len(b)], b) {
			return true
		}
	}
	return false
}
