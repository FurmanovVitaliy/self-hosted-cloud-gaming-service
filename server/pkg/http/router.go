package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/julienschmidt/httprouter"
)

type handler interface {
	SingIn(w http.ResponseWriter, r *http.Request)
	SignUp(w http.ResponseWriter, r *http.Request)
	Logout(w http.ResponseWriter, r *http.Request)

	//GetGames(w http.ResponseWriter, r *http.Request)

	//GetRooms(w http.ResponseWriter, r *http.Request)
	//CreateRoom(w http.ResponseWriter, r *http.Request)
	//JoinRoom(w http.ResponseWriter, r *http.Request)
}

/* httprouter.Router realization */
type httpRouter struct {
	router  *httprouter.Router
	handler handler
}

func (hr *httpRouter) RegisterRoutes() {
	hr.router.HandlerFunc("POST", "/signin", hr.handler.SingIn)
	hr.router.HandlerFunc("POST", "/signup", hr.handler.SignUp)
	hr.router.HandlerFunc("GET", "/logout", hr.handler.Logout)

	//hr.router.HandlerFunc("GET", "/games", hr.handler.GetGames)

	//hr.router.HandlerFunc("GET", "/rooms", hr.handler.GetRooms)
	//hr.router.HandlerFunc("POST", "/room/create", hr.handler.CreateRoom)
	//hr.router.HandlerFunc("GET", "/room/join/:uuid", hr.handler.JoinRoom)
}

func NewHttpRouter(handler handler) *httpRouter {
	return &httpRouter{
		router:  httprouter.New(),
		handler: handler,
	}
}

func (hr *httpRouter) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	hr.router.ServeHTTP(w, r)
}

/* Gin router realization */
type ginRouter struct {
	router  *gin.Engine
	handler handler
}

func NewGinRouter(handler handler) *ginRouter {
	return &ginRouter{
		router:  gin.Default(),
		handler: handler,
	}
}

func (gr *ginRouter) RegisterRoutes() {
	gr.router.POST("/api/v1/signin", gr.wrap(gr.handler.SingIn))
	gr.router.POST("/api/v1/signup", gr.wrap(gr.handler.SignUp))
	gr.router.GET("/api/v1/logout", gr.wrap(gr.handler.Logout))
}

func (gr *ginRouter) wrap(h func(w http.ResponseWriter, r *http.Request)) gin.HandlerFunc {
	return func(c *gin.Context) {
		h(c.Writer, c.Request)
	}
}

func (gr *ginRouter) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	gr.router.ServeHTTP(w, r)
}
