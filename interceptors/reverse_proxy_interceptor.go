package interceptors

import (
	"net/http/httputil"
	"github.com/GetSimpl/go-simpl/logger"
	"github.com/GetSimpl/proxy-server/services"
	"net/url"
	"net/http"
	"bytes"
	"encoding/json"
	"fmt"
)

type ReverseProxyInterceptor struct {
	Proxy *httputil.ReverseProxy
}

func NewProxy(urlString string) *ReverseProxyInterceptor{
	uri, err := url.Parse(urlString)
	if err != nil {
		logger.E("Crashing! ")
		panic("Invalid URL string for Target.")
	}

	proxy := httputil.NewSingleHostReverseProxy(uri)
	return &ReverseProxyInterceptor{ Proxy: proxy }
}


func (proxyInterceptor *ReverseProxyInterceptor) ServeHTTP (w http.ResponseWriter, r *http.Request){
	logRequest(r)
	// Add Rules Engine here
	// Async Updating of Redis for Rate Limiting and other logic
	rulesEngineChecker := services.NewRulesEngineChecker()
	checker := services.IsIPValid(r, rulesEngineChecker)

	status, errors := checker.Success()

	fmt.Println(errors)

	if status != true {
		responseHash := make(map[string]interface{})
		responseHash["success"] = false
		responseHash["errors"] = checker.Errors
		responseHash["api_version"] = 1.0
		jsonResponse, err := json.Marshal(responseHash)
		if err != nil {
			logger.E(err)
			panic(err)
		}
		w.Header().Set("Content-Type","application/json")
		w.Write([]byte(jsonResponse))
		return;
	}else{
		proxyInterceptor.Proxy.ServeHTTP(w, r)
	}
}

func headersAsString(headers http.Header) string {
	var headerParams bytes.Buffer
	for k, v := range headers {
		headerParams.WriteString(k)
		headerParams.WriteString("=")
		headerParams.WriteString(v[0])
		headerParams.WriteString(" ")
	}
	return headerParams.String()
}

func logRequest(r *http.Request){
	headers := headersAsString(r.Header)
	logParams := "Request: " + r.RequestURI + " Headers:" + headers
	logger.I(logParams)
}
