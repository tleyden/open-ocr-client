package ocrclient

import (
	"bufio"
	"fmt"
	"os"
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

	ocrRequest := OcrRequest{
		ImgUrl:     "http://fake.io/a.png",
		EngineType: ENGINE_TESSERACT,
	}

	openOcrClient := NewHttpClient(openOcrUrl)

	ocrDecoded, err := openOcrClient.DecodeImageUrl(ocrRequest)
	assert.True(t, err == nil)
	assert.Equals(t, ocrDecoded, fakeDecodedOcr)

}

func TestDecodeImageReader(t *testing.T) {

	if os.Getenv("USER") != "tleyden" {
		t.Skip("skipping test; only meant for developer workstation")
	}

	// this integration test requires a real openocr http rest api
	// server up and running on port 8080

	port := 8080
	fakeDecodedOcr := "You can create local variables"

	openOcrUrl := fmt.Sprintf("http://localhost:%d", port)
	openOcrClient := NewHttpClient(openOcrUrl)

	file, err := os.Open("ocr_test.png")
	assert.True(t, err == nil)
	reader := bufio.NewReader(file)

	engineArgs := map[string]interface{}{
		"lang": "eng",
	}

	ocrRequest := OcrRequest{
		EngineType:    ENGINE_TESSERACT,
		InplaceDecode: true, // decode in place rather than using rabbitmq
		EngineArgs:    engineArgs,
	}

	ocrDecoded, err := openOcrClient.DecodeImageReader(reader, ocrRequest)
	assert.True(t, err == nil)
	assert.True(t, strings.HasPrefix(ocrDecoded, fakeDecodedOcr))

}
