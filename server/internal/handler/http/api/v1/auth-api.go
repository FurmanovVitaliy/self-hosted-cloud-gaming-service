package v1

import (
	"net/http"

	"github.com/FurmanovVitaliy/pixel-cloud/internal/usecase"
)

func (h *handler) SingIn(w http.ResponseWriter, r *http.Request) {
	req, err := decode[usecase.LogingUserReq](r)
	if err != nil {
		errorResponse(err, w, r)
		return
	}
	res, err := h.uc.SignIn(r.Context(), &req)

	if err != nil {
		errorResponse(err, w, r)
		return
	}
	cookie := &http.Cookie{
		Name:     "jwt",
		Value:    res.AccessToken,
		MaxAge:   60 * 60 * 24 * 7,
		Path:     "/",
		Domain:   "192.168.1.13",
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteDefaultMode,
	}
	http.SetCookie(w, cookie)
	if err := encode(w, r, http.StatusOK, res); err != nil {
		errorResponse(err, w, r)
		return
	}
}
func (h *handler) Logout(w http.ResponseWriter, r *http.Request) {
	cookie := http.Cookie{
		Name:     "jwt",
		Value:    "",
		MaxAge:   -1,
		Path:     "/",
		Domain:   "192.168.1.13",
		HttpOnly: true,
		Secure:   false,
	}
	http.SetCookie(w, &cookie)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Logout success"))
}

func (h *handler) SignUp(w http.ResponseWriter, r *http.Request) {
	req, err := decode[usecase.CreateUserReq](r)
	if err != nil {
		errorResponse(err, w, r)
	}
	res, err := h.uc.SignUp(r.Context(), &req)
	if err != nil {
		errorResponse(err, w, r)
		return
	}
	if err := encode(w, r, http.StatusCreated, res); err != nil {
		errorResponse(err, w, r)
		return
	}
}
