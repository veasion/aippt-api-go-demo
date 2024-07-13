package main

import (
	"bufio"
	"bytes"
	"io"
	"net/http"
	"os"
	"strings"
)

type StreamConsumer func(data string) error

type HttpResponse struct {
	status      int
	text        string
	headers     http.Header
	contentType string
}

func PostSse(url string, headers map[string]string, body string, consumer StreamConsumer) (*HttpResponse, error) {
	bodyReader := strings.NewReader(body)
	req, err := http.NewRequest("POST", url, bodyReader)
	if err != nil {
		return nil, err
	}
	if headers != nil {
		for key, value := range headers {
			req.Header.Set(key, value)
		}
	}
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	contentType := resp.Header.Get("Content-Type")
	defer resp.Body.Close()
	var text string
	if resp.StatusCode == 200 {
		if strings.Contains(contentType, "text/event-stream") {
			reader := bufio.NewReader(resp.Body)
			for {
				line, err := reader.ReadBytes('\n')
				lineStr := string(bytes.TrimSpace(line))
				if strings.HasPrefix(lineStr, "data:") {
					lineStr = strings.TrimSpace(lineStr[5:])
					if lineStr == "" || lineStr == "[DONE]" {
						continue
					}
					err := consumer(lineStr)
					if err != nil {
						return nil, err
					}
				}
				if err == io.EOF {
					break
				} else if err != nil {
					return nil, err
				}
			}
		} else if strings.Contains(contentType, "application/json") {
			data, err := io.ReadAll(resp.Body)
			if err != nil {
				return nil, err
			}
			text = string(data)
		}
	}
	return &HttpResponse{status: resp.StatusCode, text: text, headers: resp.Header, contentType: contentType}, nil
}

func PostJson(url string, headers map[string]string, body string) (*HttpResponse, error) {
	bodyReader := strings.NewReader(body)
	req, err := http.NewRequest("POST", url, bodyReader)
	if err != nil {
		return nil, err
	}
	if headers != nil {
		for key, value := range headers {
			req.Header.Set(key, value)
		}
	}
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	contentType := resp.Header.Get("Content-Type")
	defer resp.Body.Close()
	var text string
	if resp.StatusCode == 200 {
		data, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		text = string(data)
	}
	return &HttpResponse{status: resp.StatusCode, text: text, headers: resp.Header, contentType: contentType}, nil
}

func Download(url string, path string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	out, err := os.Create(path)
	if err != nil {
		return err
	}
	defer out.Close()
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}
	return nil
}
