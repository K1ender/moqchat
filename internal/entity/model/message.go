package model

import "github.com/google/uuid"

type MessageType string

const (
	GroupMessageType MessageType = "group"
	UserMessageType  MessageType = "user"
)

type Message struct {
	ID      uuid.UUID
	Type    MessageType
	From    uuid.UUID
	To      uuid.UUID
	Content string
}
