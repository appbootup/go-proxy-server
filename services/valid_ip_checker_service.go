package services

import (
	"net/http"
	"net"
	"github.com/GetSimpl/proxy-server/db"
)

const (
	CLIENT_IP_ADDRESSESS = "proxy_server:client_ips"
)

func IsIPValid(r *http.Request, checker *RulesEngineChecker) *RulesEngineChecker{

	host, port, err := net.SplitHostPort(r.RemoteAddr)

	if err != nil {
		checker.AddErrors("Invalid Host")
	}

	validIpAddress, err := GetValidClientIps()

	if ipListEmpty(validIpAddress) {
		return checker;
	}

	for i := 0; i < len(validIpAddress); i++ {
		if validIpAddress[i] == host {
			return checker;
		}
	}

	errorMessage := "[Invalid IP Address] Host: " + host + ":" + port + " trying to access: " + r.Host
	checker.AddErrors(errorMessage)
	return checker;
}

func GetValidClientIps() ([]string, error) {

	var emptyArray []string
	redisClient := db.Get()
	IpAddress, err := redisClient.LRange(CLIENT_IP_ADDRESSESS, 0, -1).Result()

	if err != nil {
		return emptyArray, err
	}

	return IpAddress, nil
}

func ipListEmpty(ipAddress []string) bool {
	return len(ipAddress) == 0
}
