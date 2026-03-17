package patient

import "encoding/json"

type ThemeItem struct {
	Name        string
	Count       int
	WeightClass string
}

type ThemeCloudViewModel struct {
	PatientID  string
	Themes     []ThemeItem
	TotalCount int
}

func (vm ThemeCloudViewModel) ThemesJSON() string {
	data, _ := json.Marshal(vm.Themes)
	return string(data)
}

func CalculateWeightClass(count, maxCount int) string {
	if maxCount == 0 {
		return "theme-lv1"
	}

	ratio := float64(count) / float64(maxCount)

	switch {
	case ratio >= 0.8:
		return "theme-lv5"
	case ratio >= 0.6:
		return "theme-lv4"
	case ratio >= 0.4:
		return "theme-lv3"
	case ratio >= 0.2:
		return "theme-lv2"
	default:
		return "theme-lv1"
	}
}
