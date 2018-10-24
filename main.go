package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"runtime"

	"github.com/wuxin4692/mail-provider/config"
	"github.com/wuxin4692/mail-provider/http"
)

func prepare() {
	//GOMAXPROCS设置可同时执行的最大CPU数，并返回先前的设置。 若 n < 1，它就不会更改当前设置。本地机器的逻辑CPU数可通过 NumCPU 查询。
	runtime.GOMAXPROCS(runtime.NumCPU())
	//初始化日志格式　19
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
}

func init() {
	prepare()

	cfg := flag.String("c", "cfg.json", "配置文件")
	version := flag.Bool("v", false, "查看版本")
	help := flag.Bool("h", false, "帮助")
	flag.Parse()

	handleVersion(*version)
	handleHelp(*help)
	handleConfig(*cfg)
}

func main() {
	http.Start()
}

func handleVersion(displayVersion bool) {
	if displayVersion {
		fmt.Println(config.VERSION)
		os.Exit(0)
	}
}

func handleHelp(displayHelp bool) {
	if displayHelp {
		flag.Usage()
		os.Exit(0)
	}
}

func handleConfig(configFile string) {
	err := config.Parse(configFile)
	if err != nil {
		log.Fatalln(err)
	}
}
