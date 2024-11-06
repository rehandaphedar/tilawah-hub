package models

type Timing struct {
	Text     string    `json:"text,omitempty"`
	Segments []Segment `json:"segments"`
}

type Segment struct {
	Start float64 `json:"start"`
	End   float64 `json:"end"`
	Text  string  `json:"text,omitempty"`
	Score float64 `json:"score,omitempty"`
}
