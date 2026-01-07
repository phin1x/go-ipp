package ipp

import "io"

type ResponseDecoder struct {
	reader io.Reader
}

func NewResponseDecoder(reader io.Reader) *ResponseDecoder {
	return &ResponseDecoder{reader: reader}
}

func (r *ResponseDecoder) Decode(data io.Writer) (*Response, error) {
	return newResponseStateMachine().Decode(r.reader, data)
}
