package main

import(
	"fmt"
	"net/http"
	"regexp"
)

type route struct{
	pattern *regexp.Regexp
	handler http.Handler
}

type RegexpHandler struct{
	routes []*route
}

//func (this *RegexpHandler)Handler(pattern *regexp.Regexp, handler http.Handler){
func (this *RegexpHandler)Handle(s string, handler http.Handler){
	pattern, _:= regexp.Compile(s)
	this.routes = append(this.routes, &route{pattern, handler})
}

func (this *RegexpHandler)HandleFunc(s string, handler func(http.ResponseWriter, *http.Request)){
	pattern, _:= regexp.Compile(s)
	a := http.HandlerFunc(handler)
	this.routes = append(this.routes, &route{pattern, a})
}

func (this *RegexpHandler)ServeHTTP(w http.ResponseWriter, r *http.Request){
	for _, route := range this.routes{
		fmt.Println(route.pattern.MatchString(r.URL.Path), route.pattern, r.URL.Path)
		if route.pattern.MatchString(r.URL.Path){
			route.handler.ServeHTTP(w, r)
			return
		}
	}
	http.NotFound(w, r)
}

func DefaultHandler(w http.ResponseWriter, r *http.Request){
	fmt.Fprint(w, "default")
}

func UrlHandler(w http.ResponseWriter, r *http.Request){
	fmt.Fprint(w, "test")
}

func main(){
	var _handler RegexpHandler
	_handler.HandleFunc("/[a-zA-Z0-9]+", UrlHandler)
	_handler.HandleFunc("/", DefaultHandler)
	err := http.ListenAndServe(":8080", &_handler)
	if err != nil{
		fmt.Println(err)
		return
	}
}
