package database

import (
	"errors"
	"fmt"
	"watgbridge/state"

	"go.mau.fi/whatsmeow/types"
)

func MsgIdAddNewPair(waMsgId string, participantId types.JID, waChatId types.JID, tgMsgId, tgThreadId int64) error {
	var (
		db = state.State.Database
	)

	var bridgePair MsgIdPair
	res := db.Where(&MsgIdPair{
		WaMessageId: waMsgId,
		WaChatJid:   GetDatabaseJid(waChatId),
	}).Find(&bridgePair)
	if res.Error != nil {
		return res.Error
	}

	if bridgePair.WaMessageId == waMsgId {
		bridgePair.WaSenderJid = GetDatabaseJid(participantId)
		bridgePair.WaChatJid = GetDatabaseJid(waChatId)
		bridgePair.TgMessageId = tgMsgId
		bridgePair.TgThreadId = tgThreadId
		res = db.Save(&bridgePair)
		return res.Error
	}
	// else
	res = db.Create(&MsgIdPair{
		WaMessageId: waMsgId,
		WaSenderJid: GetDatabaseJid(participantId),
		WaChatJid:   GetDatabaseJid(waChatId),
		TgMessageId: tgMsgId,
		TgThreadId:  tgThreadId,
	})
	return res.Error
}

func MsgIdGetTgFromWa(waMsgId, waChatId string) (int64, int64, error) {
	db := state.State.Database

	var bridgePair MsgIdPair
	res := db.Where(&MsgIdPair{
		WaMessageId: waMsgId,
		WaChatJid:   waChatId,
	}).Find(&bridgePair)

	return bridgePair.TgThreadId, bridgePair.TgMessageId, res.Error
}

func MsgIdGetWaFromTg(tgMsgId, tgThreadId int64) (msgId, participantId, chatId string, err error) {
	db := state.State.Database

	var bridgePair MsgIdPair
	res := db.Where(&MsgIdPair{
		TgMessageId: tgMsgId,
		TgThreadId:  tgThreadId,
	}).Find(&bridgePair)

	return bridgePair.WaMessageId, bridgePair.WaSenderJid, bridgePair.WaChatJid, res.Error
}

func MsgIdGetUnreadWa(waChatId types.JID) ([]MsgIdPair, error) {
	db := state.State.Database

	var bridgePairs []MsgIdPair
	res := db.Where(&MsgIdPair{
		WaChatJid: GetDatabaseJid(waChatId),
		WaIsRead:  false,
	}).Find(&bridgePairs)

	return bridgePairs, res.Error
}

func MsgIdMarkReadWa(waChatId types.JID, waMsgId string) error {

	db := state.State.Database

	var bridgePair MsgIdPair
	res := db.Where(&MsgIdPair{
		WaMessageId: waMsgId,
		WaChatJid:   GetDatabaseJid(waChatId),
	}).First(&bridgePair)
	if res.Error != nil {
		return res.Error
	}

	bridgePair.WaIsRead = true
	saveRes := db.Save(&bridgePair)
	return saveRes.Error
}

func MsgIdDeletePair(tgMsgId int64) error {

	db := state.State.Database
	res := db.Where(&MsgIdPair{
		TgMessageId: tgMsgId,
	}).Delete(&MsgIdPair{})

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
		cocoChatThread.ThreadId = tgThreadId
		if pushName != "" {
			cocoContact, _ = FindCocoContactById(cocoChatThread.CocoContactId)
			cocoContact.PushName = pushName
		}
		var contactRes = db.Save(&cocoContact)
		var threadRes = db.Save(&chatThread)
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
		if pushName != "" {
			cocoContact.PushName = pushName
		}

	}

	var contactRes = db.Save(&cocoContact)
	var res = db.Create(&CocoChatThread{
		CocoContactId: cocoContact.ID,
		ThreadId:      tgThreadId,
	})
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

func ChatThreadDropPairByTg(tgThreadId int64) error {
	db := state.State.Database

	res := db.Where(&CocoChatThread{
		ThreadId: tgThreadId,
	}).Delete(&CocoChatThread{})

	return res.Error
}

func ChatThreadGetWaFromTg(tgThreadId int64) (CocoContact, bool) {
	db := state.State.Database

	var chatPair CocoChatThread
	var cocoContact CocoContact
	res := db.Where(&CocoChatThread{
		ThreadId: tgThreadId,
	}).First(&chatPair)
	if res.Error == nil {
		res = db.Where(&CocoContact{ID: chatPair.CocoContactId}).First(&cocoContact)
	}

	return cocoContact, res.Error == nil
}

func ChatThreadGetAllPairs() ([]CocoChatThread, error) {

	db := state.State.Database

	var chatPairs []CocoChatThread
	res := db.Find(&chatPairs)

	return chatPairs, res.Error
}

type CocoContactInfo struct {
	*types.ContactInfo
	Lid types.JID
}

func ContactNameBulkAddOrUpdate(contacts map[types.JID]CocoContactInfo) error {
	var (
		db             = state.State.Database
		logger         = state.State.Logger
		newContacts    []CocoContact
		updateContacts []CocoContact
	)

	for k, manualContactData := range contacts {
		jid := GetDatabaseJid(k)
		lid := GetDatabaseJid(manualContactData.Lid)

		jidContact, found := FindCocoContactSingleId(k)
		if !found {
			lidContact, foundLid := FindCocoContactSingleId(manualContactData.Lid)
			if !foundLid {
				newContacts = append(newContacts, CocoContact{
					Jid:          jid,
					Lid:          lid,
					Name:         manualContactData.FirstName,
					PushName:     manualContactData.PushName,
					BusinessName: manualContactData.BusinessName,
					FullName:     manualContactData.FullName,
				})
			} else {
				updateContacts = append(updateContacts, CocoContact{
					ID:           lidContact.ID,
					Jid:          jid,
					Lid:          lid,
					Name:         manualContactData.FirstName,
					PushName:     manualContactData.PushName,
					BusinessName: manualContactData.BusinessName,
					FullName:     manualContactData.FullName,
				})
			}
		} else {
			updateContacts = append(updateContacts, CocoContact{
				ID:           jidContact.ID,
				Jid:          jid,
				Lid:          lid,
				Name:         manualContactData.FirstName,
				PushName:     manualContactData.PushName,
				BusinessName: manualContactData.BusinessName,
				FullName:     manualContactData.FullName,
			})
		}
	}

	var finalError error
	// Create new ones
	if len(newContacts) > 0 {
		newError := db.Create(&newContacts).Error
		if newError != nil {
			logger.Error(fmt.Sprintf("Error creating new contacts: %v", newError))
		}
		finalError = errors.Join(newError, finalError)
	}

	// Update existing ones
	if len(updateContacts) > 0 {
		for _, c := range updateContacts {
			err := db.Save(&c).Error
			finalError = errors.Join(err, finalError)
			if err != nil {
				logger.Error(fmt.Sprintf("Error updating contacts: %v", err))
			}
		}
	}

	return finalError
}

func ContactNameGet(waUserId types.JID) (string, string, string, string, string, error) {
	contact, found := FindCocoContactSingleId(waUserId)

	if !found {
		return "", "", "", "", "", errors.New("contact not found")
	}

	return GetStringJidAsPhoneNumber(contact.Jid), contact.Name, contact.FullName, contact.PushName, contact.BusinessName, nil
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

func CocoContactUpdateJid(cocoId int32, id types.JID) error {
	db := state.State.Database

	contact, found := FindCocoContactById(cocoId)
	if !found {
		return nil
	}

	contact.Jid = GetDatabaseJid(id)
	var res = db.Save(&contact)

	return res.Error
}

func CocoContactUpdateLid(cocoId int32, id types.JID) error {
	db := state.State.Database

	contact, found := FindCocoContactById(cocoId)
	if !found {
		return nil
	}

	contact.Lid = GetDatabaseJid(id)
	var res = db.Save(&contact)

	return res.Error
}

func UpdateEphemeralSettings(waChatId string, isEphemeral bool, ephemeralTimer uint32) error {
	db := state.State.Database

	var settings ChatEphemeralSettings
	res := db.Where(&ChatEphemeralSettings{ID: waChatId}).Find(&settings)

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
	res := db.Where(&ChatEphemeralSettings{ID: waChatId}).Find(&settings)

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
		PushName: "User (" + id.User + ")",
	}
	result := db.Create(&userContact)

	return userContact, result.Error == nil
}

func CreateCocoContactLid(id types.JID) (CocoContact, bool) {
	db := state.State.Database

	userContact := CocoContact{
		Lid:      GetDatabaseJid(id),
		PushName: "User (" + id.User + ")",
	}
	result := db.Create(&userContact)

	return userContact, result.Error == nil
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

func GetStringJidAsPhoneNumber(j string) string {
	jid, _ := types.ParseJID(j)
	if jid.Server == "lid" {
		return "Unknown Number from WhatsApp"
	}
	return jid.ToNonAD().User
}
