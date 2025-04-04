package domain

type Dimensions struct {
	Height float64 `json:"height"`
	Width  float64 `json:"width"`
	Length float64 `json:"length"`
	Unit   string  `json:"unit"`
}
