package network

import (
	"net"
)

const (
	connType  = "udp"
	DefaultIP = "0.0.0.0"
)

func NewUDPConnection(connHost string, connPort string) (*net.UDPConn, error) {
	udpAddr, err := net.ResolveUDPAddr(connType, connHost+":"+connPort)
	if err != nil {
		return nil, err
	}

	conn, err := net.DialUDP(connType, nil, udpAddr)
	if err != nil {
		return nil, err
	}

	return conn, nil
}

func NewUDPListenConn(ipAddr string, port int) (*net.UDPConn, error) {
	ip := net.ParseIP(ipAddr)
	listenConn, err := net.ListenUDP(connType, &net.UDPAddr{IP: ip, Port: port})
	if err != nil {
		return nil, err
	}

	return listenConn, nil
}
