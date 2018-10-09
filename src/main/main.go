package main

import (
	"github.com/iwind/TeaWebCode/teaweb"
	"net/http"
	_ "net/http/pprof"
)

func main() {
	go http.ListenAndServe("0.0.0.0:8080", nil)

	teaweb.Start()
}
