package db

import (
	"github.com/gofrs/uuid"
	"sync"
	"time"
)

// Заказ на доставку товаров
type Order struct {
	ID              uuid.UUID
	IsOpen          bool
	DeliveryTime    int64
	DeliveryAddress string
	Products        []Product
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

type Product struct {
	ID    uuid.UUID
	Name  string
	Price float64
}

type DB struct {
	m     sync.Mutex          // мьютекс для синхронизации достуа
	id    uuid.UUID           // текущее значение ID для нового заказа
	store map[uuid.UUID]Order // БД заказа
}

// Конструктор БД
func New() *DB {
	db := DB{store: make(map[uuid.UUID]Order)}
	return &db
}

// Orders возвращает все заказы
func (db *DB) Orders() []Order {
	db.m.Lock()
	defer db.m.Unlock()
	var data []Order
	for _, v := range db.store {
		data = append(data, v)
	}

	return data
}

// NewOrder создает новый заказ.
func (db *DB) NewOrder(o Order) uuid.UUID {
	db.m.Lock()
	defer db.m.Unlock()
	o.ID = uuid.Must(uuid.NewV4())
	o.CreatedAt = time.Now()
	o.UpdatedAt = time.Now()
	db.store[o.ID] = o
	return o.ID
}

// UpdateOrder обновляет данные заказа по ID.
func (db *DB) UpdateOrder(o Order) {
	db.m.Lock()
	defer db.m.Unlock()

	// Проверяем, существует ли заказ с данным ID
	if _, ok := db.store[o.ID]; !ok {
		return
	}

	// Обновляем данные заказа в картотеке
	order := db.store[o.ID]
	order.UpdatedAt = time.Now()
	order = o

	// Сохраняем обновленный заказ обратно в картотеку
	db.store[o.ID] = order
}

// DeleteOrder удаляет заказ по ID.
func (db *DB) DeleteOrder(id uuid.UUID) {
	db.m.Lock()
	defer db.m.Unlock()

	delete(db.store, id)
}
