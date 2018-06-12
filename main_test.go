package main

import (
	"testing"
)

func TestSendMail(t *testing.T) {
	r, err := sendMail(Event{
		Subject:   "テストメール",
		ToName:    "テスト太郎",
		ToAddress: "foo@bar.com",
		BodyText:  "本文です",
	})
	if err != nil {
		t.Error(err)
	}

	t.Log(r)
}
