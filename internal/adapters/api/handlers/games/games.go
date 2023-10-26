package games

import (
	"cloud/internal/domain/games"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

const (
	gamesURL = "/games"
	gameURL  = "/games/:id"
)

type handler struct {
	gameService games.Service
}

func NewHandler(gameService games.Service) *handler {
	return &handler{gameService: gameService}
}

func (h *handler) Register(router *httprouter.Router) {
	router.GET(gamesURL, h.GetAll)
}

func (h *handler) GetAll(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	cookies := r.Cookies()
	var tokenString string
	if len(cookies) == 0 {
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Println("no cookie")
		return

	}

	for _, cookie := range cookies {
		if cookie.Name == "jwt" {
			tokenString = cookie.Value
		}
	}
	fmt.Println("tocken :", tokenString)
	_, err := h.gameService.CheckToken(tokenString)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	//w.Header().Set("Access-Control-Allow-Origin", "")
	w.Header().Set("Content-Type", "application/json")
	games, _ := h.gameService.GetAll()
	gamesJSON, _ := json.Marshal(games)
	w.Write(gamesJSON)
	w.WriteHeader(http.StatusOK)

}
