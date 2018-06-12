package main

import (
	"errors"
	"os"

	"github.com/asaskevich/govalidator"
	"github.com/aws/aws-lambda-go/lambda"
	sendgrid "github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

type Event struct {
	Subject   string `json:"subject" valid:"required"`
	ToName    string `json:"to_name" valid:"required"`
	ToAddress string `json:"to_address" valid:"required"`
	BodyText  string `json:"body_text" valid:"required"`
}

type Response struct {
	StatusCode int `json:"status_code"`
}

func sendMail(event Event) (Response, error) {
	_, err := govalidator.ValidateStruct(event)
	if err != nil {
		return Response{StatusCode: 500}, errors.New(err.Error())
	}

	from := mail.NewEmail(os.Getenv("SENDER_NAME"), os.Getenv("SENDER_ADDRESS"))
	subject := event.Subject
	to := mail.NewEmail(event.ToName, event.ToAddress)
	plainTextContent := event.BodyText
	htmlContent := event.BodyText
	message := mail.NewSingleEmail(from, subject, to, plainTextContent, htmlContent)
	client := sendgrid.NewSendClient(os.Getenv("SENDGRID_API_KEY"))
	response, err := client.Send(message)

	if err != nil {
		return Response{StatusCode: 500}, errors.New(err.Error())
	}
	return Response{StatusCode: response.StatusCode}, nil
}

func main() {
	lambda.Start(sendMail)
}
