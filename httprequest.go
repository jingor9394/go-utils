package utils

import (
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

type HttpRequest struct {
	timeout int
}

func NewHttpRequest(timeout int) *HttpRequest {
	return &HttpRequest{timeout}
}

func (r *HttpRequest) HttpRequest(reqUrl string, method string, data string, headers map[string]string) ([]byte, error) {
	body := strings.NewReader(data)
	req, err := http.NewRequest(method, reqUrl, body)
	if err != nil {
		return nil, err
	}
	if headers != nil {
		for header, val := range headers {
			if header == "Host" {
				req.Host = val
			} else {
				req.Header.Add(header, val)
			}
		}
	}

	reqTimeout := time.Duration(r.timeout) * time.Second
	client := &http.Client{
		Timeout: reqTimeout,
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	ret, err := ioutil.ReadAll(resp.Body)
	return ret, err
}
