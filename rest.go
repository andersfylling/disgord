package disgord

import (
	"bytes"
	"encoding/json"
	"github.com/andersfylling/disgord/httd"
	"io"
	"mime/multipart"
	"strconv"

	"github.com/json-iterator/go"
)

// URLParameters converts a struct of values to a valid URL query string
type URLParameters interface {
	GetQueryString() string
}

// unmarshalJSONIterator https://github.com/json-iterator/go
func unmarshalJSONIterator(data []byte, v interface{}) (err error) {
	err = jsoniter.Unmarshal(data, v)
	return
}

// unmarshalSTD standard GoLang implementation
func unmarshalSTD(data []byte, v interface{}) (err error) {
	err = json.Unmarshal(data, v)
	return
}

func unmarshal(data []byte, v interface{}) error {
	return unmarshalJSONIterator(data, v)
}

// Helpers for file uploading in messages
func (f *CreateChannelMessageFileParams) write(i int, mp *multipart.Writer) error {
	w, err := mp.CreateFormFile("file"+strconv.FormatInt(int64(i), 10), f.FileName)
	if err != nil {
		return err
	}

	if _, err = io.Copy(w, f.Reader); err != nil {
		return err
	}

	return nil
}

func (p *CreateChannelMessageParams) prepare() (postBody interface{}, contentType string, err error) {
	if len(p.Files) == 0 {
		postBody = p
		contentType = httd.ContentTypeJSON
		return
	}

	// Set up a new multipart writer, as we'll be using this for the POST body instead
	buf := new(bytes.Buffer)
	mp := multipart.NewWriter(buf)

	// Write the existing JSON payload
	var payload []byte
	payload, err = json.Marshal(p)
	if err != nil {
		return
	}
	if err = mp.WriteField("payload_json", string(payload)); err != nil {
		return
	}

	// Iterate through all the files and write them to the multipart blob
	for i, file := range p.Files {
		if err = file.write(i, mp); err != nil {
			return
		}
	}

	mp.Close()

	postBody = buf
	contentType = mp.FormDataContentType()

	return
}
