package main

import (
	"net/http"
	"ordersApi/pkg/api"
	"ordersApi/pkg/db"
	"time"
)

func main() {
	dbase := db.New()
	p := []db.Product{
		{Name: "Яблоки", Price: 20},
		{Name: "Груши", Price: 30},
	}

	o := db.Order{
		IsOpen:       true,
		DeliveryTime: time.Now().Unix(),
		Products:     p,
	}

	dbase.NewOrder(o)
	a := api.New(dbase)

	http.ListenAndServe(":80", a.Router())
}

//http://localhost/orders
