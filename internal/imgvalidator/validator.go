package imgvalidator

import (
	"fmt"
	"net/http"
)

var allowedMIMETypes = map[string]string{
	"image/png":  ".png",
	"image/jpeg": ".jpg",
	"image/webp": ".webp",
}

func ValidateMIME(header []byte) (string, error) {
	mime := detectMIME(header)
	if _, ok := allowedMIMETypes[mime]; !ok {
		return "", fmt.Errorf("unsupported file type: %s", mime)
	}
	return mime, nil
}

func detectMIME(header []byte) string {
	if len(header) >= 12 &&
		header[0] == 0x52 && header[1] == 0x49 && header[2] == 0x46 && header[3] == 0x46 &&
		header[8] == 0x57 && header[9] == 0x45 && header[10] == 0x42 && header[11] == 0x50 {
		return "image/webp"
	}
	return http.DetectContentType(header)
}

func AllowedExtension(mime string) string {
	return allowedMIMETypes[mime]
}
