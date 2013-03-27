package main


import(
	"github.com/astaxie/beego"
	"strings"
	"net/mail"
	"time"
	"fmt"
)

var defaultFromEmailAddr string
var defaultFromName string
var defaultFrom mail.Address

type Msg struct{
	To string
	Subject string
	Body string
}

var msgChan chan *Msg
var timeOutChan chan bool


func init(){
	defaultFromEmailAddr = "service@tlt.cn"
	defaultFromName = "太灵通"
	defaultFrom = mail.Address{defaultFromName, defaultFromEmailAddr}
	msgChan = make(chan *Msg, 5000)
	timeOutChan = make(chan bool)
}

/***********/

func timeOut(){
	time.Sleep(1e9 * 10)
	timeOutChan <- true
}

func sendMail(){
	for{
		msg := <- msgChan
		fmt.Printf("to=%s,subject=%s,body=%s\n", msg.To, msg.Subject, msg.Body)
		time.Sleep(1e9*20)
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
	go sendMail()
	go sendMail()
	go sendMail()
	go sendMail()
	go sendMail()
	go timeOut()
	beego.RegisterController("/sendmail", &SendMailController{})
	beego.Run()
}
