package vault

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"

	"github.com/google/uuid"
)

type role struct {
	roleId   string
	secretId string
}

type secret struct {
	isDeleted bool
	mount     string
	path      string
	secret    map[string]any
}

func NewServer(token string) *httptest.Server {
	users := map[string]string{}
	roles := map[string]role{}
	secrets := []secret{}

	mux := http.NewServeMux()

	// enable userpass auth
	mux.HandleFunc("POST /v1/sys/auth/userpass", func(rw http.ResponseWriter, r *http.Request) {
		if r.Header.Get("X-Vault-Token") != token {
			rw.WriteHeader(http.StatusForbidden)
			return
		}

		rw.WriteHeader(http.StatusCreated)
	})

	// disable userpass auth
	mux.HandleFunc("DELETE /v1/sys/auth/userpass", func(rw http.ResponseWriter, r *http.Request) {
		if r.Header.Get("X-Vault-Token") != token {
			rw.WriteHeader(http.StatusForbidden)
			return
		}

		rw.WriteHeader(http.StatusOK)
	})

	// enable approle auth
	mux.HandleFunc("POST /v1/sys/auth/approle", func(rw http.ResponseWriter, r *http.Request) {
		if r.Header.Get("X-Vault-Token") != token {
			rw.WriteHeader(http.StatusForbidden)
			return
		}

		rw.WriteHeader(http.StatusCreated)
	})

	// disable approle auth
	mux.HandleFunc("DELETE /v1/sys/auth/approle", func(rw http.ResponseWriter, r *http.Request) {
		if r.Header.Get("X-Vault-Token") != token {
			rw.WriteHeader(http.StatusForbidden)
			return
		}

		rw.WriteHeader(http.StatusOK)
	})

	// create policy
	mux.HandleFunc("POST /v1/sys/policy/{name}", func(rw http.ResponseWriter, r *http.Request) {
		if r.Header.Get("X-Vault-Token") != token {
			rw.WriteHeader(http.StatusForbidden)
			return
		}

		rw.WriteHeader(http.StatusCreated)
	})

	// delete policy
	mux.HandleFunc("DELETE /v1/sys/policy/{name}", func(rw http.ResponseWriter, r *http.Request) {
		if r.Header.Get("X-Vault-Token") != token {
			rw.WriteHeader(http.StatusForbidden)
			return
		}

		rw.WriteHeader(http.StatusOK)
	})

	// lookup self token
	mux.HandleFunc("GET /v1/auth/token/lookup-self", func(rw http.ResponseWriter, r *http.Request) {
		if r.Header.Get("X-Vault-Token") != token {
			rw.WriteHeader(http.StatusForbidden)
			return
		}

		rw.WriteHeader(http.StatusOK)
	})

	// create user
	mux.HandleFunc("POST /v1/auth/userpass/users/{username}", func(rw http.ResponseWriter, r *http.Request) {
		if r.Header.Get("X-Vault-Token") != token {
			rw.WriteHeader(http.StatusForbidden)
			return
		}

		username := r.PathValue("username")
		jsonBody := map[string]string{}

		err := json.NewDecoder(r.Body).Decode(&jsonBody)
		if err != nil {
			rw.WriteHeader(http.StatusInternalServerError)
			rw.Write([]byte(err.Error()))
			return
		}

		password, ok := jsonBody["password"]

		if !ok {
			rw.WriteHeader(http.StatusInternalServerError)
			rw.Write([]byte("empty password"))
			return
		}

		users[username] = password

		rw.WriteHeader(http.StatusCreated)
	})

	// login user
	mux.HandleFunc("PUT /v1/auth/userpass/login/{username}", func(rw http.ResponseWriter, r *http.Request) {
		username := r.PathValue("username")

		var b map[string]string
		err := json.NewDecoder(r.Body).Decode(&b)
		if err != nil && err != io.EOF {
			rw.WriteHeader(http.StatusInternalServerError)
			rw.Write([]byte(err.Error()))
			return
		}

		if users[username] != b["password"] {
			rw.WriteHeader(http.StatusForbidden)
			return
		}

		rw.WriteHeader(http.StatusOK)
		rw.Header().Add("Content-Type", "application/json")
		rw.Write([]byte(fmt.Sprintf(`{
			"auth": {
			  "client_token": "%s",
			  "metadata": {
				"username": "%s"
			  }
			}
		  }`, token, username)))
	})

	// login approle
	mux.HandleFunc("PUT /v1/auth/approle/login", func(rw http.ResponseWriter, r *http.Request) {
		var b map[string]string
		err := json.NewDecoder(r.Body).Decode(&b)
		if err != nil && err != io.EOF {
			rw.WriteHeader(http.StatusInternalServerError)
			rw.Write([]byte(err.Error()))
			return
		}

		roleId := b["role_id"]
		secretId := b["secret_id"]

		for _, role := range roles {
			if role.roleId == roleId && role.secretId == secretId {
				rw.WriteHeader(http.StatusOK)
				rw.Write([]byte(fmt.Sprintf(`{
					"auth": {
					  "client_token": "%s"
					}
				  }`, token)))
				return
			}
		}

		rw.WriteHeader(http.StatusNotFound)
	})

	// delete user
	mux.HandleFunc("DELETE /v1/auth/userpass/users/{username}", func(rw http.ResponseWriter, r *http.Request) {
		if r.Header.Get("X-Vault-Token") != token {
			rw.WriteHeader(http.StatusForbidden)
			return
		}

		username := r.PathValue("username")

		delete(users, username)

		rw.WriteHeader(http.StatusOK)
	})

	// create approle
	mux.HandleFunc("POST /v1/auth/approle/role/{rolename}", func(rw http.ResponseWriter, r *http.Request) {
		if r.Header.Get("X-Vault-Token") != token {
			rw.WriteHeader(http.StatusForbidden)
			return
		}

		rolename := r.PathValue("rolename")

		roles[rolename] = role{
			roleId:   uuid.NewString(),
			secretId: "",
		}

		rw.WriteHeader(http.StatusCreated)
	})

	// get approle roleid
	mux.HandleFunc("GET /v1/auth/approle/role/{rolename}/role-id", func(rw http.ResponseWriter, r *http.Request) {
		if r.Header.Get("X-Vault-Token") != token {
			rw.WriteHeader(http.StatusForbidden)
			return
		}

		rolename := r.PathValue("rolename")

		if r, ok := roles[rolename]; ok {
			rw.WriteHeader(http.StatusOK)
			rw.Write([]byte(fmt.Sprintf(`{
				"data": {
				  "role_id": "%s"
				}
			  }`, r.roleId)))
			return
		} else {
			rw.WriteHeader(http.StatusNotFound)
			return
		}
	})

	// generate secretid
	mux.HandleFunc("POST /v1/auth/approle/role/{rolename}/secret-id", func(rw http.ResponseWriter, r *http.Request) {
		if r.Header.Get("X-Vault-Token") != token {
			rw.WriteHeader(http.StatusForbidden)
			return
		}

		rolename := r.PathValue("rolename")

		secretId := uuid.NewString()
		if r, ok := roles[rolename]; ok {
			roles[rolename] = role{
				roleId:   r.roleId,
				secretId: secretId,
			}
		} else {
			rw.WriteHeader(http.StatusNotFound)
			return
		}

		rw.WriteHeader(http.StatusCreated)
		rw.Write([]byte(fmt.Sprintf(`{
			"data": {
			  "secret_id": "%s"
			}
		  }`, secretId)))
	})

	// delete approle
	mux.HandleFunc("DELETE /v1/auth/approle/role/{rolename}", func(rw http.ResponseWriter, r *http.Request) {
		if r.Header.Get("X-Vault-Token") != token {
			rw.WriteHeader(http.StatusForbidden)
			return
		}

		rolename := r.PathValue("rolename")

		_, ok := roles[rolename]

		if !ok {
			rw.WriteHeader(http.StatusNotFound)
			return
		}

		delete(roles, rolename)

		rw.WriteHeader(http.StatusOK)
	})

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
