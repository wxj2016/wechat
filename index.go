package wechat

import (
	"crypto/sha1"
	"fmt"
	"github.com/go-redis/redis"
	"log"
	"net/http"
	"os"
	"sort"
	"wechat/config"
)

var R *redis.Client
func init()  {
	// redis ================================================
	client := redis.NewClient(&redis.Options{
		Addr:   config.Redishost,
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	_, err := client.Ping().Result()
	if err != nil {
		panic(err)
	} else {
		R = client
	}
}

type weixinQuery struct {
	Signature    string `json:"signature"`
	Timestamp    string `json:"timestamp"`
	Nonce        string `json:"nonce"`
	EncryptType  string `json:"encrypt_type"`
	MsgSignature string `json:"msg_signature"`
	Echostr      string `json:"echostr"`
}

type Wechat struct {
	appid, Token, appsecret, SaveFileDir string
	Query                                weixinQuery
	Message                              map[string]interface{}
	Request                              *http.Request
	ResponseWriter                       http.ResponseWriter
	FuncTxt                              func(string) string
	FuncEven                             func(wxEvent, EventKey string) string
}

type jsonToken struct {
	Expire_time  int
	Access_token string
}

type jsonTicket struct {
	Expire_time  int
	Jsapi_ticket string
}

func NewWx(appid, token, appsecret, SaveFileDir string) (wx *Wechat, err error) {
	wx = &Wechat{}
	wx.appid = appid
	wx.Token = token
	wx.appsecret = appsecret

	// 初始化默认文本信息处理方法
	wx.FuncTxt = func(inMsg string) string {
		inMsg = "收到信息：" + inMsg
		return inMsg
	}
	// 初始化默认事件信息处理方法
	wx.FuncEven = func(wxEvent, EventKey string) string {
		inMsg := ""
		if wxEvent == "subscribe" || wxEvent == "SCAN" {
			if EventKey == "123" || EventKey == "qrscene_123" {
				inMsg = "追梦人 CBD家居"
			} else {
				inMsg = "CBD家居欢迎您！"
			}
		}
		return inMsg
	}

	if SaveFileDir != "" {
		wx.SaveFileDir = SaveFileDir
		err = os.Mkdir(wx.SaveFileDir, os.ModePerm)

		if err != nil {
			if os.IsNotExist(err) {
				log.Println(err)
			}
		}
	}
	return wx, nil
}

// 启用接管
func (this *Wechat) Run(r *http.Request, w http.ResponseWriter) {
	this.Request = r
	this.ResponseWriter = w
	this.initWxQuery()
	if this.Query.Signature != this.signature() {
		w.WriteHeader(403)
		return
	}

	if r.Method == http.MethodGet {
		if len(this.Query.Echostr) > 0 {
			w.Write([]byte(this.Query.Echostr))
			return
		}
	} else if r.Method == http.MethodPost {
		this.listen()
	}
}

func (this *Wechat) initWxQuery() {
	var q weixinQuery
	q.Nonce = this.Request.URL.Query().Get("nonce")
	q.Echostr = this.Request.URL.Query().Get("echostr")
	q.Signature = this.Request.URL.Query().Get("signature")
	q.Timestamp = this.Request.URL.Query().Get("timestamp")
	q.EncryptType = this.Request.URL.Query().Get("encrypt_type")
	q.MsgSignature = this.Request.URL.Query().Get("msg_signature")
	this.Query = q
	//log.Println("===", this.Query, "===")
}

func (this *Wechat) signature() string {
	strs := sort.StringSlice{this.Token, this.Query.Timestamp, this.Query.Nonce}
	sort.Strings(strs)
	str := ""
	for _, s := range strs {
		str += s
	}
	h := sha1.New()
	h.Write([]byte(str))
	return fmt.Sprintf("%x", h.Sum(nil))
}
