package vault_mock

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
)

type secret struct {
	isDeleted bool
	mount     string
	path      string
	secret    map[string]any
}

func NewVaultServer(token string) *httptest.Server {
	secrets := []secret{}

	mux := http.NewServeMux()

	// create secret
	mux.HandleFunc("POST /v1/{mount}/data/{path}", func(rw http.ResponseWriter, r *http.Request) {
		if r.Header.Get("X-Vault-Token") != token {
			rw.WriteHeader(http.StatusForbidden)
			return
		}

		mount := r.PathValue("mount")
		path := r.PathValue("path")

		var data map[string]any

		err := json.NewDecoder(r.Body).Decode(&data)
		if err != nil && err != io.EOF {
			rw.WriteHeader(http.StatusInternalServerError)
			rw.Write([]byte(err.Error()))
			return
		}

		secrets = append(secrets, secret{
			mount:  mount,
			path:   path,
			secret: data["data"].(map[string]any),
		})

		rw.WriteHeader(http.StatusCreated)
	})

	// get secret
	mux.HandleFunc("GET /v1/{mount}/data/{path}", func(rw http.ResponseWriter, r *http.Request) {
		if r.Header.Get("X-Vault-Token") != token {
			rw.WriteHeader(http.StatusForbidden)
			return
		}

		mount := r.PathValue("mount")
		path := r.PathValue("path")

		for _, secret := range secrets {
			if !secret.isDeleted && secret.mount == mount && secret.path == path {
				b, err := json.Marshal(secret.secret)
				if err != nil && err != io.EOF {
					rw.WriteHeader(http.StatusInternalServerError)
					rw.Write([]byte(err.Error()))
					return
				}

				rw.WriteHeader(http.StatusOK)
				rw.Header().Add("Cache-Control", "no-cache")
				rw.Header().Add("Content-Type", "application/json")
				rw.Write([]byte(fmt.Sprintf(`{"data": {"data": %s}}`, string(b))))
				return
			}
		}

		rw.WriteHeader(http.StatusNotFound)
	})

	// delete secret
	mux.HandleFunc("DELETE /v1/{mount}/data/{path}", func(rw http.ResponseWriter, r *http.Request) {
		if r.Header.Get("X-Vault-Token") != token {
			rw.WriteHeader(http.StatusForbidden)
			return
		}

		mount := r.PathValue("mount")
		path := r.PathValue("path")

		for _, secret := range secrets {
			if !secret.isDeleted && secret.mount == mount && secret.path == path {
				secret.isDeleted = true
				rw.WriteHeader(http.StatusOK)
				return
			}
		}

		rw.WriteHeader(http.StatusNotFound)
	})

	return httptest.NewServer(mux)
}
