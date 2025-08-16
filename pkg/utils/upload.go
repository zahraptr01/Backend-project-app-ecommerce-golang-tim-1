package utils

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
)

var cdnUploadURL = "https://cdn-lumoshive-academy.vercel.app/api/v1/upload"

// UploadImageToCDN meng-upload file ke CDN dan mengembalikan data.url
func UploadImageToCDN(ctx context.Context, r io.Reader, filename string, folder string) (string, error) {
	var body bytes.Buffer
	writer := multipart.NewWriter(&body)

	// field: image (file)
	fw, err := writer.CreateFormFile("image", filename)
	if err != nil {
		return "", err
	}
	if _, err := io.Copy(fw, r); err != nil {
		return "", err
	}

	// field: folder (text)
	if err := writer.WriteField("folder", folder); err != nil {
		return "", err
	}

	if err := writer.Close(); err != nil {
		return "", err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, cdnUploadURL, &body)
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		b, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("cdn upload failed: %s", string(b))
	}

	var payload struct {
		Success bool   `json:"success"`
		Message string `json:"message"`
		Data    struct {
			URL string `json:"url"`
		} `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		return "", err
	}
	if !payload.Success || payload.Data.URL == "" {
		return "", fmt.Errorf("invalid cdn response")
	}
	return payload.Data.URL, nil
}
