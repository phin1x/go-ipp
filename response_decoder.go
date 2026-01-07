package ipp

import "io"

type ResponseDecoder struct {
	reader io.Reader
}

func NewResponseDecoder(reader io.Reader) *ResponseDecoder {
	return &ResponseDecoder{reader: reader}
}

func (r *ResponseDecoder) Decode(reader io.Reader) (*Response, error) {
	return newResponseStateMachine().Decode(reader)
}
