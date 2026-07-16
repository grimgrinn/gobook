package storage

import (
	"fmt"
	"log"
	"net/smtp"
)

func bytesInUse(username string) int64 { return 0 /*...*/ }

// Настройка отправителя электронных писем.
// Примечение: никогда не помещайте пароль в исходный текст!
const sender = "notification@example.com"
const password = "correcthorsebatterystaple"
const hostname = "smtp.example.com"
const template = `Внимание, вы использовали %d байт хранилища, %d%% вашей квоты.`

var notifyUser = func(username, msg string) {
	auth := smtp.PlainAuth("", sender, password, hostname)
	err := smtp.SendMail(hostname+":587", auth, sender, []string{username}, []byte(msg))
	if err != nil {
		log.Printf("Сбой smpt.SendEmail(%s): %s", username, err)
	}
}

func CheckQuota(username string) {
	used := bytesInUse(username)
	const quota = 1000000000 // 1GB
	percent := 100 * used / quota
	if percent < 90 {
		return // OK
	}
	msg := fmt.Sprintf(template, used, percent)
	notifyUser(username, msg)
}
