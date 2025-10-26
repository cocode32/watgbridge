package telegram

import "time"

type BotConversationState struct {
	Command    string
	Step       string
	Data       map[string]string
	Started    time.Time
	Selectable []string // list of options for inline buttons
}
