package main

import (
	//"io/ioutil"
	r "github.com/dancannon/gorethink"
	"log"
)

var session *r.Session

type Product struct {
	ProductId     int     `gorethink:"ProductId"`
	StoreName     string  `gorethink:"StoreName"`
	Brand         string  `gorethink:"Brand"`
	Name          string  `gorethink:"Name"`
	OriginalPrice float64 `gorethink:"OriginalPrice"`
	SalePrice     float64 `gorethink:"SalePrice"`
	Sizes         []Size  `gorethink:"Sizes"`
	Type          string  `gorethink:"Type"`
	Color         string  `gorethink:"Color"`
	ImageUrl      string  `gorethink:"ImageUrl"`
	Url           string  `gorethink:"Url"`
	//UpdateId      int     `gorethink:"UpdateId"`
}

type Size struct {
	SizeName string `gorethink:"SizeName"`
	Quantity int    `gorethink:"Quantity"`
}

func init() {
	var err error
	session, err = r.Connect(r.ConnectOpts{
		Address:  "192.168.1.4:28015",
		Database: "Benign",
	})

	if err != nil {
		log.Fatalln(err.Error())
	}
}

func main() {
	r.Table("Products").Insert(Product{
		StoreName:     "KarmaLoop",
		Brand:         "NikeSB",
		Name:          "HyperDunks",
		OriginalPrice: 104.99,
		SalePrice:     99.99,
		Sizes: []Size{{
			SizeName: "12",
			Quantity: 15,
		}},
		Type:     "Shoes",
		Color:    "Blue",
		ImageUrl: "/boo/boo",
		UpdateId: 1,
	}).Run(session)
}
