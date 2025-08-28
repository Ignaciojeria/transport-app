package request

type GeneratePreSignedUrlRequest struct {
	Lpn              string       `json:"lpn" example:"1234567890"`
	OrderReferenceID string       `json:"orderReferenceID" example:"1234567890"`
	UploadUrls       []UploadUrls `json:"uploadUrls"`
}
type UploadUrls struct {
	Type        string `json:"type" example:"proof_of_delivery"`
	Count       int    `json:"count" example:"1"`
	ContentType string `json:"contentType" example:"image/webp"`
}
