package messages

import (
	"errors"
	"fmt"
	"net/http"
	//"github.com/julienschmidt/httprouter"
)

type appHandler func(w http.ResponseWriter, r *http.Request /*params httprouter.Param*/) error

func Middleware(h appHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var appErr *AppError
		err := h(w, r)
		if err != nil {
			w.Header().Set("Content-Type", "application/type")
			if errors.As(err, &appErr) {
				if errors.Is(err, ErrNotFound) {
					w.WriteHeader(http.StatusNotFound)
					w.Write(ErrNotFound.Marshal())
					return
				} /*else if errors.Is(err, NoAuthErr){
					w.WriteHeader(http.StatusUnauthorized)
					w.Write(ErrNoAuth.Marshal())
				}*/
				err = err.(*AppError)
				w.WriteHeader(http.StatusBadRequest)
				w.Write(ErrNotFound.Marshal())
			}
			fmt.Println("err", err)
			w.Write(wrapSystemError(err).Marshal())
			w.WriteHeader(http.StatusTeapot)

		}
	}
}
