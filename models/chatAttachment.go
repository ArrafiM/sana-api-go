package models

import (
	"time"
)

type ChatAttachment struct {
	ID        uint   `gorm:"primarykey"`
	ChatID    uint   `json:"chat_id" binding:"required"`
	Url       string `json:"url" binding:"required"`
	CreatedAt time.Time
}

type ChatAttachmentCreate struct {
	ChatId uint   `form:"chat_id" json:"chat_id" binding:"required"`
	Url    string `form:"files" json:"url" binding:"required"`
}
