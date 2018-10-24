package config

import (
	"encoding/json"
	"fmt"
	"log"
	"sync"

	"io/ioutil"
)

type HttpConfig struct {
	Listen string `json:"listen"`
	Token  string `json:"token"`
}

type SmtpConfig struct {
	Addr     string `json:"addr"`
	Username string `json:"username"`
	Password string `json:"password"`
	From     string `json:"from"`
}

type GlobalConfig struct {
	Debug bool        `json:"debug"`
	Http  *HttpConfig `json:"http"`
	Smtp  *SmtpConfig `json:"smtp"`
}

var (
	ConfigFile string
	config     *GlobalConfig
	configLock = new(sync.RWMutex)
)

//读锁
func Config() *GlobalConfig {
	configLock.RLock()
	defer configLock.RUnlock()
	return config
}

func Parse(cfg string) error {
	//如果cfg为空，返回错误信息
	if cfg == "" {
		return fmt.Errorf("使用 -c 参数来指定配置文件")
	}
	//如果ｃｆｇ路径的文件不存在则返回错误信息
	if !file.IsExist(cfg) {
		return fmt.Errorf("文件 %s 不存在", cfg)
	}

	ConfigFile = cfg
	//读取文件内容并转换为ｓｔｒｉｎｇ类型
	b, err := ioutil.ReadFile(cfg)
	if err != nil {
		return fmt.Errorf("读取文件 %s 错误: %s", cfg, err.Error())
	}
	ｓｔｒ := string(b)
	//去掉空格
	cfg_str := strings.TrimSpace(str)
	//将上面步骤读取的配置文件中的json反序列化并存到　ＧｌｏｂａｌＣｏｎｆｉｇ　结构体中
	var c GlobalConfig
	err = json.Unmarshal([]byte(ｃｆｇ_str), &c)
	if err != nil {
		return fmt.Errorf("反序列化 %s 错误: %s", cfg, err.Error())
	}
	//锁定写操作
	configLock.Lock()
	//跳转到函数最后并解除锁定
	defer configLock.Unlock()
	config = &c

	log.Println("读取配置文件", cfg, "成功....")
	return nil
}
