package main

import (
	"net/http"
	"fmt"
	"runtime"
	"github.com/GetSimpl/go-simpl/logger"
	"github.com/GetSimpl/proxy-server/proxies"
	"github.com/GetSimpl/proxy-server/db"
	"os"

)

var port string = os.Getenv("PORT")
var proxy *proxies.ReverseProxyInterceptor = proxies.NewProxy("http://localhost:4000")

func main(){
	logger.Init(1)
	db.Init()
	runtime.GOMAXPROCS(20)
	fmt.Println("Starting app on port 9000")
	http.ListenAndServe(port, proxy)
}


