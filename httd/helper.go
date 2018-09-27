package httd

import (
	"bytes"
	"compress/zlib"
	"io"
)

func BinaryToText(packet []byte) (text []byte, err error) {
	b := bytes.NewReader(packet)
	var r io.ReadCloser

	r, err = zlib.NewReader(b)
	if err != nil {
		return
	}
	defer r.Close()

	buffer := new(bytes.Buffer)
	_, err = buffer.ReadFrom(r)
	if err != nil {
		return
	}

	text = buffer.Bytes()
	return
}
