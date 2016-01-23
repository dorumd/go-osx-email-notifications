package utils

import (
	"../models"
	"github.com/deckarep/gosx-notifier"
	"time"
)

// NotifyMessage - triggers an notification
func NotifyMessage(message models.Message) {
	note := gosxnotifier.NewNotification(message.From)
	note.Title = message.Subject
	note.Sound = gosxnotifier.Default
	note.Group = "com.go-email.notification"
	note.Sender = "com.apple.Mail"
	note.Link = message.Link
	note.Push()
	time.Sleep(2 * time.Second)
}

// SystemNotification used for system notifications
func SystemNotification(text string) {
	note := gosxnotifier.NewNotification(text)
	note.Push()
	time.Sleep(2 * time.Second)
}
