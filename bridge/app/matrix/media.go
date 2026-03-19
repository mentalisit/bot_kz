package matrix

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func (m *Matrix) UploadMediaFromURL(url, filename string) (string, string, error) {
	// 1. Check cache (Thread-safe)
	m.mu.RLock()
	if mxc, ok := m.AvatarCache[url]; ok {
		m.mu.RUnlock()
		return mxc, "image/jpeg", nil // Cache doesn't store contentType, assuming avatar for now
	}
	m.mu.RUnlock()

	// 2. Download from HTTP
	resp, err := http.Get(url)
	if err != nil {
		return "", "", fmt.Errorf("failed to download media: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", "", fmt.Errorf("download media returned status: %d", resp.StatusCode)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", "", fmt.Errorf("failed to read media data: %v", err)
	}

	contentType := resp.Header.Get("Content-Type")
	if contentType == "" {
		contentType = http.DetectContentType(data)
	}

	// 3. Upload to Matrix
	mxc, err := m.UploadMedia(data, contentType, filename)
	if err != nil {
		return "", "", err
	}

	// 4. Update cache (only for presumably avatars or small repetitive images)
	if filename == "avatar" {
		m.mu.Lock()
		m.AvatarCache[url] = mxc
		m.mu.Unlock()
	}

	return mxc, contentType, nil
}

func (m *Matrix) UploadMedia(data []byte, contentType, filename string) (string, error) {
	url := fmt.Sprintf("%s/_matrix/media/v3/upload?access_token=%s&user_id=%s&filename=%s",
		m.Config.Matrix.HomeserverURL, m.Config.Matrix.ASToken, m.Config.Matrix.Username, filename)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(data))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", contentType)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)
	if resp.StatusCode >= 300 {
		return "", fmt.Errorf("media upload failed: %d %s", resp.StatusCode, string(respBody))
	}

	var result struct {
		ContentURI string `json:"content_uri"`
	}
	if err := json.Unmarshal(respBody, &result); err != nil {
		return "", err
	}

	if result.ContentURI == "" {
		return "", fmt.Errorf("empty content_uri in response")
	}

	return result.ContentURI, nil
}
