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
	runtime.GOMAXPROCS(runtime.NumCPU())
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
