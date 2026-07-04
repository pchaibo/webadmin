package binan

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"time"

	"github.com/tidwall/gjson"
)

type PositionRisk struct {
	Symbol           string  //交易对
	positionSide     string  //方向
	PositionAmt      float64 //持仓数量 头寸数量，符号代表多空方向, 正数为多，负数为空
	EntryPrice       float64 //开仓均价
	MarkPrice        float64 //价格
	InitialMargin    float64 //保证金
	UnRealizedProfit float64 //盈亏

}

// 账户余额
func Getbalance(ApiKey, SecretKey string) (res int, body []byte, err error) {
	res = 0
	endpoint := "/fapi/v3/account"

	params := url.Values{}
	params.Set("timestamp", fmt.Sprintf("%d", time.Now().UnixMilli()))
	params.Set("recvWindow", "10000")
	queryString := params.Encode()
	signature := Sign(queryString, SecretKey)

	fullURL := fmt.Sprintf("%s%s?%s&signature=%s", BaseURL, endpoint, queryString, signature)
	req, err := http.NewRequest("GET", fullURL, nil)
	if err != nil {
		fmt.Println(err.Error())
		return

	}

	req.Header.Set("X-MBX-APIKEY", ApiKey)

	client := Userclient()

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	defer resp.Body.Close()

	body, err = io.ReadAll(resp.Body)

	if resp.StatusCode != 200 {
		fmt.Println("Error:", resp.Status)
		fmt.Println(string(body))
		return
	}
	res = 1

	return

}

// 账户信息
func Getaccount(ApiKey, SecretKey string) (res int, body []byte, err error) {
	res = 0
	endpoint := "/fapi/v3/account"

	params := url.Values{}
	params.Set("timestamp", fmt.Sprintf("%d", time.Now().UnixMilli()))
	params.Set("recvWindow", "10000")
	queryString := params.Encode()
	signature := Sign(queryString, SecretKey)

	fullURL := fmt.Sprintf("%s%s?%s&signature=%s", BaseURL, endpoint, queryString, signature)
	req, err := http.NewRequest("GET", fullURL, nil)
	if err != nil {
		fmt.Println(err.Error())
		return

	}

	req.Header.Set("X-MBX-APIKEY", ApiKey)

	client := Userclient()

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	defer resp.Body.Close()

	body, err = io.ReadAll(resp.Body)

	if resp.StatusCode != 200 {
		fmt.Println("Error:", resp.Status)
		fmt.Println(string(body))
		return
	}
	res = 1

	return

}

// https://developers.binance.com/docs/zh-CN/derivatives/usds-margined-futures/trade/rest-api/Position-Information-V2
func GetPositionRisk(ApiKey, SecretKey string) (riskarr []PositionRisk, err error) {
	endpoint := "/fapi/v3/positionRisk"

	params := url.Values{}
	params.Set("timestamp", fmt.Sprintf("%d", time.Now().UnixMilli()))
	//params.Set("symbol", "ethusdt")
	params.Set("recvWindow", "10000")
	queryString := params.Encode()
	signature := Sign(queryString, SecretKey)

	fullURL := fmt.Sprintf("%s%s?%s&signature=%s", BaseURL, endpoint, queryString, signature)

	req, err := http.NewRequest("GET", fullURL, nil)
	if err != nil {
		fmt.Println(err.Error())
		return

	}

	req.Header.Set("X-MBX-APIKEY", ApiKey)

	client := Userclient()

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	if resp.StatusCode != 200 {
		fmt.Println("Error:", resp.Status)
		fmt.Println(string(body))
		return
	}

	// var positions []map[string]interface{}
	// json.Unmarshal(body, &positions)
	// out, _ := json.MarshalIndent(positions, "", "  ")
	// fmt.Println(string(out))

	//var riskarr []PositionRisk
	fmt.Println("\n ")
	gjson.Parse(string(body)).ForEach(func(_, v gjson.Result) bool {
		if v.Get("positionAmt").String() != "0" {
			// fmt.Printf(" %s 方向: %s 价格：%f 持仓: %f  保证金: %f 盈亏:%f  \n",
			// 	v.Get("symbol").String(),
			// 	v.Get("positionSide").String(),
			// 	v.Get("markPrice").Float(),
			// 	v.Get("positionAmt").Float(),
			// 	v.Get("initialMargin").Float(),
			// 	v.Get("unRealizedProfit").Float(),
			// )
			var riskdata PositionRisk
			riskdata.Symbol = v.Get("symbol").String()
			riskdata.positionSide = v.Get("positionSide").String()
			riskdata.PositionAmt = v.Get("positionAmt").Float()
			riskdata.EntryPrice = v.Get("entryPrice").Float()
			riskdata.MarkPrice = v.Get("markPrice").Float()
			riskdata.InitialMargin = v.Get("initialMargin").Float()
			riskdata.UnRealizedProfit = v.Get("unRealizedProfit").Float()

			riskarr = append(riskarr, riskdata)
		}

		return true
	})

	return
}

// 生成签名
func Sign(data, SecretKey string) string {
	h := hmac.New(sha256.New, []byte(SecretKey))
	h.Write([]byte(data))
	return hex.EncodeToString(h.Sum(nil))
}

func Userclient() (client *http.Client) {
	// ======== 配置 HTTP Client（支持代理） =========
	proxyAddr := "http://127.0.0.1:1080"
	//proxyAddr := Proxyurl
	transport := &http.Transport{
		DialContext: (&net.Dialer{
			Timeout: 10 * time.Second,
		}).DialContext,
		TLSHandshakeTimeout: 10 * time.Second,
	}

	// 如果传入代理地址，则启用代理
	if ProxyEnabled {
		proxyURL, err := url.Parse(proxyAddr)
		if err != nil {
			fmt.Printf("proxy parse error: %v", err)
			return
		}
		transport.Proxy = http.ProxyURL(proxyURL)
	}

	//var client *http.Client
	client = &http.Client{
		Transport: transport,
		Timeout:   15 * time.Second,
	}

	//空不走代理
	if !ProxyEnabled {
		client = &http.Client{}
	}
	return

}
