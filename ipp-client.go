package ipp

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"strconv"
)

const (
	ippContentType = "application/ipp"
)

type IPPClient struct {
	host string
	port int
	username string
	password string
	useTLS bool

	client *http.Client
}

func NewIPPClient(host string, port int, username, password string, useTLS bool) *IPPClient {
	httpClient := http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
	}

	return &IPPClient{host, port, username, password, useTLS, &httpClient}
}

func (c *IPPClient) getHttpUri(namespace, object string) string {
	proto := "http"
	if c.useTLS {
		proto = "https"
	}
	return fmt.Sprintf("%s://%s:%d/%s/%s", proto, c.host, c.port, namespace, object)
}

func (c *IPPClient) call(url string, req *Request) (*Response, error) {
	payload, err := req.Encode()
	if err != nil {
		return nil, err
	}

	var body io.Reader
	size := len(payload)

	if req.File != "" {
		fileStats, err := os.Stat(req.File)
		if os.IsNotExist(err) {
			return nil, err
		}
		size += int(fileStats.Size())

		fileReader, err := os.Open(req.File)
		if err != nil {
			return nil, err
		}
		defer fileReader.Close()

		body = io.MultiReader(bytes.NewBuffer(payload), fileReader)
	} else {
		body = bytes.NewBuffer(payload)
	}

	httpReq, err := http.NewRequest("POST", url, body)
	if err != nil {
		return nil, err
	}

	httpReq.Header.Set("Content-Length", strconv.Itoa(size))
	httpReq.Header.Set("Content-Type", ippContentType)

	if c.username != "" && c.password != "" {
		httpReq.SetBasicAuth(c.username, c.password)
	}

	httpResp, err := c.client.Do(httpReq)
	if err != nil {
		return nil, err
	}
	defer httpResp.Body.Close()

	if httpResp.StatusCode != 200 {
		return nil, fmt.Errorf("ipp server returned with http status code %d", httpResp.StatusCode)
	}

	return NewResponseDecoder(httpResp.Body).Decode()
}

func (c *IPPClient) PrintFile(printer, filePath, jobName string, copies, priority int) (int, error) {
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return -1, err
	}

	printerURI := fmt.Sprintf("ipp://localhost/printers/%s", printer)

	if jobName == "" {
		jobName = path.Base(filePath)
	}

	req := NewRequest(OperationCreateJob, 1)
	req.OperationAttributes["printer-uri"] = printerURI
	req.OperationAttributes["requesting-user-name"] = c.username
	req.OperationAttributes["job-name"] = c.username
	req.JobAttributes["copies"] = copies
	req.JobAttributes["job-priority"] = priority

	resp, err := c.call(c.getHttpUri("printers", printerURI), req)
	if err != nil {
		return -1, err
	}

	jobID := resp.Jobs[0]["job-id"][0].Value.(int)

	req = NewRequest(OperationSendDocument, 2)
	req.OperationAttributes["printer-uri"] = printerURI
	req.OperationAttributes["requesting-user-name"] = c.username
	req.OperationAttributes["job-id"] = jobID
	req.OperationAttributes["document-name1"] = jobName
	req.OperationAttributes["document-format"] = "application/octet-stream"
	req.OperationAttributes[""] = true
	req.File = filePath

	resp, err = c.call(c.getHttpUri("printers", printerURI), req)
	if err != nil {
		return -1, err
	}

	return jobID, nil
}