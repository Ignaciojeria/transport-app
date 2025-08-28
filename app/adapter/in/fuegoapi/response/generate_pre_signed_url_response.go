package response

type GeneratePreSignedUrlResponse struct {
	UploadUrls []struct {
		UploadURL        string `json:"uploadUrl" example:"https://gateway.storjshare.io/bucket/file.jpg?signature=upload"`
		DownloadURL      string `json:"downloadUrl" example:"https://gateway.storjshare.io/bucket/file.jpg?signature=download"`
		Method           string `json:"method" example:"PUT"`
		ContentType      string `json:"contentType" example:"image/webp"`
		UploadExpiresAt  string `json:"uploadExpiresAt" example:"2025-08-30T12:00:00Z"`
		DownloadExpiresAt string `json:"downloadExpiresAt" example:"2035-08-28T12:00:00Z"`
		FileKey          string `json:"fileKey" example:"orders/123/lpn/456/delivery_1_abc123.webp"`
	} `json:"uploadUrls"`
}
