package email

import (
	"context"
	"log"

	"github.com/resend/resend-go/v2"
)

func New(apiKey string) *resend.Client {
	client := resend.NewClient(apiKey)
	return client
}

func Send(ctx context.Context, client *resend.Client, params *resend.SendEmailRequest) (*resend.SendEmailResponse, error) {
	sent, err := client.Emails.SendWithContext(ctx, params)
	if err != nil {
		return nil, err
	}

	log.Printf("Email sent, ID: %s\n", sent.Id)
	return sent, nil
}
