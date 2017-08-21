package sms

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"net/url"
	"sort"
	"strings"
)

const HeaderMNSPrefix = "x-mns-"

//授权签名
func (client *Client) SignRequest(req *Request, payload []byte) {

	//	SignString = VERB + "\n"
	//	+ CONTENT-MD5 + "\n"
	//	+ CONTENT-TYPE + "\n"
	//	+ DATE + "\n"
	//	+ CanonicalizedLOGHeaders + "\n"
	//	+ CanonicalizedResource

	if _, ok := req.headers["Authorization"]; ok {
		return
	}

	contentType := req.headers["Content-Type"]
	contentMd5 := ""
	//Header中未加Content-MD5
	//if payload != nil {
	//	contentMd5 = Md5(payload)
	//	req.headers["Content-MD5"] = contentMd5
	//}
	date := req.headers["Date"]
	canonicalizedHeader := canonicalizeHeader(req.headers)
	canonicalizedResource := canonicalizeResource(req)

	signString := req.method + "\n" + contentMd5 + "\n" + contentType + "\n" + date + "\n" + canonicalizedHeader + "\n" + canonicalizedResource
	signature := CreateSignature(signString, client.AccessKeySecret)
	req.headers["Authorization"] = "MNS " + client.AccessKeyId + ":" + signature
}

func canonicalizeResource(req *Request) string {
	canonicalizedResource := req.path
	var paramNames []string
	if req.params != nil && len(req.params) > 0 {
		for k, _ := range req.params {
			paramNames = append(paramNames, k)
		}
		sort.Strings(paramNames)

		var query []string
		for _, k := range paramNames {
			query = append(query, url.QueryEscape(k)+"="+url.QueryEscape(req.params[k]))
		}
		canonicalizedResource = canonicalizedResource + "?" + strings.Join(query, "&")
	}
	return canonicalizedResource
}

//Have to break the abstraction to append keys with lower case.
func canonicalizeHeader(headers map[string]string) string {
	var canonicalizedHeaders []string

	for k, _ := range headers {
		if lower := strings.ToLower(k); strings.HasPrefix(lower, HeaderMNSPrefix) {
			canonicalizedHeaders = append(canonicalizedHeaders, lower)
		}
	}

	sort.Strings(canonicalizedHeaders)

	var headersWithValue []string

	for _, k := range canonicalizedHeaders {
		headersWithValue = append(headersWithValue, k+":"+headers[k])
	}
	return strings.Join(headersWithValue, "\n")
}

//CreateSignature creates signature for string following Aliyun rules
func CreateSignature(stringToSignature, accessKeySecret string) string {
	// Crypto by HMAC-SHA1
	hmacSha1 := hmac.New(sha1.New, []byte(accessKeySecret))
	hmacSha1.Write([]byte(stringToSignature))
	sign := hmacSha1.Sum(nil)

	// Encode to Base64
	base64Sign := base64.StdEncoding.EncodeToString(sign)

	return base64Sign
}

func percentReplace(str string) string {
	str = strings.Replace(str, "+", "%20", -1)
	str = strings.Replace(str, "*", "%2A", -1)
	str = strings.Replace(str, "%7E", "~", -1)

	return str
}

// CreateSignatureForRequest creates signature for query string values
func CreateSignatureForRequest(method string, values *url.Values, accessKeySecret string) string {

	canonicalizedQueryString := percentReplace(values.Encode())

	stringToSign := method + "&%2F&" + url.QueryEscape(canonicalizedQueryString)

	return CreateSignature(stringToSign, accessKeySecret)
}
