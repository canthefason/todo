package todo

import (
	"errors"
	"time"

	"github.com/satori/go.uuid"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

var ErrInvalidId = errors.New("invalid id")

type Service interface {
	Upsert(*Todo) error
	Get() ([]Todo, error)
	Delete(id uuid.UUID) error
}

type MongoService struct {
	*mgo.Collection
}

func NewMongoService(db *mgo.Database) *MongoService {
	return &MongoService{
		Collection: db.C("todo"),
	}
}

func (ms *MongoService) Upsert(t *Todo) error {
	_, err := ms.UpsertId(t.Id, t)

	return err
}

func (ms *MongoService) Get() ([]Todo, error) {
	todos := make([]Todo, 0)
	err := ms.Find(bson.M{}).All(&todos)

	return todos, err
}

func (ms *MongoService) Delete(id uuid.UUID) error {
	return ms.Collection.RemoveId(id)
}

type Todo struct {
	Id        uuid.UUID `bson:"_id" json:"id"`
	Title     string    `bson:"title" json:"title"`
	Completed bool      `bson:"completed" json:"completed"`
	UpdatedAt time.Time `bson:"updatedAt" json:"updatedAt"`
	//UserId    bson.ObjectId `bson:"userId" json:"userId"`
}

type TodoService struct {
	Service
}

func NewTodoService(s Service) *TodoService {
	return &TodoService{
		Service: s,
	}
}

func (ts *TodoService) Upsert(t *Todo) error {
	// Check user id
	emptyId, _ := uuid.FromString("")
	if uuid.Equal(t.Id, emptyId) {
		return ErrInvalidId
	}
	t.UpdatedAt = time.Now()

	return ts.Service.Upsert(t)
}
