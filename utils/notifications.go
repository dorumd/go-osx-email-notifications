package utils

import (
	"../models"
	"github.com/deckarep/gosx-notifier"
)

// Notify - triggers an notification
func Notify(message models.Message) {
	note := gosxnotifier.NewNotification(message.From)
	note.Title = message.Subject
	note.Sound = gosxnotifier.Default
	note.Group = "com.go-email.notification"
	note.Sender = "com.apple.Mail"
	note.Link = message.Link
	note.Push()
}
