package handlers

import "github.com/julienschmidt/httprouter"

// TODO : fix dependency on httprouter
type Handler interface {
	Register(router *httprouter.Router)
}
