package main

import (
	"net/smtp"
	"log"
	"crypto/tls"

	"github.com/google/gxui"
	"github.com/google/gxui/drivers/gl"
	"github.com/google/gxui/samples/flags"
)

func appMain(driver gxui.Driver) {
	theme := flags.CreateTheme(driver)

	from := theme.CreateTextBox()
	from.SetDesiredWidth(500)
	from.SetText("from")
	from.OnGainedFocus(func() {
		if from.Text() == "from" {
			from.SetText("")
		}
		})
	from.OnLostFocus(func() {
		if from.Text() == "" {
			from.SetText("from")
		}
		})

	password := theme.CreateTextBox()
	password.SetDesiredWidth(500)
	password.SetText("password")
	password.OnGainedFocus(func() {
		if password.Text() == "password" {
			password.SetText("")
		}
		})
	password.OnLostFocus(func() {
		if password.Text() == "" {
			password.SetText("password")
		}
		})

	to := theme.CreateTextBox()
	to.SetDesiredWidth(500)
	to.SetText("to")
	to.OnGainedFocus(func() {
		if to.Text() == "to" {
			to.SetText("")
		}
		})
	to.OnLostFocus(func() {
		if to.Text() == "" {
			to.SetText("to")
		}
		})

	host := theme.CreateTextBox()
	host.SetDesiredWidth(500)
	host.SetText("host")
	host.OnGainedFocus(func() {
		if host.Text() == "host" {
			host.SetText("")
		}
		})
	host.OnLostFocus(func() {
		if host.Text() == "" {
			host.SetText("host")
		}
		})

	port := theme.CreateTextBox()
	port.SetDesiredWidth(500)
	port.SetText("port")
	port.OnGainedFocus(func() {
		if port.Text() == "port" {
			port.SetText("")
		}
		})
	port.OnLostFocus(func() {
		if port.Text() == "" {
			port.SetText("port")
		}
		})

	result := theme.CreateLabel()
	result.SetText("")

	message := theme.CreateTextBox()
	message.SetDesiredWidth(500)
	message.SetMultiline(true)
	message.SetText("message")
	message.OnGainedFocus(func(){
		if message.Text() == "message" {
			message.SetText("")
		}
		})
	message.OnLostFocus(func() {
		if message.Text() == "" {
			message.SetText("message")
		}
		})

	sendButton := theme.CreateButton()
	sendButton.SetText("Send")
	sendButton.OnClick(func(ev gxui.MouseEvent) {
		result.SetText("Processing...")
		_from := from.Text()
		_pass := password.Text()
		_to := to.Text()
		_host := host.Text()
		_port := port.Text()

		_msg := "From: " + _from + "\r\n" +
			"To: " + _to + "\r\n" +
			"Subject: Hello there\r\n\r\n" +
			message.Text()+"\r\n"

	    auth := smtp.PlainAuth("",_from, _pass, _host)

	    tlsconfig := &tls.Config {
	        InsecureSkipVerify: true,
	        ServerName: _host,
	    }

	    conn, err := tls.Dial("tcp", _host+":"+_port, tlsconfig)
	    if err != nil {
	        result.SetText("Error:" + err.Error())
	        return
	    }
	    log.Printf("dial fine")

	    c, err := smtp.NewClient(conn, _host)
	    if err != nil {
	        result.SetText("Error:" + err.Error())
	        return
	    }
	    log.Printf("client fine")

	    if err = c.Auth(auth); err != nil {
	        result.SetText("Error:" + err.Error())
	        return
	    }
	    log.Printf("auth fine")

	    if err = c.Mail(_from); err != nil {
	        result.SetText("Error:" + err.Error())
	        return
	    }
	    log.Printf("from fine")

	    if err = c.Rcpt(_to); err != nil {
	        result.SetText("Error:" + err.Error())
	        return
	    }
	    log.Printf("to fine")

	    w, err := c.Data()
	    if err != nil {
	        result.SetText("Error:" + err.Error())
	        return
	    }
	    log.Printf("data fine")

	    _, err = w.Write([]byte(_msg))
	    if err != nil {
	        result.SetText("Error:" + err.Error())
	        return
	    }
	    log.Printf("write fine")

	    err = w.Close()
	    if err != nil {
	        result.SetText("Error:" + err.Error())
	        return
	    }
	    log.Printf("close fine")

	    c.Quit()
		log.Printf("fine")
		result.SetText("Message sent")
		return
	})

	layout := theme.CreateLinearLayout()
	layout.SetSizeMode(gxui.Fill)
	layout.AddChild(from)
	layout.AddChild(password)
	layout.AddChild(to)
	layout.AddChild(host)
	layout.AddChild(port)
	layout.AddChild(message)
	layout.AddChild(sendButton)
	layout.AddChild(result)

	layout.SetHorizontalAlignment(gxui.AlignCenter)

	scroll := theme.CreateScrollLayout()
	scroll.SetScrollAxis(false, true)
	scroll.SetChild(layout)


	window := theme.CreateWindow(800, 600, "SMTP sender")
	window.SetScale(flags.DefaultScaleFactor)
	window.AddChild(scroll)
	window.OnClose(driver.Terminate)
}

func main() {
	gl.StartDriver(appMain)
}
