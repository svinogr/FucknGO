package ui

type UiConfigStr struct {
	WWW WWWConfig `json:"www"`
}

type WWWConfig struct {
	Js     string `json:"js"`
	Html   string `json:"html"`
	Static string `json:"static"`
}
