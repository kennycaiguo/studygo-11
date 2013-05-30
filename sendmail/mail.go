package main


import(
	"github.com/astaxie/beego"
	"encoding/base64"
	"strings"
	"net/mail"
	"net/smtp"
	"time"
	"flag"
	"fmt"
)

var defaultFrom mail.Address
var b64 *base64.Encoding 

type Msg struct{
	To string
	Subject string
	Body string
}

var msgChan chan *Msg
var timeOutChan chan bool

var (
	paramHost = flag.String("host", "", "email server ip")
	paramPort = flag.String("port", "25", "email server port")
	paramFromEmail = flag.String("from", "service@tlt.cn", "from")
	paramFromName = flag.String("name", "TLT", "from")
)

/***********/

func timeOut(){
	time.Sleep(1e9 * 10)
	timeOutChan <- true
}

func sendMail(){
	for{
		msg := <- msgChan
		fmt.Printf("to=%s,subject=%s\n", msg.To, msg.Subject)
		header := make(map[string]string)
		header["From"] = defaultFrom.String()
		header["To"] = msg.To
		header["Subject"] = fmt.Sprintf("=?UTF-8?B?%s?=", b64.EncodeToString([]byte(msg.Subject)))
		header["MIME-Version"] = "1.0"
		header["Content-Type"] = "text/html; charset=UTF-8"
		header["Content-Transfer-Encoding"] = "base64"

		message := ""

		for k, v := range header {
			message += fmt.Sprintf("%s: %s\r\n", k, v)
		}
		message += "\r\n" + b64.EncodeToString([]byte(msg.Body))
		//auth := smtp.PlainAuth("", "", "", host)
		//err := smtp.SendMail(host, auth, "service@tlt.cn", []string{to.Address}, []byte(message))
		err := smtp.SendMail(*paramHost+":"+*paramPort, nil, *paramFromEmail, []string{msg.To}, []byte(message))
		if err == nil{
			fmt.Println("send success")
		}else{
			fmt.Println(err)
		}
	}
}

/***********/

type BaseController struct{
	beego.Controller
}

func (this *BaseController)Write(str string){
	this.Ctx.WriteString(str)
}

func (this *BaseController)Param(key string, defaultValue string) string {
	input := this.Input()
	//fmt.Printf("%v\n", input)
	value := input.Get(key)
	if value == "" {
		value = defaultValue
	}
	return value
}

func (this *BaseController)GetParam(key string, defaultValue string) string{
	return this.Param(key, defaultValue)
}

func (this *BaseController)PostParam(key string, defaultValue string) string{
	return this.Param(key, defaultValue)
}


/*************/

type SendMailController struct{
	BaseController
}

func (this *SendMailController)Post(){
	to := this.PostParam("to", "")
	toLen := len(to)
	subject := this.PostParam("subject", "")
	body := this.PostParam("body", "")
	if to=="" {
		this.Write("mailto is empty, to=" + to)
		return
	}
	toFlag := strings.Index(to, "@")
	if toFlag<=0 || toFlag>=toLen-1 {
		this.Write("mailto is error, to=" + to)
		return
	}

	msg := &Msg{to, subject, body}
	select{
		case msgChan <- msg :
			this.Write("ok")
		case <- timeOutChan:
			this.Write("timeout")
		default:
			this.Write("error")
	}
}

func main(){

	flag.Parse()
	if *paramHost== ""{
		flag.PrintDefaults()
		return
	}
	defaultFrom = mail.Address{*paramFromName, *paramFromEmail}
	msgChan = make(chan *Msg, 5000)
	timeOutChan = make(chan bool)
	b64 = base64.NewEncoding("ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/")

	go sendMail()
	go sendMail()
	go sendMail()
	go sendMail()
	go sendMail()
	go timeOut()

	beego.RegisterController("/", &SendMailController{})
	beego.Run()
}
