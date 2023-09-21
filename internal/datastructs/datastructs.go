package datastructs

import (
	"github.com/erupshis/effective_mobile/internal/msgbroker"
)

//go:generate easyjson -all datastructs.go
type PersonData struct {
	Name        string `json:"name"`
	Surname     string `json:"surname"`
	Patronymic  string `json:"patronymic,omitempty"`
	Age         int64  `json:"age,omitempty"`
	Sex         string `json:"gender,omitempty"`
	Nationality string `json:"country,omitempty"`
}

type ErrorMessage struct {
	Error           string `json:"error"`
	OriginalMessage string `json:"original"`
}

type ExtraDataFilling struct {
	Raw  msgbroker.Message
	Data PersonData
}

//Parse data from Remote API.

type Error struct {
	Data string `json:"error"`
}

type Age struct {
	Data int64 `json:"age"`
}

type Gender struct {
	Data string `json:"gender"`
}

type CountryData struct {
	Id string `json:"country_id"`
}

type Nationality struct {
	Data []CountryData `json:"country"`
}
