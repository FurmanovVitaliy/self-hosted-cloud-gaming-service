package metric

import (
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

const (
	URL     = "/api/metric/heartbeat"
	URLtest = "/api/test"
)

type Handler struct{}

func (h *Handler) Register(router *httprouter.Router) {
	router.HandlerFunc(http.MethodGet, URL, h.Heartbeat)
	router.HandlerFunc(http.MethodGet, URLtest, h.Test)
}
func (h *Handler) Test(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(204)
}
func (h *Handler) Heartbeat(w http.ResponseWriter, r *http.Request) {
	// Выводим HTTP-метод и URL
	fmt.Printf("HTTP метод: %s\n", r.Method)
	fmt.Printf("URL: %s\n", r.URL.String())

	// Выводим заголовки запроса
	fmt.Println("Заголовки:")
	for key, values := range r.Header {
		for _, value := range values {
			fmt.Printf("%s: %s\n", key, value)
		}
	}
	fmt.Println("Тело запроса:")
	buf := make([]byte, 1024)
	for {
		n, err := r.Body.Read(buf)
		if err != nil {
			break
		}
		fmt.Print(string(buf[:n]))
	}
	w.WriteHeader(200)
}
