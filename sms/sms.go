package sms

import (
	"encoding/xml"
	"net/http"
)

const MNSAPIVersion = "2015-06-06"

type Client struct {
	AccessKeyId     string
	AccessKeySecret string
	Endpoint        string
	Version         string
	httpClient      *http.Client
}

func NewClient(accessKeyId, accessKeySecret, endpoint string) (client *Client) {
	client = &Client{
		AccessKeyId:     accessKeyId,
		AccessKeySecret: accessKeySecret,
		Endpoint:        endpoint,
		Version:         MNSAPIVersion,
		httpClient:      &http.Client{},
	}
	return client
}

type Message struct {
	XMLNS             string `xml:"xmlns,attr"`
	MessageBody       string `xml:"MessageBody"`
	MessageAttributes MessageAttributes
}

type MessageAttributes struct {
	DirectSMS string
}

type DirectSMS struct {
	FreeSignName string
	TemplateCode string
	Type         string
	Receiver     string
	SmsParams    string
}

type MsgSendRes struct {
	XMLName        xml.Name `xml:"Message"`
	MessageId      string   `xml:"MessageId"`
	MessageBodyMD5 string   `xml:"MessageBodyMD5"`
}

type ReturnMessage struct {
	XMLName        xml.Name `xml:"Message"`
	XMLNS          string   `xml:"xmlns,attr"`
	MessageId      string   `xml:"MessageId"`
	MessageBodyMD5 string   `xml:"MessageBodyMD5"`
}
