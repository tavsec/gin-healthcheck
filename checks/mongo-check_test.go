package checks

import (
	"fmt"
	"testing"

	"go.mongodb.org/mongo-driver/mongo/integration/mtest"
)

func TestMongoCheck_Pass(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().CreateClient(true).ClientType(mtest.Mock))
	defer mt.Close()

	fmt.Println(mt.Client)
	check := NewMongoCheck(2, mt.Client)
	if !check.Pass() {
		t.Errorf("Expected MongoCheck.Pass to return true, got false")
	}
}

func TestMongoCheck_Name(t *testing.T) {
	check := NewMongoCheck(2, nil)
	if check.Name() != "mongodb" {
		t.Errorf("Expected MongoCheck.Name to return 'mongodb', got '%s'", check.Name())
	}
}

func TestMongoCheck_Fail(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	mt.Close()

	// Client not set
	check := NewMongoCheck(2, nil)
	if check.Pass() {
		t.Errorf("Expected MongoCheck.Pass to return false, got true")
	}

	// Connection closed
	check = NewMongoCheck(2, mt.Client)
	if check.Pass() {
		t.Errorf("Expected MongoCheck.Pass to return false, got true")
	}
}
