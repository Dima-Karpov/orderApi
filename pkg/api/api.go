package api

import (
	"encoding/json"
	"github.com/gofrs/uuid"
	uuid2 "github.com/google/uuid"
	"github.com/gorilla/mux"
	"net/http"
	"ordersApi/pkg/db"
)

// API приложения.
type API struct {
	r  *mux.Router // маршрутизатор запросов
	db *db.DB      // база данных
}

// Конструктор API.
func New(db *db.DB) *API {
	api := API{}
	api.db = db
	api.r = mux.NewRouter()
	api.endpoints()
	return &api
}

// Router возвращает маршрутизатор запросов.
func (api *API) Router() *mux.Router {
	return api.r
}

// Регистрация методов API в маршрутизаторе запросов.
func (api *API) endpoints() {
	api.r.Use(api.headersMiddleware)
	api.r.HandleFunc("/orders", api.ordersHandler).Methods(http.MethodGet)
	api.r.HandleFunc("/orders", api.newOrderHandler).Methods(http.MethodPost)
	api.r.HandleFunc("/orders/{id}", api.updateOrderHandler).Methods(http.MethodPatch)
	api.r.HandleFunc("/orders/{id}", api.deleteOrderHandler).Methods(http.MethodDelete)
}

// ordersHandler возвращает все заказы
func (api *API) ordersHandler(w http.ResponseWriter, r *http.Request) {
	orders := api.db.Orders()
	json.NewEncoder(w).Encode(orders)
}

func (api *API) newOrderHandler(w http.ResponseWriter, r *http.Request) {
	var o db.Order
	err := json.NewDecoder(r.Body).Decode(&o)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	id := api.db.NewOrder(o)
	w.Write([]byte(id.String()))
}

func (api *API) updateOrderHandler(w http.ResponseWriter, r *http.Request) {
	s := mux.Vars(r)["id"]

	var o db.Order
	err := json.NewDecoder(r.Body).Decode(&o)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Попробуем распарсить строку как UUID
	idUUID, err := uuid2.Parse(s)
	if err != nil {
		http.Error(w, "Invalid UUID format: "+err.Error(), http.StatusBadRequest)
		return
	}

	o.ID = uuid.UUID(idUUID)
	api.db.UpdateOrder(o)
	w.WriteHeader(http.StatusOK)
}

func (api *API) deleteOrderHandler(w http.ResponseWriter, r *http.Request) {
	s := mux.Vars(r)["id"]

	// Попробуем распарсить строку как UUID
	idUUID, err := uuid2.Parse(s)
	if err != nil {
		http.Error(w, "Invalid UUID format: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Удаляем заказ
	api.db.DeleteOrder(uuid.UUID(idUUID))
	w.WriteHeader(http.StatusOK)
}

// headersMiddleware устанавливает заголовки ответа сервера.
func (api *API) headersMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}
