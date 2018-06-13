package main

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"os"

	"github.com/antonholmquist/jason"
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
	StatusCode int    `json:"status_code"`
	Message    string `json:"message"`
}

func sendMail(ctx context.Context, event json.RawMessage) (Response, error) {
	j, err := json.Marshal(&event)
	if err != nil {
		return Response{StatusCode: 500}, errors.New(err.Error())
	}
	v, err := jason.NewObjectFromBytes(j)
	if err != nil {
		return Response{StatusCode: 500}, errors.New(err.Error())
	}

	m := mail.NewV3Mail()
	from := mail.NewEmail(os.Getenv("SENDER_NAME"), os.Getenv("SENDER_ADDRESS"))
	if s, err := v.GetString("subject"); err == nil {
		m.Subject = s
	} else {
		return Response{StatusCode: 500}, errors.New(err.Error())
	}
	m.SetFrom(from)
	m.SetTemplateID(os.Getenv("TEMPLATE_ID"))

	p := mail.NewPersonalization()
	toName, err := v.GetString("to_name")
	if err != nil {
		return Response{StatusCode: 500}, errors.New(err.Error())
	}
	toAddress, err := v.GetString("to_address")
	if err != nil {
		return Response{StatusCode: 500}, errors.New(err.Error())
	}
	to := mail.NewEmail(toName, toAddress)
	p.AddTos(to)

	root, err := v.GetObject()
	if err != nil {
		return Response{StatusCode: 500}, errors.New(err.Error())
	}

	for key, value := range root.Map() {
		strValue, _ := value.String()
		log.Printf("%v = %v", key, strValue)
		p.SetSubstitution("%"+key+"%", strValue)
	}
	m.AddPersonalizations(p)

	// p.SetCustomArg
	request := sendgrid.GetRequest(os.Getenv("SENDGRID_API_KEY"), "/v3/mail/send", "https://api.sendgrid.com")
	request.Method = "POST"
	request.Body = mail.GetRequestBody(m)
	response, err := sendgrid.API(request)
	if err != nil {
		return Response{StatusCode: 500}, errors.New(err.Error())
	}
	return Response{
		StatusCode: response.StatusCode,
		Message:    response.Body,
	}, nil
}

func main() {
	lambda.Start(sendMail)
}
