package helper

import (
	"bytes"
	"fmt"
	"golang.org/x/image/webp"
	"image/png"
	"io"
	"net/http"
	"time"
)

func DownloadFile(url string) ([]byte, error) {
	return DownloadFileAuth(url, "")
}

func DownloadFileAuth(url string, auth string) ([]byte, error) {
	var buf bytes.Buffer
	client := &http.Client{
		Timeout: time.Second * 5,
	}
	req, err := http.NewRequest("GET", url, nil)
	if auth != "" {
		req.Header.Add("Authorization", auth)
	}
	if err != nil {
		return nil, err
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	io.Copy(&buf, resp.Body)
	data := buf.Bytes()
	return data, nil
}
func ConvertWebPToPNG(data *[]byte) error {
	// Check if data is empty
	if len(*data) == 0 {
		return fmt.Errorf("empty file data")
	}

	// Check WebP file signature
	if len(*data) < 12 {
		return fmt.Errorf("file too small to be valid WebP")
	}

	// WebP files should start with "RIFF" and have "WEBP" at bytes 8-11
	if string((*data)[0:4]) != "RIFF" || string((*data)[8:12]) != "WEBP" {
		return fmt.Errorf("invalid WebP file signature")
	}

	r := bytes.NewReader(*data)
	m, err := webp.Decode(r)
	if err != nil {
		return fmt.Errorf("webp decode failed: %w", err)
	}
	var output []byte
	w := bytes.NewBuffer(output)
	if err = png.Encode(w, m); err != nil {
		return fmt.Errorf("png encode failed: %w", err)
	}
	*data = w.Bytes()
	return nil
}
