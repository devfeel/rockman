package netx

import "net"

func CheckTcpConnect(endPoint string) bool {
	conn, err := net.Dial("tcp", endPoint)
	if err == nil {
		conn.Close()
	}
	return err == nil
}
