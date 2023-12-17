package metric

import (
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

}
func (h *Handler) Heartbeat(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(200)
}
