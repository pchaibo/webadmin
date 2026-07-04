package binan

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"webadmin/config"
	"webadmin/model"

	"github.com/gorilla/websocket" //这里使用的是 gorilla 的 websocket 库
	redis "github.com/redis/go-redis/v9"
	"github.com/tidwall/gjson" //处理json
	"golang.org/x/net/proxy"
)

type Symbol struct {
	Symbol        string
	C             float64
	H             float64
	I             float64
	Change        float64
	Changepercent float64
}

var ch = make(chan int, 1)    //通道
var Redisclinet *redis.Client //创建redis连接池
var Coinname string
var Ctx = context.Background()

func init() {
	Redisclinet = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // 没有密码，默认值
		DB:       0,  // 默认DB 0
	})

}

func BinanPrice() {
	time.Sleep(time.Second * 2) //休眠2秒

	var conindata []model.Coin
	if err := model.Db.Model(&model.Coin{}).Where("status = ?", 1).Find(&conindata).Error; err != nil {
		fmt.Println("error ", err.Error())
		return
	}
	var co []string
	for _, v := range conindata {
		//coadd := v.Symbol + "@miniTicker"
		coadd := v.Symbol + "@ticker"
		co = append(co, coadd)

	}
	coinlist := strings.Join(co, "/")
	fmt.Println("coinlist:", coinlist)

	Coinname = coinlist
	ProxyEnabled = config.Get("ProxyEnabled") != ""
	go apisocket()
	for {
		<-ch                        //通道取数据
		time.Sleep(time.Second * 2) //休眠2秒
		apisocket()

	}

}

func apisocket() {
	dialer := websocket.Dialer{}
	//代理
	if ProxyEnabled {
		netDialer, err := proxy.SOCKS5("tcp", "127.0.0.1:1080", nil, nil)
		if err != nil {
			log.Println(err)
			Logs.Println("proxy tcp erro", err.Error())
		}
		dialer = websocket.Dialer{NetDial: netDialer.Dial}
	}
	//

	//创建一个拨号器，也可以用默认的 websocket.DefaultDialer
	//dialer := websocket.Dialer{}
	//向服务器发送连接请求，websocket 统一使用 ws://，默认端口和http一样都是80
	//coinstring := "btcusdt@miniTicker/ethusdt@miniTicker/bnbusdt@miniTicker/neousdt@miniTicker/ltcusdt@miniTicker/qtumusdt@miniTicker/adausdt@miniTicker/xrpusdt@miniTicker/eosusdt@miniTicker/tusdusdt@miniTicker/iotausdt@miniTicker/xlmusdt@miniTicker/ontusdt@miniTicker/trxusdt@miniTicker/etcusdt@miniTicker/icxusdt@miniTicker/venusdt@miniTicker/nulsusdt@miniTicker/vetusdt@miniTicker/paxusdt@miniTicker/bchabcusdt@miniTicker/bchsvusdt@miniTicker/usdcusdt@miniTicker/linkusdt@miniTicker/wavesusdt@miniTicker/bttusdt@miniTicker/usdsusdt@miniTicker/ongusdt@miniTicker/hotusdt@miniTicker/zilusdt@miniTicker/zrxusdt@miniTicker/fetusdt@miniTicker/batusdt@miniTicker/xmrusdt@miniTicker/zecusdt@miniTicker/iostusdt@miniTicker/celrusdt@miniTicker/dashusdt@miniTicker/nanousdt@miniTicker/omgusdt@miniTicker/thetausdt@miniTicker/enjusdt@miniTicker/mithusdt@miniTicker/maticusdt@miniTicker/atomusdt@miniTicker/tfuelusdt@miniTicker/oneusdt@miniTicker/ftmusdt@miniTicker/algousdt@miniTicker/usdsbusdt@miniTicker/gtousdt@miniTicker/erdusdt@miniTicker/dogeusdt@miniTicker/duskusdt@miniTicker/ankrusdt@miniTicker/winusdt@miniTicker/cosusdt@miniTicker/npxsusdt@miniTicker/cocosusdt@miniTicker/mtlusdt@miniTicker/tomousdt@miniTicker/perlusdt@miniTicker/dentusdt@miniTicker/mftusdt@miniTicker/keyusdt@miniTicker/stormusdt@miniTicker/dockusdt@miniTicker/wanusdt@miniTicker/funusdt@miniTicker/cvcusdt@miniTicker/chzusdt@miniTicker/bandusdt@miniTicker/busdusdt@miniTicker/beamusdt@miniTicker/xtzusdt@miniTicker/renusdt@miniTicker/rvnusdt@miniTicker/hcusdt@miniTicker/hbarusdt@miniTicker/nknusdt@miniTicker/stxusdt@miniTicker/kavausdt@miniTicker/arpausdt@miniTicker/iotxusdt@miniTicker/rlcusdt@miniTicker/mcousdt@miniTicker/ctxcusdt@miniTicker/bchusdt@miniTicker/troyusdt@miniTicker/viteusdt@miniTicker/fttusdt@miniTicker/eurusdt@miniTicker/ognusdt@miniTicker/drepusdt@miniTicker/bullusdt@miniTicker/bearusdt@miniTicker/ethbullusdt@miniTicker/ethbearusdt@miniTicker/tctusdt@miniTicker/wrxusdt@miniTicker/btsusdt@miniTicker/lskusdt@miniTicker/bntusdt@miniTicker/ltousdt@miniTicker/eosbullusdt@miniTicker/eosbearusdt@miniTicker/xrpbullusdt@miniTicker/xrpbearusdt@miniTicker/stratusdt@miniTicker/aionusdt@miniTicker/mblusdt@miniTicker/cotiusdt@miniTicker/bnbbullusdt@miniTicker/bnbbearusdt@miniTicker/stptusdt@miniTicker/wtcusdt@miniTicker/datausdt@miniTicker/xzcusdt@miniTicker/solusdt@miniTicker/ctsiusdt@miniTicker/hiveusdt@miniTicker/chrusdt@miniTicker/btcupusdt@miniTicker/btcdownusdt@miniTicker/gxsusdt@miniTicker/ardrusdt@miniTicker/lendusdt@miniTicker/mdtusdt@miniTicker/stmxusdt@miniTicker/kncusdt@miniTicker/repusdt@miniTicker/lrcusdt@miniTicker/pntusdt@miniTicker/compusdt@miniTicker/bkrwusdt@miniTicker/scusdt@miniTicker/zenusdt@miniTicker/snxusdt@miniTicker/ethupusdt@miniTicker/ethdownusdt@miniTicker/adaupusdt@miniTicker/adadownusdt@miniTicker/linkupusdt@miniTicker/linkdownusdt@miniTicker/vthousdt@miniTicker/dgbusdt@miniTicker/gbpusdt@miniTicker/sxpusdt@miniTicker/mkrusdt@miniTicker/daiusdt@miniTicker/dcrusdt@miniTicker/storjusdt@miniTicker/bnbupusdt@miniTicker/bnbdownusdt@miniTicker/xtzupusdt@miniTicker/xtzdownusdt@miniTicker/manausdt@miniTicker/audusdt@miniTicker/yfiusdt@miniTicker/balusdt@miniTicker/blzusdt@miniTicker/irisusdt@miniTicker/kmdusdt@miniTicker/jstusdt@miniTicker/srmusdt@miniTicker/antusdt@miniTicker/crvusdt@miniTicker/sandusdt@miniTicker/oceanusdt@miniTicker/nmrusdt@miniTicker/dotusdt@miniTicker/lunausdt@miniTicker/rsrusdt@miniTicker/paxgusdt@miniTicker/wnxmusdt@miniTicker/trbusdt@miniTicker/bzrxusdt@miniTicker/sushiusdt@miniTicker/yfiiusdt@miniTicker/ksmusdt@miniTicker/egldusdt@miniTicker/diausdt@miniTicker/runeusdt@miniTicker/fiousdt@miniTicker/umausdt@miniTicker/eosupusdt@miniTicker/eosdownusdt@miniTicker/trxupusdt@miniTicker/trxdownusdt@miniTicker/xrpupusdt@miniTicker/xrpdownusdt@miniTicker/dotupusdt@miniTicker/dotdownusdt@miniTicker/belusdt@miniTicker/wingusdt@miniTicker/ltcupusdt@miniTicker/ltcdownusdt@miniTicker/uniusdt@miniTicker/nbsusdt@miniTicker/oxtusdt@miniTicker/sunusdt@miniTicker/avaxusdt@miniTicker/hntusdt@miniTicker/flmusdt@miniTicker/uniupusdt@miniTicker/unidownusdt@miniTicker/ornusdt@miniTicker/utkusdt@miniTicker/xvsusdt@miniTicker/alphausdt@miniTicker/aaveusdt@miniTicker/nearusdt@miniTicker/sxpupusdt@miniTicker/sxpdownusdt@miniTicker/filusdt@miniTicker/filupusdt@miniTicker/fildownusdt@miniTicker/yfiupusdt@miniTicker/yfidownusdt@miniTicker/injusdt@miniTicker/audiousdt@miniTicker/ctkusdt@miniTicker/bchupusdt@miniTicker/bchdownusdt@miniTicker/akrousdt@miniTicker/axsusdt@miniTicker/hardusdt@miniTicker/dntusdt@miniTicker/straxusdt@miniTicker/unfiusdt@miniTicker/roseusdt@miniTicker/avausdt@miniTicker/xemusdt@miniTicker/aaveupusdt@miniTicker/aavedownusdt@miniTicker/sklusdt@miniTicker/susdusdt@miniTicker/sushiupusdt@miniTicker/sushidownusdt@miniTicker/xlmupusdt@miniTicker/xlmdownusdt@miniTicker/grtusdt@miniTicker/juvusdt@miniTicker/psgusdt@miniTicker/1inchusdt@miniTicker/reefusdt@miniTicker/ogusdt@miniTicker/atmusdt@miniTicker/asrusdt@miniTicker/celousdt@miniTicker/rifusdt@miniTicker/btcstusdt@miniTicker/truusdt@miniTicker/ckbusdt@miniTicker/twtusdt@miniTicker/firousdt@miniTicker/litusdt@miniTicker/sfpusdt@miniTicker/dodousdt@miniTicker/cakeusdt@miniTicker/acmusdt@miniTicker/badgerusdt@miniTicker/fisusdt@miniTicker/omusdt@miniTicker/pondusdt@miniTicker/degousdt@miniTicker/aliceusdt@miniTicker/linausdt@miniTicker/perpusdt@miniTicker/rampusdt@miniTicker/superusdt@miniTicker/cfxusdt@miniTicker/epsusdt@miniTicker/autousdt@miniTicker/tkousdt@miniTicker/pundixusdt@miniTicker/tlmusdt@miniTicker/1inchupusdt@miniTicker/1inchdownusdt@miniTicker/btgusdt@miniTicker/mirusdt@miniTicker/barusdt@miniTicker/forthusdt@miniTicker/bakeusdt@miniTicker/burgerusdt@miniTicker/slpusdt@miniTicker/shibusdt@miniTicker/icpusdt@miniTicker/arusdt@miniTicker/polsusdt@miniTicker/mdxusdt@miniTicker/maskusdt@miniTicker/lptusdt@miniTicker/nuusdt@miniTicker/xvgusdt@miniTicker/atausdt@miniTicker/gtcusdt@miniTicker/tornusdt@miniTicker/keepusdt@miniTicker/ernusdt@miniTicker/klayusdt@miniTicker/phausdt@miniTicker/bondusdt@miniTicker/mlnusdt@miniTicker/dexeusdt@miniTicker/c98usdt@miniTicker/clvusdt@miniTicker/qntusdt@miniTicker/flowusdt@miniTicker/tvkusdt@miniTicker"
	//coinstring = "!ticker@arr"
	//wss := "wss://stream.binance.com:9443/stream?streams="
	wssurl := "wss://fstream.binance.com/market/stream?streams=" + Coinname
	//wssurl := "wss://fstream.binance.com/market/stream?streams=bnbusdt@miniTicker/btcusdt@miniTicker"
	connect, _, err := dialer.Dial(wssurl, nil)
	if err != nil {
		ch <- 10
		fmt.Print("connecterro")
		//log.Println(err)
		return
	}
	//离开作用域关闭连接，go 的常规操作
	defer connect.Close()

	//定时向客户端发送数据
	//go TickWriter(connect)

	//启动数据读取循环，读取客户端发送来的数据
	for {
		//从 websocket 中读取数据
		//messageType 消息类型，websocket 标准
		//messageData 消息数据

		messageType, messageData, err := connect.ReadMessage()
		if err != nil {
			log.Println(err)
			break
		}
		switch messageType {
		case websocket.TextMessage: //文本数据
			Redis_set(string(messageData))
		case websocket.BinaryMessage: //二进制数据
			//fmt.Println(messageData)
		case websocket.PingMessage: //Ping
			Logs.Println("ping")
			connect.WriteMessage(websocket.PongMessage, nil)
		case websocket.PongMessage: //Pong
		case websocket.CloseMessage: //关闭
			Logs.Println("close")
			return
		default:

		}

	}
	ch <- 10 //退出
}

func Redis_set(str string) {
	//Logs.Println("str:", str)
	//处理json
	coinname := gjson.Get(str, "data.s")
	cointype := strings.ToLower(coinname.String()) //转小写
	var symbol Symbol
	symbol.Symbol = cointype
	symbol.C, _ = strconv.ParseFloat(gjson.Get(str, "data.c").String(), 64)
	symbol.H, _ = strconv.ParseFloat(gjson.Get(str, "data.h").String(), 64)
	symbol.I, _ = strconv.ParseFloat(gjson.Get(str, "data.l").String(), 64)
	symbol.Change, _ = strconv.ParseFloat(gjson.Get(str, "data.p").String(), 64)
	symbol.Changepercent, _ = strconv.ParseFloat(gjson.Get(str, "data.P").String(), 64)

	jsonstr, _ := json.Marshal(symbol) //返回json
	err := Redisclinet.HSet(Ctx, "coinprice", cointype, jsonstr).Err()
	if err != nil {
		fmt.Println("redis: ", err.Error())
		return
	}

}

func TickWriter(connect *websocket.Conn) {
	student := make(map[string]interface{})
	student["pong"] = 1492420473027
	b, err := json.Marshal(student)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(b)
	for {
		//向客户端发送类型为文本的数据
		err := connect.WriteMessage(websocket.PongMessage, []byte(b))
		if nil != err {
			log.Println(err)
			break
		}
		//休息一秒
		time.Sleep(time.Second * 5)
	}
}
