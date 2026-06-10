package bankofthailand

type Locale string

const (
	LocaleThai    Locale = "th"
	LocaleEnglish Locale = "en"
)

func normalizeLocale(loc Locale) Locale {
	switch loc {
	case LocaleThai, "TH", "tha":
		return LocaleThai
	case LocaleEnglish, "EN", "eng":
		return LocaleEnglish
	default:
		return LocaleEnglish
	}
}

func pickString(loc Locale, thai, english string) string {
	loc = normalizeLocale(loc)
	if loc == LocaleThai {
		if thai != "" {
			return thai
		}
		return english
	}
	if english != "" {
		return english
	}
	return thai
}
