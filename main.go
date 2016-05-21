package main

import (
	"net/http"
	"fmt"
	"runtime"
	"github.com/GetSimpl/go-simpl/logger"
	"github.com/GetSimpl/proxy-server/interceptors"

)

var port string = ":9000"
var proxy *interceptors.ReverseProxyInterceptor = interceptors.NewProxy("http://localhost:4000")

func main(){
	logger.Init(1)
	runtime.GOMAXPROCS(20)
	fmt.Println("Starting app on port 9000")
	http.ListenAndServe(port, proxy)
}


