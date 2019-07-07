/**
 * Author:飘~
 * Date:2019/6/9
 * Description:场景二维码生成
 */
package wechat

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"wechat/config"
)

//var str = `{"expire_seconds": 604800,"action_name": "QR_SCENE", "action_info": {"scene": {"scene_id": 123}}}`
//action_name	二维码类型，
// QR_SCENE为临时的整型参数值，
// QR_STR_SCENE为临时的字符串参数值，
// QR_LIMIT_SCENE为永久的整型参数值，
// QR_LIMIT_STR_SCENE为永久的字符串参数值

func (this *Wechat) GetQrcode(str string) string {
	token := this.GetAccessToken()
	url := config.Links["qrcode"] + "?access_token=" + token
	reader := bytes.NewBuffer([]byte(str))
	bodytype := "application/json;charset=utf-8"
	res_post, err := http.Post(url, bodytype, reader)
	if err == nil {
		body_post, _ := ioutil.ReadAll(res_post.Body)
		defer res_post.Body.Close()

		res := map[string]string{}
		json.Unmarshal(body_post, &res)

		url = config.Links["showqrcode"] + "?ticket=" + res["ticket"]
		return url
	} else {
		return "err"
	}
}
