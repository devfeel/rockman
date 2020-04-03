package _http

import (
	"testing"
	"time"
)

func TestHttpGet(t *testing.T) {
	url := "http://www.baidu.com"
	result := HttpGet(url, 30*time.Second)
	if result.Error != nil {
		t.Error(result.Error)
	} else {
		t.Log("HttpGet success")
	}
}

func TestHttpPost(t *testing.T) {
	url := "http://www.baidu.com"
	result := HttpPost(url, "", "", 30*time.Second)
	if result.Error != nil {
		t.Error(result.Error)
	} else {
		t.Log("HttpPost success")
	}
}
