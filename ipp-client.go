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
	host     string
	port     int
	username string
	password string
	useTLS   bool

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

	if req.File != nil && req.FileSize != -1 {
		size += int(req.FileSize)

		body = io.MultiReader(bytes.NewBuffer(payload), req.File)
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

	//httpyBody, _ := ioutil.ReadAll(httpResp.Body)
	//fmt.Println(httpyBody)

	return NewResponseDecoder(httpResp.Body).Decode()
}

func (c *IPPClient) Print(document io.Reader, size int, printer, jobName string, copies, priority int) (int, error) {
	printerURI := fmt.Sprintf("ipp://localhost/printers/%s", printer)

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
	fmt.Println(jobID)

	req = NewRequest(OperationSendDocument, 2)
	req.OperationAttributes["printer-uri"] = printerURI
	req.OperationAttributes["requesting-user-name"] = c.username
	req.OperationAttributes["job-id"] = jobID
	req.OperationAttributes["document-name"] = jobName
	req.OperationAttributes["document-format"] = "application/octet-stream"
	req.OperationAttributes["last-document"] = true
	req.File = document
	req.FileSize = size

	resp, err = c.call(c.getHttpUri("printers", printerURI), req)
	if err != nil {
		return -1, err
	}

	return jobID, nil
}

func (c *IPPClient) PrintFile(filePath, printer, jobName string, copies, priority int) (int, error) {
	fileStats, err := os.Stat(filePath)
	if os.IsNotExist(err) {
		return -1, err
	}

	if jobName == "" {
		jobName = path.Base(filePath)
	}

	document, err := os.Open(filePath)
	if err != nil {
		return 0, err
	}
	defer document.Close()

	return c.Print(document, int(fileStats.Size()), printer, jobName, copies, priority)
}
