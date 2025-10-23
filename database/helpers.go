package database

import (
	"database/sql"
	"errors"
	"fmt"

	"watgbridge/state"

	"go.mau.fi/whatsmeow/types"
)

func MsgIdAddNewPair(waMsgId, participantId, waChatId string, tgChatId, tgMsgId, tgThreadId int64) error {

	db := state.State.Database

	var bridgePair MsgIdPair
	res := db.Where("id = ? AND wa_chat_id = ?", waMsgId, waChatId).Find(&bridgePair)
	if res.Error != nil {
		return res.Error
	}

	if bridgePair.ID == waMsgId {
		bridgePair.ParticipantId = participantId
		bridgePair.WaChatId = waChatId
		bridgePair.TgChatId = tgChatId
		bridgePair.TgMsgId = tgMsgId
		bridgePair.TgThreadId = tgThreadId
		bridgePair.MarkRead = sql.NullBool{Valid: true, Bool: false}
		res = db.Save(&bridgePair)
		return res.Error
	}
	// else
	res = db.Create(&MsgIdPair{
		ID:            waMsgId,
		ParticipantId: participantId,
		WaChatId:      waChatId,
		TgChatId:      tgChatId,
		TgMsgId:       tgMsgId,
		TgThreadId:    tgThreadId,
		MarkRead:      sql.NullBool{Valid: true, Bool: false},
	})
	return res.Error
}

func MsgIdGetTgFromWa(waMsgId, waChatId string) (int64, int64, int64, error) {

	db := state.State.Database

	var bridgePair MsgIdPair
	res := db.Where("id = ? AND wa_chat_id = ?", waMsgId, waChatId).Find(&bridgePair)

	return bridgePair.TgChatId, bridgePair.TgThreadId, bridgePair.TgMsgId, res.Error
}

func MsgIdGetWaFromTg(tgChatId, tgMsgId, tgThreadId int64) (msgId, participantId, chatId string, err error) {

	db := state.State.Database

	var bridgePair MsgIdPair
	res := db.Where("tg_chat_id = ? AND tg_msg_id = ? AND tg_thread_id = ?", tgChatId, tgMsgId, tgThreadId).Find(&bridgePair)

	return bridgePair.ID, bridgePair.ParticipantId, bridgePair.WaChatId, res.Error
}

func MsgIdGetUnread(waChatId string) (map[string]([]string), error) {

	db := state.State.Database

	var bridgePairs []MsgIdPair
	res := db.Where("wa_chat_id = ? AND mark_read = false", waChatId).Find(&bridgePairs)

	var msgIds = make(map[string]([]string))

	for _, pair := range bridgePairs {
		if _, found := msgIds[pair.ParticipantId]; !found {
			msgIds[pair.ParticipantId] = []string{}
		}
		msgIds[pair.ParticipantId] = append(msgIds[pair.ParticipantId], pair.ID)
	}

	return msgIds, res.Error
}

func MsgIdMarkRead(waChatId, waMsgId string) error {

	db := state.State.Database

	var bridgePair MsgIdPair
	res := db.Where("id = ? AND wa_chat_id = ?", waMsgId, waChatId).Find(&bridgePair)
	if res.Error != nil {
		return res.Error
	}

	if bridgePair.ID == waMsgId {
		bridgePair.MarkRead = sql.NullBool{Valid: true, Bool: true}
		res = db.Save(&bridgePair)
		return res.Error
	}

	return nil
}

func MsgIdDeletePair(tgChatId, tgMsgId int64) error {

	db := state.State.Database
	res := db.Where("tg_chat_id = ? AND tg_msg_id = ?", tgChatId, tgMsgId).Delete(&MsgIdPair{})

	return res.Error
}

func MsgIdDropAllPairs() error {

	db := state.State.Database
	res := db.Where("1 = 1").Delete(&MsgIdPair{})

	return res.Error
}

func AddNewChatThread(waChatId string, tgThreadId int64) error {

	db := state.State.Database

	var chatThread CocoChatThread
	cocoContact, found := GetChatThread(waChatId)
	if found {
		chatThread.ThreadId = fmt.Sprintf("%d", tgThreadId)
		var res = db.Save(&chatThread)
		return res.Error
	}

	var res = db.Create(&CocoChatThread{
		CocoContactId: cocoContact.ID,
		ThreadId:      fmt.Sprintf("%d", tgThreadId),
	})
	return res.Error
}

func GetChatThread(waChatId string) (CocoChatThread, bool) {
	db := state.State.Database

	var chatThread CocoChatThread
	cocoContact, found := FindCocoContactSingleId(waChatId)

	if !found {
		return chatThread, false
	}

	var result = db.Where(&CocoChatThread{
		CocoContactId: cocoContact.ID,
	}).First(&chatThread)

	return chatThread, result.Error == nil
}

// TODO maybe not needed - confirm
//func ChatThreadDropPairByTg(tgChatId, tgThreadId int64) error {
//	db := state.State.Database
//
//	res := db.Where("tg_chat_id = ? AND tg_thread_id = ?", tgChatId, tgThreadId).Delete(&ChatThreadPair{})
//
//	return res.Error
//}
//
//func ChatThreadGetWaFromTg(tgChatId, tgThreadId int64) (string, error) {
//	db := state.State.Database
//
//	var chatPair ChatThreadPair
//	res := db.Where("tg_chat_id = ? AND tg_thread_id = ?", tgChatId, tgThreadId).Find(&chatPair)
//
//	return chatPair.ID, res.Error
//}

func ChatThreadGetAllPairs() ([]CocoChatThread, error) {

	db := state.State.Database

	var chatPairs []CocoChatThread
	res := db.Find(&chatPairs)

	return chatPairs, res.Error
}

// TODO not referenced it seems
//func ChatThreadDropAllPairs() error {
//
//	db := state.State.Database
//	res := db.Where("1 = 1").Delete(&ChatThreadPair{})
//
//	return res.Error
//}

func ContactNameAddNew(waUserId, firstName, fullName, pushName, businessName string) error {
	db := state.State.Database

	contact, found := FindCocoContactSingleId(waUserId)
	if !found {
		panic("to be honest, if we didn't find or create the contact before, now it's all fucking broken")
	}

	contact.Name = firstName
	contact.FullName = fullName
	contact.PushName = pushName
	contact.BusinessName = businessName
	var res = db.Save(&contact)
	return res.Error
}

func ContactNameBulkAddOrUpdate(contacts map[types.JID]types.ContactInfo) error {

	var (
		db           = state.State.Database
		contactNames []CocoContact
	)

	for k, manualContactData := range contacts {
		jid, lid := GetJidOrLid(k.User)
		if lid == "" {
			contactNames = append(contactNames, CocoContact{
				Jid:          jid,
				Name:         manualContactData.FirstName,
				PushName:     manualContactData.PushName,
				BusinessName: manualContactData.BusinessName,
				FullName:     manualContactData.FullName,
			})
			continue
		}

		contactNames = append(contactNames, CocoContact{
			Lid:          lid,
			Name:         manualContactData.FirstName,
			PushName:     manualContactData.PushName,
			BusinessName: manualContactData.BusinessName,
			FullName:     manualContactData.FullName,
		})
	}

	res := db.Save(&contactNames)

	return res.Error
}

func ContactNameGet(waUserId string) (string, string, string, string, error) {
	contact, found := FindCocoContactSingleId(waUserId)

	if !found {
		return "", "", "", "", errors.New("contact not found")
	}

	return contact.Name, contact.FullName, contact.PushName, contact.BusinessName, nil
}

func ContactGetAll() (map[int]CocoContact, error) {
	db := state.State.Database

	var contacts []CocoContact
	res := db.Find(&contacts)

	results := make(map[int]CocoContact)
	for _, contact := range contacts {
		results[int(contact.ID)] = contact
	}
	return results, res.Error
}

func CocoContactUpdatePushName(chatId string, senderAltId string, pushName string) error {
	if pushName == "" {
		return nil
	}

	db := state.State.Database

	jid, lid := GetJidLid(chatId, senderAltId)

	contact, found := FindCocoContact(jid, lid)
	if !found {
		return nil
	}

	contact.PushName = pushName
	var res = db.Save(&contact)

	return res.Error
}

func UpdateEphemeralSettings(waChatId string, isEphemeral bool, ephemeralTimer uint32) error {
	db := state.State.Database

	var settings ChatEphemeralSettings
	res := db.Where("id = ?", waChatId).Find(&settings)

	if res.Error != nil {
		return res.Error
	}

	if settings.ID != waChatId {
		res = db.Create(&ChatEphemeralSettings{
			ID:             waChatId,
			IsEphemeral:    isEphemeral,
			EphemeralTimer: ephemeralTimer,
		})
		return res.Error
	}

	settings.IsEphemeral = isEphemeral
	settings.EphemeralTimer = ephemeralTimer

	res = db.Save(&settings)

	return res.Error
}

func GetEphemeralSettings(waChatId string) (bool, uint32, bool, error) {
	db := state.State.Database

	var settings ChatEphemeralSettings
	res := db.Where("id = ?", waChatId).Find(&settings)

	if res.Error != nil {
		return false, 0, false, res.Error
	}

	if settings.ID != waChatId {
		return false, 0, false, nil
	}

	return settings.IsEphemeral, settings.EphemeralTimer, true, nil
}

func FindCocoContact(jid string, lid string) (CocoContact, bool) {
	db := state.State.Database

	var userContact CocoContact
	var result = db.Where(&CocoContact{
		Jid: jid,
		Lid: lid,
	}).First(&userContact)

	if result.Error == nil {
		return userContact, true
	}
	db.Create(&CocoContact{
		Jid: jid,
		Lid: lid,
	})

	result = db.Where(&CocoContact{
		Jid: jid,
		Lid: lid,
	}).First(&userContact)

	if result.Error != nil {
		return userContact, false
	}
	return userContact, true
}

func FindCocoContactSingleId(idFromWhatsmeow string) (CocoContact, bool) {
	db := state.State.Database

	var userContact CocoContact
	var result = db.Where(&CocoContact{
		Jid: idFromWhatsmeow,
	}).Or(&CocoContact{
		Lid: idFromWhatsmeow,
	}).First(&userContact)

	if result.Error == nil {
		return userContact, true
	}

	// create the damn thing
	jid, lid := GetJidOrLid(idFromWhatsmeow)
	if lid == "" {
		db.Create(&CocoContact{
			Jid: jid,
		})

		result = db.Where(&CocoContact{
			Jid: jid,
		}).First(&userContact)

		return userContact, result.Error == nil
	}

	db.Create(&CocoContact{
		Lid: lid,
	})

	result = db.Where(&CocoContact{
		Lid: lid,
	}).First(&userContact)

	return userContact, result.Error == nil
}

// GetJidLid IDK how to put this in utils, so it's gonna live here
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

func GetJidOrLid(id string) (string, string) {
	lid := ""
	jid := ""
	chat, _ := types.ParseJID(id)
	if chat.Server == "lid" {
		lid = chat.User
	} else {
		jid = chat.User
	}

	return jid, lid
}
