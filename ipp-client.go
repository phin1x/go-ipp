package ipp

import (
	"bytes"
	"crypto/tls"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"path"
	"strconv"
)

const (
	ippContentType = "application/ipp"
)

type Document struct {
	Document io.Reader
	Size     int
	Name     string
	MimeType string
}

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

func (c *IPPClient) getHttpUri(namespace string, object interface{}) string {
	proto := "http"
	if c.useTLS {
		proto = "https"
	}

	uri := fmt.Sprintf("%s://%s:%d", proto, c.host, c.port)

	if namespace != "" {
		uri = fmt.Sprintf("%s/%s", uri, namespace)
	}

	if object != nil {
		uri = fmt.Sprintf("%s/%v", uri, object)
	}

	return uri
}

func (c *IPPClient) getPrinterUri(printer string) string {
	return fmt.Sprintf("ipp://localhost/printers/%s", printer)
}

func (c *IPPClient) getJobUri(jobID int) string {
	return fmt.Sprintf("ipp://localhost/jobs/%d", jobID)
}

func (c *IPPClient) getClassUri(printer string) string {
	return fmt.Sprintf("ipp://localhost/classes/%s", printer)
}

func (c *IPPClient) SendRequest(url string, req *Request) (*Response, error) {
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

	//read the response into a temp buffer due to some wired EOF errors
	httpBody, _ := ioutil.ReadAll(httpResp.Body)
	//fmt.Println(httpBody)
	return NewResponseDecoder(bytes.NewBuffer(httpBody)).Decode()

	//return NewResponseDecoder(httpResp.Body).Decode()
}

func (c *IPPClient) Print(docs []Document, printer, jobName string, copies, priority int) (int, error) {
	printerURI := c.getPrinterUri(printer)

	req := NewRequest(OperationCreateJob, 1)
	req.OperationAttributes["printer-uri"] = printerURI
	req.OperationAttributes["requesting-user-name"] = c.username
	req.OperationAttributes["job-name"] = c.username
	req.JobAttributes["copies"] = copies
	req.JobAttributes["job-priority"] = priority

	resp, err := c.SendRequest(c.getHttpUri("printers", printer), req)
	if err != nil {
		return -1, err
	}

	if len(resp.Jobs) == 0 {
		return 0, errors.New("server doesn't returned a job id")
	}

	jobID := resp.Jobs[0]["job-id"][0].Value.(int)

	documentCount := len(docs) - 1

	for docID, doc := range docs {
		req = NewRequest(OperationSendDocument, 2)
		req.OperationAttributes["printer-uri"] = printerURI
		req.OperationAttributes["requesting-user-name"] = c.username
		req.OperationAttributes["job-id"] = jobID
		req.OperationAttributes["document-name"] = doc.Name
		req.OperationAttributes["document-format"] = doc.MimeType
		req.OperationAttributes["last-document"] = docID == documentCount
		req.File = doc.Document
		req.FileSize = doc.Size

		resp, err = c.SendRequest(c.getHttpUri("printers", printer), req)
		if err != nil {
			return -1, err
		}
	}

	return jobID, nil
}

func (c *IPPClient) PrintFile(filePath, printer string, copies, priority int) (int, error) {
	fileStats, err := os.Stat(filePath)
	if os.IsNotExist(err) {
		return -1, err
	}

	fileName := path.Base(filePath)

	document, err := os.Open(filePath)
	if err != nil {
		return 0, err
	}
	defer document.Close()

	return c.Print([]Document{
		{
			Document: document,
			Name:     fileName,
			Size:     int(fileStats.Size()),
			MimeType: MimeTypeOctetStream,
		},
	}, printer, fileName, copies, priority)
}

func (c *IPPClient) GetPrinterAttributes(printer string, attributes []string) (Attributes, error) {
	req := NewRequest(OperationGetPrinterAttributes, 1)
	req.OperationAttributes["printer-uri"] = c.getPrinterUri(printer)
	req.OperationAttributes["requesting-user-name"] = c.username

	if attributes == nil {
		req.OperationAttributes["requested-attributes"] = DefaultPrinterAttributes
	} else {
		req.OperationAttributes["requested-attributes"] = attributes
	}

	resp, err := c.SendRequest(c.getHttpUri("printers", printer), req)
	if err != nil {
		return nil, err
	}

	if len(resp.Printers) == 0 {
		return nil, errors.New("server doesn't return any printer attributes")
	}

	return resp.Printers[0], nil
}

func (c *IPPClient) ResumePrinter(printer string) error {
	req := NewRequest(OperationResumePrinter, 1)
	req.OperationAttributes["printer-uri"] = c.getPrinterUri(printer)

	_, err := c.SendRequest(c.getHttpUri("admin", ""), req)
	return err
}

func (c *IPPClient) PausePrinter(printer string) error {
	req := NewRequest(OperationPausePrinter, 1)
	req.OperationAttributes["printer-uri"] = c.getPrinterUri(printer)

	_, err := c.SendRequest(c.getHttpUri("admin", ""), req)
	return err
}

func (c *IPPClient) GetJobAttributes(jobID int, attributes []string) (Attributes, error) {
	req := NewRequest(OperationGetJobAttributes, 1)
	req.OperationAttributes["job-uri"] = c.getJobUri(jobID)

	if attributes == nil {
		req.OperationAttributes["requested-attributes"] = DefaultJobAttributes
	} else {
		req.OperationAttributes["requested-attributes"] = attributes
	}

	resp, err := c.SendRequest(c.getHttpUri("jobs", jobID), req)
	if err != nil {
		return nil, err
	}

	if len(resp.Printers) == 0 {
		return nil, errors.New("server doesn't return any job attributes")
	}

	return resp.Printers[0], nil
}

func (c *IPPClient) GetJobs(printer string, whichJobs JobStateFilter, myJobs bool, attributes []string) (map[int]Attributes, error) {
	req := NewRequest(OperationGetJobs, 1)
	req.OperationAttributes["printer-uri"] = c.getPrinterUri(printer)
	req.OperationAttributes["which-jobs"] = string(whichJobs)
	req.OperationAttributes["my-jobs"] = myJobs

	if attributes == nil {
		req.OperationAttributes["requested-attributes"] = DefaultJobAttributes
	} else {
		req.OperationAttributes["requested-attributes"] = append(attributes, "job-id")
	}

	resp, err := c.SendRequest(c.getHttpUri("", nil), req)
	if err != nil {
		return nil, err
	}

	jobIDMap := make(map[int]Attributes)

	for _, jobAttributes := range resp.Jobs {
		jobIDMap[jobAttributes["job-id"][0].Value.(int)] = jobAttributes
	}

	return jobIDMap, nil
}

func (c *IPPClient) CancelJob(jobID int, purge bool) error {
	req := NewRequest(OperationCancelJob, 1)
	req.OperationAttributes["job-uri"] = c.getJobUri(jobID)
	req.OperationAttributes["purge-jobs"] = purge

	_, err := c.SendRequest(c.getHttpUri("jobs", ""), req)
	return err
}

func (c *IPPClient) CancelAllJob(printer string, purge bool) error {
	req := NewRequest(OperationCancelJobs, 1)
	req.OperationAttributes["printer-uri"] = c.getPrinterUri(printer)
	req.OperationAttributes["purge-jobs"] = purge

	_, err := c.SendRequest(c.getHttpUri("admin", ""), req)
	return err
}

func (c *IPPClient) RestartJob(jobID int) error {
	req := NewRequest(OperationRestartJob, 1)
	req.OperationAttributes["job-uri"] = c.getJobUri(jobID)

	_, err := c.SendRequest(c.getHttpUri("jobs", ""), req)
	return err
}

func (c *IPPClient) HoldJobUntil(jobID int, holdUntil string) error {
	req := NewRequest(OperationRestartJob, 1)
	req.OperationAttributes["job-uri"] = c.getJobUri(jobID)
	req.JobAttributes["job-hold-until"] = holdUntil

	_, err := c.SendRequest(c.getHttpUri("jobs", ""), req)
	return err
}

func (c *IPPClient) PrintTestPage(printer string) (int, error) {
	testPage := new(bytes.Buffer)
	testPage.WriteString("#PDF-BANNER\n")
	testPage.WriteString("Template default-testpage.pdf\n")
	testPage.WriteString("Show printer-name printer-info printer-location printer-make-and-model printer-driver-name")
	testPage.WriteString("printer-driver-version paper-size imageable-area job-id options time-at-creation")
	testPage.WriteString("time-at-processing\n\n")

	return c.Print([]Document{
		{
			Document: testPage,
			Name:     "Test Page",
			Size:     testPage.Len(),
			MimeType: MimeTypePostscript,
		},
	}, printer, "Test Page", 1, DefaultJobPriority)
}

func (c *IPPClient) TestConnection() error {
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", c.host, c.port))
	if err != nil {
		return err
	}
	conn.Close()

	return nil
}
