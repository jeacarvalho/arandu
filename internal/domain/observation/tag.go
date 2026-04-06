package observation

import "time"

// TagType represents the type of classification tag
type TagType string

const (
	TagTypeEmotion      TagType = "emotion"
	TagTypeBehavior     TagType = "behavior"
	TagTypeCognition    TagType = "cognition"
	TagTypeRelationship TagType = "relationship"
	TagTypeSomatic      TagType = "somatic"
	TagTypeContext      TagType = "context"
)

// Tag represents a pre-defined classification tag
type Tag struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	TagType   TagType   `json:"tag_type"`
	Color     string    `json:"color"`
	SortOrder int       `json:"sort_order"`
	CreatedAt time.Time `json:"created_at"`
}

// ObservationTag represents the many-to-many relationship
// between observations and tags, with intensity level
type ObservationTag struct {
	ID            string    `json:"id"`
	ObservationID string    `json:"observation_id"`
	TagID         string    `json:"tag_id"`
	Tag           *Tag      `json:"tag,omitempty"`
	Intensity     int       `json:"intensity"` // 1-5
	CreatedAt     time.Time `json:"created_at"`
}

// TagSummary represents a summary of tags by type
type TagSummary struct {
	TagType TagType `json:"tag_type"`
	Count   int     `json:"count"`
}

// ClassificationData represents data for classification UI
type ClassificationData struct {
	ObservationID string           `json:"observation_id"`
	AvailableTags []Tag            `json:"available_tags"`
	SelectedTags  []ObservationTag `json:"selected_tags"`
}

// IsValidTagType checks if the tag type is valid
func IsValidTagType(t string) bool {
	switch TagType(t) {
	case TagTypeEmotion, TagTypeBehavior, TagTypeCognition,
		TagTypeRelationship, TagTypeSomatic, TagTypeContext:
		return true
	}
	return false
}

// IsValidIntensity checks if intensity is within valid range (1-5)
func IsValidIntensity(i int) bool {
	return i >= 1 && i <= 5
}

// TagTypeLabel returns the human-readable label for a tag type
func TagTypeLabel(tt TagType) string {
	switch tt {
	case TagTypeEmotion:
		return "Emoção"
	case TagTypeBehavior:
		return "Comportamento"
	case TagTypeCognition:
		return "Cognição"
	case TagTypeRelationship:
		return "Relação"
	case TagTypeSomatic:
		return "Soma"
	case TagTypeContext:
		return "Contexto"
	default:
		return string(tt)
	}
}

// IntensityColor returns a color based on intensity level
func IntensityColor(intensity int) string {
	switch intensity {
	case 1:
		return "bg-green-100 text-green-800 border-green-200"
	case 2:
		return "bg-green-200 text-green-900 border-green-300"
	case 3:
		return "bg-yellow-100 text-yellow-800 border-yellow-200"
	case 4:
		return "bg-orange-100 text-orange-800 border-orange-200"
	case 5:
		return "bg-red-100 text-red-800 border-red-200"
	default:
		return "bg-gray-100 text-gray-800 border-gray-200"
	}
}
