package database

import (
	"database/sql"
	"errors"
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

func AddNewChatThread(waChatId types.JID, tgThreadId int64) error {
	db := state.State.Database

	cocoChatThread, found := GetChatThread(waChatId)
	if found {
		cocoChatThread.ThreadId = tgThreadId
		var res = db.Save(&cocoChatThread)
		return res.Error
	}

	cocoContact, found := FindCocoContactSingleId(waChatId)
	if !found {
		// unknown contact, just create one to make the thread
		_, lid := GetJidOrLid(waChatId)
		if lid == "" {
			cocoContact, _ = CreateCocoContactJid(waChatId)
		} else {
			cocoContact, _ = CreateCocoContactLid(waChatId)
		}
	}

	var res = db.Create(&CocoChatThread{
		CocoContactId: cocoContact.ID,
		ThreadId:      tgThreadId,
	})
	return res.Error
}

func AddNewChatThreadWithPush(waChatId types.JID, tgThreadId int64, pushName string) error {
	db := state.State.Database

	var chatThread CocoChatThread
	var cocoContact CocoContact
	cocoChatThread, found := GetChatThread(waChatId)
	if found {
		chatThread.ThreadId = tgThreadId
		cocoContact, _ = FindCocoContactById(cocoChatThread.CocoContactId)
		cocoContact.PushName = pushName
		var threadRes = db.Save(&chatThread)
		var contactRes = db.Create(&cocoContact)
		return errors.Join(threadRes.Error, contactRes.Error)
	}

	cocoContact, found = FindCocoContactSingleId(waChatId)
	if !found {
		// unknown contact, just create one to make the thread
		_, lid := GetJidOrLid(waChatId)
		if lid == "" {
			cocoContact, _ = CreateCocoContactJid(waChatId)
		} else {
			cocoContact, _ = CreateCocoContactLid(waChatId)
		}
		cocoContact.PushName = pushName
	}

	var res = db.Create(&CocoChatThread{
		CocoContactId: cocoContact.ID,
		ThreadId:      tgThreadId,
	})
	var contactRes = db.Save(&cocoContact)
	return errors.Join(res.Error, contactRes.Error)
}

func GetChatThread(waChatId types.JID) (CocoChatThread, bool) {
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

func ChatThreadGetWaFromTg(tgThreadId int64) (CocoContact, bool) {
	db := state.State.Database

	var chatPair CocoChatThread
	var cocoContact CocoContact
	res := db.Where(&CocoChatThread{
		ThreadId: tgThreadId,
	}).Find(&chatPair)
	if res.Error != nil {
		res = db.Where(&CocoContact{ID: chatPair.CocoContactId}).Find(&cocoContact)
	}

	return cocoContact, res.Error == nil
}

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

// TODO not used?
//func ContactNameAddNew(waUserId, firstName, fullName, pushName, businessName string) error {
//	db := state.State.Database
//
//	contact, found, _, _ := FindCocoContactSingleId(waUserId)
//	if !found {
//		panic("to be honest, if we didn't find or create the contact before, now it's all fucking broken")
//	}
//
//	contact.Name = firstName
//	contact.FullName = fullName
//	contact.PushName = pushName
//	contact.BusinessName = businessName
//	var res = db.Save(&contact)
//	return res.Error
//}

func ContactNameBulkAddOrUpdate(contacts map[types.JID]types.ContactInfo) error {

	var (
		db           = state.State.Database
		contactNames []CocoContact
	)

	for k, manualContactData := range contacts {
		jid, lid := GetJidOrLid(k)
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

func ContactNameGet(waUserId types.JID) (string, string, string, string, error) {
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

func CocoContactUpdatePushName(senderId types.JID, senderAltId types.JID, pushName string) error {
	if pushName == "" {
		return nil
	}

	db := state.State.Database

	contact, found := FindCocoContact(senderId, senderAltId)
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

func FindCocoContactById(id int32) (CocoContact, bool) {
	db := state.State.Database

	var userContact CocoContact
	var result = db.Where(&CocoContact{
		ID: id,
	}).First(&userContact)

	return userContact, result.Error == nil
}

func FindCocoContact(jid types.JID, lid types.JID) (CocoContact, bool) {
	db := state.State.Database

	var userContact CocoContact
	var result = db.Where(&CocoContact{
		Jid: GetDatabaseJid(jid),
		Lid: GetDatabaseJid(lid),
	}).First(&userContact)

	return userContact, result.Error == nil
}

func CreateCocoContact(jid types.JID, lid types.JID, name string) (CocoContact, bool) {
	db := state.State.Database

	var pushName = name
	if pushName == "" {
		pushName = jid.User
	}

	userContact := CocoContact{
		Jid:      GetDatabaseJid(jid),
		Lid:      GetDatabaseJid(lid),
		PushName: name,
	}
	result := db.Create(&userContact)

	return userContact, result.Error == nil
}

func FindCocoContactSingleId(idFromWhatsmeow types.JID) (CocoContact, bool) {
	db := state.State.Database

	var userContact CocoContact
	var result = db.Where(&CocoContact{
		Jid: GetDatabaseJid(idFromWhatsmeow),
	}).Or(&CocoContact{
		Lid: GetDatabaseJid(idFromWhatsmeow),
	}).First(&userContact)

	return userContact, result.Error == nil
}

func CreateCocoContactJid(id types.JID) (CocoContact, bool) {
	db := state.State.Database

	userContact := CocoContact{
		Jid:      GetDatabaseJid(id),
		PushName: id.User,
	}
	result := db.Create(&userContact)

	return userContact, result.Error == nil
}

func CreateCocoContactLid(id types.JID) (CocoContact, bool) {
	db := state.State.Database

	userContact := CocoContact{
		Lid:      GetDatabaseJid(id),
		PushName: id.User,
	}
	result := db.Create(&userContact)

	return userContact, result.Error == nil
}

// GetJidLid IDK how to put this in utils, so it's gonna live here
func GetJidLid(chatId types.JID, altId types.JID) (string, string) {
	lid := ""
	jid := ""
	if chatId.Server == "lid" {
		lid = GetDatabaseJid(chatId)
	} else {
		jid = GetDatabaseJid(chatId)
	}

	if altId.Server == "lid" {
		if lid != "" {
			panic("both id's were lids...")
		}
		lid = GetDatabaseJid(altId)
	} else {
		jid = GetDatabaseJid(altId)
	}

	return jid, lid
}

func GetJidOrLid(id types.JID) (string, string) {
	lid := ""
	jid := ""
	if id.Server == "lid" {
		lid = GetDatabaseJid(id)
	} else {
		jid = GetDatabaseJid(id)
	}

	return jid, lid
}

func GetDatabaseJid(j types.JID) string {
	return j.ToNonAD().String()
}
