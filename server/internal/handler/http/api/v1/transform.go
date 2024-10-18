package v1

import (
	"encoding/json"
	"net/http"

	"github.com/FurmanovVitaliy/pixel-cloud/pkg/errors"
	"github.com/gorilla/websocket"
)

func encode[T any](w http.ResponseWriter, r *http.Request, status int, v T) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(v); err != nil {
		return errors.New(http.StatusInternalServerError, "API", "000001", "failed to encode response")
	}
	return nil
}

func decode[T any](r *http.Request) (T, error) {
	var v T
	if err := json.NewDecoder(r.Body).Decode(&v); err != nil {
		return v, errors.New(http.StatusBadRequest, "API", "000002", "failed to decode request")
	}
	return v, nil
}

// TODO: add  return status code and error message with out AppError class check
// errorResponse returns an error response in JSON format
func errorResponse(err error, w http.ResponseWriter, r *http.Request) {
	// Handle AppError type specifically
	if appErr, ok := err.(*errors.AppError); ok {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(appErr.TransportCode) // Use the error code as status
		response := map[string]string{
			"message": appErr.Message,
			"code":    appErr.Code,
		}
		// Encode error message as JSON
		if err := json.NewEncoder(w).Encode(response); err != nil {
			// If encoding fails, log and send fallback error
			http.Error(w, "failed to encode error response", http.StatusInternalServerError)
		}
		return
	}

	// Fallback for non-AppError errors
	http.Error(w, err.Error(), http.StatusTeapot)
}

func errorResponseWebsocket(err error, w *websocket.Conn) {
	if appErr, ok := err.(*errors.AppError); ok {
		var e = map[string]interface{}{
			"tag":     "error",
			"content": appErr.Error(),
		}
		w.WriteJSON(e)
		w.Close()
		return
	}
	w.WriteMessage(websocket.TextMessage, []byte(err.Error()))
	w.Close()
}
