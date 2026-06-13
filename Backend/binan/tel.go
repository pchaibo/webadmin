//测试：go run ./test.go

package binan

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

func Binantel() {
	t := time.Now().Unix()
	s := fmt.Sprintf("%d", t)
	fmt.Println("time:", s)
	const strUrl = "https://api.ihuyi.com/vm/Submit.json"
	v := url.Values{}
	v.Set("account", "VM83102050")                        //APIID(用户中心【云语音】-【语音通知】-【产品总览】查看)
	v.Set("password", "88b44048eed0beecad58ad7edb926d6e") //1、APIKEY(用户中心【云语音】-【语音通知】-【产品总览】查看)2、动态密码(生成动态密码方式请看该文档末尾的说明)
	v.Set("mobile", "15059862915")                        //支持中国内地手机号码及固话号码。手机号：11位，示例：139****0000。固话号码：{区号}{号码}，示例：02151****29。接收手机号码，只能提交一个号码
	v.Set("content", "您的订单号是：9666。已由顺风快递发出，请注意查收。")       //根据发送方式不同：1、完整内容方式提交完整的短信内容，如：您的订单号是：9633。已由顺风快递发出，请注意查收。2、模板变量方式模板中的变量内容，多个变量以英文竖线（|）隔开①单变量示例模板内容：您的订单已发出，订单号：【变量】，请注意查收。参数写法：content=1234最终短信为：您的订单已发出，订单号：1234，请注意查收。②多变量示例模板内容：订单号：【变量1】，联系人：【变量2】，手机号：【变量3】，金额：【变量4】。参数写法：content=20180515006|张三|136xxxxxxxx|100元最终短信为：订单号：20180515006，联系人：张三，手机号：136xxxxxxxx，金额：100元。
	//v.Set("templateid", "1361")                           //语音模板ID（使用模板变量方式发送时必填）调试阶段可使用系统默认模板ID：1361（模板内容为：您的订单号是：【变量】。已由【变量】发出，请注意查收。）
	v.Set("time", s) //Unix时间戳（10位整型数字，当使用动态密码方式时为必填）

	body := strings.NewReader(v.Encode()) //把form数据编码
	client := &http.Client{}
	req, _ := http.NewRequest("POST", strUrl, body)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := client.Do(req) //发送
	if err != nil {
		fmt.Println(err)
		return
	}

	defer resp.Body.Close() //一定要关闭resp.Body
	res, _ := io.ReadAll(resp.Body)
	fmt.Println("Binantel:", string(res))
	//{"code":2,"msg":"提交成功","voiceid":"197362477"}
}
