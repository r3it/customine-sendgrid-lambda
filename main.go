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

type Response struct {
	StatusCode int    `json:"status_code"`
	Message    string `json:"message"`
}

func sendMail(ctx context.Context, event json.RawMessage) (Response, error) {
	j, err := json.Marshal(&event)
	if err != nil {
		return errResponse(err)
	}
	v, err := jason.NewObjectFromBytes(j)
	if err != nil {
		return errResponse(err)
	}

	m := mail.NewV3Mail()
	from := mail.NewEmail(os.Getenv("SENDER_NAME"), os.Getenv("SENDER_ADDRESS"))
	if s, err := v.GetString("subject"); err == nil {
		m.Subject = s
	} else {
		return errResponse(err)
	}
	m.SetFrom(from)

	templateID, err := v.GetString("template_id")
	if err != nil {
		return errResponse(err)
	}
	m.SetTemplateID(templateID)

	p := mail.NewPersonalization()
	toName, err := v.GetString("to_name")
	if err != nil {
		return errResponse(err)
	}
	toAddress, err := v.GetString("to_address")
	if err != nil {
		return errResponse(err)
	}
	to := mail.NewEmail(toName, toAddress)
	p.AddTos(to)
	bcc := mail.NewEmail(os.Getenv("BCC_NAME"), os.Getenv("BCC_ADDRESS"))
	if bcc.Name != "" && bcc.Address != "" {
		p.AddBCCs(bcc)
	}

	root, err := v.GetObject()
	if err != nil {
		return errResponse(err)
	}

	for key, value := range root.Map() {
		strValue, _ := value.String()
		log.Printf("%v = %v", key, strValue)
		p.SetSubstitution("%"+key+"%", strValue)
	}
	m.AddPersonalizations(p)

	request := sendgrid.GetRequest(os.Getenv("SENDGRID_API_KEY"), "/v3/mail/send", "https://api.sendgrid.com")
	request.Method = "POST"
	request.Body = mail.GetRequestBody(m)
	response, err := sendgrid.API(request)
	if err != nil {
		return errResponse(err)
	}
	return Response{
		StatusCode: response.StatusCode,
		Message:    response.Body,
	}, nil
}

func errResponse(err error) (Response, error) {
	return Response{StatusCode: 500}, errors.New(err.Error())
}

func main() {
	lambda.Start(sendMail)
}
