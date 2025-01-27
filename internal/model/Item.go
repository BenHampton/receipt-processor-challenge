package model

type Item struct {
	ShortDescription string `json:"shortDescription"`
	Price            string `json:"price"`
}

var Items []Item
