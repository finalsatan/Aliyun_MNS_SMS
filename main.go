package main

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"github.com/finalsatan/aliyun_mns/sms"
	"io/ioutil"
	"net/http"
)

const (
	AliyunAccessKeyId     = "xxxxxxxxx"
	AliyunAccessKeySecret = "xxxxxxxxx"
	EndPoint              = "xxxxxxxxx.mns.cn-hangzhou.aliyuncs.com"
	Path                  = "/topics/xxxxxxxxx/messages"
)

func main() {
	client := sms.NewClient(AliyunAccessKeyId, AliyunAccessKeySecret, EndPoint)
	dsms := sms.DirectSMS{
		FreeSignName: "xxxxxx",
		TemplateCode: "SMS_xxxxxxxx",
		Type:         "singleContent",
		Receiver:     "13xxxxxxxxx,17xxxxxxxxx",
		SmsParams:    `{"code":"211051"}`,
	}

	dsmsJson, err := json.Marshal(dsms)
	if err != nil {
		panic(err)
	}
	fmt.Println("DirectSMS: ", string(dsmsJson))

	msgAttr := sms.MessageAttributes{
		DirectSMS: string(dsmsJson),
	}

	message := sms.Message{
		XMLNS:             "http://mns.aliyuncs.com/doc/v1/",
		MessageBody:       "smscontent",
		MessageAttributes: msgAttr,
	}
	data, err := xml.Marshal(message)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(data))

	req := sms.NewRequest(EndPoint, http.MethodPost, Path, data, map[string]string{})

	resp, err := client.DoRequest(req)
	if err != nil {
		panic(err)
	}

	fmt.Println(resp)
	defer resp.Body.Close()
	//err = xml.NewDecoder(response.Body).Decode(msg)

	retdata, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	fmt.Println(string(retdata))

	v := sms.MsgSendRes{}
	err = xml.Unmarshal(retdata, &v)
	if err != nil {
		fmt.Printf("error: %v", err)
		return
	}
	fmt.Println("message_id: ", v.MessageId)
	fmt.Println("message_body_md5: ", v.MessageBodyMD5)
}
