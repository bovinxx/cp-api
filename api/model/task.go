package model

type Task struct {
	Id         string `json:"id"`
	Code       string `json:"code"`
	Translator string `json:"translator"`
	Result     string `json:"result"`
}
