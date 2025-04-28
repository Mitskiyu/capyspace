package email

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sesv2"
	"github.com/aws/aws-sdk-go-v2/service/sesv2/types"
)

type Email struct {
	To       []string
	From     string
	Subject  string
	HTMLBody string
	RawBody  string
}

func New() (*sesv2.Client, error) {
	region := os.Getenv("AWS_REGION")

	if region == "" {
		region = "eu-central-1"
	}

	cfg, err := config.LoadDefaultConfig(context.Background(),
		config.WithRegion(region),
	)
	if err != nil {
		return nil, fmt.Errorf("unable to load AWS SDK config: %v", err)
	}

	svc := sesv2.NewFromConfig(cfg)
	return svc, nil
}

func Send(ctx context.Context, svc *sesv2.Client, email Email) error {
	input := &sesv2.SendEmailInput{
		Destination: &types.Destination{
			ToAddresses: email.To,
		},
		Content: &types.EmailContent{
			Simple: &types.Message{
				Body: &types.Body{
					Html: &types.Content{
						Charset: aws.String("UTF-8"),
						Data:    aws.String(email.HTMLBody),
					},
					Text: &types.Content{
						Charset: aws.String("UTF-8"),
						Data:    aws.String(email.RawBody),
					},
				},
				Subject: &types.Content{
					Charset: aws.String("UTF-8"),
					Data:    aws.String(email.Subject),
				},
			},
		},
		FromEmailAddress: aws.String(email.From),
	}

	output, err := svc.SendEmail(ctx, input)
	if err != nil {
		return fmt.Errorf("error sending email: %v", err)
	}

	log.Printf("Email sent: %v", aws.ToString(output.MessageId))
	return nil
}
