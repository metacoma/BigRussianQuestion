#!/bin/bash


send_message() {
  message="$1"
  curl -X POST \
     https://api.telegram.org/bot$TELEGRAM_BOT_TOKEN/sendMessage  \
     -F text="${message}"                                         \
     -F parse_mode="HTML"                                         \
     -F chat_id="${TELEGRAM_CHAT_ID}"
}


message=`echo "$*" | sed 's/["\(\)]//g'`
curl --data-urlencode "q=$message" telegrambot:9090/question
