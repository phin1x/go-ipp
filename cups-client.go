package ipp

import "strings"

type CUPSClient struct {
	*IPPClient
}

func NewCUPSClient(host string, port int, username, password string, useTLS bool) *CUPSClient {
	ippClient := NewIPPClient(host, port, username, password, useTLS)
	return &CUPSClient{ippClient}
}

func (c *CUPSClient) GetDevices() (map[string]Attributes, error) {
	req := NewRequest(OperationCupsGetDevices, 1)

	resp, err := c.SendRequest(c.getHttpUri("", nil), req, nil)
	if err != nil {
		return nil, err
	}

	printerNameMap := make(map[string]Attributes)

	for _, printerAttributes := range resp.PrinterAttributes {
		printerNameMap[printerAttributes[AttributeDeviceURI][0].Value.(string)] = printerAttributes
	}

	return printerNameMap, nil
}

func (c *CUPSClient) MoveJob(jobID int, destPrinter string) error {
	req := NewRequest(OperationCupsMoveJob, 1)
	req.OperationAttributes[AttributeJobURI] = c.getJobUri(jobID)
	req.PrinterAttributes[AttributeJobPrinterURI] = c.getPrinterUri(destPrinter)

	_, err := c.SendRequest(c.getHttpUri("jobs", ""), req, nil)
	return err
}

func (c *CUPSClient) MoveAllJob(srcPrinter, destPrinter string) error {
	req := NewRequest(OperationCupsMoveJob, 1)
	req.OperationAttributes[AttributePrinterURI] = c.getPrinterUri(srcPrinter)
	req.PrinterAttributes[AttributeJobPrinterURI] = c.getPrinterUri(destPrinter)

	_, err := c.SendRequest(c.getHttpUri("jobs", ""), req, nil)
	return err
}

func (c *CUPSClient) GetPPDs() (map[string]Attributes, error) {
	req := NewRequest(OperationCupsGetPPDs, 1)

	resp, err := c.SendRequest(c.getHttpUri("", nil), req, nil)
	if err != nil {
		return nil, err
	}

	ppdNameMap := make(map[string]Attributes)

	for _, printerAttributes := range resp.PrinterAttributes {
		ppdNameMap[printerAttributes[AttributePPDName][0].Value.(string)] = printerAttributes
	}

	return ppdNameMap, nil
}

func (c *CUPSClient) AcceptJobs(printer string) error {
	req := NewRequest(OperationCupsAcceptJobs, 1)
	req.OperationAttributes[AttributePrinterURI] = c.getPrinterUri(printer)

	_, err := c.SendRequest(c.getHttpUri("admin", ""), req, nil)
	return err
}

func (c *CUPSClient) RejectJobs(printer string) error {
	req := NewRequest(OperationCupsRejectJobs, 1)
	req.OperationAttributes[AttributePrinterURI] = c.getPrinterUri(printer)

	_, err := c.SendRequest(c.getHttpUri("admin", ""), req, nil)
	return err
}

func (c *CUPSClient) AddPrinterToClass(class, printer string) error {
	attributes, err := c.GetPrinterAttributes(class, []string{AttributeMemberURIs})
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

	_, err = c.SendRequest(c.getHttpUri("admin", ""), req, nil)
	return err
}

func (c *CUPSClient) DeletePrinterFromClass(class, printer string) error {
	attributes, err := c.GetPrinterAttributes(class, []string{AttributeMemberURIs})
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
		return c.DeleteClass(class)
	}

	req := NewRequest(OperationCupsAddModifyClass, 1)
	req.OperationAttributes[AttributePrinterURI] = c.getClassUri(class)
	req.PrinterAttributes[AttributeMemberURIs] = memberURIList

	_, err = c.SendRequest(c.getHttpUri("admin", ""), req, nil)
	return err
}

func (c *CUPSClient) DeleteClass(class string) error {
	req := NewRequest(OperationCupsDeleteClass, 1)
	req.OperationAttributes[AttributePrinterURI] = c.getClassUri(class)

	_, err := c.SendRequest(c.getHttpUri("admin", ""), req, nil)
	return err
}

func (c *CUPSClient) CreatePrinter(name, deviceURI, ppd string, shared bool, errorPolicy string, information, location string) error {
	req := NewRequest(OperationCupsAddModifyPrinter, 1)
	req.OperationAttributes[AttributePrinterURI] = c.getPrinterUri(name)
	req.OperationAttributes[AttributePPDName] = ppd
	req.OperationAttributes[AttributePrinterIsShared] = shared
	req.PrinterAttributes[AttributePrinterStateReason] = "none"
	req.PrinterAttributes[AttributeDeviceURI] = deviceURI
	req.PrinterAttributes[AttributePrinterInfo] = information
	req.PrinterAttributes[AttributePrinterLocation] = location
	req.PrinterAttributes[AttributePrinterErrorPolicy] = string(errorPolicy)

	_, err := c.SendRequest(c.getHttpUri("admin", ""), req, nil)
	return err
}

func (c *CUPSClient) SetPrinterPPD(printer, ppd string) error {
	req := NewRequest(OperationCupsAddModifyPrinter, 1)
	req.OperationAttributes[AttributePrinterURI] = c.getPrinterUri(printer)
	req.OperationAttributes[AttributePPDName] = ppd

	_, err := c.SendRequest(c.getHttpUri("admin", ""), req, nil)
	return err
}

func (c *CUPSClient) SetPrinterDeviceURI(printer, deviceURI string) error {
	req := NewRequest(OperationCupsAddModifyPrinter, 1)
	req.OperationAttributes[AttributePrinterURI] = c.getPrinterUri(printer)
	req.PrinterAttributes[AttributeDeviceURI] = deviceURI

	_, err := c.SendRequest(c.getHttpUri("admin", ""), req, nil)
	return err
}

func (c *CUPSClient) SetPrinterIsShared(printer string, shared bool) error {
	req := NewRequest(OperationCupsAddModifyPrinter, 1)
	req.OperationAttributes[AttributePrinterURI] = c.getPrinterUri(printer)
	req.OperationAttributes[AttributePrinterIsShared] = shared

	_, err := c.SendRequest(c.getHttpUri("admin", ""), req, nil)
	return err
}

func (c *CUPSClient) SetPrinterErrorPolicy(printer string, errorPolicy string) error {
	req := NewRequest(OperationCupsAddModifyPrinter, 1)
	req.OperationAttributes[AttributePrinterURI] = c.getPrinterUri(printer)
	req.PrinterAttributes[AttributePrinterErrorPolicy] = string(errorPolicy)

	_, err := c.SendRequest(c.getHttpUri("admin", ""), req, nil)
	return err
}

func (c *CUPSClient) SetPrinterInformation(printer, information string) error {
	req := NewRequest(OperationCupsAddModifyPrinter, 1)
	req.OperationAttributes[AttributePrinterURI] = c.getPrinterUri(printer)
	req.PrinterAttributes[AttributePrinterInfo] = information

	_, err := c.SendRequest(c.getHttpUri("admin", ""), req, nil)
	return err
}

func (c *CUPSClient) SetPrinterLocation(printer, location string) error {
	req := NewRequest(OperationCupsAddModifyPrinter, 1)
	req.OperationAttributes[AttributePrinterURI] = c.getPrinterUri(printer)
	req.PrinterAttributes[AttributePrinterLocation] = location

	_, err := c.SendRequest(c.getHttpUri("admin", ""), req, nil)
	return err
}

func (c *CUPSClient) DeletePrinter(printer string) error {
	req := NewRequest(OperationCupsDeletePrinter, 1)
	req.OperationAttributes[AttributePrinterURI] = c.getPrinterUri(printer)

	_, err := c.SendRequest(c.getHttpUri("admin", ""), req, nil)
	return err
}

func (c *CUPSClient) GetPrinters(attributes []string) (map[string]Attributes, error) {
	req := NewRequest(OperationCupsGetPrinters, 1)

	if attributes == nil {
		req.OperationAttributes[AttributeRequestedAttributes] = DefaultPrinterAttributes
	} else {
		req.OperationAttributes[AttributeRequestedAttributes] = append(attributes, AttributePrinterName)
	}

	resp, err := c.SendRequest(c.getHttpUri("", nil), req, nil)
	if err != nil {
		return nil, err
	}

	printerNameMap := make(map[string]Attributes)

	for _, printerAttributes := range resp.PrinterAttributes {
		printerNameMap[printerAttributes[AttributePrinterName][0].Value.(string)] = printerAttributes
	}

	return printerNameMap, nil
}

func (c *CUPSClient) GetClasses(attributes []string) (map[string]Attributes, error) {
	req := NewRequest(OperationCupsGetClasses, 1)

	if attributes == nil {
		req.OperationAttributes[AttributeRequestedAttributes] = DefaultClassAttributes
	} else {
		req.OperationAttributes[AttributeRequestedAttributes] = append(attributes, AttributePrinterName)
	}

	resp, err := c.SendRequest(c.getHttpUri("", nil), req, nil)
	if err != nil {
		return nil, err
	}

	printerNameMap := make(map[string]Attributes)

	for _, printerAttributes := range resp.PrinterAttributes {
		printerNameMap[printerAttributes[AttributePrinterName][0].Value.(string)] = printerAttributes
	}

	return printerNameMap, nil
}
