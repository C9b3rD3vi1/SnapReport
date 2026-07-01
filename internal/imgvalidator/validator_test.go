package imgvalidator

import "testing"

func TestValidateMIME(t *testing.T) {
	tests := []struct {
		name    string
		header  []byte
		want    string
		wantErr bool
	}{
		{"png", []byte{0x89, 0x50, 0x4E, 0x47, 0x0D, 0x0A, 0x1A, 0x0A}, "image/png", false},
		{"jpeg", []byte{0xFF, 0xD8, 0xFF, 0xE0}, "image/jpeg", false},
		{"webp", []byte{0x52, 0x49, 0x46, 0x46, 0x00, 0x00, 0x00, 0x00, 0x57, 0x45, 0x42, 0x50}, "image/webp", false},
		{"empty", []byte{}, "", true},
		{"gif rejected", []byte{0x47, 0x49, 0x46, 0x38}, "", true},
		{"pdf rejected", []byte{0x25, 0x50, 0x44, 0x46}, "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ValidateMIME(tt.header)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateMIME() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ValidateMIME() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAllowedExtension(t *testing.T) {
	tests := []struct {
		mime string
		want string
	}{
		{"image/png", ".png"},
		{"image/jpeg", ".jpg"},
		{"image/webp", ".webp"},
		{"image/gif", ""},
		{"", ""},
	}
	for _, tt := range tests {
		t.Run(tt.mime, func(t *testing.T) {
			if got := AllowedExtension(tt.mime); got != tt.want {
				t.Errorf("AllowedExtension() = %v, want %v", got, tt.want)
			}
		})
	}
}
