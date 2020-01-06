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

func HttpHead(url string, timeout time.Duration) (body string, contentType string, intervalTime int64, errReturn error) {
	startTime := time.Now()
	intervalTime = 0
	contentType = ""
	body = ""
	errReturn = nil

	resp, err := GetHttpClient(timeout).Head(url)
	if err != nil {
		intervalTime = int64(time.Now().Sub(startTime) / time.Millisecond)
		errReturn = err
		return
	}
	defer resp.Body.Close()

	bytebody, err := ioutil.ReadAll(resp.Body)
	intervalTime = int64(time.Now().Sub(startTime) / time.Millisecond)
	if err != nil {
		intervalTime = int64(time.Now().Sub(startTime) / time.Millisecond)
		errReturn = err
		return
	}
	body = string(bytebody)
	contentType = resp.Header.Get("Content-Type")
	intervalTime = int64(time.Now().Sub(startTime) / time.Millisecond)
	return
}

func HttpGet(url string, timeout time.Duration) (body string, contentType string, intervalTime int64, errReturn error) {
	startTime := time.Now()
	intervalTime = 0
	contentType = ""
	body = ""
	errReturn = nil

	resp, err := GetHttpClient(timeout).Get(url)
	if err != nil {
		intervalTime = int64(time.Now().Sub(startTime) / time.Millisecond)
		errReturn = err
		return
	}
	defer resp.Body.Close()

	bytebody, err := ioutil.ReadAll(resp.Body)
	intervalTime = int64(time.Now().Sub(startTime) / time.Millisecond)
	if err != nil {
		intervalTime = int64(time.Now().Sub(startTime) / time.Millisecond)
		errReturn = err
		return
	}
	body = string(bytebody)
	contentType = resp.Header.Get("Content-Type")
	intervalTime = int64(time.Now().Sub(startTime) / time.Millisecond)
	return
}

func HttpPost(url string, postBody string, bodyType string, timeout time.Duration) (body string, contentType string, intervalTime int64, errReturn error) {
	startTime := time.Now()
	intervalTime = 0
	contentType = ""
	body = ""
	errReturn = nil
	postBytes := bytes.NewBuffer([]byte(postBody))
	if bodyType == "" {
		bodyType = "application/x-www-form-urlencoded"
	}
	resp, err := GetHttpClient(timeout).Post(url, bodyType, postBytes)
	if err != nil {
		intervalTime = int64(time.Now().Sub(startTime) / time.Millisecond)
		errReturn = err
		return
	}
	defer resp.Body.Close()

	byteBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		intervalTime = int64(time.Now().Sub(startTime) / time.Millisecond)
		errReturn = err
		return
	}
	body = string(byteBody)
	contentType = resp.Header.Get("Content-Type")
	intervalTime = int64(time.Now().Sub(startTime) / time.Millisecond)
	return

}

//从指定query集合获取指定key的值
func GetQuery(querys url.Values, key string) string {
	if len(querys[key]) > 0 {
		return querys[key][0]
	}
	return ""
}
