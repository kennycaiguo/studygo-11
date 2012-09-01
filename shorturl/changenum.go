/**
 * 任意进制置换(62进制以内)
 **/
package main

import "math"
import "strings"

const chars string = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
const charLen int = 62

//将10进制整数转换为62进制
func IntToNum(value int) string{
	num := ""
	i := value
	for ; i>0; i=i/charLen{
		a := string(chars[i%charLen])
		num = a + num
	}
	return num
}

//将62进制转换为10进制整数
func NumToInt(value string) int{
	l := len(value) - 1
	if l < 0{
		return 0
	}
	num := 0
	index := 0
	for i,v:= range value{
		index = strings.Index(chars, string(v))
		num += index * int(math.Pow(float64(charLen), float64(l-i)))
	}
	return num
}

/*func main(){
	v := IntToNum(98760)
	fmt.Println(v)
	fmt.Println(NumToInt(v))
}*/
