package http

import (
	"fmt"
	"net/http"
	"strings"

	"crypto/tls"
	"log"
	"net/mail"
	"net/smtp"

	"github.com/toolkits/web/param"
	"github.com/wuxin4692/mail-provider/config"
)

func configProcRoutes() {

	http.HandleFunc("/sender/mail", func(w http.ResponseWriter, r *http.Request) {
		cfg := config.Config()
		//用token来限制发送权限
		token := param.String(r, "token", "")
		if cfg.Http.Token != token {
			//如果请求中的token和ｊｓｏｎ中的不同返回４０３
			http.Error(w, "no privilege", http.StatusForbidden)
			return
		}
		//取请求参数
		//收件人邮箱
		tos := param.MustString(r, "tos")
		//邮件主题
		subject := param.MustString(r, "subject")
		//邮件正文
		content := param.MustString(r, "content")
		tos = strings.Replace(tos, ",", ";", -1)

		from := mail.Address{"", cfg.Smtp.From}
		to := mail.Address{"", tos}
		subj := subject
		body := content

		//-----------------------------修改　添加html模板的发送方式
		headers := make(map[string]string)
		headers["From"] = from.String()
		headers["To"] = to.String()
		headers["Subject"] = subj

		// Setup message
		message := ""
		for k, v := range headers {
			message += fmt.Sprintf("%s: %s\r\n", k, v)
		}
		message += "\r\n" + body

		auth := smtp.PlainAuth("", cfg.Smtp.Username, cfg.Smtp.Password, cfg.Smtp.Addr)

		// 建立ＴＬＳ加密链接
		tlsconfig := &tls.Config{
			// InsecureSkipVerify控制客户端是否认证服务端的证书链和主机名。
			// 如果InsecureSkipVerify为真，TLS连接会接受服务端提供的任何证书和该证书中的任何主机名。
			// 此时，TLS连接容易遭受中间人攻击，这种设置只应用于测试。
			InsecureSkipVerify: false,
			// 用于认证返回证书的主机名（除非设置了InsecureSkipVerify）。
			// 也被用在客户端的握手里，以支持虚拟主机。
			ServerName: cfg.Smtp.Addr,
		}

		// Here is the key, you need to call tls.Dial instead of smtp.Dial
		// for smtp servers running on 465 that require an ssl connection
		// from the very beginning (no starttls)
		//Dial使用net.Dial连接指定的网络和地址，然后发起TLS握手，返回生成的TLS连接。Dial会将nil的配置视为零值的配置；
		conn, err := tls.Dial("tcp", cfg.Smtp.Addr, tlsconfig)
		if err != nil {
			log.Panic(err)
		}
		//NewClient使用已经存在的连接conn和作为服务器名的host（用于身份认证）来创建一个*Client。
		c, err := smtp.NewClient(conn, cfg.Smtp.Addr)
		if err != nil {
			log.Panic(err)
		}

		// Auth使用提供的认证机制进行认证。失败的认证会关闭该连接。只有服务端支持AUTH时，本方法才有效。（但是不支持时，调用会默默的成功）
		if err = c.Auth(auth); err != nil {
			log.Panic(err)
		}

		// Mail发送MAIL命令和邮箱地址from到服务器。如果服务端支持8BITMIME扩展，本方法会添加BODY=8BITMIME参数。方法初始化一次邮件传输，后应跟1到多个Rcpt方法的调用。
		if err = c.Mail(cfg.Smtp.From); err != nil {
			log.Panic(err)
		}
		//Rcpt发送RCPT命令和邮箱地址to到服务器。调用Rcpt方法之前必须调用了Mail方法，之后可以再一次调用Rcpt方法，也可以调用Data方法。
		if err = c.Rcpt(tos); err != nil {
			log.Panic(err)
		}

		// Data
		//Data发送DATA指令到服务器并返回一个io.WriteCloser，用于写入邮件信息。调用者必须在调用c的下一个方法之前关闭这个io.WriteCloser。方法必须在一次或多次Rcpt方法之后调用。
		p, err := c.Data()
		if err != nil {
			log.Panic(err)
		}
		//Write方法len(p) 字节数据从p写入底层的数据流。它会返回写入的字节数(0 <= n <= len(p))
		//和遇到的任何导致写入提取结束的错误。Write必须返回非nil的错误，如果它返回的 n < len(p)。Write不能修改切片p中的数据，即使临时修改也不行。
		_, err = p.Write([]byte(message))
		if err != nil {
			log.Panic(err)
		}
		//关闭io.WriteCloser
		err = p.Close()
		if err != nil {
			log.Panic(err)
		}
		//Quit发送QUIT命令并关闭到服务端的连接。
		c.Quit()

		http.Error(w, "complete!", http.StatusOK)
	})

}
