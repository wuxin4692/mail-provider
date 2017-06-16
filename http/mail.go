package http

import (
        "fmt"
	"net/http"
	"strings"
	"github.com/open-falcon/mail-provider/config"
	//"github.com/toolkits/smtp"
	"github.com/toolkits/web/param"
        //"net"
        "net/mail"
        "net/smtp"
        "crypto/tls"
        "log"
)

func configProcRoutes() {

	http.HandleFunc("/sender/mail", func(w http.ResponseWriter, r *http.Request) {
		cfg := config.Config()
		token := param.String(r, "token", "")
		if cfg.Http.Token != token {
			http.Error(w, "no privilege", http.StatusForbidden)
			return
		}

		tos := param.MustString(r, "tos")
		subject := param.MustString(r, "subject")
		content := param.MustString(r, "content")
		tos = strings.Replace(tos, ",", ";", -1)

		//s := smtp.New(cfg.Smtp.Addr, cfg.Smtp.Username, cfg.Smtp.Password)
		//err := s.SendMail(cfg.Smtp.From, tos, subject, content)
		//if err != nil {
		//	http.Error(w, err.Error(), http.StatusInternalServerError)
		//} else {
		//	http.Error(w, "success", http.StatusOK)
		//}

// 参考 https://gist.github.com/chrisgillis/10888032 ， 对上面一段的smtp不加密传输进行弹头改装，如下。
// 测试的时候，使用的是fastmail，注意两点，一个是生成一个app使用的发送的密码，另外一个就是要使用端口587（starttls，使用465不work。）
// 参考 https://www.fastmail.com/help/technical/servernamesandports.html#email
    from := mail.Address{"", cfg.Smtp.From}
    to   := mail.Address{"", tos}
    subj := subject
    body := content

    // Setup headers
    headers := make(map[string]string)
    headers["From"] = from.String()
    headers["To"] = to.String()
    headers["Subject"] = subj

    // Setup message
    message := ""
    for k,v := range headers {
        message += fmt.Sprintf("%s: %s\r\n", k, v)
    }
    message += "\r\n" + body

    auth := smtp.PlainAuth("",cfg.Smtp.Username, cfg.Smtp.Password, cfg.Smtp.Addr)

    // TLS config
    tlsconfig := &tls.Config {
        InsecureSkipVerify: false,
        ServerName: cfg.Smtp.Addr,
    }

    // Here is the key, you need to call tls.Dial instead of smtp.Dial
    // for smtp servers running on 465 that require an ssl connection
    // from the very beginning (no starttls)
    conn, err := tls.Dial("tcp", cfg.Smtp.Addr, tlsconfig)
    if err != nil {
        log.Panic(err)
    }

    c, err := smtp.NewClient(conn, cfg.Smtp.Addr)
    if err != nil {
        log.Panic(err)
    }

    // Auth
    if err = c.Auth(auth); err != nil {
        log.Panic(err)
    }

    // To && From
    if err = c.Mail(cfg.Smtp.From); err != nil {
        log.Panic(err)
    }

    if err = c.Rcpt(tos); err != nil {
        log.Panic(err)
    }

    // Data
    //w, err = c.Data()
    p, err := c.Data()
    if err != nil {
        log.Panic(err)
    }

    _, err = p.Write([]byte(message))
    if err != nil {
        log.Panic(err)
    }

    err = p.Close()
    if err != nil {
        log.Panic(err)
    }

    c.Quit()







                http.Error(w,"complete!",http.StatusOK)
	})

}
