/**
 * Author:飘~
 * Date:2019/6/9
 * Description:
 */
package wechat

import (
	"bufio"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"wechat/config"
)

// 获取code
func (this *Wechat) GetOauthCode(redirect, scope, state string) (redirectUrl string) {
	return config.Links["oauth_code"] + "?appid=" + this.appid + "&redirect_uri=" + redirect + "&response_type=code&scope=" + scope + "&state=" + state + "#wechat_redirect"
}

// 通过code获取access_token
func (this *Wechat) GetOauthAccessToken(code string) (token []byte) {
	url := config.Links["oauth_access_token"] + "?appid=" + this.appid + "&secret=" + this.appsecret + "&code=" + code + "&grant_type=authorization_code"
	resp, err := http.Get(url)
	defer resp.Body.Close()
	if err != nil {
		log.Println("GetOauthAccessToken Error", err)
	}

	if resp.StatusCode != http.StatusOK {
		log.Println("Error:GetOauthAccessToken status code", resp.StatusCode)
	}
	bodyReader := bufio.NewReader(resp.Body)
	bytes, _ := ioutil.ReadAll(bodyReader)

	return bytes
}

// 通过access_token获取用户信息，仅限scope为SCOPE_POP
func (this *Wechat) OauthUserInfo(access_token, openid string) (userInfo []byte, err error) {
	url := config.Links["oauth_userinfo"] + "?access_token=" + access_token + "&openid=" + openid + "&lang=zh_CN"

	resp, err := http.Get(url)
	defer resp.Body.Close()
	if err != nil {
		return []byte{}, err
	}
	if resp.StatusCode != http.StatusOK {
		return []byte{}, errors.New("OauthUserInfo Error")
	}
	bodyReader := bufio.NewReader(resp.Body)
	bytes, err := ioutil.ReadAll(bodyReader)

	if err != nil {
		return []byte{}, err
	}
	return bytes, nil
}
