package rpc

import "testing"

func TestRpcClient_CallEcho(t *testing.T) {
	client := NewRpcClient("127.0.0.1", "2398")
	message := "echo message"
	wantMessage := "echo message"
	err, result := client.CallEcho(message)
	if err != nil {
		t.Error(err)
	} else {
		if result == wantMessage {
			t.Log("success:", wantMessage, result)
		} else {
			t.Error("failed:", wantMessage, result)
		}
	}
}
