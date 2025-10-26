package whatsapp

import (
	"context"
	"watgbridge/database"

	"go.mau.fi/whatsmeow"
	waTypes "go.mau.fi/whatsmeow/types"
)

func SyncContactsWithBridge(waClient *whatsmeow.Client) error {
	contacts, err := waClient.Store.Contacts.GetAllContacts(context.Background())

	wrappedContacts := make(map[waTypes.JID]database.CocoContactInfo, len(contacts))
	for jid, info := range contacts {
		lid, _ := waClient.Store.LIDs.GetLIDForPN(context.Background(), jid.ToNonAD())
		wrappedContacts[jid] = database.CocoContactInfo{
			ContactInfo: &info,
			Lid:         lid,
		}
	}

	if err == nil {
		err = database.ContactNameBulkAddOrUpdate(wrappedContacts)
	}

	return err
}
