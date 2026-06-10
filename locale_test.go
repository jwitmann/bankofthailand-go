package bankofthailand

import "testing"

func TestPickString(t *testing.T) {
	tests := []struct {
		name     string
		loc      Locale
		thai     string
		english  string
		expected string
	}{
		{"english preferred", LocaleEnglish, "ไทย", "English", "English"},
		{"thai preferred", LocaleThai, "ไทย", "English", "ไทย"},
		{"english fallback to thai", LocaleEnglish, "", "", ""},
		{"thai fallback to english", LocaleThai, "", "English", "English"},
		{"english fallback when thai empty", LocaleEnglish, "", "English", "English"},
		{"normalize TH", "TH", "ไทย", "English", "ไทย"},
		{"normalize EN", "EN", "ไทย", "English", "English"},
		{"unknown defaults to english", "xyz", "ไทย", "English", "English"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := pickString(tt.loc, tt.thai, tt.english)
			if got != tt.expected {
				t.Errorf("pickString(%q, %q, %q) = %q, want %q", tt.loc, tt.thai, tt.english, got, tt.expected)
			}
		})
	}
}
