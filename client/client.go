package main

import (
	"./../"
	"./models"
	"database/sql"
	"fmt"
	"github.com/mattn/go-sqlite3"
	"gopkg.in/qml.v1"
	"log"
	"os"
	"time"
)

func main() {
	if err := qml.Run(run); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

func run() error {
	engine := qml.NewEngine()

	client, err := engine.LoadFile("client.qml")
	if err != nil {
		return err
	}

	ctrl := initDatabase()
	context := engine.Context()
	context.SetVar("controller", &ctrl)

	window := client.CreateWindow(nil)
	ctrl.Root = window.Root()
	ctrl.Ticker = time.NewTicker(time.Duration(ctrl.Config.Timeout) * time.Second)
	go func(ctrl models.Control) {
		provider.Init()
		provider.Run(ctrl.Config.Timeout, ctrl.Config.NotificationsLimit, ctrl.Ticker)
	}(ctrl)

	window.Show()
	window.Wait()
	return nil
}

func initDatabase() models.Control {
	var dbDriver string
	sql.Register(dbDriver, &sqlite3.SQLiteDriver{})

	database, err := sql.Open(dbDriver, "db")

	if err != nil {
		fmt.Println("Failed to create the handle")
	}

	if err := database.Ping(); err != nil {
		fmt.Println("Failed to keep connection alive")
	}

	if err != nil {
		log.Fatal(err)
	}

	_, err = database.Exec(
		"CREATE TABLE IF NOT EXISTS Config (id integer PRIMARY KEY, notifications_limit integer NOT NULL, timeout integer NOT NULL)",
	)

	if err != nil {
		log.Fatal(err)
	}

	config := models.Config{NotificationsLimit: 3, Timeout: 10}
	_ = database.QueryRow(
		"SELECT notifications_limit, timeout FROM Config ORDER BY id DESC LIMIT 1",
	).Scan(&config.NotificationsLimit, &config.Timeout)

	controller := models.Control{Config: config, Database: database}
	return controller
}
