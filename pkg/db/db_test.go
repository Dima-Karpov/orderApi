package db

import (
	"github.com/gofrs/uuid"
	"reflect"
	"sync"
	"testing"
	"time"
)

func TestNewOrder(t *testing.T) {
	db := New()
	o := Order{
		IsOpen:          true,
		DeliveryTime:    time.Now().Unix(),
		DeliveryAddress: "123 Main St",
		Products:        []Product{{Name: "Product 1", Price: 10.0}},
	}
	id := db.NewOrder(o)
	if _, ok := db.store[id]; !ok {
		t.Errorf("Failed to create new order")
	}
}

func TestUpdateOrder(t *testing.T) {
	db := New()
	o := Order{
		IsOpen:          true,
		DeliveryTime:    time.Now().Unix(),
		DeliveryAddress: "123 Main St",
		Products:        []Product{{Name: "Product 1", Price: 10.0}},
	}
	id := db.NewOrder(o)

	o.ID = id
	o.DeliveryAddress = "456 Elm St"

	db.UpdateOrder(o)

	updatedOrder := db.store[id]

	if updatedOrder.DeliveryAddress != "456 Elm St" {
		t.Errorf("Failed to update order")
	}
}

func TestDeleteOrder(t *testing.T) {
	db := New()
	o := Order{
		IsOpen:          true,
		DeliveryTime:    time.Now().Unix(),
		DeliveryAddress: "123 Main St",
		Products:        []Product{{Name: "Product 1", Price: 10.0}},
	}
	id := db.NewOrder(o)

	db.DeleteOrder(id)

	if _, ok := db.store[id]; ok {
		t.Errorf("Failed to delete order")
	}
}

func TestDB_Orders(t *testing.T) {
	// Создаем тестовые данные
	o1 := Order{
		ID:              uuid.Must(uuid.NewV4()),
		IsOpen:          true,
		DeliveryTime:    time.Now().Unix(),
		DeliveryAddress: "123 Main St",
		Products:        []Product{{Name: "Product 1", Price: 10.0}},
	}

	type fields struct {
		m     sync.Mutex
		id    uuid.UUID
		store map[uuid.UUID]Order
	}
	tests := []struct {
		name   string
		fields fields
		want   []Order
	}{
		{
			name: "One order",
			fields: fields{
				m:  sync.Mutex{},
				id: uuid.Must(uuid.NewV4()),
				store: map[uuid.UUID]Order{
					o1.ID: o1,
				},
			},
			want: []Order{o1},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := &DB{
				m:     tt.fields.m,
				id:    tt.fields.id,
				store: tt.fields.store,
			}
			if got := db.Orders(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Orders() = %v, want %v", got, tt.want)
			}
		})
	}
}
