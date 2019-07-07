/**
 * Author:飘~
 * Date:2019/6/10
 * Description:
 */
package wechat

import (
	"encoding/xml"
	"errors"
	"fmt"
	"github.com/clbanning/mxj"
	"io/ioutil"
	"log"
	"strconv"
	"time"
)

func (this *Wechat) listen() {
	err := this.initMessage()
	if err != nil {
		log.Println(err)
		this.ResponseWriter.WriteHeader(403)
		return
	}

	MsgType, ok := this.Message["MsgType"].(string)
	log.Println("===", MsgType, "===")

	if !ok {
		this.ResponseWriter.WriteHeader(403)
		return
	}

	switch MsgType {
	case "text":
		this.text(this.FuncTxt)
		break
	case "event":
		this.event(this.FuncEven)
		break
	default:
		break
	}

	return
}

func (this *Wechat) initMessage() error {

	body, err := ioutil.ReadAll(this.Request.Body)

	if err != nil {
		return err
	}

	m, err := mxj.NewMapXml(body)

	if err != nil {
		return err
	}

	if _, ok := m["xml"]; !ok {
		return errors.New("Invalid Message.")
	}

	message, ok := m["xml"].(map[string]interface{})

	if !ok {
		return errors.New("Invalid Field `xml` Type.")
	}

	this.Message = message

	log.Println(this.Message)

	return nil
}

type Base struct {
	FromUserName CDATAText
	ToUserName   CDATAText
	MsgType      CDATAText
	CreateTime   CDATAText
}

func (b *Base) InitBaseData(w *Wechat, msgtype string) {
	b.FromUserName = Value2CDATA(w.Message["ToUserName"].(string))
	b.ToUserName = Value2CDATA(w.Message["FromUserName"].(string))
	b.CreateTime = Value2CDATA(strconv.FormatInt(time.Now().Unix(), 10))
	b.MsgType = Value2CDATA(msgtype)
}

type CDATAText struct {
	Text string `xml:",innerxml"`
}

type TextMessage struct {
	XMLName xml.Name `xml:"xml"`
	Base
	Content CDATAText
}

func Value2CDATA(v string) CDATAText {
	return CDATAText{"<![CDATA[" + v + "]]>"}
}

//================================================================================

func (this *Wechat) text(f func(s string) string) {
	inMsg, ok := this.Message["Content"].(string)

	if !ok {
		return
	}
	var reply TextMessage
	log.Printf("%s", inMsg)

	ss := f(inMsg)

	reply.InitBaseData(this, "text")
	reply.Content = Value2CDATA(fmt.Sprintf("%s", ss))

	replyXml, err := xml.Marshal(reply)
	if err != nil {
		log.Println(err)
		this.ResponseWriter.WriteHeader(403)
		return
	}

	this.ResponseWriter.Header().Set("Content-Type", "text/xml")
	this.ResponseWriter.Write(replyXml)
}

func (this *Wechat) event(f func(wxEvent, EventKey string) string) {
	wxEvent, _ := this.Message["Event"].(string)
	EventKey, _ := this.Message["EventKey"].(string)
	log.Println("####wxEvent is:", wxEvent, "####EventKey is:", EventKey, "end")
	inMsg := f(wxEvent, EventKey)
	if inMsg != "" {
		var reply TextMessage
		reply.InitBaseData(this, "text")
		reply.Content = Value2CDATA(fmt.Sprintf("%s", inMsg))

		replyXml, err := xml.Marshal(reply)

		if err != nil {
			log.Println(err)
			this.ResponseWriter.WriteHeader(403)
			return
		}

		this.ResponseWriter.Header().Set("Content-Type", "text/xml")
		this.ResponseWriter.Write(replyXml)
	}
}
