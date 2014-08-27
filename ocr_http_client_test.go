package ocrclient

import (
	"bufio"
	"fmt"
	"os"
	"testing"

	"github.com/couchbaselabs/go.assert"
	"github.com/couchbaselabs/logg"
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

	ocrRequest := OcrRequest{
		ImgUrl:     "http://fake.io/a.png",
		EngineType: ENGINE_TESSERACT,
	}

	openOcrClient := NewHttpClient(openOcrUrl)

	ocrDecoded, err := openOcrClient.DecodeImageUrl(ocrRequest)
	assert.True(t, err == nil)
	assert.Equals(t, ocrDecoded, fakeDecodedOcr)

}

func IntegrationTestDecodeImageReader(t *testing.T) {

	// this integration test requires a real openocr http rest api
	// server up and running on port 8080

	port := 8080
	fakeDecodedOcr := "fake ocr"

	openOcrUrl := fmt.Sprintf("http://localhost:%d", port)
	openOcrClient := NewHttpClient(openOcrUrl)

	file, err := os.Open("ocr_test.png")
	assert.True(t, err == nil)
	reader := bufio.NewReader(file)

	ocrRequest := OcrRequest{
		EngineType:    ENGINE_TESSERACT,
		InplaceDecode: true, // decode in place rather than using rabbitmq
	}

	ocrDecoded, err := openOcrClient.DecodeImageReader(reader, ocrRequest)
	logg.Log("err: %v", err)
	assert.True(t, err == nil)
	assert.Equals(t, ocrDecoded, fakeDecodedOcr)

}
