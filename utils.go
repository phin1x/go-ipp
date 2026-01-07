package ipp

import (
	"fmt"
	"os"
	"path"
)

// ParseControlFile reads and decodes a cups control file into a response
func ParseControlFile(jobID int, spoolDirectory string) (*Response, error) {
	if spoolDirectory == "" {
		spoolDirectory = "/var/spool/cups"
	}

	controlFilePath := path.Join(spoolDirectory, fmt.Sprintf("c%d", jobID))

	if _, err := os.Stat(controlFilePath); err != nil {
		return nil, err
	}

	controlFile, err := os.Open(controlFilePath)
	if err != nil {
		return nil, err
	}
	defer controlFile.Close()

	return NewResponseDecoder(controlFile).Decode(nil)
}
