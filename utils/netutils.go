package utils

import (
	"io"
	"net"
	"net/http"
	"time"
)

var tp = &http.Transport{
	DialContext: (&net.Dialer{
		KeepAlive: 10 * time.Minute,
	}).DialContext,
	ResponseHeaderTimeout: 60 * time.Second,
	MaxIdleConnsPerHost:   100,
	IdleConnTimeout:       60 * time.Second,
	TLSHandshakeTimeout:   10 * time.Second,
}
var client = &http.Client{
	Timeout:   10 * time.Second,
	Transport: tp,
}

func DoGet(url string, retry int) ([]byte, error) {
	resp, err := client.Get(url)
	if err != nil {
		if retry > 0 {
			time.Sleep(500 * time.Millisecond)
			return DoGet(url, retry-1)
		}
		return nil, err
	}
	defer resp.Body.Close()
	return io.ReadAll(resp.Body)
}
func DoPost(url, contentType string, reqBody io.Reader, retry int) ([]byte, error) {
	resp, err := client.Post(url, contentType, reqBody)
	if err != nil {
		if retry > 0 {
			time.Sleep(500 * time.Millisecond)
			return DoPost(url, contentType, reqBody, retry-1)
		}
		return nil, err
	}
	defer resp.Body.Close()
	return io.ReadAll(resp.Body)
}

func DoGetAndHeader(url string, retry int, headers map[string]string) ([]byte, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return []byte{}, err
	}
	for s := range headers {
		req.Header.Add(s, headers[s])
	}
	resp, err := client.Do(req)
	if err != nil {
		if retry > 0 {
			time.Sleep(500 * time.Millisecond)
			return DoGet(url, retry-1)
		}
		return nil, err
	}
	defer resp.Body.Close()
	return io.ReadAll(resp.Body)
}
func DoPostAndHeader(url string, reqBody io.Reader, headers map[string]string, retry int) ([]byte, error) {
	req, err := http.NewRequest("POST", url, reqBody)
	if err != nil {
		return []byte{}, err
	}
	for s := range headers {
		req.Header.Add(s, headers[s])
	}
	resp, err := client.Do(req)
	if err != nil {
		if retry > 0 {
			time.Sleep(500 * time.Millisecond)
			return DoPostAndHeader(url, reqBody, headers, retry-1)
		}
		return nil, err
	}
	defer resp.Body.Close()
	return io.ReadAll(resp.Body)
}
