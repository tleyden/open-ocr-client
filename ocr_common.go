package ocrclient

// These are unforunately duplicated with code in open-ocr core package

type OcrEngineType int

const (
	ENGINE_TESSERACT = OcrEngineType(iota)
	ENGINE_MOCK
)

type OcrRequest struct {
	ImgUrl     string        `json:"img_url"`
	EngineType OcrEngineType `json:"engine"`

	// decode ocr in http handler rather than putting in queue
	InplaceDecode bool `json:"inplace_decode"`
}
