package firestore

import (
	"crypto/sha1"
	"fmt"
	"hash/fnv"
	"time"

	"github.com/google/uuid"
)

func entityMockRef(collection string, id interface{}) string {
	return fmt.Sprintf("%s/%s", collection, id)
}

func newUUID() string {
	return uuid.New().String()
}

func newID() int {
	hash32 := fnv.New32()
	hash32.Write([]byte(newUUID()))
	return int(hash32.Sum32())
}

func Sha1(v string) string {
	return fmt.Sprintf("%x", sha1.Sum([]byte(v)))
}

func Sha1f(f string, a ...interface{}) string {
	return Sha1(fmt.Sprintf(f, a...))
}

type dynamicData map[string]interface{}

type entityMock struct {
	ID        string      `firestore:"id"`
	Name      string      `firestore:"name"`
	Age       int         `firestore:"age"`
	Data      dynamicData `firestore:"data"`
	Deleted   bool        `firestore:"deleted"`
	CreatedAt time.Time   `firestore:"created_at"`
	UpdatedAt time.Time   `firestore:"updated_at"`
}

type entityKeyMock struct {
	collection string
	name       string
	value      interface{}
}

func (k entityKeyMock) Name() string {
	return k.name
}

func (k entityKeyMock) Value() interface{} {
	return k.value
}

func (k entityKeyMock) EntityName() string {
	return k.collection
}
