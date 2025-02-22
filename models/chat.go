package models

import (
	"time"

	"github.com/lib/pq"
	"gorm.io/gorm"
)

type ChatRoom struct {
	gorm.Model
	Member pq.Int64Array `json:"member" gorm:"type:integer[]"`
}

type ChatRoomList struct {
	ChatRoom
	NewchatAt    time.Time      `json:"newchat_at"`
	SenderId     uint           `json:"sender_id"`
	Message      string         `json:"message"`
	IsRead       bool           `json:"is_read"`
	MyData       map[string]any `json:"mydata" gorm:"-"`
	ReceiverData map[string]any `json:"receiverdata" gorm:"-"`
	UnreadMsg    int64          `json:"unreadmsg" gorm:"-"`
}

type Chat struct {
	gorm.Model
	ChatroomId   uint            `json:"chatroom_id"`
	SenderId     uint            `json:"sender_id"`
	ReceiverId   uint            `json:"receiver_id"`
	Message      string          `json:"message"`
	IsRead       bool            `json:"is_read"`
	Chatroom     *ChatRoom       `json:"chatroom" gorm:"-"`
	UserSender   *User           `json:"user_sender" gorm:"foreignKey:sender_id"`
	UserReceiver *User           `json:"user_receiver" gorm:"foreignKey:receiver_id"`
	File         *ChatAttachment `json:"file" gorm:"foreignKey:ChatID"`
}

type ChatCreate struct {
	ReceiverId int    `form:"receiver_id" json:"receiver_id" binding:"required"`
	Message    string `form:"message" json:"message" binding:"required"`
	File       string `form:"files" json:"url"`
}

func (ChatRoomList) TableName() string {
	return "chat_rooms"
}
