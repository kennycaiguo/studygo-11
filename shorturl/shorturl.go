package main

import(
	"sqlite"
	"sync"
	"errors"
	"net/http"
	"fmt"
	"strings"
)

var con *sqlite.Conn

func ConnectDb(filename string) error{
	_con, err := sqlite.Open(filename)
	con = _con
	return err
}

type RandNum struct{
	num int
	lock sync.Mutex
}
var randNum RandNum

type ShortUrl struct{
	id int
	url string
	num string
}

func (this ShortUrl)RandNum() (error){
	randNum.lock.Lock()
	defer randNum.lock.Unlock()
	rs, err := con.Prepare("select id from short_num limit 1")
	if err != nil{
		return err
	}
	defer rs.Finalize()
	rs.Exec()
	err = rs.Scan(&randNum.num)
	if err != nil{
		return err
	}
	con.Exec("update short_num set id=id+1 limit 1")
	this.num = IntToNum(randNum.num)
	return nil
}

func (this ShortUrl)Insert() error{
	return con.Exec("insert into short_url(url, num) values (?,?)", this.url, this.num)
}

func (this ShortUrl)Load() error{
	if this.num == "" {
		return errors.New("num is empty")
	}
	rs, err := con.Prepare("select * from short_url where url=?")
	if err != nil{
		return err
	}
	defer rs.Finalize()
	rs.Exec(this.url)
	if !rs.Next(){
		return errors.New("no data")
	}
	err = rs.Scan(&this.id, &this.url, &this.num)
	if err == nil{
		return nil
	}
	return err
}

const defaultHtml string = `
<!DOCTYPE html>
<html>
	<head>
		<meta http-equiv="Content-Type" content="text/html; charset=utf-8">
		<title>短网址服务</title>
	</head>
	<body style="color:#666;">
		<div style="width:700px;margin:136px auto 0;">
			<div style="font-size:24px;font-weight:bold;text-align:center;color:#999;margin-bottom:50px;">短网址服务</div>
			<div>
			<form method="post" action="/">
				<span style="font-weight:bold;font-size:24px;line-height38px;">网址:</span>
				<input type="text" name="url" style="width:500px;height:30px;line-height:30px;font-size:14px;"/>
				<input type="submit" name="submit" value="生成短网址" style="width:106px;height:38px;line-height:38px;text-align:center;font-size:16px;font-weight:bold;color:#666;"/>
			</form>
			</div>
		</div>
	</body>
</html>
`

func DefaultHandler(w http.ResponseWriter, r *http.Request){
	if r.Method == "GET"{
		fmt.Fprint(w, defaultHtml)
		return
	}
	if r.Method != "POST"{
		return
	}
	url := r.FormValue("url")
	if strings.Index(url, "http://") != 0{
		if strings.Index(url, "https://") != 0{
			return
		}
	}
	shortUrl := ShortUrl{0, url, ""}
	err := shortUrl.Load()
	if err != nil{
		fmt.Fprintf(w, "http://%s/%s", r.Host, shortUrl.num)
		return
	}
	shortUrl.num, _ = shortUrl.RandNum()
	shortUrl.url = url
	shortUrl.Insert()
	fmt.Fprintf(w, "http://%s/%s", r.Host, shortUrl.num)
}

func main(){
	err := ConnectDb("shorturl.sqlite")
	if err != nil{
		fmt.Println(err)
		return
	}
	defer con.Close()
	http.HandleFunc("/", DefaultHandler)
	err = http.ListenAndServe(":8080", nil)
	if err != nil{
		fmt.Println(err)
		return
	}
}

