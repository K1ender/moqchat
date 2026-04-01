package entity

import "github.com/google/uuid"

type MessageType string

const (
	GroupMessageType MessageType = "group"
	UserMessageType  MessageType = "user"
)

type Message struct {
	id      uuid.UUID
	typ     MessageType
	from    uuid.UUID
	to      uuid.UUID
	content string
}
