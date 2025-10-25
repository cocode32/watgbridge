ALTER TABLE msg_id_pairs
    RENAME COLUMN wa_message_id TO id;

ALTER TABLE msg_id_pairs
    RENAME COLUMN wa_sender_jid TO participant_id;

ALTER TABLE msg_id_pairs
    RENAME COLUMN wa_chat_jid TO wa_chat_id;

ALTER TABLE msg_id_pairs
    RENAME COLUMN tg_message_id TO tg_msg_id;
