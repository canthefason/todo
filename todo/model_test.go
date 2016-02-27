package todo

import (
	"fmt"
	"math/rand"
	"os"
	"testing"
	"time"

	"github.com/satori/go.uuid"
	"gopkg.in/mgo.v2"
)

func tearUp(t *testing.T, f func(*mgo.Database)) {
	mgoAddr := os.Getenv("MONGO_ADDR")
	if mgoAddr == "" {
		mgoAddr = "localhost:27017"
	}

	session, err := mgo.Dial(mgoAddr)
	if err != nil {
		t.Fatalf("Could not connect to mongo: %s", err)
	}
	defer session.Close()

	rand.Seed(time.Now().UnixNano())

	db := session.DB(fmt.Sprintf("test-%d", rand.Int31()))
	defer db.DropDatabase()

	f(db)
}

func TestMongoService(t *testing.T) {
	tearUp(t, func(db *mgo.Database) {
		ms := NewMongoService(db)

		to := &Todo{}
		to.Id = uuid.NewV4()
		to.Title = "first thing"
		to.Completed = false

		if err := ms.Upsert(to); err != nil {
			t.Fatalf("Expected nil, got %s", err)
		}
		todos, err := ms.Get()
		if err != nil {
			t.Fatalf("expected nil, got %s", err)
		}
		if len(todos) != 1 {
			t.Fatalf("expected 1 todo item, got %d", len(todos))
		}
		actualTodo := todos[0]
		if actualTodo.Title != to.Title {
			t.Fatalf("expected %s as title, got %s", to.Title, actualTodo.Title)
		}

		if err := ms.Delete(to.Id); err != nil {
			t.Fatalf("expected nil, got %s", err)
		}

		todos, err = ms.Get()
		if err != nil {
			t.Fatalf("expected nil, got %s", err)
		}
		if len(todos) != 0 {
			t.Fatalf("expected 0 todo item, got %d", len(todos))
		}

	})
}

func TestTodoService(t *testing.T) {
	tearUp(t, func(db *mgo.Database) {
		ts := NewTodoService(NewMongoService(db))

		to := &Todo{}
		to.Title = "code review"
		if err := ts.Upsert(to); err != ErrInvalidId {
			t.Fatalf("expected %s, got %v", ErrInvalidId, err)
		}

		to.Id = uuid.NewV4()
		if err := ts.Upsert(to); err != nil {
			t.Fatalf("expected nil, got %s", err)
		}

		todos, err := ts.Get()
		if err != nil {
			t.Fatalf("expected nil, got %s", err)
		}

		actualTodo := todos[0]
		if actualTodo.Title != to.Title {
			t.Fatalf("expected %s as title, got %s", to.Title, actualTodo.Title)
		}
	})
}
