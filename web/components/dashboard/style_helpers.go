package dashboard

// GetStatIconClasses retorna classes completas baseadas nas cores
// Uso de mapa para garantir que todas as classes existam no CSS
func GetStatIconClasses(bgColor string, textColor string) string {
	// Mapeamento de cores para classes completas
	bgClasses := map[string]string{
		"bg-arandu-soft": "bg-arandu-soft",
		"bg-primary-100": "bg-primary-100",
		"bg-primary-50":  "bg-primary-50",
		"bg-accent-100":  "bg-accent-100",
		"bg-accent-50":   "bg-accent-50",
		"bg-neutral-100": "bg-neutral-100",
		"bg-neutral-50":  "bg-neutral-50",
		"bg-blue-50":     "bg-blue-50",
		"bg-blue-100":    "bg-blue-100",
		"bg-green-50":    "bg-green-50",
		"bg-green-100":   "bg-green-100",
		"bg-amber-50":    "bg-amber-50",
		"bg-amber-100":   "bg-amber-100",
		"bg-purple-50":   "bg-purple-50",
		"bg-purple-100":  "bg-purple-100",
		"bg-red-50":      "bg-red-50",
		"bg-red-100":     "bg-red-100",
		"bg-arandu-bg":   "bg-arandu-bg",
	}

	textClasses := map[string]string{
		"text-arandu-primary": "text-arandu-primary",
		"text-primary-600":    "text-primary-600",
		"text-primary-500":    "text-primary-500",
		"text-accent-600":     "text-accent-600",
		"text-accent-500":     "text-accent-500",
		"text-neutral-600":    "text-neutral-600",
		"text-neutral-500":    "text-neutral-500",
		"text-blue-600":       "text-blue-600",
		"text-blue-500":       "text-blue-500",
		"text-green-600":      "text-green-600",
		"text-green-500":      "text-green-500",
		"text-amber-600":      "text-amber-600",
		"text-amber-500":      "text-amber-500",
		"text-purple-600":     "text-purple-600",
		"text-purple-500":     "text-purple-500",
		"text-red-600":        "text-red-600",
		"text-red-500":        "text-red-500",
		"text-arandu-dark":    "text-arandu-dark",
	}

	baseClasses := "w-12 h-12 rounded-xl flex items-center justify-center text-xl"

	bgClass := bgClasses[bgColor]
	if bgClass == "" {
		bgClass = "bg-neutral-100" // fallback
	}

	textClass := textClasses[textColor]
	if textClass == "" {
		textClass = "text-neutral-600" // fallback
	}

	return baseClasses + " " + bgClass + " " + textClass
}

// GetStatValueClasses retorna classes para o valor do stat
func GetStatValueClasses(textColor string) string {
	validColors := map[string]string{
		"text-arandu-primary": "text-arandu-primary",
		"text-primary-600":    "text-primary-600",
		"text-accent-600":     "text-accent-600",
		"text-neutral-600":    "text-neutral-600",
		"text-blue-600":       "text-blue-600",
		"text-green-600":      "text-green-600",
		"text-amber-600":      "text-amber-600",
		"text-purple-600":     "text-purple-600",
		"text-red-600":        "text-red-600",
		"text-arandu-dark":    "text-arandu-dark",
	}

	if color, ok := validColors[textColor]; ok {
		return "text-2xl font-bold " + color
	}
	return "text-2xl font-bold text-neutral-600"
}
