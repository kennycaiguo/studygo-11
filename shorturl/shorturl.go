package main

import (
	"fmt"
	"sqlite"
)

func main() {
	db,err := sqlite.Open("shorturl.sqlite")
	if err != nil{
		fmt.Println(err)
		return
	}
	insert(db)
	db.Close()
}

func insert(db *sqlite.Conn){
	err := db.Exec("insert into short_url(id,url, num) values (?, ?, ?)", 3, "http://tg.tlt.cn", 3)
	if err != nil{
		fmt.Println(err)
		return
	}
}

func show(db *sqlite.Conn){
	var id int
	var url string
	var num string
	rs, err := db.Prepare("select * from short_url")
	if err != nil{
		fmt.Println(err)
		return
	}
	rs.Exec()
	for rs.Next(){
		err = rs.Scan(&id, &url, &num)
		if err == nil{
			fmt.Println(id, url ,num)
		}
	}
}
