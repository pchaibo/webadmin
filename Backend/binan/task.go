package binan

import (
	"encoding/json"
	"fmt"
	"log"
	"math"
	"strings"
	"sync"
	"time"

	"webadmin/config"
	"webadmin/model"

	"github.com/tidwall/gjson"
)

const (
	BaseURL = "https://fapi.binance.com"
)

var Look sync.Mutex
var wg sync.WaitGroup
var Logs = log.Default()
var Proxytcp string
var ProxyEnabled bool

type userHeyue struct {
	User   *model.User
	Heyues []*model.Heyue
}

var Users []*userHeyue

func Taskuser() {
	Updateuser()
	for true {
		Updtecoinpric()
		time.Sleep(10 * time.Second)
	}

}

func Updtecoinpric() {
	var coin []model.Coin
	if err := model.Db.Model(&model.Coin{}).Find(&coin).Error; err != nil {
		Logs.Println(err.Error())
		return
	}

	for _, v := range coin {
		//cointype := strings.ToLower(.String()) //转小写
		reddate, err := Redisclinet.HGet(Ctx, "coinprice", v.Symbol).Result()
		if err != nil {
			continue
		}
		var symbol Symbol
		err = json.Unmarshal([]byte(reddate), &symbol)
		if err != nil {
			continue
		}
		var coinup model.Coin
		coinup.Close = symbol.C
		coinup.High = symbol.H
		coinup.Low = symbol.I
		model.Db.Model(&model.Coin{}).Where("id=?", v.Id).Updates(coinup)

	}
}
func init() {
	ProxyEnabled = config.Get("Proxyurl") != ""
}
func Task() {
	//Proxyurl = "http://" + config.Get("Proxyurl")
	//Updateuser()
	time.Sleep(2 * time.Second)
	for true {

		for _, v := range Users {
			if len(v.Heyues) < 1 {
				continue
			}
			wg.Add(1)
			go Userinfo(v)

		}
		wg.Wait()

		time.Sleep(6 * time.Second) //
		Updateuser()
	}

}

// 更新用户信息
func Updateuser() {
	var users []*model.User
	err := model.Db.Where("status=1 and bnaccess !='' and bnasecret !=''").Find(&users).Error
	if err != nil {
		Logs.Println(err.Error())
		return
	}

	Users = nil
	for _, v := range users {
		var heyuns []*model.Heyue
		model.Db.Where("status=1 and user_id=?", v.Id).Find(&heyuns)
		for _, heyun := range heyuns {
			v.Heyue = append(v.Heyue, heyun)
		}
		Look.Lock()
		Users = append(Users, &userHeyue{User: v, Heyues: heyuns})
		Look.Unlock()

	}

}

// check user 合约
func Userinfo(u *userHeyue) {
	defer wg.Done()
	//查询金额更新
	userbalance(u.User)
	//Logs.Println("userinfo Margin:  \n ", user.Margin)
	resdata, err := GetPositionRisk(u.User.Bnaccess, u.User.Bnasecret)
	if err != nil {
		Logs.Println("resdata error: ", err.Error())
		return
	}
	if len(resdata) < 1 {
		Logs.Println("resdata not null ") //查询无交易停用
		return
	}

	user := u.User
	heyues := u.Heyues
	for _, hey := range heyues {
		if len(resdata) > 0 {

			//查询持仓
			res := Checkheyun(user, hey, resdata)
			if res != 1 {
				Logs.Println("symbol  null: ", hey.Symbol)
				Addpositon(user, hey)
				continue

				//网格平仓
			} else if hey.Rangeclosing == 2 && hey.Is_num > 1 {
				Rangclosing(user, hey, resdata)
				continue
			}

		} else {
			//无持仓
			Checkadd(user, hey)
			continue
		}
	}

}

type Balance struct {
	//Asset              string // 资产 usdt 类型
	TotalWalletBalance      float64
	TotalCrossWalletBalance float64 //全仓余额
	TotalCrossUnPnl         float64 //全仓持仓未实现盈利
	AvailableBalance        float64 //下单可用余额
	MaxWithdrawAmount       float64 //最大可转出余额
}

// 账户余额
func userbalance(user *model.User) (rest int) {
	rest = 0
	if user.Margin < 1 {
		user.Margin = 1
	}
	//查账户
	res, body, err := Getbalance(user.Bnaccess, user.Bnasecret)
	if err != nil {
		Logs.Println("userinfo error: ", err.Error())
		return
	}
	if res != 1 {
		Logs.Println("body error: ", string(body))
		return
	}
	bodys := string(body)

	var balance Balance
	balance.TotalWalletBalance = gjson.Get(bodys, "totalWalletBalance").Float()
	balance.TotalCrossWalletBalance = gjson.Get(bodys, "totalCrossWalletBalance").Float()
	balance.TotalCrossUnPnl = gjson.Get(bodys, "totalCrossUnPnl").Float()
	balance.AvailableBalance = gjson.Get(bodys, "availableBalance").Float()
	balance.MaxWithdrawAmount = gjson.Get(bodys, "maxWithdrawAmount").Float()

	//fmt.Println("balance:", balance)
	jsonstr, _ := json.Marshal(balance) //返回json
	id := fmt.Sprintf("%d", user.Id)
	err = Redisclinet.HSet(Ctx, "user", id, jsonstr).Err()
	if err != nil {
		fmt.Println("redis: ", err.Error())
		return
	}

	return
}

func Checkuserinfo(user *model.User) (rest int) {
	rest = 0
	if user.Margin < 1 {
		user.Margin = 1
	}
	//查账户
	res, body, err := Getaccount(user.Bnaccess, user.Bnasecret)
	if err != nil {
		Logs.Println("userinfo error: ", err.Error())
		return
	}
	if res != 1 {
		Logs.Println("body error: ", string(body))
		return
	}
	bodys := string(body)
	totalInitialMargin := gjson.Get(bodys, "totalInitialMargin").Float()
	totalMarginBalance := gjson.Get(bodys, "totalMarginBalance").Float()
	marginmode := totalInitialMargin / totalMarginBalance * 10
	if marginmode < float64(user.Margin) {
		Logs.Printf("totalInitialMargin:%f totalMarginBalance:%f  marginmode:%f \n", totalInitialMargin, totalMarginBalance, marginmode)
		rest = 1
		return
	}
	return
}

// 风控
func Ckeckrisk(heyue *model.Heyue) (rest int) {
	rest = 0
	if heyue.Is_num > 1 && heyue.Risk == 2 && heyue.NewTime > 1000 {
		if (heyue.NewTime + int64(heyue.RiskTime)*60) > time.Now().Unix() {
			return
		} else {
			rest = 1
			return
		}
	}

	return
}

// checkuserinfo 加仓
func Checkadd(user *model.User, heyue *model.Heyue) {

	//保证金比例限制
	rest := Checkuserinfo(user)
	if rest != 1 {
		return
	}
	//调用风控
	if heyue.Risk == 2 && heyue.Is_num > 1 && heyue.NewTime > 100 {
		risk := Ckeckrisk(heyue)
		if risk != 1 {
			Logs.Println("风控时间写 ", user.Username, heyue.Symbol, heyue.Repeatprice, heyue.Side)
			return
		}
	}

	Logs.Println("Checkadd  : ", user.Username, heyue.Symbol, heyue.Oneprice, heyue.Repeatprice)
	Addpositon(user, heyue)

}

// 查询持仓
func Checkheyun(user *model.User, heyue *model.Heyue, resdata []PositionRisk) (res int) {
	var positionSide string
	res = 0
	for _, v := range resdata {
		if heyue.Side == 1 {
			positionSide = "LONG"
		} else {
			positionSide = "SHORT"
		}

		if v.Symbol == strings.ToUpper(heyue.Symbol) && v.positionSide == positionSide {
			res = 1
			sell := heyue.Sellprice * 1e-2
			Rangepercent := float64(heyue.Rangepercent) * 1e-2 * float64(heyue.Is_num) //百分比* 0.01 * 次数
			Marginpercentage := math.Abs(v.UnRealizedProfit) / v.InitialMargin
			//收益平仓
			if v.UnRealizedProfit > 0 && (v.UnRealizedProfit/v.InitialMargin) > sell {
				Logs.Println("平仓 id: ", heyue.Id)
				rest := Checkuserinfo(user) //查询保证金比例
				if rest == 1 {
					posite_red := Closeposition(user, heyue, v)
					if posite_red == 1 {
						heyue.Is_num = 0
						heyue.Newprice = 0
						Addpositon(user, heyue) //马上加仓
					}
					return
				} else {
					Logs.Println("仓位过重，平仓 id: ", heyue.Id, heyue.Symbol)
					return
				}

				//网格类型 1:差价usdt计算
			} else if heyue.Rangetype == 1 && heyue.Newprice > 0 && heyue.Rangeprice > 0 && heyue.Is_num < heyue.Num {
				//网格加仓做多
				if positionSide == "LONG" && (heyue.Newprice-heyue.Rangeprice) > v.MarkPrice {
					Checkadd(user, heyue)
					//网格做空
				} else if positionSide == "SHORT" && (heyue.Newprice+heyue.Rangeprice) < v.MarkPrice {
					Checkadd(user, heyue)
				}

				//网格类型 2:保证金百分比计算
			} else if heyue.Rangetype == 2 && heyue.Is_num < heyue.Num && v.UnRealizedProfit < 0 && Marginpercentage > Rangepercent {
				if Rangepercent <= 0 {
					Logs.Println("保证金百分比小于0:", Rangepercent, v.UnRealizedProfit, Marginpercentage)
					continue
				}
				//v.MarkPrice < heyue.Newprice
				if positionSide == "LONG" {
					Logs.Println("保证金百分比计算LONG:", Rangepercent, v.UnRealizedProfit, Marginpercentage, v.MarkPrice)
					Checkadd(user, heyue)
					//&& v.MarkPrice > heyue.Newprice
				} else if positionSide == "SHORT" {
					Logs.Println("保证金百分比计算SHORT:", Rangepercent, v.UnRealizedProfit, Marginpercentage, v.MarkPrice)
					Checkadd(user, heyue)
				}

			}

			return
		}

	}
	return

}

// 网格平仓
func Rangclosing(user *model.User, heyue *model.Heyue, resdata []PositionRisk) {
	var positionSide string

	for _, v := range resdata {
		if heyue.Side == 1 {
			positionSide = "LONG"
		} else {
			positionSide = "SHORT"
		}
		if v.Symbol == strings.ToUpper(heyue.Symbol) && v.positionSide == positionSide {
			//取网格日志
			var order = model.Heyueorder{}
			model.Db.Model(&model.Heyueorder{}).Where("ordertype=1 and  symbol=? and  side =? and user_id=?  and num=? ", heyue.Symbol, heyue.Side, heyue.UserId, heyue.Is_num).Order("id desc").First(&order)
			//Logs.Println("网格 log: ", order)
			if order.Price <= 0 || heyue.Newprice <= 0 {
				continue
			}

			Quantityprice := order.Quantity * order.Price                                     //原价
			user_LONG := Quantityprice + Quantityprice/20*float64(heyue.Rangeclosingpct)*1e-2 //除0倍* 百分比
			user_SHORT := Quantityprice - Quantityprice/20*float64(heyue.Rangeclosingpct)*1e-2

			newprcie := v.MarkPrice * order.Quantity
			v.PositionAmt = order.Quantity
			v.UnRealizedProfit = 0
			//做多
			if heyue.Side == 1 && newprcie > user_LONG && v.MarkPrice > order.Price && v.MarkPrice > heyue.Newprice {
				Logs.Println("newprcie:", newprcie)
				Logs.Println("user_LONG:", user_LONG)
				Logs.Println("网格平多 log: ", order)
				rest := Checkuserinfo(user)
				if rest == 1 {
					RangCloseposition(user, heyue, v)
				}

			} else if heyue.Side == 2 && newprcie < user_SHORT && v.MarkPrice < order.Price && v.MarkPrice < heyue.Newprice {
				Logs.Println("newprcie:", newprcie)
				Logs.Println("user_SHORT:", user_SHORT)
				Logs.Println("网格平空 log: ", order)
				rest := Checkuserinfo(user)
				if rest == 1 {
					RangCloseposition(user, heyue, v)
				}
			}
		}

	}
}
