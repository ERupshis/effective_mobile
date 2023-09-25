package datastructs // import "github.com/erupshis/effective_mobile/internal/datastructs"


TYPES

type Age struct {
	Data int64 `json:"age"`
}
    Age support struct to extract person's age from remote API message.

func (v Age) MarshalEasyJSON(w *jwriter.Writer)
    MarshalEasyJSON supports easyjson.Marshaler interface

func (v Age) MarshalJSON() ([]byte, error)
    MarshalJSON supports json.Marshaler interface

func (v *Age) UnmarshalEasyJSON(l *jlexer.Lexer)
    UnmarshalEasyJSON supports easyjson.Unmarshaler interface

func (v *Age) UnmarshalJSON(data []byte) error
    UnmarshalJSON supports json.Unmarshaler interface

type Countries struct {
	Data []CountryData `json:"country"`
}
    Countries support struct to extract person's country from remote API
    message.

func (v Countries) MarshalEasyJSON(w *jwriter.Writer)
    MarshalEasyJSON supports easyjson.Marshaler interface

func (v Countries) MarshalJSON() ([]byte, error)
    MarshalJSON supports json.Marshaler interface

func (v *Countries) UnmarshalEasyJSON(l *jlexer.Lexer)
    UnmarshalEasyJSON supports easyjson.Unmarshaler interface

func (v *Countries) UnmarshalJSON(data []byte) error
    UnmarshalJSON supports json.Unmarshaler interface

type CountryData struct {
	Id          string  `json:"country_id"`
	Probability float64 `json:"probability"`
}
    CountryData support struct to extract the most relevant person's country
    from remote API message.

func (v CountryData) MarshalEasyJSON(w *jwriter.Writer)
    MarshalEasyJSON supports easyjson.Marshaler interface

func (v CountryData) MarshalJSON() ([]byte, error)
    MarshalJSON supports json.Marshaler interface

func (v *CountryData) UnmarshalEasyJSON(l *jlexer.Lexer)
    UnmarshalEasyJSON supports easyjson.Unmarshaler interface

func (v *CountryData) UnmarshalJSON(data []byte) error
    UnmarshalJSON supports json.Unmarshaler interface

type Error struct {
	Data string `json:"error"`
}
    Error support struct to extract error from remote API message.

func (v Error) MarshalEasyJSON(w *jwriter.Writer)
    MarshalEasyJSON supports easyjson.Marshaler interface

func (v Error) MarshalJSON() ([]byte, error)
    MarshalJSON supports json.Marshaler interface

func (v *Error) UnmarshalEasyJSON(l *jlexer.Lexer)
    UnmarshalEasyJSON supports easyjson.Unmarshaler interface

func (v *Error) UnmarshalJSON(data []byte) error
    UnmarshalJSON supports json.Unmarshaler interface

type ErrorMessage struct {
	Error           string `json:"error"`
	OriginalMessage string `json:"original"`
}
    ErrorMessage kafka's error message for response on incorrect input.

func (v ErrorMessage) MarshalEasyJSON(w *jwriter.Writer)
    MarshalEasyJSON supports easyjson.Marshaler interface

func (v ErrorMessage) MarshalJSON() ([]byte, error)
    MarshalJSON supports json.Marshaler interface

func (v *ErrorMessage) UnmarshalEasyJSON(l *jlexer.Lexer)
    UnmarshalEasyJSON supports easyjson.Unmarshaler interface

func (v *ErrorMessage) UnmarshalJSON(data []byte) error
    UnmarshalJSON supports json.Unmarshaler interface

type ExtraDataFilling struct {
	Raw  msgbroker.Message
	Data PersonData
}
    ExtraDataFilling extended struct with raw message(for add it in error
    message) in case of errors.

func (v ExtraDataFilling) MarshalEasyJSON(w *jwriter.Writer)
    MarshalEasyJSON supports easyjson.Marshaler interface

func (v ExtraDataFilling) MarshalJSON() ([]byte, error)
    MarshalJSON supports json.Marshaler interface

func (v *ExtraDataFilling) UnmarshalEasyJSON(l *jlexer.Lexer)
    UnmarshalEasyJSON supports easyjson.Unmarshaler interface

func (v *ExtraDataFilling) UnmarshalJSON(data []byte) error
    UnmarshalJSON supports json.Unmarshaler interface

type Gender struct {
	Data string `json:"gender"`
}
    Gender support struct to extract person's gender from remote API message.

func (v Gender) MarshalEasyJSON(w *jwriter.Writer)
    MarshalEasyJSON supports easyjson.Marshaler interface

func (v Gender) MarshalJSON() ([]byte, error)
    MarshalJSON supports json.Marshaler interface

func (v *Gender) UnmarshalEasyJSON(l *jlexer.Lexer)
    UnmarshalEasyJSON supports easyjson.Unmarshaler interface

func (v *Gender) UnmarshalJSON(data []byte) error
    UnmarshalJSON supports json.Unmarshaler interface

type PersonData struct {
	Id         int64  `json:"id"`
	Name       string `json:"name"`
	Surname    string `json:"surname"`
	Patronymic string `json:"patronymic,omitempty"`
	Age        int64  `json:"age,omitempty"`
	Gender     string `json:"gender,omitempty"`
	Country    string `json:"country,omitempty"`
}
    PersonData main struct containing person's data.

func (v PersonData) MarshalEasyJSON(w *jwriter.Writer)
    MarshalEasyJSON supports easyjson.Marshaler interface

func (v PersonData) MarshalJSON() ([]byte, error)
    MarshalJSON supports json.Marshaler interface

func (v *PersonData) UnmarshalEasyJSON(l *jlexer.Lexer)
    UnmarshalEasyJSON supports easyjson.Unmarshaler interface

func (v *PersonData) UnmarshalJSON(data []byte) error
    UnmarshalJSON supports json.Unmarshaler interface

