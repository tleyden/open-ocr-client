
[![Build Status](https://drone.io/github.com/tleyden/open-ocr-client/status.png)](https://drone.io/github.com/tleyden/open-ocr-client/latest) [![GoDoc](http://godoc.org/github.com/tleyden/open-ocr-client?status.png)](http://godoc.org/github.com/tleyden/open-ocr-client) 

This is a client library for accessing an [OpenOCR](http://www.openocr.net) server.

## Quick Start

This assumes you have [Go installed](https://golang.org/doc/install).

```
$ go get -u -v github.com/tleyden/open-ocr-client
$ cd cli-main
$ go build
$ wget http://bit.ly/ocrimage
$ ./cli-main --openOcrUrl <your-openocr-url> --openOcrPort <your-openocr-port> upload --file ocrimage
```
