package models

import (
	"github.com/gorilla/websocket"
	"time"
)

const (
	MPrivate = iota
	MGroup
	MFriendReq
	MFriendAgree
	MGroupReq
	MGroupAgree
)
const (
	CString = iota
	CPic
)

type Message struct {
	CreatedAt   *time.Time
	SenderId    string
	ReceiverId  string
	Content     interface{}
	ContentType int
	MsgType     int
}

func (m *Message) TableName() string {
	return "Massages"
}

type ClientNode struct {
	Conn     *websocket.Conn
	DataChan chan *Message
	Online   bool
}
