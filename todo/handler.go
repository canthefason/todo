package todo

import (
	"net/http"

	"github.com/ant0ine/go-json-rest/rest"
	"github.com/satori/go.uuid"
)

func MustInitRouter() rest.App {
	router, err := rest.MakeRouter(
		rest.Post("/", UpdateHandler),
		rest.Get("/", GetHandler),
		rest.Delete("/#id", DeleteHandler),
	)
	if err != nil {
		panic(err)
	}

	return router
}

func UpdateHandler(w rest.ResponseWriter, r *rest.Request) {
	t := new(Todo)
	if err := r.DecodeJsonPayload(t); err != nil {
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	serv := r.Env["todo"].(*TodoService)

	err := serv.Upsert(t)
	if err != nil {
		if err.Error() == ErrInvalidId.Error() {
			rest.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.WriteJson(t)
}

func GetHandler(w rest.ResponseWriter, r *rest.Request) {
	serv := r.Env["todo"].(*TodoService)
	todos, err := serv.Get()
	if err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.WriteJson(todos)
}

func DeleteHandler(w rest.ResponseWriter, r *rest.Request) {
	serv := r.Env["todo"].(*TodoService)
	id, err := uuid.FromString(r.PathParam("id"))
	if err != nil {
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := serv.Delete(id); err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
