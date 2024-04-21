package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"orchestrator/internal/core/enteties"
)

func (h *Handlers) Login(res http.ResponseWriter, req *http.Request) {
	var buf bytes.Buffer
	if _, err := buf.ReadFrom(req.Body); err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}

	var user enteties.User
	err := json.Unmarshal(buf.Bytes(), &user)
	if err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}

	var userDB *enteties.User
	userDB, err = h.service.Login(req.Context(), &user)
	if err != nil {
		res.WriteHeader(http.StatusBadRequest)
		return
	}

	authToken, err := buildJWTString(userDB.ID.String(), secretKey)
	if err != nil {
		res.WriteHeader(http.StatusBadRequest)
		return
	}

	http.SetCookie(res, &http.Cookie{
		Name:  authorizationTokenCookie,
		Value: authToken,
	})
	res.WriteHeader(http.StatusOK)
}
