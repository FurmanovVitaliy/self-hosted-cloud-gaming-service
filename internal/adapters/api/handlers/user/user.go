package user

import (
	"cloud/internal/adapters/api/handlers"
	"cloud/internal/apperror"
	"cloud/internal/domain/user"
	"cloud/pkg/logger"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

var _ handlers.Handler = &handler{}

const (
	loginURL  = "/login"
	singupURL = "/singup"
	logoutURL = "/logout"

	usersURL = "/users"
	//userURL  = "/users/:uuid"
)

type handler struct {
	service *user.Service
	logger  *logger.Logger
}

func NewHandler(service *user.Service, logger *logger.Logger) handlers.Handler {
	return &handler{
		service: service,
		logger:  logger,
	}

}

func (h *handler) Register(router *httprouter.Router) {
	router.HandlerFunc(http.MethodGet, usersURL, apperror.Middleware(h.GetList))
	//router.HandlerFunc(http.MethodGet, userURL, apperror.Middleware(h.GetUserByUUID))
	router.HandlerFunc(http.MethodPost, singupURL, apperror.Middleware(h.CreateUser))
	router.HandlerFunc(http.MethodPost, loginURL, apperror.Middleware(h.Login))
	router.HandlerFunc(http.MethodGet, logoutURL, apperror.Middleware(h.Logout))
	//router.HandlerFunc(http.MethodPut, userURL, apperror.Middleware(h.UpdateUser))
	//router.HandlerFunc(http.MethodPatch, userURL, apperror.Middleware(h.PartiallyUpdateUser))
	//router.HandlerFunc(http.MethodDelete, userURL, apperror.Middleware(h.DeleteUser))
}
func (h *handler) GetList(w http.ResponseWriter, r *http.Request) error {
	users, err := h.service.GetList(r.Context())
	if err != nil {
		return err
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(users)
	return nil
}
func (h *handler) CreateUser(w http.ResponseWriter, r *http.Request) error {
	var req user.CreateUserReq
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		return err
	}
	res, err := h.service.Create(r.Context(), &req)
	if err != nil {
		return err
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(res)

	return nil
}
func (h *handler) Login(w http.ResponseWriter, r *http.Request) error {
	var user user.LogingUserReq
	decoder := json.NewDecoder(r.Body)

	if err := decoder.Decode(&user); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return err
	}
	fmt.Println(user)
	u, err := h.service.Login(r.Context(), &user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return err
	}
	cookie := &http.Cookie{
		Name:     "jwt",
		Value:    u.AccessToken,
		MaxAge:   60 * 60 * 24 * 7,
		Path:     "/",
		Domain:   "192.168.1.13",
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteDefaultMode,
	}

	http.SetCookie(w, cookie)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(user)
	return nil
}

func (h *handler) Logout(w http.ResponseWriter, r *http.Request) error {
	cookie := http.Cookie{
		Name:     "jwt",
		Value:    "",
		MaxAge:   -1,
		Path:     "",
		Domain:   "",
		HttpOnly: true,
		Secure:   false,
	}
	http.SetCookie(w, &cookie)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Logout success"))

	return nil
}
