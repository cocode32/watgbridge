package database

import (
	"watgbridge/state"
)

type CocoContact struct {
	ID           int32 `gorm:"primaryKey;autoIncrement"`
	Lid          string
	Jid          string
	Name         string
	FullName     string
	PushName     string
	BusinessName string
}

type CocoChatThread struct {
	ID            int32 `gorm:"primaryKey;autoIncrement"`
	CocoContactId int32
	ThreadId      int64
}

type MsgIdPair struct {
	// WhatsApp
	WaMessageId      string `gorm:"primaryKey;"` // Message ID
	WaParticipantJid string // Sender JID
	WaChatJid        string // Chat JID
	WaIsRead         bool   // keep track of messages we sent to whatsapp to mark as read

	// Telegram
	TgThreadId  int64
	TgMessageId int64
}

type ChatEphemeralSettings struct {
	ID             string `gorm:"primaryKey;"` // WhatsApp Chat ID
	IsEphemeral    bool
	EphemeralTimer uint32
}

func AutoMigrate() error {
	db := state.State.Database
	return db.AutoMigrate(
		&MsgIdPair{},
		&ChatEphemeralSettings{},
		&CocoContact{},
		&CocoChatThread{},
	)
}
