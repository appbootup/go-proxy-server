package interceptors

import (
	"net/http/httputil"
	"github.com/GetSimpl/go-simpl/logger"
	"net/url"
	"net/http"
	"bytes"
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
	proxyInterceptor.Proxy.ServeHTTP(w, r)
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
