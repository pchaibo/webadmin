package binan

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"webadmin/model"

	"github.com/tidwall/gjson"
	"gorm.io/gorm"
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

// 更新价格
func Updateheyue(heyue *model.Heyue, v PositionRisk) {
	var newheyue model.Heyue
	newheyue.Is_num = 1
	var uewprice float64
	if v.EntryPrice > 10 {
		uewprice = RoundTo(v.EntryPrice, 3)
	} else {
		uewprice = RoundTo(v.EntryPrice, 6)
	}
	newheyue.Newprice = uewprice
	newheyue.NewTime = time.Now().Unix()
	if v.EntryPrice > 0 && heyue.Id > 0 {
		model.Db.Model(&model.Heyue{}).Where("id = ?", heyue.Id).Select("is_num", "newprice", "new_time").Updates(&newheyue)
	}

}
func getcoinpric(heyue *model.Heyue) (model.Coin, error) {
	var coinup model.Coin
	err := model.Db.Where("status=1 and symbol= ?", heyue.Symbol).First(&coinup).Error
	if err != nil {
		Logs.Println("getcoinpric Coin error ", err.Error())
		return coinup, err
	}
	reddate, err := Redisclinet.HGet(Ctx, "coinprice", heyue.Symbol).Result()
	if err != nil {
		return coinup, err
	}
	var symbol Symbol
	err = json.Unmarshal([]byte(reddate), &symbol)
	if err != nil {
		return coinup, err
	}

	coinup.Close = symbol.C
	coinup.High = symbol.H
	coinup.Low = symbol.I
	return coinup, err
}

// 加仓
func Addpositon(user *model.User, heyue *model.Heyue) {
	//自动开仓
	// if heyue.Autoadd != 1 {
	// 	return
	// }
	var Coin model.Coin
	err := model.Db.Where("status=1 and symbol= ?", heyue.Symbol).First(&Coin).Error
	if err != nil {
		Logs.Println("Addpositon Coin error ", err.Error())
		return
	}
	var num float64
	var total float64
	if heyue.Is_num == 0 {

		//预估价格-做多减仓
		if heyue.Side == 1 && Coin.Close > heyue.TopPrice && heyue.TopPrice > 0 {
			//不开仓
			if heyue.ReductionRatio >= 100 {
				//Logs.Println("ReductionRatio gt  100 ")
				return
			}
			num = (heyue.Oneprice - heyue.Oneprice*float64(heyue.ReductionRatio)*0.01) / Coin.Close
			total = heyue.Oneprice

			//预估价格-做空减仓
		} else if heyue.Side == 2 && Coin.Close < heyue.TopPrice && heyue.TopPrice > 0 {
			//不开仓
			if heyue.ReductionRatio >= 100 {
				//Logs.Println("ReductionRatio gt  100 ")
				return
			}
			num = (heyue.Oneprice - heyue.Oneprice*float64(heyue.ReductionRatio)*0.01) / Coin.Close
			total = heyue.Oneprice

		} else {
			num = heyue.Oneprice / Coin.Close
			total = heyue.Oneprice
		}

	} else {
		num = heyue.Repeatprice / Coin.Close
		total = heyue.Repeatprice
	}
	//计算数量
	//quantity := RoundTo(num, int(Coin.Priceprecision)) //数量
	quantity := RoundTo(num, int(Coin.Quantityprecision)) //数量

	//https://fapi.binance.com/fapi/v1/exchangeInfo 小数位数
	if quantity == 0 {
		Logs.Println("add order quantity null : ", heyue.Symbol, quantity)
		return
	}
	symbol := strings.ToUpper(heyue.Symbol)

	var newheyue model.Heyue
	newheyue.Is_num = heyue.Is_num + 1
	newheyue.Newprice = Coin.Close
	newheyue.NewTime = time.Now().Unix()
	newheyue.Closeposition = 1
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
			model.Db.Model(&model.Heyue{}).Where("id = ?", heyue.Id).Select("is_num", "newprice", "new_time", "closeposition").Updates(&newheyue)
			var order model.Heyueorder
			order.Total = total
			order.Ordertype = 1
			order.Quantity = quantity
			price := math.Round(Coin.Close*10000) / 10000
			order.Num = newheyue.Is_num
			order.Price = price
			order.Side = heyue.Side
			order.Orderid = orderId
			order.UserId = uint(user.Id)
			order.Username = user.Username
			order.Symbol = heyue.Symbol
			order.Log = response
			order.Status = 1
			order.AddTime = time.Now().Unix()
			model.Db.Create(&order)
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
			model.Db.Model(&model.Heyue{}).Where("id = ?", heyue.Id).Select("is_num", "newprice", "new_time", "closeposition").Updates(&newheyue)
			var order model.Heyueorder
			order.Total = total
			order.Ordertype = 1
			order.Quantity = quantity
			price := math.Round(Coin.Close*10000) / 10000
			order.Num = newheyue.Is_num
			order.Price = price
			order.Side = heyue.Side
			order.Orderid = orderId
			order.UserId = uint(user.Id)
			order.Username = user.Username
			order.Symbol = heyue.Symbol
			order.Log = response
			order.Status = 1
			order.AddTime = time.Now().Unix()
			model.Db.Create(&order)
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
			model.Db.Model(&model.Heyue{}).Where("id = ?", heyue.Id).UpdateColumn("is_num", gorm.Expr("is_num - ?", 1)).UpdateColumn("newprice", newpric)
			var order model.Heyueorder
			order.Ordertype = 2
			order.Quantity = resdata.PositionAmt
			price := newpric
			order.Price = price
			order.Total = price * resdata.PositionAmt
			order.Side = heyue.Side
			order.Orderid = orderId
			order.UserId = uint(user.Id)
			order.Username = user.Username
			order.Symbol = heyue.Symbol
			order.Log = response
			order.Usdt = 0
			order.Status = 1
			order.Num = heyue.Is_num
			order.AddTime = time.Now().Unix()
			model.Db.Create(&order)
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
			model.Db.Model(&model.Heyue{}).Where("id = ?", heyue.Id).UpdateColumn("is_num", gorm.Expr("is_num - ?", 1)).UpdateColumn("newprice", newpric)
			var order model.Heyueorder
			order.Ordertype = 2
			order.Quantity = resdata.PositionAmt
			price := newpric
			order.Price = price
			order.Total = price * resdata.PositionAmt
			order.Side = heyue.Side
			order.Orderid = orderId
			order.UserId = uint(user.Id)
			order.Username = user.Username
			order.Symbol = heyue.Symbol
			order.Status = 1
			order.Log = response_short
			order.Usdt = 0
			order.Num = heyue.Is_num
			order.AddTime = time.Now().Unix()
			model.Db.Create(&order)

		}
		return
	}
	return
}
func DecimalPlaces(num float64) int {
	s := strconv.FormatFloat(num, 'f', -1, 64)

	if !strings.Contains(s, ".") {
		return 0
	}

	return len(strings.Split(s, ".")[1])
}

// RoundTo 保留指定小数位
func Roundto(num float64, precision int) float64 {
	factor := math.Pow(10, float64(precision))
	return math.Round(num*factor) / factor
}

// 平仓
func Closeposition(user *model.User, heyue *model.Heyue, resdata PositionRisk) (rest int) {
	rest = 0
	var quantity float64
	quantity = resdata.PositionAmt
	quantity_num := DecimalPlaces(quantity)
	if heyue.Side == 1 {
		//按比例平仓
		if heyue.Closepositionrate < 100 {
			PositionAmt := quantity * float64(heyue.Closepositionrate) * 0.01
			quantity = Roundto(PositionAmt, quantity_num)
		}
		response, err := placeOrder(resdata.Symbol, "SELL", "LONG", user.Bnaccess, user.Bnasecret, quantity)
		if err != nil {
			Logs.Println("Closeposition Error  order: ", err.Error())
			return
		}
		Logs.Println("平仓 LONG: ", response)
		//更新数据
		orderId := gjson.Get(response, "orderId").Int()
		if orderId > 0 {
			rest = 1
			//重置网格 repeatprice
			if heyue.Resset == 1 && heyue.Resnum > 0 && heyue.Resrangepercent > 0 {
				model.Db.Model(&model.Heyue{}).Where("id = ?", heyue.Id).UpdateColumn("is_num", 0).UpdateColumn("newprice", 0).UpdateColumn("num", heyue.Resnum).UpdateColumn("rangepercent", heyue.Resrangepercent).UpdateColumn("repeatprice", heyue.Resrepeatprice)
			} else {
				model.Db.Model(&model.Heyue{}).Where("id = ?", heyue.Id).UpdateColumn("is_num", 0).UpdateColumn("closeposition", 2).UpdateColumn("newprice", 0)
			}

			var order model.Heyueorder
			order.Ordertype = 2
			order.Quantity = resdata.PositionAmt
			price := math.Round(resdata.MarkPrice*10000) / 10000
			order.Price = price
			order.Total = price * resdata.PositionAmt
			order.Side = heyue.Side
			order.Orderid = orderId
			order.UserId = uint(user.Id)
			order.Username = user.Username
			order.Symbol = heyue.Symbol
			order.Log = response
			order.Usdt = resdata.UnRealizedProfit
			order.Status = 1
			order.Num = 1
			order.AddTime = time.Now().Unix()
			model.Db.Create(&order)
		}

		return
	} else if heyue.Side == 2 {
		if heyue.Closepositionrate < 100 {
			PositionAmt := quantity * float64(heyue.Closepositionrate) * 0.01
			quantity = Roundto(PositionAmt, quantity_num)
		}
		response_short, err := placeOrder(resdata.Symbol, "BUY", "SHORT", user.Bnaccess, user.Bnasecret, quantity)
		if err != nil {
			Logs.Printf("Error placing buy order: %v", err)
			return
		}
		Logs.Println("平仓 SHORT：", response_short)
		orderId := gjson.Get(response_short, "orderId").Int()
		if orderId > 0 {
			rest = 1
			//重置网格 repeatprice
			if heyue.Resset == 1 && heyue.Resnum > 0 && heyue.Resrangepercent > 0 {
				model.Db.Model(&model.Heyue{}).Where("id = ?", heyue.Id).UpdateColumn("is_num", 0).UpdateColumn("newprice", 0).UpdateColumn("num", heyue.Resnum).UpdateColumn("rangepercent", heyue.Resrangepercent).UpdateColumn("repeatprice", heyue.Resrepeatprice)
			} else {
				model.Db.Model(&model.Heyue{}).Where("id = ?", heyue.Id).UpdateColumn("is_num", 0).UpdateColumn("closeposition", 2).UpdateColumn("newprice", 0)
			}
			var order model.Heyueorder
			order.Ordertype = 2
			order.Quantity = resdata.PositionAmt
			price := math.Round(resdata.MarkPrice*10000) / 10000
			order.Price = price
			order.Total = price * resdata.PositionAmt
			order.Side = heyue.Side
			order.Orderid = orderId
			order.UserId = uint(user.Id)
			order.Username = user.Username
			order.Symbol = heyue.Symbol
			order.Status = 1
			order.Num = 1
			order.Log = response_short
			order.Usdt = resdata.UnRealizedProfit
			order.AddTime = time.Now().Unix()
			model.Db.Create(&order)

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
