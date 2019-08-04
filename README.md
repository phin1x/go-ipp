# go-ipp

[![Version](https://img.shields.io/github/release-pre/phin1x/go-ipp.svg)](https://github.com/phin1x/go-ipp/releases/tag/v1.1.0)
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

## Examples

Print a file
```go
client := ipp.NewIPPClient("printserver", 631, "user", "password", true)
client.PrintFile("/path/to/file", "my-printer", map[string]interface{}{})
```

## Licence

Apache Licence Version 2.0

