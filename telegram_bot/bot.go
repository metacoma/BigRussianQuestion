package main

import (
  "database/sql"
  "strings"
  "github.com/go-telegram-bot-api/telegram-bot-api"
  "net/http"
  "os/exec"
  "os"
  "log"
  "fmt"
  "strconv"
  "time"
  _ "github.com/mattn/go-sqlite3"

)

var TELEGRAM_TOKEN = os.Getenv("TELEGRAM_TOKEN")
var GOLD_CHAT_ID = getenvInt("GOLD_CHAT_ID")
var FLOW_CHAT_ID = getenvInt("FLOW_CHAT_ID")
var PREMODERATION_CHAT_ID = getenvInt("PREMODERATION_CHAT_ID")
var MIRROR_TXT_CHANNEL = getenvInt("MIRROR_TXT_CHANNEL")
var bot *tgbotapi.BotAPI
var sqliteDatabase *sql.DB
var SQLITE_DB_PATH = os.Getenv("SQLITE_DB_PATH")



var messages_map = make(map[string]string)

func makeTimestamp() int64 {
    return time.Now().UnixNano() / int64(time.Millisecond)
}

func getenvInt(key string) (int64) {
    var v int64;
    v, _ = strconv.ParseInt(os.Getenv(key), 10, 64)
    return v
}

func questionHandler(w http.ResponseWriter, r *http.Request) {

  r.ParseForm()

  parsed := make(map[string]string);

  for k, v := range r.Form {
    parsed[k] = strings.Join(v, "")
    fmt.Println("key:", k)
    fmt.Println("val:", strings.Join(v, ""))
  }

  sendButton(FLOW_CHAT_ID, parsed["q"], "")
  return
}

func gold(w http.ResponseWriter, r *http.Request) {

  r.ParseForm()

  parsed := make(map[string]string);

  for k, v := range r.Form {
    parsed[k] = strings.Join(v, "")
    fmt.Println("key:", k)
    fmt.Println("val:", strings.Join(v, ""))
  }

  telegram_msg := tgbotapi.NewMessage(PREMODERATION_CHAT_ID, parsed["q"])
  telegram_msg.ParseMode = tgbotapi.ModeHTML;
  bot.Send(telegram_msg)
  return
}

func httpServer() {
  http.HandleFunc("/question", questionHandler) // set router
  http.HandleFunc("/gold", gold) // 
  http.ListenAndServe(":9090", nil) // set listen port
}

func sendButton(dst_chat int64, text string, message_id string) {
  var unixtime int64
  var row []tgbotapi.InlineKeyboardButton
  keyboard := tgbotapi.InlineKeyboardMarkup{}
  if (message_id == "") {
    unixtime = makeTimestamp()
    message_id = fmt.Sprintf("%d", unixtime)
  }

  if (dst_chat == PREMODERATION_CHAT_ID || dst_chat == GOLD_CHAT_ID) {

    if (dst_chat != GOLD_CHAT_ID) {
      row = append(row, tgbotapi.NewInlineKeyboardButtonData("Отправить в золотой фонд", message_id))
      keyboard.InlineKeyboard = append(keyboard.InlineKeyboard, row)
      msg := tgbotapi.NewMessage(dst_chat, text)
      msg.ReplyMarkup = keyboard
      bot.Send(msg)
    } else {
      tmp_file := fmt.Sprintf("/tmp/image_%d.png", time.Now().Unix())
      GenerateImage(text, tmp_file)
      msg := tgbotapi.NewPhotoUpload(dst_chat, tmp_file)
      bot.Send(msg)
    }

  } else {
    StoreAnswer(sqliteDatabase, unixtime, text)
    row = append(row, tgbotapi.NewInlineKeyboardButtonData("Отправить в золотой фонд", message_id))
    keyboard.InlineKeyboard = append(keyboard.InlineKeyboard, row)
    msg := tgbotapi.NewMessage(dst_chat, text)
    msg.ReplyMarkup = keyboard
    bot.Send(msg)
  }


}

func GenerateImage(text string, dst_file string) string {
        //cmd := fmt.Sprintf("cd /image_generator && cat /images/russianflow.jpg | ./image_generator.sh '%s' > %s", text, dst_file)
        cmd := fmt.Sprintf("/image_generator/generate.sh '%s' > %s", text, dst_file)
        log.Printf("GenerateImage: %s\n", cmd)
        exec.Command("sh","-c",cmd).Output()
        return dst_file
}

func createTable(db *sql.DB) {
  createTableSQL := `CREATE TABLE answers (
    "message_id" integer NOT NULL PRIMARY KEY,   
    "answer" TEXT    
    );` // SQL Statement for Create Table

  log.Println("Create answers table...")
  statement, err := db.Prepare(createTableSQL) // Prepare SQL Statement
  if err != nil {
    log.Fatal(err.Error())
  }
  statement.Exec() // Execute SQL Statements
  log.Println("Answers table created")
}

func getAnswer(db *sql.DB, message_id int64) string {
  q := fmt.Sprintf("SELECT answer FROM answers WHERE message_id = %d", message_id)
  row, err := db.Query(q)
  log.Printf("getAnswer: query = %s\n", q)
  if err != nil {
    log.Fatal(err)
  }
  defer row.Close()
  for row.Next() { // Iterate and fetch the records from result cursor
    var answer string
    row.Scan(&answer)
    return answer
  }

  return ""
}

func getAnswerByTxtID(db *sql.DB, message_txt_id string) string {
  message_id, _ := strconv.ParseInt(message_txt_id, 10, 64)
  return getAnswer(db, message_id)
}

func StoreAnswer(db *sql.DB, message_id int64, answer string) {
  insertSQL := `INSERT INTO answers(message_id, answer) VALUES (?, ?)`
  statement, err := db.Prepare(insertSQL)
  if err != nil {
    log.Fatalln(err.Error())
  }
  _, err = statement.Exec(message_id, answer)
  if err != nil {
    log.Fatalln(err.Error())
  }
}

func InitDb() *sql.DB {
  needCreateTable := false
  if _, err := os.Stat(SQLITE_DB_PATH); os.IsNotExist(err) {
    file, err := os.Create(SQLITE_DB_PATH) // Create SQLite file
    if err != nil {
      log.Fatal(err.Error())
    }
    file.Close()
    log.Println("db created")
    needCreateTable = true
  }

  sqliteDatabase, _ = sql.Open("sqlite3", SQLITE_DB_PATH)

  if (needCreateTable) {
    createTable(sqliteDatabase)
  }
  return sqliteDatabase
}


func main() {
  var err error
  bot, err = tgbotapi.NewBotAPI(TELEGRAM_TOKEN)

  InitDb()

  if (err != nil) {
    log.Panic(err)
  }
  go httpServer()

  u := tgbotapi.NewUpdate(0)
  u.Timeout = 60

  updates, err := bot.GetUpdatesChan(u)



  for update := range updates {
    if update.CallbackQuery != nil {

      log.Printf("------\n%+v\n-------\n", update.CallbackQuery)

      callback_data := string(update.CallbackQuery.Data)
      if (update.CallbackQuery.ChatInstance == "786515482557635255") {
        sendButton(PREMODERATION_CHAT_ID, getAnswerByTxtID(sqliteDatabase, callback_data), callback_data)
        deleteMessageConfig := tgbotapi.NewDeleteMessage(FLOW_CHAT_ID, update.CallbackQuery.Message.MessageID)
        bot.DeleteMessage(deleteMessageConfig)
      }

      if (update.CallbackQuery.ChatInstance == "8475962581961792526") {

        sendButton(GOLD_CHAT_ID, getAnswerByTxtID(sqliteDatabase, callback_data), callback_data)
        deleteMessageConfig := tgbotapi.NewDeleteMessage(PREMODERATION_CHAT_ID, update.CallbackQuery.Message.MessageID)
        bot.DeleteMessage(deleteMessageConfig)

        msg := tgbotapi.NewMessage(MIRROR_TXT_CHANNEL, getAnswerByTxtID(sqliteDatabase, callback_data))
        bot.Send(msg)
      }

      // always delete the message that was clicked
    }

    if update.Message == nil {
      continue
    }

  }

  log.Printf("Token is %s", TELEGRAM_TOKEN)
}
