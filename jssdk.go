/**
 * Author:飘~
 * Date:2019/6/9
 * Description:
 */
package wechat

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"
	"wechat/config"
	"wechat/tools"
	"wechat/wxType"
)

type signPackage struct {
	Debug     bool     `json:"debug"`
	AppId     string   `json:"appId"`
	NonceStr  string   `json:"nonceStr"`
	Timestamp int      `json:"timestamp"`
	Url       string   `json:"url"`
	Signature string   `json:"signature"`
	JsApiList []string `json:"jsApiList"`
}

func (this *Wechat) fetchJsApiTicket(accessToken string) (ti string, err error) {
	url := config.Links["js_api_tocket"] + "?access_token=" + accessToken + "&type=jsapi"
	resp, err := http.Get(url)
	defer resp.Body.Close()
	if err != nil {
		log.Println(err, "jsApiTicket Error---")
		return "", err
	}

	if resp.StatusCode != http.StatusOK {
		fmt.Println("Error:status code aaa", resp.StatusCode)
		return "", err
	}

	bodyReader := bufio.NewReader(resp.Body)

	bytes, _ := ioutil.ReadAll(bodyReader)

	ticket := wxType.Ticket{}
	err = json.Unmarshal(bytes, &ticket)

	if err != nil {
		log.Println(err, "Unmarshal Error")
		return "", err
	}
	// 保存到文件中
	jj := jsonTicket{int(time.Now().Unix()) + 7180, ticket.Ticket}
	data, err := json.Marshal(jj)
	if err != nil {
		log.Println("fetchJsApiTicket Marshal error", err)
	}
	R.Set("ticket",data,time.Duration(time.Second*7180)).Err()
	//tools.SetFile(string(data), this.SaveFileDir+string(os.PathSeparator)+"ticket.txt")
	return ticket.Ticket, nil
}

func (this *Wechat) getJsApiTicket() (ticket string) {
	//json_str := tools.ReadFile(this.SaveFileDir + string(os.PathSeparator) + "ticket.txt")
	json_str,_:=R.Get("ticket").Result()
	if json_str==""{
		access_token, _ := this.fetchAccessToken()
		s, _ := this.fetchJsApiTicket(access_token.AccessToken)
		return s
	}
	// 解析json串
	t := jsonTicket{}
	err := json.Unmarshal([]byte(json_str), &t)
	if err != nil {
		fmt.Println(err.Error())
		access_token, _ := this.fetchAccessToken()
		s, _ := this.fetchJsApiTicket(access_token.AccessToken)
		return s
	}
	expire_time := t.Expire_time
	jsapi_ticket := t.Jsapi_ticket
	if expire_time < int(time.Now().Unix()) {
		// ticket过期需要从新获取
		// 获取access_token生成ticket
		access_token, _ := this.fetchAccessToken()
		s, _ := this.fetchJsApiTicket(access_token.AccessToken)
		jsapi_ticket = s

	}
	return jsapi_ticket
}

func (this *Wechat) GetWechatConfig(url string, debug bool) string {
	ticket := this.getJsApiTicket()
	nonceStr := tools.GetRandomString(16)
	my_string := fmt.Sprintf("jsapi_ticket=%s&noncestr=%s&timestamp=%s&url=%s", ticket, nonceStr, strconv.Itoa(int(time.Now().Unix())), url)
	//log.Println("my_string:", my_string)
	signature := tools.MySha1(my_string)
	//log.Println("signature:", signature)
	signPackage := signPackage{debug, this.appid, nonceStr, int(time.Now().Unix()), url, signature, []string{"updateAppMessageShareData", "updateTimelineShareData", "checkJsApi", "scanQRCode", "onMenuShareTimeline", "onMenuShareAppMessage","chooseImage","previewImage","uploadImage"}}
	data, err := json.Marshal(signPackage)
	if err != nil {
		fmt.Println(err)
	}
	return string(data)
}

// 返回几个重要的值，
//func (this *Wechat) Getjsapisign(url string) (appid,timestamp,nonceStr,signature string) {
//	appid=this.appid
//	ticket := this.getJsApiTicket()
//	timestamp=strconv.Itoa(int(time.Now().Unix()))
//	nonceStr = tools.GetRandomString(16)
//	my_string := fmt.Sprintf("jsapi_ticket=%s&noncestr=%s&timestamp=%s&url=%s", ticket, nonceStr, timestamp, url)
//	//log.Println("my_string:", my_string)
//	signature = tools.MySha1(my_string)
//	//log.Println("signature:", signature)
//	return
//}

func (this *Wechat) Getjsapisign(url string) (string,int,string,string) {
	ticket := this.getJsApiTicket()
	nonceStr := tools.GetRandomString(16)
	tttt:=int(time.Now().Unix())
	my_string := fmt.Sprintf("jsapi_ticket=%s&noncestr=%s&timestamp=%s&url=%s", ticket, nonceStr, strconv.Itoa(tttt), url)
	//log.Println("my_string:", my_string)
	signature := tools.MySha1(my_string)
	////log.Println("signature:", signature)
	//signPackage := signPackage{debug, this.appid, nonceStr, int(time.Now().Unix()), url, signature, []string{"updateAppMessageShareData", "updateTimelineShareData", "checkJsApi", "scanQRCode", "onMenuShareTimeline", "onMenuShareAppMessage","chooseImage","previewImage","uploadImage"}}
	//data, err := json.Marshal(signPackage)
	//if err != nil {
	//	fmt.Println(err)
	//}
	//return string(data)
	return this.appid,tttt,nonceStr,signature
}

//https://api.weixin.qq.com/cgi-bin/shorturl?access_token=ACCESS_TOKEN
func (this *Wechat) Getsorturl(logurl string) (string) {
	token := this.GetAccessToken()
	url := config.Links["shorturl"] + "?access_token=" + token

	ns:=fmt.Sprintf(`{
  "action":"long2short",
  "long_url":"%s",
}`,logurl)
	log.Println("HHHHHHHH:",ns)
	reader := bytes.NewBuffer([]byte(ns))
	bodytype := "application/json;charset=utf-8"
	res_post, err := http.Post(url, bodytype, reader)
	if err == nil {
		body_post, _ := ioutil.ReadAll(res_post.Body)
		defer res_post.Body.Close()

		return string(body_post)
	}else{
		log.Println("++++++1589103976+++++++",err)
		return ""
	}
}

//template_send
func (this *Wechat) Template_send(datastr string) (string) {
	token := this.GetAccessToken()
	url2 := config.Links["template_send"] + "?access_token=" + token

	reader := bytes.NewBuffer([]byte(datastr))
	bodytype := "application/json;charset=utf-8"
	res_post, err := http.Post(url2, bodytype, reader)
	if err == nil {
		body_post, _ := ioutil.ReadAll(res_post.Body)
		defer res_post.Body.Close()

		return string(body_post)
	}else{
		return ""
	}
}

// 发送一次订阅消息
func (this *Wechat) Subscribemsg(touser,template_id,dsturl,appid,pagepath,scene,title,msg string) (string) {
	token := this.GetAccessToken()
	url := config.Links["subscribe"] + "?access_token=" + token
	ns:=fmt.Sprintf(`{
    "touser":"%s",
    "template_id":"%s",
    "url":"%s",
    "miniprogram":{
    "appid":"%s",
    "pagepath":"%s"    
},
    "scene":"%s",
    "title":"%s",
    "data":{
    "content":{
    "value":"%s",
    "color":""
}
}
}`,touser,template_id,dsturl,appid,pagepath,scene,title,msg)


	reader := bytes.NewBuffer([]byte(ns))
	bodytype := "application/json;charset=utf-8"
	res_post, err := http.Post(url, bodytype, reader)
	if err == nil {
		body_post, _ := ioutil.ReadAll(res_post.Body)
		defer res_post.Body.Close()
		return string(body_post)
	}else{
		log.Println("++++++1589103976+++++++",err)
		return ""
	}
}