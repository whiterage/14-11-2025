package service

import "testing"

func TestNormalizeURL(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		input   string
		wantURL string
		wantErr bool
	}{
		{"https input", "https://example.com", "https://example.com", false},
		{"http input", "http://example.com/page", "http://example.com/page", false},
		{"without scheme", "example.com", "https://example.com", false},
		{"with spaces", "   yandex.ru  ", "https://yandex.ru", false},
		{"empty", "", "", true},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, err := normalizeURL(tt.input)
			if tt.wantErr {
				if err == nil {
					t.Fatalf("expected error for %q", tt.input)
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got != tt.wantURL {
				t.Fatalf("unexpected url: got %q want %q", got, tt.wantURL)
			}
		})
	}
}
