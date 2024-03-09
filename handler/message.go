package handler

import (
	"demo/dao"
	Model "demo/models"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"sync"
	"time"
)

var (
	clients = make(map[string]Model.ClientNode)
	mux     sync.Mutex
)

var upgrade = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func WsHandler(id string, ctx *gin.Context) {
	defer func() {
		closeClient(id)
		err := recover()
		if err != nil {
			log.Println("handler.WsHandler:", err)
		}
	}()
	//创建节点
	conn, err := upgrade.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		panic(err)
	}
	c := addClient(id, conn)
	//读取未发送消息
	go receiveMsgFromRdb(id, c.DataChan)
	//转发消息
	go receiveMsgFromClient(conn)

	//接收消息并发送至客户端
	pingTicker := time.NewTicker(time.Second * 30)
	for {
		select {
		case msg := <-c.DataChan:
			//log.Println("msg:", msg)
			err := conn.WriteJSON(msg)
			if err != nil {
				panic(err)
			}
		//心跳
		case <-pingTicker.C:
			log.Println("发送ticker")
			err := conn.SetReadDeadline(time.Now().Add(time.Second * 10))
			if err != nil {
				panic(err)
			}
			err = conn.WriteMessage(websocket.PingMessage, []byte("你好"))
			if err != nil {
				panic(err)
			}
		}
	}
}

func addClient(id string, conn *websocket.Conn) *Model.ClientNode {
	mux.Lock()
	defer mux.Unlock()
	log.Println("addClient:", id)
	c := Model.ClientNode{
		Conn:     conn,
		DataChan: make(chan *Model.Message),
	}
	clients[id] = c
	return &c
}
func closeClient(id string) {
	mux.Lock()
	defer mux.Unlock()
	c, ok := clients[id]
	if ok {
		close(c.DataChan)
		conn := c.Conn
		_ = conn.Close()
		delete(clients, id)
	}
}
func getClient(id string) *Model.ClientNode {
	mux.Lock()
	defer mux.Unlock()
	c, ok := clients[id]
	if !ok {
		return nil
	}
	return &c
}

// 接收消息进程
func receiveMsgFromClient(conn *websocket.Conn) {
	defer func() {
		err := recover()
		if err != nil {
			_ = conn.Close()
			log.Println("handler.receiveMsgFromClient:", err)
		}
	}()
	//获取连接成功后第一次消息
	//_, _, err := conn.ReadMessage()
	//if err != nil {
	//	panic(err)
	//}
	for {
		msg := &Model.Message{}
		err := conn.ReadJSON(&msg)
		if err != nil {
			panic(err)
		}
		switch msg.MsgType {
		//转发私聊消息
		case Model.MPrivate:
			err := SendMsg(msg)
			if err != nil {
				panic(err)
			}
		//转发群聊消息
		case Model.MGroup:
			group, err := GetGroupById(msg.ReceiverId)
			if err != nil {
				panic(err)
			}
			for _, member := range group.Members {
				if member.ID == msg.SenderId {
					continue
				}
				msg.ReceiverId = member.ID
				err := SendMsg(msg)
				if err != nil {
					panic(err)
				}
			}
		default:

		}
	}
}
func receiveMsgFromRdb(id string, m chan *Model.Message) {
	defer func() {
		if r := recover(); r != nil {
			log.Println("handler.handle_message.GetMsgFromRdb:", r)
		}
	}()
	Messages, err := dao.Rdb.LRange(dao.BgCtx, dao.MsgPre+id, 0, -1).Result()
	if err != nil {
		panic(err)
	}
	for _, message := range Messages {
		msg := &Model.Message{}
		//log.Println("formRdb:", message)
		err := json.Unmarshal([]byte(message), msg)
		if err != nil {
			panic(err)
		}
		m <- msg
		err = dao.Rdb.LPop(dao.BgCtx, dao.MsgPre+id).Err()
		if err != nil {
			panic(err)
		}
	}
}

// sendMsgToRdb 将信息加入缓存
func sendMsgToRdb(m *Model.Message) error {
	bytes, err := json.Marshal(m)
	if err != nil {

		return err
	}
	log.Println("sendMsgToRdb:", bytes)
	err = dao.Rdb.RPush(dao.BgCtx, dao.MsgPre+m.ReceiverId, bytes).Err()
	return err
}

// SendMsg 若客户端未连接则发送至缓存
func SendMsg(m *Model.Message) error {
	defer func() {
		if r := recover(); r != nil {
			log.Println("handler.handle_message.SendMsg:", r)
		}
	}()
	client := getClient(m.ReceiverId)
	if client == nil {
		err := sendMsgToRdb(m)
		return err
	}
	client.DataChan <- m
	return nil
}
