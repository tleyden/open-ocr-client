package ocrclient

import (
	"fmt"
	"strings"
	"testing"

	"github.com/couchbaselabs/go.assert"
	"github.com/tleyden/fakehttp"
)

func TestDecodeImageUrl(t *testing.T) {

	port := 8080
	fakeDecodedOcr := "fake ocr"
	sourceServer := fakehttp.NewHTTPServerWithPort(port)
	sourceServer.Start()
	headers := map[string]string{"Content-Type": "text/plain"}
	sourceServer.Response(200, headers, fakeDecodedOcr)

	openOcrUrl := fmt.Sprintf("http://localhost:%d", port)
	openOcrClient := NewHttpClient(openOcrUrl)
	attachmentUrl := "http://fake.io/a.png"
	ocrDecoded, err := openOcrClient.DecodeImageUrl(attachmentUrl, ENGINE_TESSERACT)
	assert.True(t, err == nil)
	assert.Equals(t, ocrDecoded, fakeDecodedOcr)

}

func TestDecodeImageReader(t *testing.T) {

	port := 8080
	fakeDecodedOcr := "fake ocr"
	sourceServer := fakehttp.NewHTTPServerWithPort(port)
	sourceServer.Start()
	headers := map[string]string{"Content-Type": "text/plain"}
	sourceServer.Response(200, headers, fakeDecodedOcr)

	openOcrUrl := fmt.Sprintf("http://localhost:%d", port)
	openOcrClient := NewHttpClient(openOcrUrl)

	testImageContents := "blah"
	reader := strings.NewReader(testImageContents)

	ocrDecoded, err := openOcrClient.DecodeImageReader(reader, ENGINE_TESSERACT)
	assert.True(t, err == nil)
	assert.Equals(t, ocrDecoded, fakeDecodedOcr)

}
