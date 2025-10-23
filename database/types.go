package database

import (
	"database/sql"

	"watgbridge/state"
)

type CocoContact struct {
	ID           int32  `gorm:"primaryKey;autoIncrement"`
	Lid          string `gorm:"index:jid_lid,unique"`
	Jid          string `gorm:"index:jid_lid,unique"`
	Name         string
	FullName     string
	PushName     string
	BusinessName string
}

type CocoChatThread struct {
	ID       int32 `gorm:"primaryKey;autoIncrement"`
	ThreadId string
}

type MsgIdPair struct {
	// WhatsApp
	ID            string `gorm:"primaryKey;"` // Message ID
	ParticipantId string // Sender JID
	WaChatId      string // Chat JID

	// Telegram
	TgChatId   int64
	TgThreadId int64
	TgMsgId    int64

	MarkRead sql.NullBool
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
