package api

import (
	"bytes"
	"encoding/json"
	uuid2 "github.com/gofrs/uuid"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"ordersApi/pkg/db"
	"testing"
	"time"
)

func TestAPI_ordersHandler(t *testing.T) {
	// Создаём чистый объект API для теста.
	dbase := db.New()
	dbase.NewOrder(db.Order{})
	api := New(dbase)
	// Создаём HTTP-запрос.
	req := httptest.NewRequest(http.MethodGet, "/orders", nil)
	// Создаём объект для записи ответа обработчика.
	rr := httptest.NewRecorder()
	// Вызываем маршрутизатор. Маршрутизатор для пути и метода запроса
	// вызовет обработчик. Обработчик запишет ответ в созданный объект.
	api.r.ServeHTTP(rr, req)
	// Проверяем код ответа.
	if !(rr.Code == http.StatusOK) {
		t.Errorf("код неверен: получили %d, а хотели %d", rr.Code, http.StatusOK)
	}
	// Читаем тело ответа.
	b, err := ioutil.ReadAll(rr.Body)
	if err != nil {
		t.Fatalf("не удалось раскодировать ответ сервера: %v", err)
	}
	// Раскодируем JSON в массив заказов.
	var data []db.Order
	err = json.Unmarshal(b, &data)
	if err != nil {
		t.Fatalf("не удалось раскодировать ответ сервера: %v", err)
	}
	// Проверяем, что в массиве ровно один элемент.
	const wantLen = 1
	if len(data) != wantLen {
		t.Fatalf("получено %d записей, ожидалось %d", len(data), wantLen)
	}
	// Также можно проверить совпадение заказов в результате
	// с добавленными в БД для теста.
}

func TestAPI_newOrderHandler(t *testing.T) {
	dbase := db.New()
	api := New(dbase)

	order := db.Order{
		ID:              uuid2.Must(uuid2.NewV4()),
		IsOpen:          true,
		DeliveryTime:    time.Now().Add(24 * time.Hour).Unix(),
		DeliveryAddress: "123 Test Street",
		Products: []db.Product{
			{ID: uuid2.Must(uuid2.NewV4()), Name: "Test Product 1", Price: 9.99},
			{ID: uuid2.Must(uuid2.NewV4()), Name: "Test Product 2", Price: 19.99},
		},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	body, err := json.Marshal(order)
	if err != nil {
		t.Fatalf("Не удалось сериализовать объект: %v", err)
	}

	req := httptest.NewRequest(http.MethodPost, "/orders", ioutil.NopCloser(bytes.NewReader(body)))
	rr := httptest.NewRecorder()

	api.r.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Errorf("Код невере: получили %d, а хотели %d", rr.Code, http.StatusOK)
	}

	respBody, err := ioutil.ReadAll(rr.Body)
	if err != nil {
		t.Fatalf("Не удалось прочитать ответ: %v", err)
	}

	if len(respBody) == 0 {
		t.Fatalf("Ожидался идентификатор созданного заказа")
	}
}

func TestAPI_updateOrderHandler(t *testing.T) {
	dbase := db.New()
	api := New(dbase)

	order := db.Order{
		ID:              uuid2.Must(uuid2.NewV4()),
		IsOpen:          true,
		DeliveryTime:    time.Now().Add(24 * time.Hour).Unix(),
		DeliveryAddress: "123 Test Street",
		Products: []db.Product{
			{ID: uuid2.Must(uuid2.NewV4()), Name: "Test Product 1", Price: 9.99},
			{ID: uuid2.Must(uuid2.NewV4()), Name: "Test Product 2", Price: 19.99},
		},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	id := dbase.NewOrder(order)

	updatedOrder := db.Order{
		ID:              id,
		IsOpen:          false,
		DeliveryTime:    time.Now().Add(48 * time.Hour).Unix(),
		DeliveryAddress: "456 Updated Street",
		Products: []db.Product{
			{ID: uuid2.Must(uuid2.NewV4()), Name: "Updated Product 1", Price: 29.99},
		},
		CreatedAt: order.CreatedAt,
		UpdatedAt: time.Now(),
	}

	body, err := json.Marshal(updatedOrder)
	if err != nil {
		t.Fatalf("Не удалось сериализовать объект: %v", err)
	}

	req := httptest.NewRequest(http.MethodPatch, "/orders/"+id.String(), ioutil.NopCloser(bytes.NewReader(body)))
	rr := httptest.NewRecorder()

	api.r.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("код неверен: получили %d, а хотели %d", rr.Code, http.StatusOK)
	}
}

func TestAPI_deleteOrderHandler(t *testing.T) {
	dbase := db.New()
	api := New(dbase)

	order := db.Order{
		ID:              uuid2.Must(uuid2.NewV4()),
		IsOpen:          true,
		DeliveryTime:    time.Now().Add(24 * time.Hour).Unix(),
		DeliveryAddress: "123 Test Street",
		Products: []db.Product{
			{ID: uuid2.Must(uuid2.NewV4()), Name: "Test Product 1", Price: 9.99},
			{ID: uuid2.Must(uuid2.NewV4()), Name: "Test Product 2", Price: 19.99},
		},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	id := dbase.NewOrder(order)

	req := httptest.NewRequest(http.MethodDelete, "/orders/"+id.String(), nil)
	rr := httptest.NewRecorder()

	api.r.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("код неверен: получили %d, а хотели %d", rr.Code, http.StatusOK)
	}
}
