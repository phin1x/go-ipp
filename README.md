# go-ipp

## Go Get

To get the package, execute:
```
go get -h github.com/phin1x/go-ipp
```

## Features

* basic ipp 2.0 compatible Client
* extended client for cups server
* create custom ipp requests
* parse ipp response and files 

## Examples

Print a file
```go
client := ipp.NewIPPClient("printserver", 631, "user", "password", true)
client.PrintFile("/path/to/file", "my-printer", 1, ipp.DefaultJobPriority)
```

## Licence

Apache Licence Version 2.0

