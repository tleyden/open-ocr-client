package main

import (
	"fmt"
	"math/rand"

	"os"

	"github.com/alecthomas/kingpin"
	"github.com/couchbaselabs/logg"
	"github.com/tleyden/open-ocr-client"
)

var (
	app           = kingpin.New("open-ocr-client", "A command-line chat application.")
	stress        = app.Command("stress", "Do a stress test")
	numIterations = stress.Arg("numIterations", "how many OCR jobs should each goroutine create?").Default("5").Int()
	numGoroutines = stress.Arg("numGoroutines", "how many goroutines should be launched?").Default("1").Int()

	numTestImages = 5
)

func init() {
	logg.LogKeys["CLI"] = true
}

func main() {
	switch kingpin.MustParse(app.Parse(os.Args[1:])) {
	case "stress":
		logg.LogTo("CLI", "do stress test")
		stressTestLauncher()
	default:
		logg.LogTo("CLI", "oops, nothing to do")
	}
}

func imageUrls() []string {
	imageUrlBase := "http://tleyden-misc.s3.amazonaws.com/ocr-test-data"

	imageUrls := []string{}
	for i := 0; i < numTestImages; i++ {
		imageUrl := fmt.Sprintf("%s/%d.png", imageUrlBase, i)
		imageUrls = append(imageUrls, imageUrl)
	}
	return imageUrls
}

func stressTestLauncher() {
	doneChannel := make(chan bool)
	for i := 0; i < *numGoroutines; i++ {
		go stressTest(doneChannel)
	}

	for i := 0; i < *numGoroutines; i++ {
		<-doneChannel
		logg.LogTo("CLI", "goroutine finished: %d", i)
	}

}

func stressTest(doneChannel chan<- bool) {

	imageUrls := imageUrls()
	logg.LogTo("CLI", "imageUrls: %v", imageUrls)
	logg.LogTo("CLI", "numIterations: %v", *numIterations)

	openOcrUrl := "http://api.openocr.net"
	client := ocrclient.NewHttpClient(openOcrUrl)

	for i := 0; i < *numIterations; i++ {
		index := randomIntInRange(0, numTestImages)
		imageUrl := imageUrls[index]
		logg.LogTo("CLI", "OCR decoding: %v.  index: %d", imageUrl, index)
		ocrDecoded, err := client.DecodeImageUrl(imageUrl, ocrclient.ENGINE_TESSERACT)
		if err != nil {
			logg.LogPanic("Error decoding image: %v", err)
		}
		logg.LogTo("CLI", "OCR decoded: %v", ocrDecoded)
	}

	doneChannel <- true

}

// return a random number between min and max - 1
// eg, if you call it with 0,1 it will always return 0
// if you call it between 0,2 it will return 0 or 1
func randomIntInRange(min, max int) int {
	if min == max {
		return min
	}
	return rand.Intn(max-min) + min
}
