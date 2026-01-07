package ipp

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"strconv"
)

var SocketNotFoundError = errors.New("unable to locate CUPS socket")
var CertNotFoundError = errors.New("unable to locate CUPS certificate")

var (
	DefaultSocketSearchPaths = []string{"/var/run/cupsd", "/var/run/cups/cups.sock", "/run/cups/cups.sock", "/private/var/run/cupsd"}
	DefaultCertSearchPaths   = []string{"/etc/cups/certs/0", "/run/cups/certs/0"}
)

const DefaultRequestRetryLimit = 3

type SocketAdapter struct {
	host              string
	useTLS            bool
	SocketSearchPaths []string
	CertSearchPaths   []string
	//RequestRetryLimit is the number of times a request will be retried when receiving an authorized status. This usually happens when a CUPs cert is expired, and a retry will use the newly generated cert. Default 3.
	RequestRetryLimit int
}

func NewSocketAdapter(host string, useTLS bool) *SocketAdapter {
	return &SocketAdapter{
		host:              host,
		useTLS:            useTLS,
		SocketSearchPaths: DefaultSocketSearchPaths,
		CertSearchPaths:   DefaultCertSearchPaths,
		RequestRetryLimit: DefaultRequestRetryLimit,
	}
}

// SendRequest performs the given IPP request to the given URL, returning the IPP response or an error if one occurred.
// Additional data will be written to an io.Writer if additionalData is not nil
func (a *SocketAdapter) SendRequest(url string, r *Request, additionalData io.Writer) (*Response, error) {
	return a.SendRequestContext(context.Background(), url, r, additionalData)
}

func (a *SocketAdapter) SendRequestContext(ctx context.Context, url string, r *Request, additionalData io.Writer) (*Response, error) {
	for i := 0; i < a.RequestRetryLimit; i++ {
		// encode request
		payload, err := r.Encode()
		if err != nil {
			return nil, fmt.Errorf("unable to encode IPP request: %w", err)
		}

		var body io.Reader
		size := len(payload)

		if r.File != nil && r.FileSize != -1 {
			size += r.FileSize

			body = io.MultiReader(bytes.NewBuffer(payload), r.File)
		} else {
			body = bytes.NewBuffer(payload)
		}

		req, err := http.NewRequestWithContext(ctx, "POST", url, body)
		if err != nil {
			return nil, fmt.Errorf("unable to create HTTP request: %w", err)
		}

		sock, err := a.GetSocket()
		if err != nil {
			return nil, err
		}

		// if cert isn't found, do a request to generate it
		cert, err := a.GetCert()
		if err != nil && !errors.Is(err, CertNotFoundError) {
			return nil, err
		}

		req.Header.Set("Content-Length", strconv.Itoa(size))
		req.Header.Set("Content-Type", ContentTypeIPP)
		req.Header.Set("Authorization", fmt.Sprintf("Local %s", cert))

		unixClient := http.Client{
			Transport: &http.Transport{
				DialContext: func(ctx context.Context, _, _ string) (net.Conn, error) {
					return net.Dial("unix", sock)
				},
			},
		}

		// send request
		httpResp, err := unixClient.Do(req)
		if err != nil {
			return nil, fmt.Errorf("unable to perform HTTP request: %w", err)
		}

		if httpResp.StatusCode == http.StatusUnauthorized {
			// retry with newly generated cert
			httpResp.Body.Close()
			continue
		}

		if httpResp.StatusCode != http.StatusOK {
			httpResp.Body.Close()
			return nil, fmt.Errorf("server did not return Status OK: %d", httpResp.StatusCode)
		}

		// buffer response to avoid read issues
		buf := new(bytes.Buffer)
		if httpResp.ContentLength > 0 {
			buf.Grow(int(httpResp.ContentLength))
		}
		if _, err := io.Copy(buf, httpResp.Body); err != nil {
			httpResp.Body.Close()
			return nil, fmt.Errorf("unable to buffer response: %w", err)
		}

		httpResp.Body.Close()

		// decode reply
		ippResp, err := NewResponseDecoder(buf).Decode(additionalData)
		if err != nil {
			return nil, fmt.Errorf("unable to decode IPP response: %w", err)
		}

		if err = ippResp.CheckForErrors(); err != nil {
			return nil, fmt.Errorf("received error IPP response: %w", err)
		}

		return ippResp, nil
	}

	return nil, errors.New("request retry limit exceeded")
}

// GetSocket returns the path to the cupsd socket by searching SocketSearchPaths
func (a *SocketAdapter) GetSocket() (string, error) {
	for _, path := range a.SocketSearchPaths {
		fi, err := os.Stat(path)
		if err != nil {
			if os.IsNotExist(err) {
				continue
			} else if os.IsPermission(err) {
				return "", errors.New("unable to access socket: Access denied")
			}
			return "", fmt.Errorf("unable to access socket: %w", err)
		}

		if fi.Mode()&os.ModeSocket != 0 {
			return path, nil
		}
	}

	return "", SocketNotFoundError
}

// GetCert returns the current CUPs authentication certificate by searching CertSearchPaths
func (a *SocketAdapter) GetCert() (string, error) {
	for _, path := range a.CertSearchPaths {
		f, err := os.Open(path)
		if err != nil {
			if os.IsNotExist(err) {
				continue
			} else if os.IsPermission(err) {
				return "", errors.New("unable to access certificate: Access denied")
			}
			return "", fmt.Errorf("unable to access certificate: %w", err)
		}
		defer f.Close()

		buf := new(bytes.Buffer)
		if _, err := io.Copy(buf, f); err != nil {
			return "", fmt.Errorf("unable to access certificate: %w", err)
		}
		return buf.String(), nil
	}

	return "", CertNotFoundError
}

func (a *SocketAdapter) GetHttpUri(namespace string, object interface{}) string {
	proto := "http"
	if a.useTLS {
		proto = "https"
	}

	uri := fmt.Sprintf("%s://%s", proto, a.host)

	if namespace != "" {
		uri = fmt.Sprintf("%s/%s", uri, namespace)
	}

	if object != nil {
		uri = fmt.Sprintf("%s/%v", uri, object)
	}

	return uri
}

func (a *SocketAdapter) TestConnection() error {
	sock, err := a.GetSocket()
	if err != nil {
		return err
	}
	conn, err := net.Dial("unix", sock)
	if err != nil {
		return err
	}
	_ = conn.Close()
	return nil
}
