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

func (c *IPPClient) SendRequest(url string, req *Request, additionalResponseData io.Writer) (*Response, error) {
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
	return NewResponseDecoder(bytes.NewBuffer(httpBody)).Decode(additionalResponseData)

	//return NewResponseDecoder(httpResp.Body).Decode()
}

// Print one or more `Document`s using IPP `Create-Job` followed by `Send-Document` request(s).
func (c *IPPClient) PrintDocuments(docs []Document, printer string, jobAttributes map[string]interface{}) (int, error) {
	printerURI := c.getPrinterUri(printer)

	req := NewRequest(OperationCreateJob, 1)
	req.OperationAttributes[OperationAttributePrinterURI] = printerURI
	req.OperationAttributes[OperationAttributeRequestingUserName] = c.username

	// set defaults for some attributes, may get overwritten
	req.OperationAttributes[OperationAttributeJobName] = docs[0].Name
	req.OperationAttributes[OperationAttributeCopies] = 1
	req.OperationAttributes[OperationAttributeJobPriority] = DefaultJobPriority

	for key, value := range jobAttributes {
		req.JobAttributes[key] = value
	}

	resp, err := c.SendRequest(c.getHttpUri("printers", printer), req, nil)
	if err != nil {
		return -1, err
	}

	if len(resp.Jobs) == 0 {
		return 0, errors.New("server doesn't returned a job id")
	}

	jobID := resp.Jobs[0][OperationAttributeJobID][0].Value.(int)

	documentCount := len(docs) - 1

	for docID, doc := range docs {
		req = NewRequest(OperationSendDocument, 2)
		req.OperationAttributes[OperationAttributePrinterURI] = printerURI
		req.OperationAttributes[OperationAttributeRequestingUserName] = c.username
		req.OperationAttributes[OperationAttributeJobID] = jobID
		req.OperationAttributes[OperationAttributeDocumentName] = doc.Name
		req.OperationAttributes[OperationAttributeDocumentFormat] = doc.MimeType
		req.OperationAttributes[OperationAttributeLastDocument] = docID == documentCount
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
	req.OperationAttributes[OperationAttributePrinterURI] = printerURI
	req.OperationAttributes[OperationAttributeRequestingUserName] = c.username
	req.OperationAttributes[OperationAttributeJobName] = doc.Name
	req.OperationAttributes[OperationAttributeDocumentFormat] = doc.MimeType

	// set defaults for some attributes, may get overwritten
	req.OperationAttributes[OperationAttributeCopies] = 1
	req.OperationAttributes[OperationAttributeJobPriority] = DefaultJobPriority

	for key, value := range jobAttributes {
		req.JobAttributes[key] = value
	}

	req.File = doc.Document
	req.FileSize = doc.Size

	resp, err := c.SendRequest(c.getHttpUri("printers", printer), req, nil)
	if err != nil {
		return -1, err
	}

	if len(resp.Jobs) == 0 {
		return 0, errors.New("server doesn't returned a job id")
	}

	jobID := resp.Jobs[0][OperationAttributeJobID][0].Value.(int)

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

	jobAttributes[OperationAttributeJobName] = fileName

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
	req.OperationAttributes[OperationAttributePrinterURI] = c.getPrinterUri(printer)
	req.OperationAttributes[OperationAttributeRequestingUserName] = c.username

	if attributes == nil {
		req.OperationAttributes[OperationAttributeRequestedAttributes] = DefaultPrinterAttributes
	} else {
		req.OperationAttributes[OperationAttributeRequestedAttributes] = attributes
	}

	resp, err := c.SendRequest(c.getHttpUri("printers", printer), req, nil)
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
	req.OperationAttributes[OperationAttributePrinterURI] = c.getPrinterUri(printer)

	_, err := c.SendRequest(c.getHttpUri("admin", ""), req, nil)
	return err
}

func (c *IPPClient) PausePrinter(printer string) error {
	req := NewRequest(OperationPausePrinter, 1)
	req.OperationAttributes[OperationAttributePrinterURI] = c.getPrinterUri(printer)

	_, err := c.SendRequest(c.getHttpUri("admin", ""), req, nil)
	return err
}

func (c *IPPClient) GetJobAttributes(jobID int, attributes []string) (Attributes, error) {
	req := NewRequest(OperationGetJobAttributes, 1)
	req.OperationAttributes[OperationAttributeJobURI] = c.getJobUri(jobID)

	if attributes == nil {
		req.OperationAttributes[OperationAttributeRequestedAttributes] = DefaultJobAttributes
	} else {
		req.OperationAttributes[OperationAttributeRequestedAttributes] = attributes
	}

	resp, err := c.SendRequest(c.getHttpUri("jobs", jobID), req, nil)
	if err != nil {
		return nil, err
	}

	if len(resp.Printers) == 0 {
		return nil, errors.New("server doesn't return any job attributes")
	}

	return resp.Printers[0], nil
}

func (c *IPPClient) GetJobs(printer, class string, whichJobs JobStateFilter, myJobs bool, firstJobId, limit int, attributes []string) (map[int]Attributes, error) {
	req := NewRequest(OperationGetJobs, 1)
	req.OperationAttributes[OperationAttributeWhichJobs] = string(whichJobs)
	req.OperationAttributes[OperationAttributeMyJobs] = myJobs

	if printer != "" {
		req.OperationAttributes[OperationAttributePrinterURI] = c.getPrinterUri(printer)
	} else if class != "" {
		req.OperationAttributes[OperationAttributePrinterURI] = c.getClassUri(printer)
	} else {
		req.OperationAttributes[OperationAttributePrinterURI] = "ipp://localhost/"
	}

	if firstJobId > 0 {
		req.OperationAttributes[OperationAttributeFirstJobID] = firstJobId
	}

	if limit > 0 {
		req.OperationAttributes[OperationAttributeLimit] = limit
	}

	if myJobs {
		req.OperationAttributes[OperationAttributeRequestingUserName] = c.username
	}

	if attributes == nil {
		req.OperationAttributes[OperationAttributeRequestedAttributes] = DefaultJobAttributes
	} else {
		req.OperationAttributes[OperationAttributeRequestedAttributes] = append(attributes, OperationAttributeJobID)
	}

	resp, err := c.SendRequest(c.getHttpUri("", nil), req, nil)
	if err != nil {
		return nil, err
	}

	jobIDMap := make(map[int]Attributes)

	for _, jobAttributes := range resp.Jobs {
		jobIDMap[jobAttributes[OperationAttributeJobID][0].Value.(int)] = jobAttributes
	}

	return jobIDMap, nil
}

func (c *IPPClient) CancelJob(jobID int, purge bool) error {
	req := NewRequest(OperationCancelJob, 1)
	req.OperationAttributes[OperationAttributeJobURI] = c.getJobUri(jobID)
	req.OperationAttributes[OperationAttributePurgeJobs] = purge

	_, err := c.SendRequest(c.getHttpUri("jobs", ""), req, nil)
	return err
}

func (c *IPPClient) CancelAllJob(printer string, purge bool) error {
	req := NewRequest(OperationCancelJobs, 1)
	req.OperationAttributes[OperationAttributePrinterURI] = c.getPrinterUri(printer)
	req.OperationAttributes[OperationAttributePurgeJobs] = purge

	_, err := c.SendRequest(c.getHttpUri("admin", ""), req, nil)
	return err
}

func (c *IPPClient) RestartJob(jobID int) error {
	req := NewRequest(OperationRestartJob, 1)
	req.OperationAttributes[OperationAttributeJobURI] = c.getJobUri(jobID)

	_, err := c.SendRequest(c.getHttpUri("jobs", ""), req, nil)
	return err
}

func (c *IPPClient) HoldJobUntil(jobID int, holdUntil string) error {
	req := NewRequest(OperationRestartJob, 1)
	req.OperationAttributes[OperationAttributeJobURI] = c.getJobUri(jobID)
	req.JobAttributes[PrinterAttributeHoldJobUntil] = holdUntil

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
		OperationAttributeJobName: "Test Page",
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
