package hello

import (
	"appengine"
	"appengine/user"
	"fmt"
	"net/http"
)

func init(){
	http.HandleFunc("/", handler)
	http.HandleFunc("/login", handleLogin)
	http.HandleFunc("/logout", handleLogout)
	http.HandleFunc("/add", handleAdd)
	http.HandleFunc("/show", handleShow)
}

func handler(w http.ResponseWriter, r *http.Request){
	u := getUser(r)
	if u == nil{
		fmt.Fprintf(w, "<html><body>test, <a href='/login'>login</a></body></html>")
	}else{
		fmt.Fprintf(w, "<html><body>hello, %s , <a href='/logout'>logout</a></body></html>", u.Email)
	}
}

func handleLogin(w http.ResponseWriter, r *http.Request){
	u := getUser(r)
	if u == nil{
		url, err := user.LoginURL(appengine.NewContext(r), "/")
		if err != nil{
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Location", url)
		w.WriteHeader(http.StatusFound)
		return
	}
	w.Header().Set("Location", "/")
	w.WriteHeader(http.StatusFound)
}


func handleLogout(w http.ResponseWriter, r *http.Request){
	u := getUser(r)
	if u == nil {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}
	url, err := user.LogoutURL(appengine.NewContext(r), "/")
	if err != nil{
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Location", url)
	w.WriteHeader(http.StatusFound)
}

var addHtml = `
	<html>
		<head></head>
		<body>
			<form action="/show" method="post">
				<textarea name="data"></textarea>
				<input type="submit" value="submit"/>
			</form>
		</body>
	</html>
`

func handleAdd(w http.ResponseWriter, r *http.Request){
	fmt.Fprintf(w, addHtml)
}

func handleShow(w http.ResponseWriter, r *http.Request){
	data := r.FormValue("data")
	fmt.Fprintf(w, "%v", data)
}

/*************util*********************/
func getUser(r *http.Request) *user.User{
	c := appengine.NewContext(r)
	return user.Current(c)
}
