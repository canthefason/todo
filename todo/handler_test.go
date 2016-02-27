package todo

import (
	"errors"
	"net/http"
	"testing"

	"github.com/ant0ine/go-json-rest/rest"
	"github.com/ant0ine/go-json-rest/rest/test"
	"github.com/satori/go.uuid"
)

type MockMiddl struct {
	Todos []Todo
	Error error
}

func (m MockMiddl) MiddlewareFunc(handler rest.HandlerFunc) rest.HandlerFunc {
	return func(w rest.ResponseWriter, r *rest.Request) {

		r.Env["todo"] = NewTodoService(
			NewTodoService(m),
		)

		handler(w, r)
	}
}

func NewMockMiddl(todos []Todo, err error) *MockMiddl {
	return &MockMiddl{
		Todos: todos,
		Error: err,
	}
}

func (ms MockMiddl) Upsert(t *Todo) error {
	return ms.Error
}

func (ms MockMiddl) Get() ([]Todo, error) {
	return ms.Todos, ms.Error
}

func (ms MockMiddl) Delete(id uuid.UUID) error {
	return ms.Error
}

func prepareApi(todos []Todo, err error) *rest.Api {
	api := rest.NewApi()
	api.Use(rest.DefaultDevStack...)

	api.Use(MockMiddl{
		Todos: todos,
		Error: err,
	})

	router := MustInitRouter()
	api.SetApp(router)

	return api
}

func TestGetHandler(t *testing.T) {
	// error case
	errorApi := prepareApi([]Todo{}, errors.New("connection error"))
	recorded := test.RunRequest(t, errorApi.MakeHandler(),
		test.MakeSimpleRequest("GET", "http://1.2.3.4/", nil))

	recorded.CodeIs(http.StatusInternalServerError)
	recorded.ContentTypeIsJson()

	todos := []Todo{
		Todo{
			Id:    uuid.NewV4(),
			Title: "testing",
		},
	}

	// success case
	successApi := prepareApi(todos, nil)
	recorded = test.RunRequest(t, successApi.MakeHandler(),
		test.MakeSimpleRequest("GET", "http://1.2.3.4/", nil))

	recorded.CodeIs(http.StatusOK)
	recorded.ContentTypeIsJson()
}

func TestUpdateHandler(t *testing.T) {
	// empty payload case
	badRequestApi := prepareApi([]Todo{}, nil)
	recorded := test.RunRequest(t, badRequestApi.MakeHandler(),
		test.MakeSimpleRequest("POST", "http://1.2.3.4/", nil))

	recorded.CodeIs(http.StatusBadRequest)
	recorded.ContentTypeIsJson()

	// invalid id case
	recorded = test.RunRequest(t, badRequestApi.MakeHandler(),
		test.MakeSimpleRequest("POST", "http://1.2.3.4/", Todo{}))

	recorded.CodeIs(http.StatusBadRequest)
	recorded.ContentTypeIsJson()

	todo := Todo{Id: uuid.NewV4()}
	// internal server error case
	serverErrorApi := prepareApi([]Todo{}, errors.New("could not connect"))
	recorded = test.RunRequest(t, serverErrorApi.MakeHandler(),
		test.MakeSimpleRequest("POST", "http://1.2.3.4/", todo))

	recorded.CodeIs(http.StatusInternalServerError)
	recorded.ContentTypeIsJson()

	successApi := prepareApi([]Todo{}, nil)
	recorded = test.RunRequest(t, successApi.MakeHandler(),
		test.MakeSimpleRequest("POST", "http://1.2.3.4/", todo))

	recorded.CodeIs(http.StatusOK)
	recorded.ContentTypeIsJson()
}

func TestDeleteHandler(t *testing.T) {
	// invalid id case
	errorApi := prepareApi([]Todo{}, errors.New("not found"))
	recorded := test.RunRequest(t, errorApi.MakeHandler(),
		test.MakeSimpleRequest("DELETE", "http://1.2.3.4/123", nil))

	recorded.CodeIs(http.StatusBadRequest)
	recorded.ContentTypeIsJson()

	// internal error case
	recorded = test.RunRequest(t, errorApi.MakeHandler(),
		test.MakeSimpleRequest("DELETE", "http://1.2.3.4/96efa48f-d506-4035-968b-447ad0f75deb", nil))

	recorded.CodeIs(http.StatusInternalServerError)
	recorded.ContentTypeIsJson()

	// success case
	successApi := prepareApi([]Todo{}, nil)
	recorded = test.RunRequest(t, successApi.MakeHandler(),
		test.MakeSimpleRequest("DELETE", "http://1.2.3.4/96efa48f-d506-4035-968b-447ad0f75deb", nil))

	recorded.CodeIs(http.StatusOK)
	recorded.ContentTypeIsJson()
}
