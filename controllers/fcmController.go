package controllers

import (
	"context"
	"log"

	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/messaging"
	"google.golang.org/api/option"
)

func InitFirebaseApp() (*messaging.Client, error) {
	// Use the context from the `context` package.
	ctx := context.Background()

	// Initialize Firebase app.
	opt := option.WithCredentialsFile("tmp/sana-mobile-783d3-firebase-adminsdk-6m1hn-cc9051e745.json")
	app, err := firebase.NewApp(ctx, nil, opt)
	if err != nil {
		return nil, err
	}

	// Get the messaging client.
	client, err := app.Messaging(ctx)
	if err != nil {
		return nil, err
	}

	return client, nil
}

func SendNotification(token string, title string, body string) error {
	client, err := InitFirebaseApp()
	if err != nil {
		log.Printf("Error sent message: %s\n", err)
		return err
	}
	// Create a message to send to the device token.
	message := &messaging.Message{
		Token: token,
		Notification: &messaging.Notification{
			Title: title,
			Body:  body,
		},
	}

	// Send the message.
	ctx := context.Background()
	response, err := client.Send(ctx, message)
	if err != nil {
		return err
	}

	log.Printf("Successfully sent message: %s\n", response)
	return nil
}
