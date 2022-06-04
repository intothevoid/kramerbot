package models

type PipupSimpleToast struct {
	Title    string `json:"title"`
	Message  string `json:"message"`
	Duration int    `json:"duration"`
	Position int    `json:"position"`
}

type PipupToast struct {
	Duration        int         `json:"duration"`
	Position        int         `json:"position"`
	Title           string      `json:"title"`
	TitleColor      string      `json:"titleColor"`
	TitleSize       int         `json:"titleSize"`
	Message         string      `json:"message"`
	MessageColor    string      `json:"messageColor"`
	MessageSize     int         `json:"messageSize"`
	BackgroundColor string      `json:"backgroundColor"`
	Media           *PipupMedia `json:"media"`
}

type PipupMedia struct {
	Image *PipupImage `json:"image"`
}

type PipupImage struct {
	Uri   string `json:"uri"`
	Width int    `json:"width"`
}
