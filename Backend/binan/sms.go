package binan

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"webadmin/config"
)

// http://utf8.api.smschinese.cn/?Uid=本站用户名&Key=接口短信密钥&smsMob=手机号码&smsText=验证码:8888
func SendSMS(mobile, text string) (string, error) {
	apiURL := "http://utf8.api.smschinese.cn/"
	smsuid := config.Get("Smsuid")
	Smskey := config.Get("Smskey")
	data := url.Values{}
	data.Set("Uid", smsuid)
	data.Set("Key", Smskey)
	data.Set("smsMob", mobile)
	data.Set("smsText", "验证码"+text+"，5分钟内有效，请勿泄露与他人")

	req, err := http.NewRequest(
		http.MethodPost,
		apiURL,
		strings.NewReader(data.Encode()),
	)
	if err != nil {
		return "", err
	}

	req.Header.Set(
		"Content-Type",
		"application/x-www-form-urlencoded",
	)

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	fmt.Println("body:", string(body))
	return string(body), nil
}

func TestSendSMS() {
	smsuid := config.Get("Smsuid")
	Smskey := config.Get("Smskey")
	str, err := SendSMS("15059862915", "100563")
	if err != nil {
		fmt.Println("err:", err.Error())
	}
	fmt.Println("str:", str)
	fmt.Println("smsuid:", smsuid)
	fmt.Println("Smskey:", Smskey)

}
