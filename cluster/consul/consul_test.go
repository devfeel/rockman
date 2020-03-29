package consul

import (
	"fmt"
	"github.com/devfeel/rockman/packets"
	"testing"
	"time"
)

var consulServer = "116.62.16.66:8500"

func TestCreateLocker(t *testing.T) {
	lockKey := "locker_rockman"
	client, err := NewConsulClient(consulServer)
	if err != nil {
		t.Error("create consul client error", err)
		return
	}

	for i := 0; i < 3; i++ {
		go func(i int) {
			lock, err := client.CreateLocker(lockKey)
			if err != nil {
				fmt.Println(time.Now(), i, "create lock err", err)
				return
			}
			_, err = lock.Locker.Lock(nil)
			if err != nil {
				fmt.Println(time.Now(), i, "lock err", err)
			} else {
				fmt.Println(time.Now(), i, "lock success")
				time.Sleep(time.Minute)
				err = lock.Locker.Unlock()
				fmt.Println(time.Now(), i, "unlock success")
			}
		}(i)
	}
	time.Sleep(time.Hour)
}

func TestConsulClient_ListKV(t *testing.T) {
	client, err := NewConsulClient(consulServer)
	if err != nil {
		t.Error("create consul client error", err)
		return
	}
	nodeKVs, _, err := client.ListKV(packets.NodeKeyPrefix)
	if err != nil {
		fmt.Println("RefreshNodes error: " + err.Error())
	}
	for _, s := range nodeKVs {
		fmt.Println(s.Key, string(s.Value), s.Session)
	}
}
