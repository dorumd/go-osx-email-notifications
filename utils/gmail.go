package utils

import (
	"../models"
	"encoding/json"
	"fmt"
	"golang.org/x/net/context"
	"google.golang.org/api/gmail/v1"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/user"
	"path/filepath"
	"strings"
	"time"

	"golang.org/x/oauth2"
)

// ProcessMessages - process unread messages
func ProcessMessages(
	mResponseMessages []*gmail.Message,
	gmailClient *gmail.Service,
	queue []models.Message,
	processedQueue map[string]bool,
	pushToQueue bool,
	lastProcessedTimestamp int64) ([]models.Message, map[string]bool, bool, int64) {
	var messagesChannel = make(chan *gmail.Message)
	var processedMessagesCount = 0
	var totalMessagesCount = 0
	var noNewMessages = false
	for _, m := range mResponseMessages {
		go func(messageId string) {
			if message, err := gmailClient.Users.Messages.Get(User, messageId).Do(); err == nil {
				if message.InternalDate > lastProcessedTimestamp {
					fmt.Println("Message fetched")
					totalMessagesCount++
					messagesChannel <- message
					noNewMessages = false
				} else {
					noNewMessages = true
				}
			} else {
				fmt.Printf("There was an error trying to get message - %s\n", err.Error())
			}
		}(m.Id)
	}

	for {
		select {
		case message := <-messagesChannel:
			processedMessagesCount++
			fmt.Printf("Inside\n")
			var messageItem = models.Message{}
			messageItem.ID = message.Id
			if ok := processedQueue[messageItem.ID]; ok != true {
				fmt.Printf("Processing headers\n")
				for _, header := range message.Payload.Headers {
					if header.Name == "From" {
						messageItem.From = header.Value
					} else if header.Name == "Subject" {
						messageItem.Subject = header.Value
					}
				}

				fmt.Printf("Processing queuee\n")
				messageItem.Link = strings.Replace(BaseGmailMessageURL, "%MESSAGE_ID%", message.Id, -1)
				if pushToQueue == true {
					queue = append(queue, messageItem)
				} else {
					NotifyMessage(messageItem)
					processedQueue[messageItem.ID] = true
					pushToQueue = true
				}
			}

			if totalMessagesCount == processedMessagesCount {
				fmt.Printf("Last processed message timestamp %#v \n", message.InternalDate)
				return queue, processedQueue, pushToQueue, message.InternalDate
			}
		case <-time.After(50 * time.Millisecond):
			if noNewMessages == true {
				return queue, processedQueue, pushToQueue, lastProcessedTimestamp
			}
			fmt.Printf(".")
		}
	}
}

// GetClient uses a Context and Config to retrieve a Token
// then generate a Client. It returns the generated Client.
func GetClient(ctx context.Context, config *oauth2.Config) *http.Client {
	cacheFile, err := tokenCacheFile()
	if err != nil {
		log.Fatalf("Unable to get path to cached credential file. %v", err)
	}
	tok, err := tokenFromFile(cacheFile)
	if err != nil {
		tok = getTokenFromWeb(config)
		saveToken(cacheFile, tok)
	}
	return config.Client(ctx, tok)
}

// getTokenFromWeb uses Config to request a Token.
// It returns the retrieved Token.
func getTokenFromWeb(config *oauth2.Config) *oauth2.Token {
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Printf("Go to the following link in your browser then type the "+
		"authorization code: \n%v\n", authURL)

	var code string
	if _, err := fmt.Scan(&code); err != nil {
		log.Fatalf("Unable to read authorization code %v", err)
	}

	tok, err := config.Exchange(oauth2.NoContext, code)
	if err != nil {
		log.Fatalf("Unable to retrieve token from web %v", err)
	}
	return tok
}

// tokenCacheFile generates credential file path/filename.
// It returns the generated credential path/filename.
func tokenCacheFile() (string, error) {
	usr, err := user.Current()
	if err != nil {
		return "", err
	}
	tokenCacheDir := filepath.Join(usr.HomeDir, ".credentials")
	os.MkdirAll(tokenCacheDir, 0700)
	return filepath.Join(tokenCacheDir,
		url.QueryEscape("gmail-go-quickstart.json")), err
}

// tokenFromFile retrieves a Token from a given file path.
// It returns the retrieved Token and any read error encountered.
func tokenFromFile(file string) (*oauth2.Token, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	t := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(t)
	defer f.Close()
	return t, err
}

// saveToken uses a file path to create a file and store the
// token in it.
func saveToken(file string, token *oauth2.Token) {
	fmt.Printf("Saving credential file to: %s\n", file)
	f, err := os.Create(file)
	if err != nil {
		log.Fatalf("Unable to cache oauth token: %v", err)
	}
	defer f.Close()
	json.NewEncoder(f).Encode(token)
}
