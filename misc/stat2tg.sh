#!/bin/sh

. ~/answers/.env

update_tg_message() {
  message_id=4967
  new_message="$*"

  curl -s -X POST \
     https://api.telegram.org/bot$TELEGRAM_BOT_TOKEN/editMessageText  \
     -F message_id="${message_id}"                                    \
     -F text="${new_message}"                                    \
     -F parse_mode="HTML"                                       \
     -F chat_id="${PREMODERATION_CHAT_ID}"
} 

db_query() {
  cd ~/answers
  sudo docker-compose exec telegrambot sqlite3 /data/db.sqlite "$*"
} 

total=`db_query 'SELECT count(message_id) from answers;'`
in_queue=`db_query 'select count(message_id) from answers WHERE is_gold = 1 AND (sent = 0 OR sent is NULL);'`

update_tg_message "Вопросов в базе: $total, в очереди: $in_queue"
