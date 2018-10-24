package http

import (
	"net/http"

	"github.com/wuxin4692/mail-provider/config"
)

func configCommonRoutes() {
	//健康检查页面　返回ok
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("ok"))
	})
	//版本查看页面　返回常量版本号
	http.HandleFunc("/version", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(config.VERSION))
	})
}
