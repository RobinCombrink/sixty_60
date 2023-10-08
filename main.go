package main

import (
	// "encoding/base64"
	// "bytes"
	"encoding/base64"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"

	// "os"

	// "parser60/authentication"
	"parser60/authentication"
	"parser60/emailparsing"
	"parser60/schema"
	webserver "parser60/web/web_server"

	"google.golang.org/api/gmail/v1"
)

const userId = "me"
const labelName = "Shopping/Sixty60"
const userLabelType = "user"

func main() {
	// messageBodies := readMessageBodiesFromGoogle(true)
	messageBodies := readMessageBodiesFromLocal(filepath.Join("secrets", "messages"))
	invoices := getInvoices(messageBodies)

	webserver.SetupHttpServer(invoices)
}

func getInvoices(messageBodies []string) (invoices []schema.Invoice) {
	for _, messageBody := range messageBodies {
		invoice := emailparsing.GetInvoiceFromHtml(string(messageBody))

		var total float64
		for _, item := range invoice.Items {
			total += float64(item.Total - item.Discount)
		}
		//Convert from cents
		total = (total / 100) + (float64(invoice.DeliveryFee) / 100)
		fmt.Printf("total: R%.2f\n", total)
		invoices = append(invoices, *invoice)
	}
	return invoices
}

func readMessageBodiesFromLocal(messageDirectory string) (messageBodies []string) {
	files, err := os.ReadDir(messageDirectory)
	if err != nil {
		log.Fatalf("Could not read directory: %v", err)
	}

	for _, file := range files {
		fileInfo, err := file.Info()
		if err != nil {
			log.Fatalf("Could not read file info: %v", err)
		}
		fileContests, err := os.ReadFile(filepath.Join(messageDirectory, fileInfo.Name()))
		if err != nil {
			log.Fatalf("Could not read file: %v", err)
		}
		messageBodies = append(messageBodies, string(string(fileContests)))
	}
	return messageBodies
}

func readMessageBodiesFromGoogle(saveLocal bool) (messageBodies []string) {
	service := authentication.Authenticate()
	labelId, err := getLabelWithName(labelName, service)
	if err != nil {
		log.Fatalf("Could not get labels:\n%v", err)
	}
	result, err := service.Users.Messages.List(userId).LabelIds(labelId).Do()
	if err != nil {
		log.Fatalf("Unable to list Messages:\n%v", err)
	}
	messageBodies = getMessageBodies(result, service)
	if saveLocal {
		saveMessagesLocally(messageBodies)
	}
	return messageBodies
}

func saveMessagesLocally(mesageBoddies []string) {
	for i, messageBody := range mesageBoddies {
		fileName := fmt.Sprint(i) + ".html"
		//TODO: Error handling
		os.WriteFile(filepath.Join("secrets", "messages", fileName), []byte(messageBody), 0644)
	}
}

func getLabelWithName(name string, service *gmail.Service) (string, error) {
	r, err := service.Users.Labels.List(userId).Do()
	if err != nil {
		return "", err
	}
	if len(r.Labels) == 0 {
		return "", nil
	}
	for _, label := range r.Labels {
		if label.Type == userLabelType && label.Name == name {
			return label.Id, nil
		}
	}
	return "", errors.New("no results found - try a different label")
}

func getMessageBodies(result *gmail.ListMessagesResponse, service *gmail.Service) (mesageBoddies []string) {
	// Loop over each message
	for _, m := range result.Messages {
		msg, err := service.Users.Messages.Get(userId, m.Id).Do()
		if err != nil {
			log.Printf("Unable to retrieve message %v: %v", m.Id, err)
			continue
		}
		var sumParts string = ""
		// Print the body of the message
		for _, part := range msg.Payload.Parts {
			if part.MimeType == "text/html" {
				data, err := base64.URLEncoding.DecodeString(part.Body.Data)
				if err != nil {
					log.Printf("Unable to decode message body: %v", err)
					continue
				}
				sumParts += string(data)
			}
		}
		if sumParts != "" {
			mesageBoddies = append(mesageBoddies, string(sumParts))
		}
	}
	return mesageBoddies
}
