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
	DefaultSocketSearchPaths = []string{"/var/run/cupsd", "/var/run/cups/cups.sock", "/run/cups/cups.sock"}
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

//DoRequest performs the given IPP request to the given URL, returning the IPP response or an error if one occurred.
//Additional data will be written to an io.Writer if additionalData is not nil
func (h *SocketAdapter) SendRequest(url string, r *Request, additionalData io.Writer) (*Response, error) {
	for i := 0; i < h.RequestRetryLimit; i++ {
		// encode request
		payload, err := r.Encode()
		if err != nil {
			return nil, fmt.Errorf("unable to encode IPP request: %v", err)
		}

		var body io.Reader
		size := len(payload)

		if r.File != nil && r.FileSize != -1 {
			size += r.FileSize

			body = io.MultiReader(bytes.NewBuffer(payload), r.File)
		} else {
			body = bytes.NewBuffer(payload)
		}

		req, err := http.NewRequest("POST", url, body)
		if err != nil {
			return nil, fmt.Errorf("unable to create HTTP request: %v", err)
		}

		sock, err := h.GetSocket()
		if err != nil {
			return nil, err
		}

		// if cert isn't found, do a request to generate it
		cert, err := h.GetCert()
		if err != nil && err != CertNotFoundError {
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
		resp, err := unixClient.Do(req)
		if err != nil {
			return nil, fmt.Errorf("unable to perform HTTP request: %v", err)
		}

		if resp.StatusCode == http.StatusUnauthorized {
			// retry with newly generated cert
			resp.Body.Close()
			continue
		}

		if resp.StatusCode != http.StatusOK {
			resp.Body.Close()
			return nil, fmt.Errorf("server did not return Status OK: %d", resp.StatusCode)
		}

		// buffer response to avoid read issues
		buf := new(bytes.Buffer)
		if _, err := io.Copy(buf, resp.Body); err != nil {
			resp.Body.Close()
			return nil, fmt.Errorf("unable to buffer response: %v", err)
		}

		resp.Body.Close()

		// decode reply
		ippResp, err := NewResponseDecoder(bytes.NewReader(buf.Bytes())).Decode(additionalData)
		if err != nil {
			return nil, fmt.Errorf("unable to decode IPP response: %v", err)
		}

		if err = ippResp.CheckForErrors(); err != nil {
			return nil, fmt.Errorf("received error IPP response: %v", err)
		}

		return ippResp, nil
	}

	return nil, errors.New("request retry limit exceeded")
}

//GetSocket returns the path to the cupsd socket by searching SocketSearchPaths
func (h *SocketAdapter) GetSocket() (string, error) {
	for _, path := range h.SocketSearchPaths {
		fi, err := os.Stat(path)
		if err != nil {
			if os.IsNotExist(err) {
				continue
			} else if os.IsPermission(err) {
				return "", errors.New("unable to access socket: Access denied")
			}
			return "", fmt.Errorf("unable to access socket: %v", err)
		}

		if fi.Mode()&os.ModeSocket != 0 {
			return path, nil
		}
	}

	return "", SocketNotFoundError
}

//GetCert returns the current CUPs authentication certificate by searching CertSearchPaths
func (h *SocketAdapter) GetCert() (string, error) {
	for _, path := range h.CertSearchPaths {
		f, err := os.Open(path)
		if err != nil {
			if os.IsNotExist(err) {
				continue
			} else if os.IsPermission(err) {
				return "", errors.New("unable to access certificate: Access denied")
			}
			return "", fmt.Errorf("unable to access certificate: %v", err)
		}
		defer f.Close()

		buf := new(bytes.Buffer)
		if _, err := io.Copy(buf, f); err != nil {
			return "", fmt.Errorf("unable to access certificate: %v", err)
		}
		return buf.String(), nil
	}

	return "", CertNotFoundError
}

func (h *SocketAdapter) GetHttpUri(namespace string, object interface{}) string {
	proto := "http"
	if h.useTLS {
		proto = "https"
	}

	uri := fmt.Sprintf("%s://%s", proto, h.host)

	if namespace != "" {
		uri = fmt.Sprintf("%s/%s", uri, namespace)
	}

	if object != nil {
		uri = fmt.Sprintf("%s/%v", uri, object)
	}

	return uri
}

func (h *SocketAdapter) TestConnection() error {
	sock, err := h.GetSocket()
	if err != nil {
		return err
	}
	conn, err := net.Dial("unix", sock)
	if err != nil {
		return err
	}
	conn.Close()
	return nil
}
