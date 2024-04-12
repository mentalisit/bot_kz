package helper

import (
	"bytes"
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
	r := bytes.NewReader(*data)
	m, err := webp.Decode(r)
	if err != nil {
		return err
	}
	var output []byte
	w := bytes.NewBuffer(output)
	if err = png.Encode(w, m); err != nil {
		return err
	}
	*data = w.Bytes()
	return nil
}
