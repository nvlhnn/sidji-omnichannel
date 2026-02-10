package services

import (
	"encoding/json"
	"fmt"
	"io"
	"mime"
	"net/http"
	"os"
	"path/filepath"

	"github.com/sidji-omnichannel/internal/config"
)

type MediaService struct {
	cfg *config.Config
}

func NewMediaService(cfg *config.Config) *MediaService {
	// Ensure uploads directory exists
	if err := os.MkdirAll("uploads", 0755); err != nil {
		fmt.Printf("Error creating uploads directory: %v\n", err)
	}
	return &MediaService{cfg: cfg}
}

// MetaMediaResponse is the response from getting media URL
type MetaMediaResponse struct {
	URL      string `json:"url"`
	MimeType string `json:"mime_type"`
	Sha256   string `json:"sha256"`
	FileSize int    `json:"file_size"`
	ID       string `json:"id"`
}

// DownloadMetaMedia downloads media from Meta and saves it locally
// Returns the local relative path (e.g., "/uploads/123.jpg")
func (s *MediaService) DownloadMetaMedia(mediaID, accessToken string) (string, error) {
	// 1. Get the download URL
	metaURL := fmt.Sprintf("https://graph.facebook.com/v19.0/%s", mediaID)
	req, err := http.NewRequest("GET", metaURL, nil)
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+accessToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to get media URL: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("meta API error getting media URL (status %d): %s", resp.StatusCode, string(body))
	}

	var mediaResp MetaMediaResponse
	if err := json.NewDecoder(resp.Body).Decode(&mediaResp); err != nil {
		return "", fmt.Errorf("failed to decode media response: %w", err)
	}

	// 2. Download the actual binary
	downloadReq, err := http.NewRequest("GET", mediaResp.URL, nil)
	if err != nil {
		return "", fmt.Errorf("failed to create download request: %w", err)
	}
	downloadReq.Header.Set("Authorization", "Bearer "+accessToken)

	downloadResp, err := client.Do(downloadReq)
	if err != nil {
		return "", fmt.Errorf("failed to download media blob: %w", err)
	}
	defer downloadResp.Body.Close()

	if downloadResp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("download error (status %d)", downloadResp.StatusCode)
	}

	// 3. Save to file
	// Determine extension from mime type
	exts, _ := mime.ExtensionsByType(mediaResp.MimeType)
	ext := ".bin"
	if len(exts) > 0 {
		ext = exts[0]
	}
    // Some mime types might return weird extensions, basic fallback logical
    if mediaResp.MimeType == "image/jpeg" { ext = ".jpg" }
    if mediaResp.MimeType == "image/png" { ext = ".png" }

	filename := fmt.Sprintf("%s%s", mediaID, ext)
	localPath := filepath.Join("uploads", filename)

	out, err := os.Create(localPath)
	if err != nil {
		return "", fmt.Errorf("failed to create local file: %w", err)
	}
	defer out.Close()

	if _, err := io.Copy(out, downloadResp.Body); err != nil {
		return "", fmt.Errorf("failed to write file: %w", err)
	}

	// Return web-accessible path (assuming we serve /uploads route)
	// Using forward slashes for URL consistency
	publicPath := fmt.Sprintf("/uploads/%s", filename)
	return publicPath, nil
}
