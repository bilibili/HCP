package ip

import (
	"net"
	"strconv"
	"strings"
)

var Ipv4Parser = ipv4Parser{}

type ipv4Parser struct {
}

func (ip4 *ipv4Parser) CheckIP(ipStr string) bool {
	address := net.ParseIP(ipStr)
	if address == nil {
		return false
	} else {
		return true
	}

}

func (ip4 *ipv4Parser) inetAToN(ipStr string) uint32 {
	bits := strings.Split(ipStr, ".")
	b0, _ := strconv.Atoi(bits[0])
	b1, _ := strconv.Atoi(bits[1])
	b2, _ := strconv.Atoi(bits[2])
	b3, _ := strconv.Atoi(bits[3])
	var sum uint32
	sum += uint32(b0) << 24
	sum += uint32(b1) << 16
	sum += uint32(b2) << 8
	sum += uint32(b3)
	return sum
}

func (ip4 *ipv4Parser) IsInnerIp(ipStr string) bool {
	defer func() {
		_ = recover()
	}()
	if !ip4.CheckIP(ipStr) {
		return false
	}
	inputIpNum := ip4.inetAToN(ipStr)
	innerIpA := ip4.inetAToN("10.255.255.255")
	innerIpB := ip4.inetAToN("172.16.255.255")
	innerIpC := ip4.inetAToN("192.168.255.255")
	return inputIpNum>>24 == innerIpA>>24 ||
		inputIpNum>>20 == innerIpB>>20 ||
		inputIpNum>>16 == innerIpC>>16
}
