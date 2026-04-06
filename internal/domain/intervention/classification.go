package intervention

import "time"

// TagType represents the type of intervention classification
type TagType string

const (
	TagTypeCognitive       TagType = "cognitive"
	TagTypeBehavioral      TagType = "behavioral"
	TagTypeEmotional       TagType = "emotional"
	TagTypePsychoeducation TagType = "psychoeducation"
	TagTypeNarrative       TagType = "narrative"
	TagTypeBody            TagType = "body"
)

// Tag represents a pre-defined classification tag for interventions
type Tag struct {
	ID        string
	Name      string
	TagType   TagType
	Color     string
	Icon      string
	SortOrder int
	CreatedAt time.Time
}

// InterventionClassification represents the relationship between an intervention and a tag
type InterventionClassification struct {
	ID             string
	InterventionID string
	TagID          string
	Tag            *Tag
	Intensity      int // 1-5, optional
	CreatedAt      time.Time
}

// GetTagTypeColor returns the color for a given tag type
func GetTagTypeColor(tagType TagType) string {
	switch tagType {
	case TagTypeCognitive:
		return "#7C3AED"
	case TagTypeBehavioral:
		return "#1D9E75"
	case TagTypeEmotional:
		return "#0F6E56"
	case TagTypePsychoeducation:
		return "#F59E0B"
	case TagTypeNarrative:
		return "#3B82F6"
	case TagTypeBody:
		return "#DC2626"
	default:
		return "#6B7280"
	}
}

// GetTagTypeIcon returns the icon class for a given tag type
func GetTagTypeIcon(tagType TagType) string {
	switch tagType {
	case TagTypeCognitive:
		return "brain"
	case TagTypeBehavioral:
		return "running"
	case TagTypeEmotional:
		return "heart"
	case TagTypePsychoeducation:
		return "book-open"
	case TagTypeNarrative:
		return "comment-dots"
	case TagTypeBody:
		return "body"
	default:
		return "tag"
	}
}

// GetTagTypeDisplayName returns the display name for a tag type
func GetTagTypeDisplayName(tagType TagType) string {
	switch tagType {
	case TagTypeCognitive:
		return "Técnica Cognitiva"
	case TagTypeBehavioral:
		return "Técnica Comportamental"
	case TagTypeEmotional:
		return "Técnica Emocional"
	case TagTypePsychoeducation:
		return "Psicoeducação"
	case TagTypeNarrative:
		return "Exploração Narrativa"
	case TagTypeBody:
		return "Intervenção Corporal"
	default:
		return "Outro"
	}
}

// TagTypeOrder defines the display order for tag types
var TagTypeOrder = []TagType{
	TagTypeCognitive,
	TagTypeBehavioral,
	TagTypeEmotional,
	TagTypePsychoeducation,
	TagTypeNarrative,
	TagTypeBody,
}
