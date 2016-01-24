package models

import (
	"./../../"
	"database/sql"
	"gopkg.in/qml.v1"
	"log"
	"time"
)

type Control struct {
	Database     *sql.DB
	Root         qml.Object
	Config       Config
	Notification Notification
	Ticker       *time.Ticker
}

type Notification struct {
	Text  string
	Color string
}

func (ctrl *Control) SaveButtonReleased(obj qml.Object) {
	ctrl.Config.NotificationsLimit, _ = obj.Property("notificationsLimit").(int)
	ctrl.Config.Timeout, _ = obj.Property("timeout").(int)
	tx, err := ctrl.Database.Begin()
	if err != nil {
		log.Fatal(err)
	}

	_, err = ctrl.Database.Exec("Insert into Config (id, notifications_limit, timeout) values (?, ?, ?)", nil, ctrl.Config.NotificationsLimit, ctrl.Config.Timeout)
	tx.Commit()
	if err != nil {
		ctrl.Notification.Text = "Unable to save your configuration"
		ctrl.Notification.Color = "red"
	} else {
		ctrl.Notification.Text = "Configuration successfully saved"
		ctrl.Notification.Color = "green"
	}

	obj.Set("notification", ctrl.Notification.Text)
	obj.Set("notificationColor", ctrl.Notification.Color)
	qml.Changed(ctrl, &ctrl.Config)
}

func (ctrl *Control) RunAgainButtonReleased() {
	ctrl.Ticker.Stop()
	ctrl.Ticker = time.NewTicker(time.Duration(ctrl.Config.Timeout) * time.Second)
	go func(config Config, ticker *time.Ticker) {
		provider.Run(ctrl.Config.Timeout, ctrl.Config.NotificationsLimit, ctrl.Ticker)
	}(ctrl.Config, ctrl.Ticker)
}
