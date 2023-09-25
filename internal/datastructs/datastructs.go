package datastructs

import (
	"github.com/erupshis/effective_mobile/internal/msgbroker"
)

// PersonData main struct containing person's data.
//
//go:generate easyjson -all datastructs.go
type PersonData struct {
	Id         int64  `json:"id"`
	Name       string `json:"name"`
	Surname    string `json:"surname"`
	Patronymic string `json:"patronymic,omitempty"`
	Age        int64  `json:"age,omitempty"`
	Gender     string `json:"gender,omitempty"`
	Country    string `json:"country,omitempty"`
}

// ErrorMessage kafka's error message for response on incorrect input.
type ErrorMessage struct {
	Error           string `json:"error"`
	OriginalMessage string `json:"original"`
}

// ExtraDataFilling extended struct with raw message(for add it in error message) in case of errors.
type ExtraDataFilling struct {
	Raw  msgbroker.Message
	Data PersonData
}

// Error support struct to extract error from remote API message.
type Error struct {
	Data string `json:"error"`
}

// Age support struct to extract person's age from remote API message.
type Age struct {
	Data int64 `json:"age"`
}

// Gender support struct to extract person's gender from remote API message.
type Gender struct {
	Data string `json:"gender"`
}

// CountryData support struct to extract the most relevant person's country from remote API message.
type CountryData struct {
	Id          string  `json:"country_id"`
	Probability float64 `json:"probability"`
}

// Countries support struct to extract person's country from remote API message.
type Countries struct {
	Data []CountryData `json:"country"`
}
