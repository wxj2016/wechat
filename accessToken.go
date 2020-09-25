/**
 * Author:飘~
 * Date:2019/6/10
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
	"time"
	"wechat/config"
	"wechat/wxType"
)

// 外部获取AccessToken
func (this *Wechat) GetAccessToken() string {
	// 获取ticket内容
	//json_str := tools.ReadFile(this.SaveFileDir + string(os.PathSeparator) + "/token.txt")
	j := jsonToken{}
	s,_:=R.Get("token").Result()
	if s==""{
		token, err := this.fetchAccessToken()
		if err != nil {
			log.Println(err)
		}
		return token.AccessToken
	}
	// 解析json串
	err := json.Unmarshal([]byte(s), &j)
	if err != nil {
		fmt.Println("GetAccessToken err:", err)
		token, err := this.fetchAccessToken()
		if err != nil {
			log.Println(err)
		}
		return token.AccessToken

		// ???????????????????
	}
	expire_time := j.Expire_time
	access_token := j.Access_token

	if expire_time < int(time.Now().Unix()) {
		log.Println("----------过期了")
		token, err := this.fetchAccessToken()
		if err != nil {
			log.Println(err)
		}
		return token.AccessToken
	}
	return access_token
}

// 去微信官方请求基础AccessToken
func (this *Wechat) fetchAccessToken() (wxType.BasicToken, error) {
	url := config.Links["access_token"] + "?grant_type=client_credential&appid=" + this.appid + "&secret=" + this.appsecret
	var token wxType.BasicToken
	resp, err := http.Get(url)
	defer resp.Body.Close()
	if err != nil {
		log.Println(err, "fetchAccessToken Error", err)
		return token, err
	}

	if resp.StatusCode != http.StatusOK {
		fmt.Println("Error:status code", resp.StatusCode)
		return token, err
	}

	bodyReader := bufio.NewReader(resp.Body)

	bytes, _ := ioutil.ReadAll(bodyReader)
	err = json.Unmarshal(bytes, &token)

	if err != nil {
		log.Println(err, "Unmarshal Error")
		return token, err
	}
	// 保存到文件中
	jj := jsonToken{int(time.Now().Unix()) + 7180, token.AccessToken}
	data, err := json.Marshal(jj)
	if err != nil {
		log.Println("fetchAccessToken Marshal error", err)
	}
	R.Set("token",data,time.Duration(time.Second*7180)).Err()
	//tools.SetFile(string(data), this.SaveFileDir+string(os.PathSeparator)+"/token.txt")
	return token, nil
}
