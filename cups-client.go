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

	resp, err := c.SendRequest(c.getHttpUri("", nil), req)
	if err != nil {
		return nil, err
	}

	printerNameMap := make(map[string]Attributes)

	for _, printerAttributes := range resp.Printers {
		printerNameMap[printerAttributes["device-uri"][0].Value.(string)] = printerAttributes
	}

	return printerNameMap, nil
}

func (c *CUPSClient) MoveJob(jobID int, destPrinter string) error {
	req := NewRequest(OperationCupsMoveJob, 1)
	req.OperationAttributes["job-uri"] = c.getJobUri(jobID)
	req.PrinterAttributes["job-printer-uri"] = c.getPrinterUri(destPrinter)

	_, err := c.SendRequest(c.getHttpUri("jobs", ""), req)
	return err
}

func (c *CUPSClient) MoveAllJob(srcPrinter, destPrinter string) error {
	req := NewRequest(OperationCupsMoveJob, 1)
	req.OperationAttributes["printer-uri"] = c.getPrinterUri(srcPrinter)
	req.PrinterAttributes["job-printer-uri"] = c.getPrinterUri(destPrinter)

	_, err := c.SendRequest(c.getHttpUri("jobs", ""), req)
	return err
}

func (c *CUPSClient) GetPPDs() (map[string]Attributes, error) {
	req := NewRequest(OperationCupsGetPpds, 1)

	resp, err := c.SendRequest(c.getHttpUri("", nil), req)
	if err != nil {
		return nil, err
	}

	ppdNameMap := make(map[string]Attributes)

	for _, printerAttributes := range resp.Printers {
		ppdNameMap[printerAttributes["ppd-name"][0].Value.(string)] = printerAttributes
	}

	return ppdNameMap, nil
}

func (c *CUPSClient) AcceptJobs(printer string) error {
	req := NewRequest(OperationCupsAcceptJobs, 1)
	req.OperationAttributes["printer-uri"] = c.getPrinterUri(printer)

	_, err := c.SendRequest(c.getHttpUri("admin", ""), req)
	return err
}

func (c *CUPSClient) RejectJobs(printer string) error {
	req := NewRequest(OperationCupsRejectJobs, 1)
	req.OperationAttributes["printer-uri"] = c.getPrinterUri(printer)

	_, err := c.SendRequest(c.getHttpUri("admin", ""), req)
	return err
}

func (c *CUPSClient) AddPrinterToClass(class, printer string) error {
	attributes, err := c.GetPrinterAttributes(class, []string{"member-uris"})
	if err != nil && !IsNotExistsError(err) {
		return err
	}

	memberURIList := make([]string, 0)

	if !IsNotExistsError(err) {
		for _, member := range attributes["member-uris"] {
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
	req.OperationAttributes["printer-uri"] = c.getClassUri(class)
	req.PrinterAttributes["member-uris"] = memberURIList

	_, err = c.SendRequest(c.getHttpUri("admin", ""), req)
	return err
}

func (c *CUPSClient) DeletePrinterFromClass(class, printer string) error {
	attributes, err := c.GetPrinterAttributes(class, []string{"member-uris"})
	if err != nil {
		return err
	}

	memberURIList := make([]string, 0)

	for _, member := range attributes["member-uris"] {
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
	req.OperationAttributes["printer-uri"] = c.getClassUri(class)
	req.PrinterAttributes["member-uris"] = memberURIList

	_, err = c.SendRequest(c.getHttpUri("admin", ""), req)
	return err
}

func (c *CUPSClient) DeleteClass(class string) error {
	req := NewRequest(OperationCupsDeleteClass, 1)
	req.OperationAttributes["printer-uri"] = c.getClassUri(class)

	_, err := c.SendRequest(c.getHttpUri("admin", ""), req)
	return err
}

func (c *CUPSClient) CreatePrinter(name, deviceURI, ppd string, shared bool, errorPolicy ErrorPolicy, information, location string) error {
	req := NewRequest(OperationCupsAddModifyPrinter, 1)
	req.OperationAttributes["printer-uri"] = c.getPrinterUri(name)
	req.OperationAttributes["ppd-name"] = ppd
	req.OperationAttributes["printer-is-shared"] = shared
	req.PrinterAttributes["printer-state-reason"] = "none"
	req.PrinterAttributes["device-uri"] = deviceURI
	req.PrinterAttributes["printer-info"] = information
	req.PrinterAttributes["printer-location"] = location
	req.PrinterAttributes["printer-error-policy"] = string(errorPolicy)

	_, err := c.SendRequest(c.getHttpUri("admin", ""), req)
	return err
}

func (c *CUPSClient) SetPrinterPPD(printer, ppd string) error {
	req := NewRequest(OperationCupsAddModifyPrinter, 1)
	req.OperationAttributes["printer-uri"] = c.getPrinterUri(printer)
	req.OperationAttributes["ppd-name"] = ppd

	_, err := c.SendRequest(c.getHttpUri("admin", ""), req)
	return err
}

func (c *CUPSClient) SetPrinterDeviceURI(printer, deviceURI string) error {
	req := NewRequest(OperationCupsAddModifyPrinter, 1)
	req.OperationAttributes["printer-uri"] = c.getPrinterUri(printer)
	req.PrinterAttributes["device-uri"] = deviceURI

	_, err := c.SendRequest(c.getHttpUri("admin", ""), req)
	return err
}

func (c *CUPSClient) SetPrinterIsShared(printer string, shared bool) error {
	req := NewRequest(OperationCupsAddModifyPrinter, 1)
	req.OperationAttributes["printer-uri"] = c.getPrinterUri(printer)
	req.OperationAttributes["printer-is-shared"] = shared

	_, err := c.SendRequest(c.getHttpUri("admin", ""), req)
	return err
}

func (c *CUPSClient) SetPrinterErrorPolicy(printer string, errorPolicy ErrorPolicy) error {
	req := NewRequest(OperationCupsAddModifyPrinter, 1)
	req.OperationAttributes["printer-uri"] = c.getPrinterUri(printer)
	req.PrinterAttributes["printer-error-policy"] = string(errorPolicy)

	_, err := c.SendRequest(c.getHttpUri("admin", ""), req)
	return err
}

func (c *CUPSClient) SetPrinterInformation(printer, information string) error {
	req := NewRequest(OperationCupsAddModifyPrinter, 1)
	req.OperationAttributes["printer-uri"] = c.getPrinterUri(printer)
	req.PrinterAttributes["printer-info"] = information

	_, err := c.SendRequest(c.getHttpUri("admin", ""), req)
	return err
}

func (c *CUPSClient) SetPrinterLocation(printer, location string) error {
	req := NewRequest(OperationCupsAddModifyPrinter, 1)
	req.OperationAttributes["printer-uri"] = c.getPrinterUri(printer)
	req.PrinterAttributes["printer-location"] = location

	_, err := c.SendRequest(c.getHttpUri("admin", ""), req)
	return err
}

func (c *CUPSClient) DeletePrinter(printer string) error {
	req := NewRequest(OperationCupsDeletePrinter, 1)
	req.OperationAttributes["printer-uri"] = c.getPrinterUri(printer)

	_, err := c.SendRequest(c.getHttpUri("admin", ""), req)
	return err
}

func (c *CUPSClient) GetPrinters(attributes []string) (map[string]Attributes, error) {
	req := NewRequest(OperationCupsGetPrinters, 1)

	if attributes == nil {
		req.OperationAttributes["requested-attributes"] = DefaultPrinterAttributes
	} else {
		req.OperationAttributes["requested-attributes"] = append(attributes, "printer-name")
	}

	resp, err := c.SendRequest(c.getHttpUri("", nil), req)
	if err != nil {
		return nil, err
	}

	printerNameMap := make(map[string]Attributes)

	for _, printerAttributes := range resp.Printers {
		printerNameMap[printerAttributes["printer-name"][0].Value.(string)] = printerAttributes
	}

	return printerNameMap, nil
}

func (c *CUPSClient) GetClasses(attributes []string) (map[string]Attributes, error) {
	req := NewRequest(OperationCupsGetClasses, 1)

	if attributes == nil {
		req.OperationAttributes["requested-attributes"] = DefaultClassAttributes
	} else {
		req.OperationAttributes["requested-attributes"] = append(attributes, "printer-name")
	}

	resp, err := c.SendRequest(c.getHttpUri("", nil), req)
	if err != nil {
		return nil, err
	}

	printerNameMap := make(map[string]Attributes)

	for _, printerAttributes := range resp.Printers {
		printerNameMap[printerAttributes["printer-name"][0].Value.(string)] = printerAttributes
	}

	return printerNameMap, nil
}
