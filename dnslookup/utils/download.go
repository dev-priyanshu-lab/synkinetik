package utils

import (
	"crypto/tls"
	"io"
	"net/http"
	"os"
	"time"
)

func DownloadLink(file, url string) bool {
	tr := &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}}
	client := &http.Client{Transport: tr, Timeout: 180 * time.Second}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return false
	}
	req.Header.Set("User-Agent", "Mozilla/5.0")
	resp, err := client.Do(req)
	if err != nil {
		return false
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return false
	}
	out, err := os.Create(file)
	if err != nil {
		return false
	}
	defer out.Close()
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return false
	}
	return true
}
