// Code generated by easyjson for marshaling/unmarshaling. DO NOT EDIT.

package trunks

import (
	json "encoding/json"
	easyjson "github.com/mailru/easyjson"
	jlexer "github.com/mailru/easyjson/jlexer"
	jwriter "github.com/mailru/easyjson/jwriter"
	http "net/http"
	time "time"
)

// suppress unused package warning
var (
	_ *json.RawMessage
	_ *jlexer.Lexer
	_ *jwriter.Writer
	_ easyjson.Marshaler
)

func easyjsonBd1621b8DecodeGithubComStraightdaveTrunksLib(in *jlexer.Lexer, out *jsonResult) {
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
		case "attack":
			out.Attack = string(in.String())
		case "seq":
			out.Seq = uint64(in.Uint64())
		case "code":
			out.Code = uint16(in.Uint16())
		case "timestamp":
			if data := in.Raw(); in.Ok() {
				in.AddError((out.Timestamp).UnmarshalJSON(data))
			}
		case "latency":
			out.Latency = time.Duration(in.Int64())
		case "bytes_out":
			out.BytesOut = uint64(in.Uint64())
		case "bytes_in":
			out.BytesIn = uint64(in.Uint64())
		case "error":
			out.Error = string(in.String())
		case "body":
			if in.IsNull() {
				in.Skip()
				out.Body = nil
			} else {
				out.Body = in.Bytes()
			}
		case "method":
			out.Method = string(in.String())
		case "url":
			out.URL = string(in.String())
		case "headers":
			if in.IsNull() {
				in.Skip()
			} else {
				in.Delim('{')
				out.Headers = make(http.Header)
				for !in.IsDelim('}') {
					key := string(in.String())
					in.WantColon()
					var v2 []string
					if in.IsNull() {
						in.Skip()
						v2 = nil
					} else {
						in.Delim('[')
						if v2 == nil {
							if !in.IsDelim(']') {
								v2 = make([]string, 0, 4)
							} else {
								v2 = []string{}
							}
						} else {
							v2 = (v2)[:0]
						}
						for !in.IsDelim(']') {
							var v3 string
							v3 = string(in.String())
							v2 = append(v2, v3)
							in.WantComma()
						}
						in.Delim(']')
					}
					(out.Headers)[key] = v2
					in.WantComma()
				}
				in.Delim('}')
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
func easyjsonBd1621b8EncodeGithubComStraightdaveTrunksLib(out *jwriter.Writer, in jsonResult) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"attack\":"
		out.RawString(prefix[1:])
		out.String(string(in.Attack))
	}
	{
		const prefix string = ",\"seq\":"
		out.RawString(prefix)
		out.Uint64(uint64(in.Seq))
	}
	{
		const prefix string = ",\"code\":"
		out.RawString(prefix)
		out.Uint16(uint16(in.Code))
	}
	{
		const prefix string = ",\"timestamp\":"
		out.RawString(prefix)
		out.Raw((in.Timestamp).MarshalJSON())
	}
	{
		const prefix string = ",\"latency\":"
		out.RawString(prefix)
		out.Int64(int64(in.Latency))
	}
	{
		const prefix string = ",\"bytes_out\":"
		out.RawString(prefix)
		out.Uint64(uint64(in.BytesOut))
	}
	{
		const prefix string = ",\"bytes_in\":"
		out.RawString(prefix)
		out.Uint64(uint64(in.BytesIn))
	}
	{
		const prefix string = ",\"error\":"
		out.RawString(prefix)
		out.String(string(in.Error))
	}
	{
		const prefix string = ",\"body\":"
		out.RawString(prefix)
		out.Base64Bytes(in.Body)
	}
	{
		const prefix string = ",\"method\":"
		out.RawString(prefix)
		out.String(string(in.Method))
	}
	{
		const prefix string = ",\"url\":"
		out.RawString(prefix)
		out.String(string(in.URL))
	}
	{
		const prefix string = ",\"headers\":"
		out.RawString(prefix)
		if in.Headers == nil && (out.Flags&jwriter.NilMapAsEmpty) == 0 {
			out.RawString(`null`)
		} else {
			out.RawByte('{')
			v6First := true
			for v6Name, v6Value := range in.Headers {
				if v6First {
					v6First = false
				} else {
					out.RawByte(',')
				}
				out.String(string(v6Name))
				out.RawByte(':')
				if v6Value == nil && (out.Flags&jwriter.NilSliceAsEmpty) == 0 {
					out.RawString("null")
				} else {
					out.RawByte('[')
					for v7, v8 := range v6Value {
						if v7 > 0 {
							out.RawByte(',')
						}
						out.String(string(v8))
					}
					out.RawByte(']')
				}
			}
			out.RawByte('}')
		}
	}
	out.RawByte('}')
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v jsonResult) MarshalEasyJSON(w *jwriter.Writer) {
	easyjsonBd1621b8EncodeGithubComStraightdaveTrunksLib(w, v)
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *jsonResult) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjsonBd1621b8DecodeGithubComStraightdaveTrunksLib(l, v)
}
