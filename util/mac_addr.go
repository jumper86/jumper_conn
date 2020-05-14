package util

import (
	"net"

	"github.com/jumper86/jumper_conn/cst"
)

func GetMacAddr() (string, error) {
	inters, err := net.Interfaces()
	if err != nil {
		return "", cst.ErrGetMacAddr
	}

	for _, inter := range inters {
		macAddr := inter.HardwareAddr.String()
		if len(macAddr) != 0 {
			return macAddr, nil
		}
	}

	return "", nil
}

func GetMacAddrs() ([]string, error) {
	inters, err := net.Interfaces()
	if err != nil {
		return nil, cst.ErrGetMacAddr
	}

	macAddrs := make([]string, 0)
	for _, inter := range inters {
		macAddr := inter.HardwareAddr.String()
		if len(macAddr) != 0 {
			macAddrs = append(macAddrs, macAddr)
		}
	}

	return macAddrs, nil
}
