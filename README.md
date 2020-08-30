# go-ipp

[![Version](https://img.shields.io/github/release-pre/phin1x/go-ipp.svg)](https://github.com/phin1x/go-ipp/releases/tag/v1.4.4)
[![Documentation](https://img.shields.io/badge/go.dev-reference-007d9c?logo=go&logoColor=white)](https://pkg.go.dev/github.com/phin1x/go-ipp?tab=doc)
[![Go Report Card](https://goreportcard.com/badge/github.com/phin1x/go-ipp)](https://goreportcard.com/report/github.com/phin1x/go-ipp)
[![Licence](https://img.shields.io/github/license/phin1x/go-ipp.svg)](https://github.com/phin1x/go-ipp/blob/master/LICENSE)

## Go Get

To get the package, execute:
```
go get -u github.com/phin1x/go-ipp
```

## Features

* basic ipp 2.0 compatible Client
* extended client for cups server
* create custom ipp requests
* parse ipp responses and ipp control files

## Example

Print a file with the ipp client
```go
package main

import "github.com/phin1x/go-ipp"

func main() {
    // create a new ipp client
    client := ipp.NewIPPClient("printserver", 631, "user", "password", true)
    // print file
    client.PrintFile("/path/to/file", "my-printer", map[string]interface{}{})
}
```

Craft and send a custom request
```go

package main

import "github.com/phin1x/go-ipp"

func main() {             
    // define a ipp request
    req := ipp.NewRequest(OperationGetJobs, 1)
    req.OperationAttributes[ipp.AttributeWhichJobs] = "completed"
    req.OperationAttributes[ipp.AttributeMyJobs] = myJobs
    req.OperationAttributes[ipp.AttributeFirstJobID] = 42
    req.OperationAttributes[ipp.AttributeRequestingUserName] = "fabian"
    
    // encode request to bytes
    payload, err := req.Encode()
    if err != nil {
        panic(err)
    }
    
    // send ipp request to remote server via http
    httpReq, err := http.NewRequest("POST", "http://my-print-server:631/printers/my-printer", bytes.NewBuffer(payload))
    if err != nil {
        panic(err)
    }
    
    // set ipp headers
    httpReq.Header.Set("Content-Length", len(payload))
    httpReq.Header.Set("Content-Type", ipp.ContentTypeIPP)
    
    httpClient := &http.Client()
    httpResp, err := httpClient.Do(httpReq)
    if err != nil {
        panic(err)
    }
    defer httpResp.Body.Close()
    
    // response must be 200 for a successful operation
    // other possible http codes are: 
    // - 500 -> server error
    // - 426 -> sever requests a encrypted connection
    // - 401 -> forbidden -> need authorization header or user is not permitted
    if httpResp.StatusCode != 200 {
        panic("non 200 response from server")
    }
    
    // decode ipp response
    resp, err := ipp.NewResponseDecoder(httpResp.Body).Decode(nil)
    if err != nil {
        panic(err)
    }
    
    // check if the response status is "ok"
    if resp.StatusCode == ipp.StatusOk {
        panic(resp.StatusCode)
    }
    
    // do something with the returned data
    for _, job := resp.JobAttributes {
        // ...
    }
}
```

## Licence

Apache Licence Version 2.0

