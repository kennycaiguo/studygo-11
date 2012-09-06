package main

import(
	"sqlite"
	"sync"
	"errors"
	"net/http"
	"fmt"
	"strings"
	"regexp"
	"os"
	"flag"
)

var con *sqlite.Conn

func ConnectDb(filename string) error{
	_con, err := sqlite.Open(filename)
	con = _con
	return err
}

var lock sync.Mutex

type ShortUrl struct{
	id int
	url string
	num string
}

func (this *ShortUrl)GetNum() (error){
	lock.Lock()
	num := 0
	defer lock.Unlock()
	rs, err := con.Prepare("select id from short_num limit 1")
	if err != nil{
		return err
	}
	defer rs.Finalize()
	rs.Exec(); rs.Next()
	err = rs.Scan(&num)
	if err != nil{
		fmt.Println("b")
		return err
	}
	err = con.Exec("update short_num set id=id+1")
	//fmt.Println(err)
	this.num = IntToNum(num)
	return nil
}

func (this *ShortUrl)Insert() error{
	return con.Exec("insert into short_url(url, num) values (?,?)", this.url, this.num)
}

func (this *ShortUrl)LoadByUrl() error{
	if this.url == "" {
		return errors.New("url is empty")
	}
	rs, err := con.Prepare("select * from short_url where url=? limit 1")
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

func (this *ShortUrl)LoadByNum() error{
	if this.num== "" {
		return errors.New("num is empty")
	}
	rs, err := con.Prepare("select * from short_url where num=? limit 1")
	if err != nil{
		return err
	}
	defer rs.Finalize()
	rs.Exec(this.num)
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
	<div><span><a href='http://www.tlt.cn/' target='_blank' style='color:#999; font-size:13px;font-weight:normal;text-decoration:none;margin:10px;'>访问太灵通</a></span></div>
		<div style="width:700px;margin:80px auto 0;">
			<div style="font-weight:bold;text-align:center;color:#999;margin-bottom:50px;">短网址服务</div>
			<div>
			<form method="post" action="/">
				<span style="line-height38px;">网址:</span>
				<input type="text" name="url" style="width:500px;height:30px;line-height:30px;font-size:14px;border:solid 1px #666"/>
				<input type="submit" name="submit" value="生成短网址" style="height:38px;line-height:38px;text-align:center;color:#666;border:solid 1px #666"/>
			</form>
			</div>
			<div style="font-size:18px; text-align:center;font-weight:bold;margin-top:20px;">%s</div>
		</div>
		<div style='text-align:center;font-size:12px;position:absolute;bottom:0;right:0;'>
		Copyright © 2009 - 2012 太灵通. All Rights Reserved
		</div>
	</body>
</html>
`

/*****************路由处理*******************/
type route struct{
	pattern *regexp.Regexp
	handler http.Handler
}

type RegexpHandler struct{
	routes []*route
}

//func (this *RegexpHandler)Handler(pattern *regexp.Regexp, handler http.Handler){
func (this *RegexpHandler)Handler(s string, handler http.Handler){
	pattern, _:= regexp.Compile(s)
	this.routes = append(this.routes, &route{pattern, handler})
}

func (this *RegexpHandler)HandleFunc(s string, handler func(http.ResponseWriter, *http.Request)){
	pattern, _:= regexp.Compile(s)
	this.routes = append(this.routes, &route{pattern, http.HandlerFunc(handler)})
}

func (this *RegexpHandler)ServeHTTP(w http.ResponseWriter, r *http.Request){
	for _, route := range this.routes{
		if route.pattern.MatchString(r.URL.Path){
			route.handler.ServeHTTP(w, r)
			return
		}
	}
	http.NotFound(w, r)
}
/***********end**************/

func DefaultHandler(w http.ResponseWriter, r *http.Request){
	/*if r.URL.Path != "/"{
		UrlHandler(w, r)
		return
	}*/
	//rootUrl := r.Header["Referer"][0]
	rootUrl := "http://" + r.Host + "/"
	if r.Method == "GET"{
		fmt.Fprintf(w, defaultHtml, "")
		return
	}
	if r.Method != "POST"{
		return
	}
	url := r.FormValue("url")
	if strings.Index(url, "http://") != 0{
		if strings.Index(url, "https://") != 0{
			fmt.Fprintf(w, defaultHtml, "")
			return
		}
	}
	shortUrl := ShortUrl{0, url, ""}
	err := shortUrl.LoadByUrl()
	if err == nil{
		fmt.Fprintf(w, defaultHtml, rootUrl + shortUrl.num)
		return
	}
	err = shortUrl.GetNum()
	if err != nil{
		fmt.Println(err)
		return
	}
	shortUrl.url = url
	shortUrl.Insert()
	fmt.Fprintf(w, defaultHtml, rootUrl + shortUrl.num)
}

func UrlHandler(w http.ResponseWriter, r *http.Request){
	num := r.URL.Path
	num_len := len(num)
	if num_len == 0{
		fmt.Fprintf(w, defaultHtml, "地址有错")
		return
	}
	if(string(num[num_len-1]) == "/"){
		num = num[1:num_len-1]
	}else{
		num = num[1:num_len]
	}
	shortUrl := ShortUrl{0, "", num}
	err := shortUrl.LoadByNum()
	if err != nil{
		fmt.Fprintf(w, defaultHtml, "地址有错")
		return
	}
	http.Redirect(w, r, shortUrl.url, http.StatusFound)
}

func main(){
	path := flag.String("path", "", "sqlite path")
	flag.Parse()
	if *path == ""{
		*path,_ = os.Getwd()
	}
	fmt.Println(*path)
	err := ConnectDb(*path + "/shorturl.sqlite")
	if err != nil{
		fmt.Println(err)
		return
	}
	defer con.Close()
	var _handler RegexpHandler
	_handler.HandleFunc("/[a-zA-Z0-9]+", UrlHandler)
	_handler.HandleFunc("/", DefaultHandler)
	//err = http.ListenAndServe(":8080", nil)
	err = http.ListenAndServe("0.0.0.0:8080", &_handler)
	if err != nil{
		fmt.Println(err)
		return
	}
}

