package backend

import (
	"context"
	"testing"
	"github.com/joho/godotenv"
)

func initEnvVar(t *testing.T) {
	if err := godotenv.Load("../../.env"); err != nil {
		t.Fatal("No .env file found")
	}
}

func initWrongSMTPVar(t *testing.T) {
	t.Setenv("SMTP_HOSTNAME", "value")
	t.Setenv("SMTP_TLS_PORT", "value")
	t.Setenv("SMTP_USER", "value")
	t.Setenv("SMTP_PW", "value")
}

func TestInitWrongSMTPAuth(t *testing.T) {
	initEnvVar(t)
	initWrongSMTPVar(t)
	_, err := InitBackend(context.TODO())
	t.Log(err)
	if err == nil {
		t.FailNow()
	}
}

func initWrongMongoDbVar(t *testing.T) {
	t.Setenv("MONGODB_URI", "value")
	t.Setenv("MONGODB_DATABASE_NAME", "value")
}

func TestInitWrongMongoDb(t *testing.T) {
	initEnvVar(t)
	initWrongMongoDbVar(t)
	_, err := InitBackend(context.TODO())
	t.Log(err)
	if err == nil {
		t.FailNow()
	}
}

func TestInitBackend(t *testing.T) {
	initEnvVar(t)
	backend, err := InitBackend(context.TODO())
	if err != nil {
		t.Log(err)
		t.FailNow()
	}
	backend.Terminate(context.TODO())
}