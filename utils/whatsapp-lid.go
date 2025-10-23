package utils

import "go.mau.fi/whatsmeow/types"

func GetJidLid(chatId string, altId string) (string, string) {
	lid := ""
	jid := ""
	chat, _ := types.ParseJID(chatId)
	if chat.Server == "lid" {
		lid = chat.User
	} else {
		jid = chat.User
	}

	alt, _ := types.ParseJID(altId)
	if alt.Server == "lid" {
		if lid != "" {
			panic("both id's were lids...")
		}
		lid = alt.User
	} else {
		jid = alt.User
	}

	return jid, lid
}
