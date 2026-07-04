package controller

import (
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"net/http"

	"webadmin/binan"
	"webadmin/config"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}
var Clients = make(map[*websocket.Conn]bool)
var mu sync.Mutex

// 拉取数据
func Queuedata() (data []byte, err error) {
	Redisclinet := binan.Redisclinet
	reddata, _ := Redisclinet.HGetAll(binan.Ctx, "coinprice").Result()
	var adddate []string
	for _, v := range reddata {
		adddate = append(adddate, v)
	}
	resdata := make(map[string]interface{})
	resdata["status"] = 1
	resdata["data"] = adddate
	// 输出当前记录
	//fmt.Println("transaction_new:", adddate)
	data, err = json.Marshal(resdata)
	return
}

// 发送数据
func Broadcaster() {
	ticker := time.NewTicker(1 * time.Second)
	for {
		<-ticker.C
		mu.Lock()

		data, _ := Queuedata() //拉取存的数据
		for conn := range Clients {
			err := conn.WriteMessage(websocket.TextMessage, data)
			if err != nil {
				conn.Close()
				delete(Clients, conn)
			}
			//
			if err := conn.WriteMessage(websocket.PingMessage, []byte("ping")); err != nil {
				// 如果发送失败，比如客户端断开，就退出
				fmt.Println("err:", err.Error())

			}
		}
		mu.Unlock()
	}
}

func WsHandler(c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		fmt.Println("WebSocket 升级失败:", err)
		return
	}

	fmt.Println("客户端已连接")

	// 添加新连接，移除旧连接
	mu.Lock()
	Clients[conn] = true
	mu.Unlock()

	// 监听关闭

	// 设置初始读超时时间（例如 30 秒）
	_ = conn.SetReadDeadline(time.Now().Add(30 * time.Second))

	// 每收到一个 Pong 就延长超时时间
	conn.SetPongHandler(func(appData string) error {
		//fmt.Println("收到 Pong，刷新超时")
		return conn.SetReadDeadline(time.Now().Add(30 * time.Second))

	})

	// 监听关闭
	go func() {
		defer func() {
			mu.Lock()
			delete(Clients, conn)
			mu.Unlock()
			conn.Close()
		}()

		for {
			messageType, message, err := conn.ReadMessage()
			if err != nil {
				fmt.Println("读取失败（超时或关闭）:", err)
				break
			}

			config.Logs.Printf("收到消息: %s\n", message)
			err = conn.WriteMessage(messageType, []byte("ok"))
			if err != nil {
				fmt.Println("写入失败:", err)
				break
			}

		}
	}()

}
