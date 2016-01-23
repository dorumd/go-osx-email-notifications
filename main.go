package main

import (
	"./models"
	"./utils"
	"fmt"
	"golang.org/x/net/context"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/gmail/v1"
	"io/ioutil"
	"log"
	"strings"
	"time"
)

func main() {
	context := context.Background()
	clientSecretFile, err := ioutil.ReadFile("config/client_secret.json")
	if err != nil {
		log.Fatalf("Unable to read client secret file: %v", err)
	}

	config, err := google.ConfigFromJSON(clientSecretFile, gmail.GmailReadonlyScope)
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
	}
	client := utils.GetClient(context, config)

	gmailClient, err := gmail.New(client)
	if err != nil {
		log.Fatalf("Unable to retrieve gmail Client %v", err)
	}

	var queue = make([]models.Message, 0)
	var processedQueue = make(map[string]bool, 0)
	var pushToQueue = false
	var firstRun = true

	ticker := time.NewTicker(20 * time.Second)
	for _ = range ticker.C {
		fmt.Printf("Processing started at %s\n", time.Now().Format(utils.DateTimeFormat))
		if mResponse, err := gmailClient.Users.Messages.List(utils.User).Q("is:unread AND is:important").Do(); err == nil {
			mResponseMessages := mResponse.Messages[0:utils.NotificationsLimit]

			if firstRun && len(mResponse.Messages) > utils.NotificationsLimit {
				utils.SystemNotification("You have more than 10 unread messages in your inbox")
				firstRun = false
			}

			for _, m := range mResponseMessages {
				if message, err := gmailClient.Users.Messages.Get(utils.User, m.Id).Do(); err == nil {
					var messageItem = models.Message{}
					messageItem.ID = message.Id
					for _, header := range message.Payload.Headers {
						if header.Name == "From" {
							messageItem.From = header.Value
						} else if header.Name == "Subject" {
							messageItem.Subject = header.Value
						}
					}
					if ok := processedQueue[messageItem.ID]; ok != true {
						messageItem.Link = strings.Replace(utils.BaseGmailMessageURL, "%MESSAGE_ID%", message.Id, -1)
						if pushToQueue == true {
							queue = append(queue, messageItem)
						} else {
							utils.NotifyMessage(messageItem)
							processedQueue[messageItem.ID] = true
							pushToQueue = true
						}
					}
				} else {
					fmt.Printf("There was an error trying to get message - %s\n", err.Error())
				}
			}
		} else {
			fmt.Printf("There was an error trying to get messages - %s\n", err.Error())
		}

		if len(queue) > 0 {
			for _, message := range queue {
				utils.NotifyMessage(message)
				processedQueue[message.ID] = true
			}
			queue = make([]models.Message, 0)
			pushToQueue = false
		}
		fmt.Printf("Processing ended at %s\n", time.Now().Format(utils.DateTimeFormat))
	}
}
