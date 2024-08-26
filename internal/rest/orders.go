package rest

import (
	"context"
	"encoding/json"
	"github.com/Back1ng/wbtech-0/internal/usecase"
	"net/http"
)

type Handler struct {
	orderUc *usecase.OrderUsecase
}

func NewHandler(orderUc *usecase.OrderUsecase) *http.ServeMux {
	h := Handler{
		orderUc: orderUc,
	}

	mux := http.NewServeMux()

	corsMiddleware := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			next.ServeHTTP(w, r)
		})
	}

	mux.Handle("/", corsMiddleware(http.HandlerFunc(h.GetAllOrdersHandler)))
	mux.Handle("/order", corsMiddleware(http.HandlerFunc(h.GetOrderHandler)))

	return mux
}

func (h Handler) GetAllOrdersHandler(w http.ResponseWriter, r *http.Request) {
	orders, err := h.orderUc.GetAllOrders(context.Background())
	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}

	json, err := json.Marshal(orders)
	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}

	w.Write(json)
	return
}

func (h Handler) GetOrderHandler(w http.ResponseWriter, r *http.Request) {
	order, err := h.orderUc.Get(context.Background(), r.URL.Query().Get("order_id"))
	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}

	json, err := json.Marshal(order)
	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}

	w.Write(json)
	return
}
