package ipp

import (
	"bytes"
	"context"
	"strings"
)

// CUPSClient implements a ipp client with specific cups operations
type CUPSClient struct {
	*IPPClient
}

// NewCUPSClient creates a new cups ipp client (used HttpAdapter internally)
func NewCUPSClient(host string, port int, username, password string, useTLS bool) *CUPSClient {
	ippClient := NewIPPClient(host, port, username, password, useTLS)
	return &CUPSClient{ippClient}
}

// NewCUPSClientWithAdapter creates a new cups ipp client with given Adapter
func NewCUPSClientWithAdapter(username string, adapter Adapter) *CUPSClient {
	ippClient := NewIPPClientWithAdapter(username, adapter)
	return &CUPSClient{ippClient}
}

// GetDevices returns a map of device uris and printer attributes
func (c *CUPSClient) GetDevices() (map[string]Attributes, error) {
	return c.GetDevicesContext(context.Background())
}

func (c *CUPSClient) GetDevicesContext(ctx context.Context) (map[string]Attributes, error) {
	req := NewRequest(OperationCupsGetDevices, 1)

	resp, err := c.SendRequestContext(ctx, c.adapter.GetHttpUri("", nil), req, nil)
	if err != nil {
		return nil, err
	}

	printerNameMap := make(map[string]Attributes)

	for _, printerAttributes := range resp.PrinterAttributes {
		printerNameMap[printerAttributes[AttributeDeviceURI][0].Value.(string)] = printerAttributes
	}

	return printerNameMap, nil
}

// MoveJob moves a job to a other printer
func (c *CUPSClient) MoveJob(jobID int, destPrinter string) error {
	return c.MoveJobContext(context.Background(), jobID, destPrinter)
}

func (c *CUPSClient) MoveJobContext(ctx context.Context, jobID int, destPrinter string) error {
	req := NewRequest(OperationCupsMoveJob, 1)
	req.OperationAttributes[AttributeJobURI] = c.getJobUri(jobID)
	req.PrinterAttributes[AttributeJobPrinterURI] = c.getPrinterUri(destPrinter)

	_, err := c.SendRequestContext(ctx, c.adapter.GetHttpUri("jobs", ""), req, nil)
	return err
}

// MoveAllJob moves all job from a printer to a other printer
func (c *CUPSClient) MoveAllJob(srcPrinter, destPrinter string) error {
	return c.MoveAllJobContext(context.Background(), srcPrinter, destPrinter)
}

func (c *CUPSClient) MoveAllJobContext(ctx context.Context, srcPrinter, destPrinter string) error {
	req := NewRequest(OperationCupsMoveJob, 1)
	req.OperationAttributes[AttributePrinterURI] = c.getPrinterUri(srcPrinter)
	req.PrinterAttributes[AttributeJobPrinterURI] = c.getPrinterUri(destPrinter)

	_, err := c.SendRequestContext(ctx, c.adapter.GetHttpUri("jobs", ""), req, nil)
	return err
}

// GetPPDs returns a map of ppd names and attributes
func (c *CUPSClient) GetPPDs() (map[string]Attributes, error) {
	return c.GetPPDsContext(context.Background())
}

func (c *CUPSClient) GetPPDsContext(ctx context.Context) (map[string]Attributes, error) {
	req := NewRequest(OperationCupsGetPPDs, 1)

	resp, err := c.SendRequestContext(ctx, c.adapter.GetHttpUri("", nil), req, nil)
	if err != nil {
		return nil, err
	}

	ppdNameMap := make(map[string]Attributes)

	for _, printerAttributes := range resp.PrinterAttributes {
		ppdNameMap[printerAttributes[AttributePPDName][0].Value.(string)] = printerAttributes
	}

	return ppdNameMap, nil
}

// AcceptJobs lets a printer accept jobs again
func (c *CUPSClient) AcceptJobs(printer string) error {
	return c.AcceptJobsContext(context.Background(), printer)
}

func (c *CUPSClient) AcceptJobsContext(ctx context.Context, printer string) error {
	req := NewRequest(OperationCupsAcceptJobs, 1)
	req.OperationAttributes[AttributePrinterURI] = c.getPrinterUri(printer)

	_, err := c.SendRequestContext(ctx, c.adapter.GetHttpUri("admin", ""), req, nil)
	return err
}

// RejectJobs does not let a printer accept jobs
func (c *CUPSClient) RejectJobs(printer string) error {
	return c.RejectJobsContext(context.Background(), printer)
}

func (c *CUPSClient) RejectJobsContext(ctx context.Context, printer string) error {
	req := NewRequest(OperationCupsRejectJobs, 1)
	req.OperationAttributes[AttributePrinterURI] = c.getPrinterUri(printer)

	_, err := c.SendRequestContext(ctx, c.adapter.GetHttpUri("admin", ""), req, nil)
	return err
}

// AddPrinterToClass adds a printer to a class, if the class does not exists it will be crated
func (c *CUPSClient) AddPrinterToClass(class, printer string) error {
	return c.AddPrinterToClassContext(context.Background(), class, printer)
}

func (c *CUPSClient) AddPrinterToClassContext(ctx context.Context, class, printer string) error {
	attributes, err := c.GetPrinterAttributesContext(ctx, class, []string{AttributeMemberURIs})
	if err != nil && !IsNotExistsError(err) {
		return err
	}

	memberURIList := make([]string, 0)

	if !IsNotExistsError(err) {
		for _, member := range attributes[AttributeMemberURIs] {
			memberString := strings.Split(member.Value.(string), "/")
			printerName := memberString[len(memberString)-1]

			if printerName == printer {
				return nil
			}

			memberURIList = append(memberURIList, member.Value.(string))
		}
	}

	memberURIList = append(memberURIList, c.getPrinterUri(printer))

	req := NewRequest(OperationCupsAddModifyClass, 1)
	req.OperationAttributes[AttributePrinterURI] = c.getClassUri(class)
	req.PrinterAttributes[AttributeMemberURIs] = memberURIList

	_, err = c.SendRequestContext(ctx, c.adapter.GetHttpUri("admin", ""), req, nil)
	return err
}

// DeletePrinterFromClass removes a printer from a class, if a class has no more printer it will be deleted
func (c *CUPSClient) DeletePrinterFromClass(class, printer string) error {
	return c.DeletePrinterFromClassContext(context.Background(), class, printer)
}

func (c *CUPSClient) DeletePrinterFromClassContext(ctx context.Context, class, printer string) error {
	attributes, err := c.GetPrinterAttributesContext(ctx, class, []string{AttributeMemberURIs})
	if err != nil {
		return err
	}

	memberURIList := make([]string, 0)

	for _, member := range attributes[AttributeMemberURIs] {
		memberString := strings.Split(member.Value.(string), "/")
		printerName := memberString[len(memberString)-1]

		if printerName != printer {
			memberURIList = append(memberURIList, member.Value.(string))
		}
	}

	if len(memberURIList) == 0 {
		return c.DeleteClassContext(ctx, class)
	}

	req := NewRequest(OperationCupsAddModifyClass, 1)
	req.OperationAttributes[AttributePrinterURI] = c.getClassUri(class)
	req.PrinterAttributes[AttributeMemberURIs] = memberURIList

	_, err = c.SendRequestContext(ctx, c.adapter.GetHttpUri("admin", ""), req, nil)
	return err
}

// DeleteClass deletes a class
func (c *CUPSClient) DeleteClass(class string) error {
	return c.DeleteClassContext(context.Background(), class)
}

func (c *CUPSClient) DeleteClassContext(ctx context.Context, class string) error {
	req := NewRequest(OperationCupsDeleteClass, 1)
	req.OperationAttributes[AttributePrinterURI] = c.getClassUri(class)

	_, err := c.SendRequestContext(ctx, c.adapter.GetHttpUri("admin", ""), req, nil)
	return err
}

// CreatePrinter creates a new printer
func (c *CUPSClient) CreatePrinter(name, deviceURI, ppd string, shared bool, errorPolicy string, information, location string) error {
	return c.CreatePrinterContext(context.Background(), name, deviceURI, ppd, shared, errorPolicy, information, location)
}

func (c *CUPSClient) CreatePrinterContext(ctx context.Context, name, deviceURI, ppd string, shared bool, errorPolicy string, information, location string) error {
	req := NewRequest(OperationCupsAddModifyPrinter, 1)
	req.OperationAttributes[AttributePrinterURI] = c.getPrinterUri(name)
	req.OperationAttributes[AttributePPDName] = ppd
	req.OperationAttributes[AttributePrinterIsShared] = shared
	req.PrinterAttributes[AttributePrinterStateReasons] = "none"
	req.PrinterAttributes[AttributeDeviceURI] = deviceURI
	req.PrinterAttributes[AttributePrinterInfo] = information
	req.PrinterAttributes[AttributePrinterLocation] = location
	req.PrinterAttributes[AttributePrinterErrorPolicy] = errorPolicy

	_, err := c.SendRequestContext(ctx, c.adapter.GetHttpUri("admin", ""), req, nil)
	return err
}

// SetPrinterPPD sets the ppd for a printer
func (c *CUPSClient) SetPrinterPPD(printer, ppd string) error {
	return c.SetPrinterPPDContext(context.Background(), printer, ppd)
}

func (c *CUPSClient) SetPrinterPPDContext(ctx context.Context, printer, ppd string) error {
	req := NewRequest(OperationCupsAddModifyPrinter, 1)
	req.OperationAttributes[AttributePrinterURI] = c.getPrinterUri(printer)
	req.OperationAttributes[AttributePPDName] = ppd

	_, err := c.SendRequestContext(ctx, c.adapter.GetHttpUri("admin", ""), req, nil)
	return err
}

// SetPrinterDeviceURI sets the device uri for a printer
func (c *CUPSClient) SetPrinterDeviceURI(printer, deviceURI string) error {
	return c.SetPrinterDeviceURIContext(context.Background(), printer, deviceURI)
}

func (c *CUPSClient) SetPrinterDeviceURIContext(ctx context.Context, printer, deviceURI string) error {
	req := NewRequest(OperationCupsAddModifyPrinter, 1)
	req.OperationAttributes[AttributePrinterURI] = c.getPrinterUri(printer)
	req.PrinterAttributes[AttributeDeviceURI] = deviceURI

	_, err := c.SendRequestContext(ctx, c.adapter.GetHttpUri("admin", ""), req, nil)
	return err
}

// SetPrinterIsShared shares or unshares a printer in the network
func (c *CUPSClient) SetPrinterIsShared(printer string, shared bool) error {
	return c.SetPrinterIsSharedContext(context.Background(), printer, shared)
}

func (c *CUPSClient) SetPrinterIsSharedContext(ctx context.Context, printer string, shared bool) error {
	req := NewRequest(OperationCupsAddModifyPrinter, 1)
	req.OperationAttributes[AttributePrinterURI] = c.getPrinterUri(printer)
	req.OperationAttributes[AttributePrinterIsShared] = shared

	_, err := c.SendRequestContext(ctx, c.adapter.GetHttpUri("admin", ""), req, nil)
	return err
}

// SetPrinterErrorPolicy sets the error policy for a printer
func (c *CUPSClient) SetPrinterErrorPolicy(printer string, errorPolicy string) error {
	return c.SetPrinterErrorPolicyContext(context.Background(), printer, errorPolicy)
}

func (c *CUPSClient) SetPrinterErrorPolicyContext(ctx context.Context, printer string, errorPolicy string) error {
	req := NewRequest(OperationCupsAddModifyPrinter, 1)
	req.OperationAttributes[AttributePrinterURI] = c.getPrinterUri(printer)
	req.PrinterAttributes[AttributePrinterErrorPolicy] = errorPolicy

	_, err := c.SendRequestContext(ctx, c.adapter.GetHttpUri("admin", ""), req, nil)
	return err
}

// SetPrinterInformation sets general printer information
func (c *CUPSClient) SetPrinterInformation(printer, information string) error {
	return c.SetPrinterInformationContext(context.Background(), printer, information)
}

func (c *CUPSClient) SetPrinterInformationContext(ctx context.Context, printer, information string) error {
	req := NewRequest(OperationCupsAddModifyPrinter, 1)
	req.OperationAttributes[AttributePrinterURI] = c.getPrinterUri(printer)
	req.PrinterAttributes[AttributePrinterInfo] = information

	_, err := c.SendRequestContext(ctx, c.adapter.GetHttpUri("admin", ""), req, nil)
	return err
}

// SetPrinterLocation sets the printer location
func (c *CUPSClient) SetPrinterLocation(printer, location string) error {
	return c.SetPrinterLocationContext(context.Background(), printer, location)
}

func (c *CUPSClient) SetPrinterLocationContext(ctx context.Context, printer, location string) error {
	req := NewRequest(OperationCupsAddModifyPrinter, 1)
	req.OperationAttributes[AttributePrinterURI] = c.getPrinterUri(printer)
	req.PrinterAttributes[AttributePrinterLocation] = location

	_, err := c.SendRequestContext(ctx, c.adapter.GetHttpUri("admin", ""), req, nil)
	return err
}

// DeletePrinter deletes a printer
func (c *CUPSClient) DeletePrinter(printer string) error {
	return c.DeletePrinterContext(context.Background(), printer)
}

func (c *CUPSClient) DeletePrinterContext(ctx context.Context, printer string) error {
	req := NewRequest(OperationCupsDeletePrinter, 1)
	req.OperationAttributes[AttributePrinterURI] = c.getPrinterUri(printer)

	_, err := c.SendRequestContext(ctx, c.adapter.GetHttpUri("admin", ""), req, nil)
	return err
}

// GetPrinters returns a map of printer names and attributes
func (c *CUPSClient) GetPrinters(attributes []string) (map[string]Attributes, error) {
	return c.GetPrintersContext(context.Background(), attributes)
}

func (c *CUPSClient) GetPrintersContext(ctx context.Context, attributes []string) (map[string]Attributes, error) {
	req := NewRequest(OperationCupsGetPrinters, 1)

	if attributes == nil {
		req.OperationAttributes[AttributeRequestedAttributes] = DefaultPrinterAttributes
	} else {
		req.OperationAttributes[AttributeRequestedAttributes] = append(attributes, AttributePrinterName)
	}

	resp, err := c.SendRequestContext(ctx, c.adapter.GetHttpUri("", nil), req, nil)
	if err != nil {
		return nil, err
	}

	printerNameMap := make(map[string]Attributes)

	for _, printerAttributes := range resp.PrinterAttributes {
		printerNameMap[printerAttributes[AttributePrinterName][0].Value.(string)] = printerAttributes
	}

	return printerNameMap, nil
}

// GetClasses returns a map of class names and attributes
func (c *CUPSClient) GetClasses(attributes []string) (map[string]Attributes, error) {
	return c.GetClassesContext(context.Background(), attributes)
}

func (c *CUPSClient) GetClassesContext(ctx context.Context, attributes []string) (map[string]Attributes, error) {
	req := NewRequest(OperationCupsGetClasses, 1)

	if attributes == nil {
		req.OperationAttributes[AttributeRequestedAttributes] = DefaultClassAttributes
	} else {
		req.OperationAttributes[AttributeRequestedAttributes] = append(attributes, AttributePrinterName)
	}

	resp, err := c.SendRequestContext(ctx, c.adapter.GetHttpUri("", nil), req, nil)
	if err != nil {
		return nil, err
	}

	printerNameMap := make(map[string]Attributes)

	for _, printerAttributes := range resp.PrinterAttributes {
		printerNameMap[printerAttributes[AttributePrinterName][0].Value.(string)] = printerAttributes
	}

	return printerNameMap, nil
}

// PrintTestPage prints a test page of type application/vnd.cups-pdf-banner
func (c *CUPSClient) PrintTestPage(printer string) (int, error) {
	return c.PrintTestPageContext(context.Background(), printer)
}

func (c *CUPSClient) PrintTestPageContext(ctx context.Context, printer string) (int, error) {
	testPage := new(bytes.Buffer)
	testPage.WriteString("#PDF-BANNER\n")
	testPage.WriteString("Template default-testpage.pdf\n")
	testPage.WriteString("Show printer-name printer-info printer-location printer-make-and-model printer-driver-name")
	testPage.WriteString("printer-driver-version paper-size imageable-area job-id options time-at-creation")
	testPage.WriteString("time-at-processing\n\n")

	return c.PrintDocumentsContext(ctx, []Document{
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
