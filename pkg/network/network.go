package network

import (
	"errors"
	"net"
)

func Connect(server string, port string) (*net.TCPConn, error) {

	addr, err := net.ResolveTCPAddr("tcp", server+":"+port)
	if err != nil {
		return nil, errors.Join(errors.New("resolve address"), err)
	}
	tcpConn, err := net.DialTCP("tcp", nil, addr)
	if err != nil {
		return nil, errors.Join(errors.New("dial tcp"), err)
	}

	return tcpConn, nil
}
