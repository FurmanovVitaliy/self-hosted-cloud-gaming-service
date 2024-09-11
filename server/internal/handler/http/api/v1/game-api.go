package v1

import "net/http"

func (h *handler) GetGames(w http.ResponseWriter, r *http.Request) {
	games, err := h.uc.GetGames()
	if err != nil {
		errorResponse(err, w, r)
		return
	}
	if err := encode(w, r, http.StatusOK, games); err != nil {
		errorResponse(err, w, r)
		return
	}
}
