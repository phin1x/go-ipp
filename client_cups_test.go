package ipp

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMapPrinterAttributesByName(t *testing.T) {
	groups := []Attributes{
		{AttributePrinterName: []Attribute{{Tag: TagName, Name: AttributePrinterName, Value: "printer-a"}}},
		{AttributePrinterName: []Attribute{{Tag: TagName, Name: AttributePrinterName, Value: "printer-b"}}},
	}

	result, err := mapPrinterAttributesByName(groups)
	assert.NoError(t, err)
	assert.Len(t, result, 2)
	assert.Contains(t, result, "printer-a")
	assert.Contains(t, result, "printer-b")
}

func TestMapPrinterAttributesByNameMissingName(t *testing.T) {
	groups := []Attributes{
		{AttributePrinterInfo: []Attribute{{Tag: TagText, Name: AttributePrinterInfo, Value: "no name here"}}},
	}

	_, err := mapPrinterAttributesByName(groups)
	assert.Error(t, err)
}

func TestMapPrinterAttributesByNameWrongType(t *testing.T) {
	groups := []Attributes{
		{AttributePrinterName: []Attribute{{Tag: TagInteger, Name: AttributePrinterName, Value: int32(42)}}},
	}

	_, err := mapPrinterAttributesByName(groups)
	assert.Error(t, err)
}
