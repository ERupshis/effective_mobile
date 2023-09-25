package datastructs // import "github.com/erupshis/effective_mobile/internal/datastructs"


TYPES

type Age struct {
	Data int64 `json:"age"`
}

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

func (v PersonData) MarshalEasyJSON(w *jwriter.Writer)
    MarshalEasyJSON supports easyjson.Marshaler interface

func (v PersonData) MarshalJSON() ([]byte, error)
    MarshalJSON supports json.Marshaler interface

func (v *PersonData) UnmarshalEasyJSON(l *jlexer.Lexer)
    UnmarshalEasyJSON supports easyjson.Unmarshaler interface

func (v *PersonData) UnmarshalJSON(data []byte) error
    UnmarshalJSON supports json.Unmarshaler interface

