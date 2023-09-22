// Code generated by easyjson for marshaling/unmarshaling. DO NOT EDIT.

package datastructs

import (
	json "encoding/json"
	msgbroker "github.com/erupshis/effective_mobile/internal/msgbroker"
	easyjson "github.com/mailru/easyjson"
	jlexer "github.com/mailru/easyjson/jlexer"
	jwriter "github.com/mailru/easyjson/jwriter"
)

// suppress unused package warning
var (
	_ *json.RawMessage
	_ *jlexer.Lexer
	_ *jwriter.Writer
	_ easyjson.Marshaler
)

func easyjson1cf0b278DecodeGithubComErupshisEffectiveMobileInternalDatastructs(in *jlexer.Lexer, out *PersonData) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		if isTopLevel {
			in.Consumed()
		}
		in.Skip()
		return
	}
	in.Delim('{')
	for !in.IsDelim('}') {
		key := in.UnsafeFieldName(false)
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		case "id":
			out.Id = string(in.String())
		case "name":
			out.Name = string(in.String())
		case "surname":
			out.Surname = string(in.String())
		case "patronymic":
			out.Patronymic = string(in.String())
		case "age":
			out.Age = int64(in.Int64())
		case "gender":
			out.Gender = string(in.String())
		case "country":
			out.Country = string(in.String())
		default:
			in.SkipRecursive()
		}
		in.WantComma()
	}
	in.Delim('}')
	if isTopLevel {
		in.Consumed()
	}
}
func easyjson1cf0b278EncodeGithubComErupshisEffectiveMobileInternalDatastructs(out *jwriter.Writer, in PersonData) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"id\":"
		out.RawString(prefix[1:])
		out.String(string(in.Id))
	}
	{
		const prefix string = ",\"name\":"
		out.RawString(prefix)
		out.String(string(in.Name))
	}
	{
		const prefix string = ",\"surname\":"
		out.RawString(prefix)
		out.String(string(in.Surname))
	}
	if in.Patronymic != "" {
		const prefix string = ",\"patronymic\":"
		out.RawString(prefix)
		out.String(string(in.Patronymic))
	}
	if in.Age != 0 {
		const prefix string = ",\"age\":"
		out.RawString(prefix)
		out.Int64(int64(in.Age))
	}
	if in.Gender != "" {
		const prefix string = ",\"gender\":"
		out.RawString(prefix)
		out.String(string(in.Gender))
	}
	if in.Country != "" {
		const prefix string = ",\"country\":"
		out.RawString(prefix)
		out.String(string(in.Country))
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v PersonData) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson1cf0b278EncodeGithubComErupshisEffectiveMobileInternalDatastructs(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v PersonData) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson1cf0b278EncodeGithubComErupshisEffectiveMobileInternalDatastructs(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *PersonData) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson1cf0b278DecodeGithubComErupshisEffectiveMobileInternalDatastructs(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *PersonData) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson1cf0b278DecodeGithubComErupshisEffectiveMobileInternalDatastructs(l, v)
}
func easyjson1cf0b278DecodeGithubComErupshisEffectiveMobileInternalDatastructs1(in *jlexer.Lexer, out *Gender) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		if isTopLevel {
			in.Consumed()
		}
		in.Skip()
		return
	}
	in.Delim('{')
	for !in.IsDelim('}') {
		key := in.UnsafeFieldName(false)
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		case "gender":
			out.Data = string(in.String())
		default:
			in.SkipRecursive()
		}
		in.WantComma()
	}
	in.Delim('}')
	if isTopLevel {
		in.Consumed()
	}
}
func easyjson1cf0b278EncodeGithubComErupshisEffectiveMobileInternalDatastructs1(out *jwriter.Writer, in Gender) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"gender\":"
		out.RawString(prefix[1:])
		out.String(string(in.Data))
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v Gender) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson1cf0b278EncodeGithubComErupshisEffectiveMobileInternalDatastructs1(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v Gender) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson1cf0b278EncodeGithubComErupshisEffectiveMobileInternalDatastructs1(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *Gender) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson1cf0b278DecodeGithubComErupshisEffectiveMobileInternalDatastructs1(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *Gender) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson1cf0b278DecodeGithubComErupshisEffectiveMobileInternalDatastructs1(l, v)
}
func easyjson1cf0b278DecodeGithubComErupshisEffectiveMobileInternalDatastructs2(in *jlexer.Lexer, out *ExtraDataFilling) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		if isTopLevel {
			in.Consumed()
		}
		in.Skip()
		return
	}
	in.Delim('{')
	for !in.IsDelim('}') {
		key := in.UnsafeFieldName(false)
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		case "Raw":
			easyjson1cf0b278DecodeGithubComErupshisEffectiveMobileInternalMsgbroker(in, &out.Raw)
		case "Data":
			(out.Data).UnmarshalEasyJSON(in)
		default:
			in.SkipRecursive()
		}
		in.WantComma()
	}
	in.Delim('}')
	if isTopLevel {
		in.Consumed()
	}
}
func easyjson1cf0b278EncodeGithubComErupshisEffectiveMobileInternalDatastructs2(out *jwriter.Writer, in ExtraDataFilling) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"Raw\":"
		out.RawString(prefix[1:])
		easyjson1cf0b278EncodeGithubComErupshisEffectiveMobileInternalMsgbroker(out, in.Raw)
	}
	{
		const prefix string = ",\"Data\":"
		out.RawString(prefix)
		(in.Data).MarshalEasyJSON(out)
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v ExtraDataFilling) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson1cf0b278EncodeGithubComErupshisEffectiveMobileInternalDatastructs2(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v ExtraDataFilling) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson1cf0b278EncodeGithubComErupshisEffectiveMobileInternalDatastructs2(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *ExtraDataFilling) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson1cf0b278DecodeGithubComErupshisEffectiveMobileInternalDatastructs2(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *ExtraDataFilling) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson1cf0b278DecodeGithubComErupshisEffectiveMobileInternalDatastructs2(l, v)
}
func easyjson1cf0b278DecodeGithubComErupshisEffectiveMobileInternalMsgbroker(in *jlexer.Lexer, out *msgbroker.Message) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		if isTopLevel {
			in.Consumed()
		}
		in.Skip()
		return
	}
	in.Delim('{')
	for !in.IsDelim('}') {
		key := in.UnsafeFieldName(false)
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		case "Key":
			if in.IsNull() {
				in.Skip()
				out.Key = nil
			} else {
				out.Key = in.Bytes()
			}
		case "Value":
			if in.IsNull() {
				in.Skip()
				out.Value = nil
			} else {
				out.Value = in.Bytes()
			}
		default:
			in.SkipRecursive()
		}
		in.WantComma()
	}
	in.Delim('}')
	if isTopLevel {
		in.Consumed()
	}
}
func easyjson1cf0b278EncodeGithubComErupshisEffectiveMobileInternalMsgbroker(out *jwriter.Writer, in msgbroker.Message) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"Key\":"
		out.RawString(prefix[1:])
		out.Base64Bytes(in.Key)
	}
	{
		const prefix string = ",\"Value\":"
		out.RawString(prefix)
		out.Base64Bytes(in.Value)
	}
	out.RawByte('}')
}
func easyjson1cf0b278DecodeGithubComErupshisEffectiveMobileInternalDatastructs3(in *jlexer.Lexer, out *ErrorMessage) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		if isTopLevel {
			in.Consumed()
		}
		in.Skip()
		return
	}
	in.Delim('{')
	for !in.IsDelim('}') {
		key := in.UnsafeFieldName(false)
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		case "error":
			out.Error = string(in.String())
		case "original":
			out.OriginalMessage = string(in.String())
		default:
			in.SkipRecursive()
		}
		in.WantComma()
	}
	in.Delim('}')
	if isTopLevel {
		in.Consumed()
	}
}
func easyjson1cf0b278EncodeGithubComErupshisEffectiveMobileInternalDatastructs3(out *jwriter.Writer, in ErrorMessage) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"error\":"
		out.RawString(prefix[1:])
		out.String(string(in.Error))
	}
	{
		const prefix string = ",\"original\":"
		out.RawString(prefix)
		out.String(string(in.OriginalMessage))
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v ErrorMessage) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson1cf0b278EncodeGithubComErupshisEffectiveMobileInternalDatastructs3(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v ErrorMessage) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson1cf0b278EncodeGithubComErupshisEffectiveMobileInternalDatastructs3(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *ErrorMessage) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson1cf0b278DecodeGithubComErupshisEffectiveMobileInternalDatastructs3(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *ErrorMessage) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson1cf0b278DecodeGithubComErupshisEffectiveMobileInternalDatastructs3(l, v)
}
func easyjson1cf0b278DecodeGithubComErupshisEffectiveMobileInternalDatastructs4(in *jlexer.Lexer, out *Error) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		if isTopLevel {
			in.Consumed()
		}
		in.Skip()
		return
	}
	in.Delim('{')
	for !in.IsDelim('}') {
		key := in.UnsafeFieldName(false)
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		case "error":
			out.Data = string(in.String())
		default:
			in.SkipRecursive()
		}
		in.WantComma()
	}
	in.Delim('}')
	if isTopLevel {
		in.Consumed()
	}
}
func easyjson1cf0b278EncodeGithubComErupshisEffectiveMobileInternalDatastructs4(out *jwriter.Writer, in Error) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"error\":"
		out.RawString(prefix[1:])
		out.String(string(in.Data))
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v Error) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson1cf0b278EncodeGithubComErupshisEffectiveMobileInternalDatastructs4(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v Error) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson1cf0b278EncodeGithubComErupshisEffectiveMobileInternalDatastructs4(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *Error) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson1cf0b278DecodeGithubComErupshisEffectiveMobileInternalDatastructs4(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *Error) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson1cf0b278DecodeGithubComErupshisEffectiveMobileInternalDatastructs4(l, v)
}
func easyjson1cf0b278DecodeGithubComErupshisEffectiveMobileInternalDatastructs5(in *jlexer.Lexer, out *CountryData) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		if isTopLevel {
			in.Consumed()
		}
		in.Skip()
		return
	}
	in.Delim('{')
	for !in.IsDelim('}') {
		key := in.UnsafeFieldName(false)
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		case "country_id":
			out.Id = string(in.String())
		case "probability":
			out.Probability = float64(in.Float64())
		default:
			in.SkipRecursive()
		}
		in.WantComma()
	}
	in.Delim('}')
	if isTopLevel {
		in.Consumed()
	}
}
func easyjson1cf0b278EncodeGithubComErupshisEffectiveMobileInternalDatastructs5(out *jwriter.Writer, in CountryData) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"country_id\":"
		out.RawString(prefix[1:])
		out.String(string(in.Id))
	}
	{
		const prefix string = ",\"probability\":"
		out.RawString(prefix)
		out.Float64(float64(in.Probability))
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v CountryData) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson1cf0b278EncodeGithubComErupshisEffectiveMobileInternalDatastructs5(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v CountryData) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson1cf0b278EncodeGithubComErupshisEffectiveMobileInternalDatastructs5(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *CountryData) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson1cf0b278DecodeGithubComErupshisEffectiveMobileInternalDatastructs5(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *CountryData) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson1cf0b278DecodeGithubComErupshisEffectiveMobileInternalDatastructs5(l, v)
}
func easyjson1cf0b278DecodeGithubComErupshisEffectiveMobileInternalDatastructs6(in *jlexer.Lexer, out *Countries) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		if isTopLevel {
			in.Consumed()
		}
		in.Skip()
		return
	}
	in.Delim('{')
	for !in.IsDelim('}') {
		key := in.UnsafeFieldName(false)
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		case "country":
			if in.IsNull() {
				in.Skip()
				out.Data = nil
			} else {
				in.Delim('[')
				if out.Data == nil {
					if !in.IsDelim(']') {
						out.Data = make([]CountryData, 0, 2)
					} else {
						out.Data = []CountryData{}
					}
				} else {
					out.Data = (out.Data)[:0]
				}
				for !in.IsDelim(']') {
					var v7 CountryData
					(v7).UnmarshalEasyJSON(in)
					out.Data = append(out.Data, v7)
					in.WantComma()
				}
				in.Delim(']')
			}
		default:
			in.SkipRecursive()
		}
		in.WantComma()
	}
	in.Delim('}')
	if isTopLevel {
		in.Consumed()
	}
}
func easyjson1cf0b278EncodeGithubComErupshisEffectiveMobileInternalDatastructs6(out *jwriter.Writer, in Countries) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"country\":"
		out.RawString(prefix[1:])
		if in.Data == nil && (out.Flags&jwriter.NilSliceAsEmpty) == 0 {
			out.RawString("null")
		} else {
			out.RawByte('[')
			for v8, v9 := range in.Data {
				if v8 > 0 {
					out.RawByte(',')
				}
				(v9).MarshalEasyJSON(out)
			}
			out.RawByte(']')
		}
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v Countries) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson1cf0b278EncodeGithubComErupshisEffectiveMobileInternalDatastructs6(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v Countries) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson1cf0b278EncodeGithubComErupshisEffectiveMobileInternalDatastructs6(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *Countries) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson1cf0b278DecodeGithubComErupshisEffectiveMobileInternalDatastructs6(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *Countries) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson1cf0b278DecodeGithubComErupshisEffectiveMobileInternalDatastructs6(l, v)
}
func easyjson1cf0b278DecodeGithubComErupshisEffectiveMobileInternalDatastructs7(in *jlexer.Lexer, out *Age) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		if isTopLevel {
			in.Consumed()
		}
		in.Skip()
		return
	}
	in.Delim('{')
	for !in.IsDelim('}') {
		key := in.UnsafeFieldName(false)
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		case "age":
			out.Data = int64(in.Int64())
		default:
			in.SkipRecursive()
		}
		in.WantComma()
	}
	in.Delim('}')
	if isTopLevel {
		in.Consumed()
	}
}
func easyjson1cf0b278EncodeGithubComErupshisEffectiveMobileInternalDatastructs7(out *jwriter.Writer, in Age) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"age\":"
		out.RawString(prefix[1:])
		out.Int64(int64(in.Data))
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v Age) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson1cf0b278EncodeGithubComErupshisEffectiveMobileInternalDatastructs7(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v Age) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson1cf0b278EncodeGithubComErupshisEffectiveMobileInternalDatastructs7(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *Age) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson1cf0b278DecodeGithubComErupshisEffectiveMobileInternalDatastructs7(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *Age) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson1cf0b278DecodeGithubComErupshisEffectiveMobileInternalDatastructs7(l, v)
}
