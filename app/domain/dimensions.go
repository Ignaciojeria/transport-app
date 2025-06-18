package domain

type Dimensions struct {
	Height int64  `json:"height"`
	Width  int64  `json:"width"`
	Length int64  `json:"length"`
	Unit   string `json:"unit"`
}
