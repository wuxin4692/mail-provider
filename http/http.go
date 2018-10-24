package http

import (
	"log"
	"net/http"
	_ "net/http/pprof"

	"github.com/wuxin4692/mail-provider/config"
)

func init() {
	configCommonRoutes()
	configProcRoutes()
}

func Start() {
	addr := config.Config().Http.Listen
	if addr == "" {
		return
	}
	s := &http.Server{
		Addr:           addr,
		MaxHeaderBytes: 1 << 20, // 请求的头域最大长度
	}
	log.Println("http listening", addr)
	log.Fatalln(s.ListenAndServe())
}
