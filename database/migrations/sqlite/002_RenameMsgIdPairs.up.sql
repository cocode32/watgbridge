ALTER TABLE msg_id_pairs
    RENAME COLUMN id TO wa_message_id;

ALTER TABLE msg_id_pairs
    RENAME COLUMN participant_id TO wa_sender_jid;

ALTER TABLE msg_id_pairs
    RENAME COLUMN wa_chat_id TO wa_chat_jid;

ALTER TABLE msg_id_pairs
    RENAME COLUMN tg_msg_id TO tg_message_id;