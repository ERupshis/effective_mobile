package datastructs

import (
	"github.com/erupshis/effective_mobile/internal/msgbroker"
)

//go:generate easyjson -all datastructs.go
type PersonData struct {
	Name        string `json:"name"`
	Surname     string `json:"surname"`
	Patronymic  string `json:"patronymic,omitempty"`
	Age         string `json:"age,omitempty"`
	Sex         string `json:"sex,omitempty"`
	Nationality string `json:"nationality,omitempty"`
}

type ErrorMessage struct {
	Error           string `json:"error"`
	OriginalMessage string `json:"original"`
}

type ExtraDataFilling struct {
	Raw  msgbroker.Message
	Data PersonData
}
