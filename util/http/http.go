package _http

import (
	"bytes"
	"crypto/tls"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"time"
)

const (
	MaxIdleConns        int = 100
	MaxIdleConnsPerHost int = 50
	IdleConnTimeout     int = 90
	InsecureSkipVerify      = true
)

type HttpResult struct {
	Status       string
	Body         string
	ContentType  string
	IntervalTime time.Duration
	Error        error
}

func getTransport(timeout time.Duration) *http.Transport {
	return &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		DialContext: (&net.Dialer{
			Timeout: timeout,
		}).DialContext,
		MaxIdleConns:        MaxIdleConns,
		MaxIdleConnsPerHost: MaxIdleConnsPerHost,
		IdleConnTimeout:     time.Duration(IdleConnTimeout) * time.Second,
		TLSClientConfig:     &tls.Config{InsecureSkipVerify: InsecureSkipVerify},
	}
}

func GetHttpClient(timeout time.Duration) *http.Client {
	return &http.Client{Transport: getTransport(timeout)}
}

func HttpHead(url string, timeout time.Duration) *HttpResult {
	startTime := time.Now()
	result := &HttpResult{}

	resp, err := GetHttpClient(timeout).Head(url)
	if err != nil {
		result.IntervalTime = time.Now().Sub(startTime)
		result.Error = err
		return result
	}
	defer resp.Body.Close()

	result.Status = resp.Status

	byteBody, err := ioutil.ReadAll(resp.Body)
	result.IntervalTime = time.Now().Sub(startTime)
	if err != nil {
		result.IntervalTime = time.Now().Sub(startTime)
		result.Error = err
		return result
	}
	result.Body = string(byteBody)
	result.ContentType = resp.Header.Get("Content-Type")
	result.IntervalTime = time.Now().Sub(startTime)
	return result
}

func HttpGet(url string, timeout time.Duration) *HttpResult {
	startTime := time.Now()
	result := &HttpResult{}

	resp, err := GetHttpClient(timeout).Get(url)
	if err != nil {
		result.IntervalTime = time.Now().Sub(startTime)
		result.Error = err
		return result
	}
	defer resp.Body.Close()

	result.Status = resp.Status

	byteBody, err := ioutil.ReadAll(resp.Body)
	result.IntervalTime = time.Now().Sub(startTime)
	if err != nil {
		result.IntervalTime = time.Now().Sub(startTime)
		result.Error = err
		return result
	}
	result.Body = string(byteBody)
	result.ContentType = resp.Header.Get("Content-Type")
	result.IntervalTime = time.Now().Sub(startTime)
	return result
}

func HttpPost(url string, postBody string, bodyType string, timeout time.Duration) *HttpResult {
	startTime := time.Now()
	result := &HttpResult{}

	postBytes := bytes.NewBuffer([]byte(postBody))
	if bodyType == "" {
		bodyType = "application/x-www-form-urlencoded"
	}
	resp, err := GetHttpClient(timeout).Post(url, bodyType, postBytes)
	if err != nil {
		result.IntervalTime = time.Now().Sub(startTime)
		result.Error = err
		return result
	}
	defer resp.Body.Close()

	result.Status = resp.Status

	byteBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		result.IntervalTime = time.Now().Sub(startTime)
		result.Error = err
		return result
	}
	result.Body = string(byteBody)
	result.ContentType = resp.Header.Get("Content-Type")
	result.IntervalTime = time.Now().Sub(startTime)
	return result

}

//从指定query集合获取指定key的值
func GetQuery(querys url.Values, key string) string {
	if len(querys[key]) > 0 {
		return querys[key][0]
	}
	return ""
}
