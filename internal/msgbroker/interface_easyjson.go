// Code generated by easyjson for marshaling/unmarshaling. DO NOT EDIT.

package msgbroker

import (
	json "encoding/json"
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

func easyjson43ae329fDecodeGithubComErupshisEffectiveMobileInternalMsgbroker(in *jlexer.Lexer, out *MessageBody) {
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
		case "name":
			out.Name = string(in.String())
		case "surname":
			out.Surname = string(in.String())
		case "patronymic":
			out.Patronymic = string(in.String())
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
func easyjson43ae329fEncodeGithubComErupshisEffectiveMobileInternalMsgbroker(out *jwriter.Writer, in MessageBody) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"name\":"
		out.RawString(prefix[1:])
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
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v MessageBody) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson43ae329fEncodeGithubComErupshisEffectiveMobileInternalMsgbroker(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v MessageBody) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson43ae329fEncodeGithubComErupshisEffectiveMobileInternalMsgbroker(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *MessageBody) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson43ae329fDecodeGithubComErupshisEffectiveMobileInternalMsgbroker(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *MessageBody) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson43ae329fDecodeGithubComErupshisEffectiveMobileInternalMsgbroker(l, v)
}
func easyjson43ae329fDecodeGithubComErupshisEffectiveMobileInternalMsgbroker1(in *jlexer.Lexer, out *Message) {
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
func easyjson43ae329fEncodeGithubComErupshisEffectiveMobileInternalMsgbroker1(out *jwriter.Writer, in Message) {
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

// MarshalJSON supports json.Marshaler interface
func (v Message) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson43ae329fEncodeGithubComErupshisEffectiveMobileInternalMsgbroker1(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v Message) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson43ae329fEncodeGithubComErupshisEffectiveMobileInternalMsgbroker1(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *Message) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson43ae329fDecodeGithubComErupshisEffectiveMobileInternalMsgbroker1(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *Message) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson43ae329fDecodeGithubComErupshisEffectiveMobileInternalMsgbroker1(l, v)
}
