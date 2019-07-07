/**
 * Author:飘~
 * Date:2019/6/9
 * Description:
 */
package wechat

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
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
	jj := jsonTicket{int(time.Now().Unix()) + 7000, ticket.Ticket}
	data, err := json.Marshal(jj)
	if err != nil {
		log.Println("fetchJsApiTicket Marshal error", err)
	}
	tools.SetFile(string(data), this.SaveFileDir+string(os.PathSeparator)+"ticket.txt")
	return ticket.Ticket, nil
}

func (this *Wechat) getJsApiTicket() (ticket string) {
	json_str := tools.ReadFile(this.SaveFileDir + string(os.PathSeparator) + "ticket.txt")
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
	signPackage := signPackage{debug, this.appid, nonceStr, int(time.Now().Unix()), url, signature, []string{"updateAppMessageShareData", "updateTimelineShareData"}}
	data, err := json.Marshal(signPackage)
	if err != nil {
		fmt.Println(err)
	}
	return string(data)
}
