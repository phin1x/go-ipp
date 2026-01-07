package ipp

import "net/http"

type HttpAdapterOption func(*HttpAdapter)

func WithHttpClient(client *http.Client) HttpAdapterOption {
	return func(adapter *HttpAdapter) {
		adapter.client = client
	}
}
