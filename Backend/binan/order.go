package binan

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"math"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/gotify/server/v2/model"
	"github.com/jinzhu/gorm"
	"github.com/tidwall/gjson"
)

func Order(ApiKey, SecretKey string) {
	// 合约市场和下单信息
	symbol := "ETHUSDT" // 合约交易对
	quantity := 0.01    // 下单数量
	//price := "30000"    // 限价单的价格
	// 做多（买入）
	response, err := placeOrder(symbol, "BUY", "LONG", ApiKey, SecretKey, quantity)
	if err != nil {
		log.Fatalf("Error placing buy order: %v", err)
	}
	fmt.Printf("Buy Order Response: %s\n", response)

	// 做空（卖出）
	response, err = placeOrder(symbol, "SELL", "SHORT", ApiKey, SecretKey, quantity)
	if err != nil {
		log.Fatalf("Error placing sell order: %v", err)
	}
	fmt.Printf("Sell Order Response: %s\n", response)
}

// RoundTo 保留 n 位小数并返回 float64
func RoundTo(v float64, n int) float64 {
	pow := math.Pow(10, float64(n))
	return math.Round(v*pow) / pow
}

// 加仓
func Addpositon(user *model.User, heyue *model.Heyue) {
	var Coin model.Coin
	err := Db.DB.Where("status=1 and symbol= ?", heyue.Symbol).First(&Coin).Error
	if err != nil {
		Logs.Println("Addpositon Coin error ", err.Error())
		return
	}
	var num float64
	var total float64
	if heyue.Is_num == 0 {
		num = heyue.Onepric / Coin.Close
		total = heyue.Onepric
	} else {
		num = heyue.Repeatpric / Coin.Close
		total = heyue.Repeatpric
	}
	//计算数量
	quantity := RoundTo(num, int(Coin.Quantityprecision)) //数量
	//https://fapi.binance.com/fapi/v1/exchangeInfo 小数位数
	if quantity == 0 {
		Logs.Println("add order quantity null : ", heyue.Symbol, quantity)
		return
	}
	symbol := strings.ToUpper(heyue.Symbol)

	var newheyue model.Heyue
	newheyue.Is_num = heyue.Is_num + 1
	newheyue.Newpric = Coin.Close
	newheyue.NewTime = time.Now().Unix()
	//LONG
	if heyue.Side == 1 {
		response, err := placeOrder(symbol, "BUY", "LONG", user.Bnaccess, user.Bnasecret, quantity)
		if err != nil {
			Logs.Println("Addpositon   Error  order: ", err.Error())
			return
		}
		Logs.Println("加仓 LONG:", response)
		orderId := gjson.Get(response, "orderId").Int()
		if orderId > 0 {
			Db.DB.Model(&model.Heyue{}).Where("id = ?", heyue.Id).Select("is_num", "newpric", "new_time").Updates(&newheyue)
			var order model.Heyuesorder
			order.Total = total
			order.Ordertype = 1
			order.Quantity = quantity
			price := math.Round(Coin.Close*10000) / 10000
			order.Num = newheyue.Is_num
			order.Price = price
			order.Side = heyue.Side
			order.Orderid = orderId
			order.UserId = user.ID
			order.Username = user.Name
			order.Symbol = heyue.Symbol
			order.Log = response
			order.Status = 1
			order.AddTime = time.Now().Unix()
			Db.HeyuesordersCreate(&order)
		}
		return
		//SHORT
	} else if heyue.Side == 2 {
		response, err := placeOrder(symbol, "SELL", "SHORT", user.Bnaccess, user.Bnasecret, quantity)
		if err != nil {
			Logs.Println("Addpositon   Error  order: ", err.Error())
			return
		}
		Logs.Println("加仓 SHORT:", response)
		orderId := gjson.Get(response, "orderId").Int()
		if orderId > 0 {
			Db.DB.Model(&model.Heyue{}).Where("id = ?", heyue.Id).Select("is_num", "newpric", "new_time").Updates(&newheyue)
			var order model.Heyuesorder
			order.Total = total
			order.Ordertype = 1
			order.Quantity = quantity
			price := math.Round(Coin.Close*10000) / 10000
			order.Num = newheyue.Is_num
			order.Price = price
			order.Side = heyue.Side
			order.Orderid = orderId
			order.UserId = user.ID
			order.Username = user.Name
			order.Symbol = heyue.Symbol
			order.Log = response
			order.Status = 1
			order.AddTime = time.Now().Unix()
			Db.HeyuesordersCreate(&order)
		}
		return
	}

}

// 网格平仓
func RangCloseposition(user *model.User, heyue *model.Heyue, resdata PositionRisk) (rest int) {
	rest = 0
	newpric := math.Round(resdata.MarkPrice*10000) / 10000
	if heyue.Side == 1 {
		response, err := placeOrder(resdata.Symbol, "SELL", "LONG", user.Bnaccess, user.Bnasecret, resdata.PositionAmt)
		if err != nil {
			Logs.Println("Closeposition Error  order: ", err.Error())
			return
		}
		Logs.Println("网格平仓 LONG: ", response)
		//更新数据
		orderId := gjson.Get(response, "orderId").Int()
		if orderId > 0 {
			rest = 1
			Db.DB.Model(&model.Heyue{}).Where("id = ?", heyue.Id).UpdateColumn("is_num", gorm.Expr("is_num - ?", 1)).UpdateColumn("newpric", newpric)
			var order model.Heyuesorder
			order.Ordertype = 2
			order.Quantity = resdata.PositionAmt
			price := newpric
			order.Price = price
			order.Total = price * resdata.PositionAmt
			order.Side = heyue.Side
			order.Orderid = orderId
			order.UserId = user.ID
			order.Username = user.Name
			order.Symbol = heyue.Symbol
			order.Log = response
			order.Usdt = 0
			order.Status = 1
			order.Num = heyue.Is_num
			order.AddTime = time.Now().Unix()
			Db.HeyuesordersCreate(&order)
		}

		return
	} else if heyue.Side == 2 {
		response_short, err := placeOrder(resdata.Symbol, "BUY", "SHORT", user.Bnaccess, user.Bnasecret, resdata.PositionAmt)
		if err != nil {
			Logs.Printf("Error placing buy order: %v", err)
			return
		}
		Logs.Println("网格平仓 SHORT：", response_short)
		orderId := gjson.Get(response_short, "orderId").Int()
		if orderId > 0 {
			rest = 1
			Db.DB.Model(&model.Heyue{}).Where("id = ?", heyue.Id).UpdateColumn("is_num", gorm.Expr("is_num - ?", 1)).UpdateColumn("newpric", newpric)
			var order model.Heyuesorder
			order.Ordertype = 2
			order.Quantity = resdata.PositionAmt
			price := newpric
			order.Price = price
			order.Total = price * resdata.PositionAmt
			order.Side = heyue.Side
			order.Orderid = orderId
			order.UserId = user.ID
			order.Username = user.Name
			order.Symbol = heyue.Symbol
			order.Status = 1
			order.Log = response_short
			order.Usdt = 0
			order.Num = heyue.Is_num
			order.AddTime = time.Now().Unix()
			Db.HeyuesordersCreate(&order)

		}
		return
	}
	return
}

// 平仓
func Closeposition(user *model.User, heyue *model.Heyue, resdata PositionRisk) (rest int) {
	rest = 0
	if heyue.Side == 1 {
		response, err := placeOrder(resdata.Symbol, "SELL", "LONG", user.Bnaccess, user.Bnasecret, resdata.PositionAmt)
		if err != nil {
			Logs.Println("Closeposition Error  order: ", err.Error())
			return
		}
		Logs.Println("平仓 LONG: ", response)
		//更新数据
		orderId := gjson.Get(response, "orderId").Int()
		if orderId > 0 {
			rest = 1
			Db.DB.Model(&model.Heyue{}).Where("id = ?", heyue.Id).UpdateColumn("is_num", 0).UpdateColumn("newpric", 0)
			var order model.Heyuesorder
			order.Ordertype = 2
			order.Quantity = resdata.PositionAmt
			price := math.Round(resdata.MarkPrice*10000) / 10000
			order.Price = price
			order.Total = price * resdata.PositionAmt
			order.Side = heyue.Side
			order.Orderid = orderId
			order.UserId = user.ID
			order.Username = user.Name
			order.Symbol = heyue.Symbol
			order.Log = response
			order.Usdt = resdata.UnRealizedProfit
			order.Status = 1
			order.Num = 1
			order.AddTime = time.Now().Unix()
			Db.HeyuesordersCreate(&order)
		}

		return
	} else if heyue.Side == 2 {
		response_short, err := placeOrder(resdata.Symbol, "BUY", "SHORT", user.Bnaccess, user.Bnasecret, resdata.PositionAmt)
		if err != nil {
			Logs.Printf("Error placing buy order: %v", err)
			return
		}
		Logs.Println("平仓 SHORT：", response_short)
		orderId := gjson.Get(response_short, "orderId").Int()
		if orderId > 0 {
			rest = 1
			Db.DB.Model(&model.Heyue{}).Where("id = ?", heyue.Id).UpdateColumn("is_num", 0).UpdateColumn("newpric", 0)
			var order model.Heyuesorder
			order.Ordertype = 2
			order.Quantity = resdata.PositionAmt
			price := math.Round(resdata.MarkPrice*10000) / 10000
			order.Price = price
			order.Total = price * resdata.PositionAmt
			order.Side = heyue.Side
			order.Orderid = orderId
			order.UserId = user.ID
			order.Username = user.Name
			order.Symbol = heyue.Symbol
			order.Status = 1
			order.Num = 1
			order.Log = response_short
			order.Usdt = resdata.UnRealizedProfit
			order.AddTime = time.Now().Unix()
			Db.HeyuesordersCreate(&order)

		}
		return
	}
	return
}

func ordersign(params url.Values, secretKey string) string {
	// 对参数进行签名
	mac := hmac.New(sha256.New, []byte(secretKey))
	mac.Write([]byte(params.Encode()))
	return hex.EncodeToString(mac.Sum(nil))
}

func placeOrder(symbol, side, position, ApiKey, SecretKey string, quantity float64) (string, error) {
	// 构造请求参数
	params := url.Values{}
	params.Set("symbol", symbol)
	params.Set("side", side)
	params.Set("type", "MARKET") //价格类型 MARKET
	quan := strconv.FormatFloat(math.Abs(quantity), 'f', -1, 64)
	Logs.Println("quan:", symbol, quan)
	params.Set("quantity", quan) //数量
	//params.Set("price", price)       // 如果是限价单，价格必填
	params.Set("positionSide", position)
	params.Set("recvWindow", "5000")
	params.Set("timestamp", fmt.Sprintf("%d", time.Now().UnixMilli()))

	// 生成签名
	signature := ordersign(params, SecretKey)
	params.Set("signature", signature)

	// 构造请求
	req, err := http.NewRequest("POST", BaseURL+"/fapi/v1/order", bytes.NewBufferString(params.Encode()))
	if err != nil {
		return "", err
	}

	// 设置请求头
	req.Header.Set("X-MBX-APIKEY", ApiKey)

	// 发送请求
	//client := &http.Client{}
	client := Userclient()
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// 读取响应
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	// 返回响应内容
	return string(body), nil
}
