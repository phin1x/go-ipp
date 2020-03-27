package ipp

import (
	"bytes"
	"crypto/tls"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"path"
	"strconv"
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

func (c *IPPClient) SendRequest(url string, req *Request, additionalResponseData io.Writer) (*Response, error) {
	payload, err := req.Encode()
	if err != nil {
		return nil, err
	}

	var body io.Reader
	size := len(payload)

	if req.File != nil && req.FileSize != -1 {
		size += req.FileSize

		body = io.MultiReader(bytes.NewBuffer(payload), req.File)
	} else {
		body = bytes.NewBuffer(payload)
	}

	httpReq, err := http.NewRequest("POST", url, body)
	if err != nil {
		return nil, err
	}

	httpReq.Header.Set("Content-Length", strconv.Itoa(size))
	httpReq.Header.Set("Content-Type", ContentTypeIPP)

	if c.username != "" && c.password != "" {
		httpReq.SetBasicAuth(c.username, c.password)
	}

	httpResp, err := c.client.Do(httpReq)
	if err != nil {
		return nil, err
	}
	defer httpResp.Body.Close()

	if httpResp.StatusCode != 200 {
		return nil, HTTPError{
			Code: httpResp.StatusCode,
		}
	}

	resp, err := NewResponseDecoder(httpResp.Body).Decode(additionalResponseData)
	if err != nil {
		return nil, err
	}

	err = resp.CheckForErrors()
	return resp, err
}

// Print one or more `Document`s using IPP `Create-Job` followed by `Send-Document` request(s).
func (c *IPPClient) PrintDocuments(docs []Document, printer string, jobAttributes map[string]interface{}) (int, error) {
	printerURI := c.getPrinterUri(printer)

	req := NewRequest(OperationCreateJob, 1)
	req.OperationAttributes[AttributePrinterURI] = printerURI
	req.OperationAttributes[AttributeRequestingUserName] = c.username

	// set defaults for some attributes, may get overwritten
	req.OperationAttributes[AttributeJobName] = docs[0].Name
	req.OperationAttributes[AttributeCopies] = 1
	req.OperationAttributes[AttributeJobPriority] = DefaultJobPriority

	for key, value := range jobAttributes {
		req.JobAttributes[key] = value
	}

	resp, err := c.SendRequest(c.getHttpUri("printers", printer), req, nil)
	if err != nil {
		return -1, err
	}

	if len(resp.JobAttributes) == 0 {
		return 0, errors.New("server doesn't returned a job id")
	}

	jobID := resp.JobAttributes[0][AttributeJobID][0].Value.(int)

	documentCount := len(docs) - 1

	for docID, doc := range docs {
		req = NewRequest(OperationSendDocument, 2)
		req.OperationAttributes[AttributePrinterURI] = printerURI
		req.OperationAttributes[AttributeRequestingUserName] = c.username
		req.OperationAttributes[AttributeJobID] = jobID
		req.OperationAttributes[AttributeDocumentName] = doc.Name
		req.OperationAttributes[AttributeDocumentFormat] = doc.MimeType
		req.OperationAttributes[AttributeLastDocument] = docID == documentCount
		req.File = doc.Document
		req.FileSize = doc.Size

		_, err = c.SendRequest(c.getHttpUri("printers", printer), req, nil)
		if err != nil {
			return -1, err
		}
	}

	return jobID, nil
}

// Print a `Document` using an IPP `Print-Job` request.
//
// `jobAttributes` can contain arbitrary key/value pairs to control the way in which the
// document is printed. [RFC 2911 § 4.2](https://tools.ietf.org/html/rfc2911#section-4.2)
// defines some useful attributes:
//
//   * [`job-priority`](https://tools.ietf.org/html/rfc2911#section-4.2.1): an integer between 1-100
//   * [`copies`](https://tools.ietf.org/html/rfc2911#section-4.2.5): a positive integer
//   * [`finishings`](https://tools.ietf.org/html/rfc2911#section-4.2.6): an enumeration
//   * [`number-up`](https://tools.ietf.org/html/rfc2911#section-4.2.9): a positive integer
//   * [`orientation-requested`](https://tools.ietf.org/html/rfc2911#section-4.2.10): an enumeration
//   * [`media`](https://tools.ietf.org/html/rfc2911#section-4.2.11): a string
//   * [`printer-resolution`](https://tools.ietf.org/html/rfc2911#section-4.2.12): a `Resolution`
//   * [`print-quality`](https://tools.ietf.org/html/rfc2911#section-4.2.13): an enumeration
//
// Your print system may provide other attributes. Define custom attributes as needed in
// `AttributeTagMapping` and provide values here.
func (c *IPPClient) PrintJob(doc Document, printer string, jobAttributes map[string]interface{}) (int, error) {
	printerURI := c.getPrinterUri(printer)

	req := NewRequest(OperationPrintJob, 1)
	req.OperationAttributes[AttributePrinterURI] = printerURI
	req.OperationAttributes[AttributeRequestingUserName] = c.username
	req.OperationAttributes[AttributeJobName] = doc.Name
	req.OperationAttributes[AttributeDocumentFormat] = doc.MimeType

	// set defaults for some attributes, may get overwritten
	req.OperationAttributes[AttributeCopies] = 1
	req.OperationAttributes[AttributeJobPriority] = DefaultJobPriority

	for key, value := range jobAttributes {
		req.JobAttributes[key] = value
	}

	req.File = doc.Document
	req.FileSize = doc.Size

	resp, err := c.SendRequest(c.getHttpUri("printers", printer), req, nil)
	if err != nil {
		return -1, err
	}

	if len(resp.JobAttributes) == 0 {
		return 0, errors.New("server doesn't returned a job id")
	}

	jobID := resp.JobAttributes[0][AttributeJobID][0].Value.(int)

	return jobID, nil
}

func (c *IPPClient) PrintFile(filePath, printer string, jobAttributes map[string]interface{}) (int, error) {
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

	jobAttributes[AttributeJobName] = fileName

	return c.PrintDocuments([]Document{
		{
			Document: document,
			Name:     fileName,
			Size:     int(fileStats.Size()),
			MimeType: MimeTypeOctetStream,
		},
	}, printer, jobAttributes)
}

func (c *IPPClient) GetPrinterAttributes(printer string, attributes []string) (Attributes, error) {
	req := NewRequest(OperationGetPrinterAttributes, 1)
	req.OperationAttributes[AttributePrinterURI] = c.getPrinterUri(printer)
	req.OperationAttributes[AttributeRequestingUserName] = c.username

	if attributes == nil {
		req.OperationAttributes[AttributeRequestedAttributes] = DefaultPrinterAttributes
	} else {
		req.OperationAttributes[AttributeRequestedAttributes] = attributes
	}

	resp, err := c.SendRequest(c.getHttpUri("printers", printer), req, nil)
	if err != nil {
		return nil, err
	}

	if len(resp.PrinterAttributes) == 0 {
		return nil, errors.New("server doesn't return any printer attributes")
	}

	return resp.PrinterAttributes[0], nil
}

func (c *IPPClient) ResumePrinter(printer string) error {
	req := NewRequest(OperationResumePrinter, 1)
	req.OperationAttributes[AttributePrinterURI] = c.getPrinterUri(printer)

	_, err := c.SendRequest(c.getHttpUri("admin", ""), req, nil)
	return err
}

func (c *IPPClient) PausePrinter(printer string) error {
	req := NewRequest(OperationPausePrinter, 1)
	req.OperationAttributes[AttributePrinterURI] = c.getPrinterUri(printer)

	_, err := c.SendRequest(c.getHttpUri("admin", ""), req, nil)
	return err
}

func (c *IPPClient) GetJobAttributes(jobID int, attributes []string) (Attributes, error) {
	req := NewRequest(OperationGetJobAttributes, 1)
	req.OperationAttributes[AttributeJobURI] = c.getJobUri(jobID)

	if attributes == nil {
		req.OperationAttributes[AttributeRequestedAttributes] = DefaultJobAttributes
	} else {
		req.OperationAttributes[AttributeRequestedAttributes] = attributes
	}

	resp, err := c.SendRequest(c.getHttpUri("jobs", jobID), req, nil)
	if err != nil {
		return nil, err
	}

	if len(resp.PrinterAttributes) == 0 {
		return nil, errors.New("server doesn't return any job attributes")
	}

	return resp.PrinterAttributes[0], nil
}

func (c *IPPClient) GetJobs(printer, class string, whichJobs string, myJobs bool, firstJobId, limit int, attributes []string) (map[int]Attributes, error) {
	req := NewRequest(OperationGetJobs, 1)
	req.OperationAttributes[AttributeWhichJobs] = string(whichJobs)
	req.OperationAttributes[AttributeMyJobs] = myJobs

	if printer != "" {
		req.OperationAttributes[AttributePrinterURI] = c.getPrinterUri(printer)
	} else if class != "" {
		req.OperationAttributes[AttributePrinterURI] = c.getClassUri(printer)
	} else {
		req.OperationAttributes[AttributePrinterURI] = "ipp://localhost/"
	}

	if firstJobId > 0 {
		req.OperationAttributes[AttributeFirstJobID] = firstJobId
	}

	if limit > 0 {
		req.OperationAttributes[AttributeLimit] = limit
	}

	if myJobs {
		req.OperationAttributes[AttributeRequestingUserName] = c.username
	}

	if attributes == nil {
		req.OperationAttributes[AttributeRequestedAttributes] = DefaultJobAttributes
	} else {
		req.OperationAttributes[AttributeRequestedAttributes] = append(attributes, AttributeJobID)
	}

	resp, err := c.SendRequest(c.getHttpUri("", nil), req, nil)
	if err != nil {
		return nil, err
	}

	jobIDMap := make(map[int]Attributes)

	for _, jobAttributes := range resp.JobAttributes {
		jobIDMap[jobAttributes[AttributeJobID][0].Value.(int)] = jobAttributes
	}

	return jobIDMap, nil
}

func (c *IPPClient) CancelJob(jobID int, purge bool) error {
	req := NewRequest(OperationCancelJob, 1)
	req.OperationAttributes[AttributeJobURI] = c.getJobUri(jobID)
	req.OperationAttributes[AttributePurgeJobs] = purge

	_, err := c.SendRequest(c.getHttpUri("jobs", ""), req, nil)
	return err
}

func (c *IPPClient) CancelAllJob(printer string, purge bool) error {
	req := NewRequest(OperationCancelJobs, 1)
	req.OperationAttributes[AttributePrinterURI] = c.getPrinterUri(printer)
	req.OperationAttributes[AttributePurgeJobs] = purge

	_, err := c.SendRequest(c.getHttpUri("admin", ""), req, nil)
	return err
}

func (c *IPPClient) RestartJob(jobID int) error {
	req := NewRequest(OperationRestartJob, 1)
	req.OperationAttributes[AttributeJobURI] = c.getJobUri(jobID)

	_, err := c.SendRequest(c.getHttpUri("jobs", ""), req, nil)
	return err
}

func (c *IPPClient) HoldJobUntil(jobID int, holdUntil string) error {
	req := NewRequest(OperationRestartJob, 1)
	req.OperationAttributes[AttributeJobURI] = c.getJobUri(jobID)
	req.JobAttributes[AttributeHoldJobUntil] = holdUntil

	_, err := c.SendRequest(c.getHttpUri("jobs", ""), req, nil)
	return err
}

func (c *IPPClient) PrintTestPage(printer string) (int, error) {
	testPage := new(bytes.Buffer)
	testPage.WriteString("#PDF-BANNER\n")
	testPage.WriteString("Template default-testpage.pdf\n")
	testPage.WriteString("Show printer-name printer-info printer-location printer-make-and-model printer-driver-name")
	testPage.WriteString("printer-driver-version paper-size imageable-area job-id options time-at-creation")
	testPage.WriteString("time-at-processing\n\n")

	return c.PrintDocuments([]Document{
		{
			Document: testPage,
			Name:     "Test Page",
			Size:     testPage.Len(),
			MimeType: MimeTypePostscript,
		},
	}, printer, map[string]interface{}{
		AttributeJobName: "Test Page",
	})
}

func (c *IPPClient) TestConnection() error {
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", c.host, c.port))
	if err != nil {
		return err
	}
	conn.Close()

	return nil
}
