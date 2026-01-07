package ipp

import (
	"bytes"
	"context"
	"crypto/tls"
	"fmt"
	"io"
	"net"
	"net/http"
	"strconv"
)

type HttpAdapter struct {
	host     string
	port     int
	username string
	password string
	useTLS   bool
	client   *http.Client
}

func NewHttpAdapter(host string, port int, username, password string, useTLS bool, opts ...HttpAdapterOption) *HttpAdapter {
	adapter := &HttpAdapter{
		host:     host,
		port:     port,
		username: username,
		password: password,
		useTLS:   useTLS,
	}

	for _, opt := range opts {
		opt(adapter)
	}

	if adapter.client == nil {
		adapter.client = &http.Client{
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{
					InsecureSkipVerify: true,
				},
			},
		}
	}

	return adapter
}

func (a *HttpAdapter) SendRequest(url string, req *Request, additionalData io.Writer) (*Response, error) {
	return a.SendRequestContext(context.Background(), url, req, additionalData)
}

func (a *HttpAdapter) SendRequestContext(ctx context.Context, url string, req *Request, additionalData io.Writer) (*Response, error) {
	payload, err := req.Encode()
	if err != nil {
		return nil, err
	}

	size := len(payload)
	var body io.Reader
	if req.File != nil && req.FileSize != -1 {
		size += req.FileSize
		body = io.MultiReader(bytes.NewBuffer(payload), req.File)
	} else {
		body = bytes.NewBuffer(payload)
	}

	httpReq, err := http.NewRequestWithContext(ctx, "POST", url, body)
	if err != nil {
		return nil, err
	}

	httpReq.Header.Set("Content-Length", strconv.Itoa(size))
	httpReq.Header.Set("Content-Type", ContentTypeIPP)

	if a.username != "" && a.password != "" {
		httpReq.SetBasicAuth(a.username, a.password)
	}

	httpResp, err := a.client.Do(httpReq)
	if err != nil {
		return nil, err
	}
	defer httpResp.Body.Close()

	if httpResp.StatusCode != 200 {
		return nil, HTTPError{
			Code: httpResp.StatusCode,
		}
	}

	// buffer response to avoid read issues
	buf := new(bytes.Buffer)
	if httpResp.ContentLength > 0 {
		buf.Grow(int(httpResp.ContentLength))
	}
	if _, err := io.Copy(buf, httpResp.Body); err != nil {
		return nil, fmt.Errorf("unable to buffer response: %w", err)
	}

	ippResp, err := NewResponseDecoder(buf).Decode(additionalData)
	if err != nil {
		return nil, err
	}

	if err = ippResp.CheckForErrors(); err != nil {
		return nil, fmt.Errorf("received error IPP response: %w", err)
	}

	return ippResp, nil
}

func (a *HttpAdapter) GetHttpUri(namespace string, object interface{}) string {
	proto := "http"
	if a.useTLS {
		proto = "https"
	}

	uri := fmt.Sprintf("%s://%s:%d", proto, a.host, a.port)

	if namespace != "" {
		uri = fmt.Sprintf("%s/%s", uri, namespace)
	}

	if object != nil {
		uri = fmt.Sprintf("%s/%v", uri, object)
	}

	return uri
}

func (a *HttpAdapter) TestConnection() error {
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", a.host, a.port))
	if err != nil {
		return err
	}
	_ = conn.Close()
	return nil
}
