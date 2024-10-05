package controllers

import (
	"fmt"
	"net/http"
	"sana-api/db"
	"sana-api/models"
	"sana-api/utils/token"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func AllChat(c *gin.Context) {
	var chats []models.Chat
	roomid := c.Query("roomid")
	page := c.Query("page")
	pageSize := c.Query("page_size")
	readMsg := c.Query("readmsg")
	if roomid != "" {
		db.CON.Where("chatroom_id = ?", roomid).
			Scopes(db.Paginate(page, pageSize)).
			Order("created_at desc").
			Find(&chats)
	} else {
		db.CON.Find(&chats)
	}
	if readMsg == "true" && roomid != "" {
		userId, _ := token.ExtractTokenID(c)
		roomId, _ := strconv.ParseUint(roomid, 10, 32)
		ReadMsg(userId, uint(roomId))
	}
	c.JSON(http.StatusOK, gin.H{"message": "all Chats", "data": chats})
}

func ReadMsg(userId uint, roomId uint) {
	db.CON.Model(&models.Chat{}).
		Where("is_read = ? and receiver_id = ? and chatroom_id = ?", false, userId, roomId).
		Update("is_read", true)
	msg := Message{
		SenderID:   fmt.Sprintf("user%s", strconv.Itoa(int(userId))),
		ReceiverID: fmt.Sprintf("user%s", strconv.Itoa(int(userId))),
		Content:    "readMsg[chat]",
	}
	BroadcastMessage(msg)
}

func ChatRoom(c *gin.Context) {
	userId, _ := token.ExtractTokenID(c)
	page := c.Query("page")
	pageSize := c.Query("page_size")
	unreadMsg := c.Query("unreadmsg")
	var chatRoom []models.ChatRoomList
	db.CON.Model(&chatRoom).
		Select(`chat_rooms.*, latest_chats.created_at as newchat_at,
			latest_chats.message, latest_chats.sender_id,
			latest_chats.is_read
		`).
		Joins(`
            LEFT JOIN (
                SELECT chatroom_id, message, sender_id, is_read, created_at 
                FROM chats AS c1 
                WHERE created_at = (
                    SELECT MAX(created_at) 
                    FROM chats AS c2 
                    WHERE c1.chatroom_id = c2.chatroom_id
                )
            ) AS latest_chats
            ON latest_chats.chatroom_id = chat_rooms.id
        `).
		Order("newchat_at desc").
		Where("? = ANY(member)", userId).
		Scopes(db.Paginate(page, pageSize)).
		Find(&chatRoom)
	totalunread := 0
	for i := range chatRoom {
		if !chatRoom[i].IsRead && chatRoom[i].SenderId != userId {
			totalunread++
		}
		for _, id := range chatRoom[i].Member {
			var userData models.User
			if id == int64(userId) {
				db.CON.First(&userData, id)
				chatRoom[i].MyData = map[string]any{
					"id":      userData.ID,
					"name":    userData.Name,
					"picture": userData.Picture,
				}
			} else {
				db.CON.First(&userData, id)
				chatRoom[i].ReceiverData = map[string]any{
					"id":      userData.ID,
					"name":    userData.Name,
					"picture": userData.Picture,
				}
			}
		}
		if unreadMsg == "true" {
			var chats models.Chat
			var count int64
			var chatroomId = chatRoom[i].ID
			db.CON.Model(&chats).
				Where("is_read = ? and sender_id != ? and chatroom_id = ?", false, userId, chatroomId).
				Count(&count)
			chatRoom[i].UnreadMsg = count
		}
	}
	c.JSON(http.StatusOK, gin.H{"message": "all Chat Rooms", "data": chatRoom, "unread": totalunread})
}

func StoreChat(c *gin.Context) {
	var chat models.ChatCreate
	userId, err := token.ExtractTokenID(c)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := c.ShouldBindJSON(&chat); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	stringMember := fmt.Sprintf("ARRAY[%d,%d]", int(userId), chat.ReceiverId)

	var chatRoom models.ChatRoom
	var roomid uint
	if err := db.CON.Where("member @> ?", gorm.Expr(stringMember)).
		First(&chatRoom).Error; err != nil {
		newChatRoom := models.ChatRoom{
			Member: []int64{int64(userId), int64(chat.ReceiverId)},
		}
		err = db.CON.Create(&newChatRoom).Error
		if err != nil {
			fmt.Println("Error inserting data:", err)
			return
		}
		roomid = newChatRoom.ID
	} else {
		roomid = chatRoom.ID
	}

	newChat := models.Chat{
		ChatroomId: uint(roomid),
		Message:    chat.Message,
		SenderId:   userId,
		ReceiverId: uint(chat.ReceiverId),
	}
	db.CON.Create(&newChat)
	msg := Message{
		SenderID:   fmt.Sprintf("user%s", strconv.Itoa(int(userId))),
		ReceiverID: fmt.Sprintf("user%s", strconv.Itoa(int(chat.ReceiverId))),
		Content:    fmt.Sprintf("chatRoomId%s", strconv.Itoa(int(roomid))),
	}
	BroadcastMessage(msg)
	//send back to sender
	msg.ReceiverID = fmt.Sprintf("user%s", strconv.Itoa(int(userId)))
	msg.SenderID = fmt.Sprintf("user%s", strconv.Itoa(int(chat.ReceiverId)))
	BroadcastMessage(msg)
	c.JSON(http.StatusOK, gin.H{"message": "chat message stored", "data": newChat})
}
