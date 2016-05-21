package services

import (
	"net/http"
	"fmt"
	"net"
)

func IsIPValid(r *http.Request, checker *RulesEngineChecker) *RulesEngineChecker{

	fmt.Println(r.RemoteAddr)
	host, _, err := net.SplitHostPort(r.RemoteAddr)

	if err != nil {
		checker.AddErrors("Invalid Host")
	}

	if (host != "127.0.0.1") {
		checker.AddErrors("Invalid IP Address")
	}

	return checker;
}
