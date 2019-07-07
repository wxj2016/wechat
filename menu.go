/**
 * Author:飘~
 * Date:2019/6/9
 * Description:
 */
package wechat

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"wechat/config"
)

var menus = `{
  "button":[
      
      {
        "type":"view",
        "name":"微官网",
        "url":"http://www.cbd88.cn"
      },
      {
        "type":"view",
        "name":"关于我们",
        "url":"http://www.cbd88.cn/wapindex.php/Index/about.html"
      }
  ]
}`

func (this *Wechat) MenuCreate(menus string) string {
	token := this.GetAccessToken()
	url := config.Links["menu_create"] + "?access_token=" + token
	reader := bytes.NewBuffer([]byte(menus))
	bodytype := "application/json;charset=utf-8"
	res_post, err := http.Post(url, bodytype, reader)
	if err == nil {
		body_post, _ := ioutil.ReadAll(res_post.Body)
		defer res_post.Body.Close()

		return string(body_post)

	} else {
		//创建菜单失败 }
		return "创建菜单失败"
	}

}

func (this *Wechat) MenuDelete() error {
	token := this.GetAccessToken()
	url := config.Links["menu_delete"] + "?access_token=" + token

	resp, err := http.Get(url)
	defer resp.Body.Close()
	if err != nil {
		log.Println(err, "MenuDelete Error")
		return err
	}

	if resp.StatusCode != http.StatusOK {
		fmt.Println("Error:status code", resp.StatusCode)
		return errors.New("StatusCode is " + strconv.Itoa(resp.StatusCode))
	}
	return nil

}
