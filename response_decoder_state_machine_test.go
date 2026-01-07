package ipp

import (
	"bytes"
	"encoding/hex"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

/**
 * Issue #??? `unsupported` fields parsing.
 * The following response from a Samsung M288x series
 * printer trigger the issue above and is included here
 * as a test case:
0000   02 00 00 01 00 00 00 01 01 47 00 12 61 74 74 72   .........G..attr
0010   69 62 75 74 65 73 2d 63 68 61 72 73 65 74 00 05   ibutes-charset..
0020   75 74 66 2d 38 48 00 1b 61 74 74 72 69 62 75 74   utf-8H..attribut
0030   65 73 2d 6e 61 74 75 72 61 6c 2d 6c 61 6e 67 75   es-natural-langu
0040   61 67 65 00 05 65 6e 2d 75 73 45 00 0b 70 72 69   age..en-usE..pri
0050   6e 74 65 72 2d 75 72 69 00 28 69 70 70 3a 2f 2f   nter-uri.(ipp://
0060   6c 6f 63 61 6c 68 6f 73 74 2f 70 72 69 6e 74 65   localhost/printe
0070   72 73 2f 53 45 43 33 30 43 44 41 37 41 38 30 43   rs/SEC30CDA7A80C
0080   30 32 05 10 00 0c 70 72 69 6e 74 65 72 2d 74 79   02....printer-ty
              ^ start of unsupported group
	         ^ value tag "unsupported" (0x10)
		    ^ length of the group (0x000c)
	         ^ decoder is falsely interpreting 0x1000 as the length.
0090   70 65 00 0c 70 72 69 6e 74 65 72 2d 74 79 70 65   pe..printer-type
00a0   10 00 0a 64 65 76 69 63 65 2d 75 72 69 00 0a 64   ...device-uri..d
00b0   65 76 69 63 65 2d 75 72 69 10 00 11 70 72 69 6e   evice-uri...prin
00c0   74 65 72 2d 69 73 2d 73 68 61 72 65 64 00 11 70   ter-is-shared..p
00d0   72 69 6e 74 65 72 2d 69 73 2d 73 68 61 72 65 64   rinter-is-shared
00e0   04 42 00 0c 70 72 69 6e 74 65 72 2d 6e 61 6d 65   .B..printer-name
00f0   00 0f 53 45 43 33 30 43 44 41 37 41 38 30 43 30   ..SEC30CDA7A80C0
0100   32 41 00 10 70 72 69 6e 74 65 72 2d 6c 6f 63 61   2A..printer-loca
0110   74 69 6f 6e 00 00 41 00 0c 70 72 69 6e 74 65 72   tion..A..printer
0120   2d 69 6e 66 6f 00 26 53 61 6d 73 75 6e 67 20 4d   -info.&Samsung M
0130   32 38 38 78 20 53 65 72 69 65 73 20 28 53 45 43   288x Series (SEC
0140   33 30 43 44 41 37 41 38 30 43 30 32 29 41 00 16   30CDA7A80C02)A..
0150   70 72 69 6e 74 65 72 2d 6d 61 6b 65 2d 61 6e 64   printer-make-and
0160   2d 6d 6f 64 65 6c 00 14 53 61 6d 73 75 6e 67 20   -model..Samsung
0170   4d 32 38 38 78 20 53 65 72 69 65 73 23 00 0d 70   M288x Series#..p
0180   72 69 6e 74 65 72 2d 73 74 61 74 65 00 04 00 00   rinter-state....
0190   00 03 41 00 15 70 72 69 6e 74 65 72 2d 73 74 61   ..A..printer-sta
01a0   74 65 2d 6d 65 73 73 61 67 65 00 19 50 72 69 6e   te-message..Prin
01b0   74 65 72 20 69 73 20 72 65 61 64 79 20 74 6f 20   ter is ready to
01c0   70 72 69 6e 74 44 00 15 70 72 69 6e 74 65 72 2d   printD..printer-
01d0   73 74 61 74 65 2d 72 65 61 73 6f 6e 73 00 04 6e   state-reasons..n
01e0   6f 6e 65 45 00 15 70 72 69 6e 74 65 72 2d 75 72   oneE..printer-ur
01f0   69 2d 73 75 70 70 6f 72 74 65 64 00 1f 69 70 70   i-supported..ipp
0200   3a 2f 2f 31 39 32 2e 31 36 38 2e 31 37 38 2e 31   ://192.168.178.1
0210   36 37 2f 69 70 70 2f 70 72 69 6e 74 45 00 00 00   67/ipp/printE...
0220   26 69 70 70 3a 2f 2f 53 45 43 33 30 43 44 41 37   &ipp://SEC30CDA7
0230   41 38 30 43 30 32 2e 6c 6f 63 61 6c 2e 2f 69 70   A80C02.local./ip
0240   70 2f 70 72 69 6e 74 03 01 02 03                  p/print.
                               ^ start of data
*/

var unsupportedResponse = `
02 00 00 01 00 00 00 07 01 47 00 12 61 74 74 72
69 62 75 74 65 73 2d 63 68 61 72 73 65 74 00 05    
75 74 66 2d 38 48 00 1b 61 74 74 72 69 62 75 74
65 73 2d 6e 61 74 75 72 61 6c 2d 6c 61 6e 67 75
61 67 65 00 05 65 6e 2d 75 73 45 00 0b 70 72 69
6e 74 65 72 2d 75 72 69 00 28 69 70 70 3a 2f 2f
6c 6f 63 61 6c 68 6f 73 74 2f 70 72 69 6e 74 65
72 73 2f 53 45 43 33 30 43 44 41 37 41 38 30 43
30 32 05 10 00 0c 70 72 69 6e 74 65 72 2d 74 79
70 65 00 0c 70 72 69 6e 74 65 72 2d 74 79 70 65
10 00 0a 64 65 76 69 63 65 2d 75 72 69 00 0a 64
65 76 69 63 65 2d 75 72 69 10 00 11 70 72 69 6e
74 65 72 2d 69 73 2d 73 68 61 72 65 64 00 11 70
72 69 6e 74 65 72 2d 69 73 2d 73 68 61 72 65 64
04 42 00 0c 70 72 69 6e 74 65 72 2d 6e 61 6d 65
00 0f 53 45 43 33 30 43 44 41 37 41 38 30 43 30
32 41 00 10 70 72 69 6e 74 65 72 2d 6c 6f 63 61
74 69 6f 6e 00 00 41 00 0c 70 72 69 6e 74 65 72
2d 69 6e 66 6f 00 26 53 61 6d 73 75 6e 67 20 4d
32 38 38 78 20 53 65 72 69 65 73 20 28 53 45 43
33 30 43 44 41 37 41 38 30 43 30 32 29 41 00 16
70 72 69 6e 74 65 72 2d 6d 61 6b 65 2d 61 6e 64
2d 6d 6f 64 65 6c 00 14 53 61 6d 73 75 6e 67 20
4d 32 38 38 78 20 53 65 72 69 65 73 23 00 0d 70
72 69 6e 74 65 72 2d 73 74 61 74 65 00 04 00 00
00 03 41 00 15 70 72 69 6e 74 65 72 2d 73 74 61
74 65 2d 6d 65 73 73 61 67 65 00 19 50 72 69 6e
74 65 72 20 69 73 20 72 65 61 64 79 20 74 6f 20
70 72 69 6e 74 44 00 15 70 72 69 6e 74 65 72 2d
73 74 61 74 65 2d 72 65 61 73 6f 6e 73 00 04 6e
6f 6e 65 45 00 15 70 72 69 6e 74 65 72 2d 75 72
69 2d 73 75 70 70 6f 72 74 65 64 00 1f 69 70 70
3a 2f 2f 31 39 32 2e 31 36 38 2e 31 37 38 2e 31
36 37 2f 69 70 70 2f 70 72 69 6e 74 45 00 00 00
26 69 70 70 3a 2f 2f 53 45 43 33 30 43 44 41 37
41 38 30 43 30 32 2e 6c 6f 63 61 6c 2e 2f 69 70
70 2f 70 72 69 6e 74 03 01 02 03
`

func hex2b(hx string) []byte {
	hx = strings.ReplaceAll(hx, " ", "")
	hx = strings.ReplaceAll(hx, "\n", "")
	data, _ := hex.DecodeString(hx)
	return data
}

func hex2reader(hex string) *bytes.Reader {
	return bytes.NewReader(hex2b(hex))
}

func TestResponseDecode(t *testing.T) {
	reader := hex2reader(unsupportedResponse)

	stateMachine := newResponseStateMachine()
	response, err := stateMachine.Decode(reader)

	if err != nil {
		t.Errorf("Expected nil, got %v", err)
	}

	if response.ProtocolVersionMajor != 2 {
		t.Errorf("Expected 2, got %v", response.ProtocolVersionMajor)
	}
	if response.StatusCode != 1 {
		t.Errorf("Expected Status:1, got %v", response.StatusCode)
	}
	if response.RequestId != 7 {
		t.Errorf("Expected RequestId:7, got %v", response.RequestId)
	}

	assert.Equal(t, 3, len(response.OperationAttributes))
	assert.Equal(t, 1, len(response.OperationAttributes[AttributeCharset]))

	charset := response.OperationAttributes[AttributeCharset][0]
	assert.Equal(t, "utf-8", charset.Value)

	unsupportedAttributes := response.UnsupportedAttributes
	assert.Equal(t, 3, len(unsupportedAttributes))

	printerType := unsupportedAttributes[AttributePrinterType][0]
	assert.Equal(t, "printer-type", printerType.Value)

	printAttributes := response.PrinterAttributes[0]
	assert.Equal(t, 8, len(printAttributes))

	printerStateMsg := printAttributes[AttributePrinterStateMessage][0]
	assert.Equal(t, "Printer is ready to print", printerStateMsg.Value)

	if response.Data != nil {
		should := []byte{0x01, 0x02, 0x03}
		assert.Equal(t, should, response.Data)
	} else {
		t.Errorf("Expected Data, got nil")
	}
}
