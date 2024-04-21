package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"orchestrator/internal/core/enteties"
)

func (h *Handlers) CreateUser(res http.ResponseWriter, req *http.Request) {
	var buf bytes.Buffer
	if _, err := buf.ReadFrom(req.Body); err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}

	var user enteties.User
	if err := json.Unmarshal(buf.Bytes(), &user); err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.service.CreateUser(req.Context(), &user); err != nil {
		if errors.Is(err, enteties.ErrorUserExist) {
			res.Header().Set("Content-Type", "application/json")
			res.WriteHeader(http.StatusBadRequest)
			res.Write([]byte(`{"error": "user already exist"}`))
			return
		}
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}

	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(http.StatusOK)
	res.Write([]byte(`{"status": "ok"}`))
}
