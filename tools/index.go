package tools

import (
	"crypto/sha1"
	"fmt"
	"io"
	"io/ioutil"
	"math/rand"
	"time"
)

// 写文件"./wechatFile/ticket.txt"
func SetFile(ticket string, filename string) bool {
	d1 := []byte(ticket)
	err := ioutil.WriteFile(filename, d1, 0644)
	if err != nil {
		fmt.Println(err)
		return false
	}
	return true
}

// 读取文件内容
func ReadFile(filename string) string {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Println("文件读取出错", err)
		return fmt.Sprintf("%s", err)
	}
	fmt.Println("文件读取成功", string(data))
	return string(data)
}

// 获取随机字符串
func GetRandomString(l int) string {
	str := "0123456789abcdefghijklmnopqrstuvwxyz"
	bytes := []byte(str)
	result := []byte{}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < l; i++ {
		result = append(result, bytes[r.Intn(len(bytes))])
	}
	return string(result)
}

// sha1加密
func MySha1(data string) string {
	t := sha1.New()
	io.WriteString(t, data)
	return fmt.Sprintf("%x", t.Sum(nil))
}
