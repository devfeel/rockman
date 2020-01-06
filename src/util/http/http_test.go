package _http

import (
	"testing"
	"time"
)

func TestHttpGet(t *testing.T) {
	url := "http://www.baidu.com"
	_, _, _, err := HttpGet(url, 30*time.Second)
	if err != nil {
		t.Error(err)
	} else {
		t.Log("HttpGet success")
	}
}

func TestHttpPost(t *testing.T) {
	url := "http://www.baidu.com"
	_, _, _, err := HttpPost(url, "", "", 30*time.Second)
	if err != nil {
		t.Error(err)
	} else {
		t.Log("HttpPost success")
	}
}
