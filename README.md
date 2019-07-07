微信SDK
======
提供微信常用的功能，如接入，各种事件，用户基本信息，菜单，场景二维码，jssdk等常用功能

## 安装
`go get github.com/wxj2016/wechat`

## 在Beego上的用法
####接入和一些消息的接管
```
const redirect = "http://wx.cbd88.com"

var wx *wechat.Wechat

var Token = "wxxxxxxxxx"
var appid = "xxxxxxxxxx"
var appsecret = "xxxxxxxxxxxxxxxxxxxxxxxxxxxx"

func init() {
	w, err := wechat.NewWx(appid, Token, appsecret, "")
	wx = w
	if err != nil {
		log.Println(err)
	}
}

type WxController struct {
	beego.Controller
}

// @router /cbdwechat/index [get,post]
func (this *WxController) Oauth() {
	this.EnableRender = false
	wx.FuncTxt = func(s string) string {
		if s == "追梦人" {
			return "追梦人 点击进入签到:http://www.cbd88.com"
		} else if s == "407787759" {
			return "是你自己"
		} else if s == "" {
			return ""
		} else {
			return "如有需求，请致电我们，400-830-1116"
		}
	}
	wx.FuncEven = func(wxEvent, EventKey string) string {
		log.Println("wxEvent:", wxEvent, "EventKey:", EventKey)
		if EventKey == "103" || EventKey == "qrscene_103" {
			return "通过场景二维码103进来的 请点击签到 http://www.cbd88.com"
		}
		return "CBD家居欢迎您！"
	}

	wx.Run(this.Ctx.Request, this.Ctx.ResponseWriter)
}
```

####场景二维码生成
```
var str = `{"action_name": "QR_LIMIT_SCENE", "action_info": {"scene": {"scene_id": 103}}}`
this.Redirect(wx.GetQrcode(str), 302)
```

####取用户基本信息
```
func (this *WxController) Index() {
	code := this.GetString("code")
	if code == "" {
		scope := "snsapi_userinfo"
		state := ""

		url := wx.GetOauthCode(redirect, scope, state)

		this.Redirect(url, 301)
	} else {
		Otoken := wxType.Otoken{}
		Otokenb := wx.GetOauthAccessToken(code)
		json.Unmarshal(Otokenb, &Otoken)

		if Otoken.Errcode != 0 {
			this.Redirect(redirect, 301)
			return
		}

		userInfo := wxType.UserInfo{}
		userInfob, err := wx.OauthUserInfo(Otoken.Access_token, Otoken.Openid)
		json.Unmarshal(userInfob, &userInfo)
	}
}
```

####创建和删除菜单方法
```
var menus = `{
  "button":[
      
      {
        "type":"view",
        "name":"微官网",
        "url":"http://www.cbd88.com"
      },
      {
        "type":"view",
        "name":"关于我们",
        "url":"http://www.cbd88.com/wapindex.php/Index/about.html"
      }
  ]
}`


res:=wx.MenuCreate(menus)
err:=wx.MenuDelete()

```


####jssdk
```
config：=wx.GetWechatConfig(url, false)
```