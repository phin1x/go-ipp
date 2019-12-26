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

	for _, printerAttributes := range resp.Printers {
		printerNameMap[printerAttributes[PrinterAttributeDeviceURI][0].Value.(string)] = printerAttributes
	}

	return printerNameMap, nil
}

func (c *CUPSClient) MoveJob(jobID int, destPrinter string) error {
	req := NewRequest(OperationCupsMoveJob, 1)
	req.OperationAttributes[OperationAttributeJobURI] = c.getJobUri(jobID)
	req.PrinterAttributes[PrinterAttributeJobPrinterURI] = c.getPrinterUri(destPrinter)

	_, err := c.SendRequest(c.getHttpUri("jobs", ""), req, nil)
	return err
}

func (c *CUPSClient) MoveAllJob(srcPrinter, destPrinter string) error {
	req := NewRequest(OperationCupsMoveJob, 1)
	req.OperationAttributes[OperationAttributePrinterURI] = c.getPrinterUri(srcPrinter)
	req.PrinterAttributes[PrinterAttributeJobPrinterURI] = c.getPrinterUri(destPrinter)

	_, err := c.SendRequest(c.getHttpUri("jobs", ""), req, nil)
	return err
}

func (c *CUPSClient) GetPPDs() (map[string]Attributes, error) {
	req := NewRequest(OperationCupsGetPpds, 1)

	resp, err := c.SendRequest(c.getHttpUri("", nil), req, nil)
	if err != nil {
		return nil, err
	}

	ppdNameMap := make(map[string]Attributes)

	for _, printerAttributes := range resp.Printers {
		ppdNameMap[printerAttributes[OperationAttributePPDName][0].Value.(string)] = printerAttributes
	}

	return ppdNameMap, nil
}

func (c *CUPSClient) AcceptJobs(printer string) error {
	req := NewRequest(OperationCupsAcceptJobs, 1)
	req.OperationAttributes[OperationAttributePrinterURI] = c.getPrinterUri(printer)

	_, err := c.SendRequest(c.getHttpUri("admin", ""), req, nil)
	return err
}

func (c *CUPSClient) RejectJobs(printer string) error {
	req := NewRequest(OperationCupsRejectJobs, 1)
	req.OperationAttributes[OperationAttributePrinterURI] = c.getPrinterUri(printer)

	_, err := c.SendRequest(c.getHttpUri("admin", ""), req, nil)
	return err
}

func (c *CUPSClient) AddPrinterToClass(class, printer string) error {
	attributes, err := c.GetPrinterAttributes(class, []string{PrinterAttributeMemberURIs})
	if err != nil && !IsNotExistsError(err) {
		return err
	}

	memberURIList := make([]string, 0)

	if !IsNotExistsError(err) {
		for _, member := range attributes[PrinterAttributeMemberURIs] {
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
	req.OperationAttributes[OperationAttributePrinterURI] = c.getClassUri(class)
	req.PrinterAttributes[PrinterAttributeMemberURIs] = memberURIList

	_, err = c.SendRequest(c.getHttpUri("admin", ""), req, nil)
	return err
}

func (c *CUPSClient) DeletePrinterFromClass(class, printer string) error {
	attributes, err := c.GetPrinterAttributes(class, []string{PrinterAttributeMemberURIs})
	if err != nil {
		return err
	}

	memberURIList := make([]string, 0)

	for _, member := range attributes[PrinterAttributeMemberURIs] {
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
	req.OperationAttributes[OperationAttributePrinterURI] = c.getClassUri(class)
	req.PrinterAttributes[PrinterAttributeMemberURIs] = memberURIList

	_, err = c.SendRequest(c.getHttpUri("admin", ""), req, nil)
	return err
}

func (c *CUPSClient) DeleteClass(class string) error {
	req := NewRequest(OperationCupsDeleteClass, 1)
	req.OperationAttributes[OperationAttributePrinterURI] = c.getClassUri(class)

	_, err := c.SendRequest(c.getHttpUri("admin", ""), req, nil)
	return err
}

func (c *CUPSClient) CreatePrinter(name, deviceURI, ppd string, shared bool, errorPolicy ErrorPolicy, information, location string) error {
	req := NewRequest(OperationCupsAddModifyPrinter, 1)
	req.OperationAttributes[OperationAttributePrinterURI] = c.getPrinterUri(name)
	req.OperationAttributes[OperationAttributePPDName] = ppd
	req.OperationAttributes[OperationAttributePrinterIsShared] = shared
	req.PrinterAttributes[PrinterAttributePrinterStateReason] = "none"
	req.PrinterAttributes[PrinterAttributeDeviceURI] = deviceURI
	req.PrinterAttributes[PrinterAttributePrinterInfo] = information
	req.PrinterAttributes[PrinterAttributePrinterLocation] = location
	req.PrinterAttributes[PrinterAttributePrinterErrorPolicy] = string(errorPolicy)

	_, err := c.SendRequest(c.getHttpUri("admin", ""), req, nil)
	return err
}

func (c *CUPSClient) SetPrinterPPD(printer, ppd string) error {
	req := NewRequest(OperationCupsAddModifyPrinter, 1)
	req.OperationAttributes[OperationAttributePrinterURI] = c.getPrinterUri(printer)
	req.OperationAttributes[OperationAttributePPDName] = ppd

	_, err := c.SendRequest(c.getHttpUri("admin", ""), req, nil)
	return err
}

func (c *CUPSClient) SetPrinterDeviceURI(printer, deviceURI string) error {
	req := NewRequest(OperationCupsAddModifyPrinter, 1)
	req.OperationAttributes[OperationAttributePrinterURI] = c.getPrinterUri(printer)
	req.PrinterAttributes[PrinterAttributeDeviceURI] = deviceURI

	_, err := c.SendRequest(c.getHttpUri("admin", ""), req, nil)
	return err
}

func (c *CUPSClient) SetPrinterIsShared(printer string, shared bool) error {
	req := NewRequest(OperationCupsAddModifyPrinter, 1)
	req.OperationAttributes[OperationAttributePrinterURI] = c.getPrinterUri(printer)
	req.OperationAttributes[OperationAttributePrinterIsShared] = shared

	_, err := c.SendRequest(c.getHttpUri("admin", ""), req, nil)
	return err
}

func (c *CUPSClient) SetPrinterErrorPolicy(printer string, errorPolicy ErrorPolicy) error {
	req := NewRequest(OperationCupsAddModifyPrinter, 1)
	req.OperationAttributes[OperationAttributePrinterURI] = c.getPrinterUri(printer)
	req.PrinterAttributes[PrinterAttributePrinterErrorPolicy] = string(errorPolicy)

	_, err := c.SendRequest(c.getHttpUri("admin", ""), req, nil)
	return err
}

func (c *CUPSClient) SetPrinterInformation(printer, information string) error {
	req := NewRequest(OperationCupsAddModifyPrinter, 1)
	req.OperationAttributes[OperationAttributePrinterURI] = c.getPrinterUri(printer)
	req.PrinterAttributes[PrinterAttributePrinterInfo] = information

	_, err := c.SendRequest(c.getHttpUri("admin", ""), req, nil)
	return err
}

func (c *CUPSClient) SetPrinterLocation(printer, location string) error {
	req := NewRequest(OperationCupsAddModifyPrinter, 1)
	req.OperationAttributes[OperationAttributePrinterURI] = c.getPrinterUri(printer)
	req.PrinterAttributes[PrinterAttributePrinterLocation] = location

	_, err := c.SendRequest(c.getHttpUri("admin", ""), req, nil)
	return err
}

func (c *CUPSClient) DeletePrinter(printer string) error {
	req := NewRequest(OperationCupsDeletePrinter, 1)
	req.OperationAttributes[OperationAttributePrinterURI] = c.getPrinterUri(printer)

	_, err := c.SendRequest(c.getHttpUri("admin", ""), req, nil)
	return err
}

func (c *CUPSClient) GetPrinters(attributes []string) (map[string]Attributes, error) {
	req := NewRequest(OperationCupsGetPrinters, 1)

	if attributes == nil {
		req.OperationAttributes[OperationAttributeRequestedAttributes] = DefaultPrinterAttributes
	} else {
		req.OperationAttributes[OperationAttributeRequestedAttributes] = append(attributes, PrinterAttributePrinterName)
	}

	resp, err := c.SendRequest(c.getHttpUri("", nil), req, nil)
	if err != nil {
		return nil, err
	}

	printerNameMap := make(map[string]Attributes)

	for _, printerAttributes := range resp.Printers {
		printerNameMap[printerAttributes[PrinterAttributePrinterName][0].Value.(string)] = printerAttributes
	}

	return printerNameMap, nil
}

func (c *CUPSClient) GetClasses(attributes []string) (map[string]Attributes, error) {
	req := NewRequest(OperationCupsGetClasses, 1)

	if attributes == nil {
		req.OperationAttributes[OperationAttributeRequestedAttributes] = DefaultClassAttributes
	} else {
		req.OperationAttributes[OperationAttributeRequestedAttributes] = append(attributes, PrinterAttributePrinterName)
	}

	resp, err := c.SendRequest(c.getHttpUri("", nil), req, nil)
	if err != nil {
		return nil, err
	}

	printerNameMap := make(map[string]Attributes)

	for _, printerAttributes := range resp.Printers {
		printerNameMap[printerAttributes[PrinterAttributePrinterName][0].Value.(string)] = printerAttributes
	}

	return printerNameMap, nil
}
