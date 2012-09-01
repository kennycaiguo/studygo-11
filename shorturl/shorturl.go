package main

import(
	"sqlite"
	"sync"
	"errors"
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

func (this ShortUrl)RandNum() (string, error){
	randNum.lock.Lock()
	defer randNum.lock.Unlock()
	rs, err := con.Prepare("select id from short_num limit 1")
	if err != nil{
		return "", err
	}
	defer rs.Finalize()
	rs.Exec()
	err = rs.Scan(&randNum.num)
	if err != nil{
		return "", err
	}
	con.Exec("update short_num set id=id+1 limit 1")
	return IntToNum(randNum.num), nil
}

func (this ShortUrl)Insert() error{
	return con.Exec("insert into short_url(url, num) values (?,?)", this.url, this.num)
}

func (this ShortUrl)Load() error{
	if this.num == "" {
		return errors.New("num is empty")
	}
	rs, err := con.Prepare("select * from short_url where num=?")
	if err != nil{
		return err
	}
	defer rs.Finalize()
	rs.Exec()
	if !rs.Next(){
		return errors.New("no data")
	}
	err = rs.Scan(&this.id, &this.url, &this.num)
	if err == nil{
		return nil
	}
	return err
}

func main(){
	ConnectDb("shorturl.sqlite")
}

