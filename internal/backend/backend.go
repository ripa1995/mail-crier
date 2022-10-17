package backend

import (
	"net/smtp"
	"go.mongodb.org/mongo-driver/mongo"
)

type SMTPConfig struct {
	Username string
	Password string
	Hostname string
}

type MongoConfig struct {
	ConnectionURI string
	DatabaseName string
}

type Backend struct {
	SMTPAuth *smtp.Auth
	MongoClient *mongo.Client
}

//TODO: Backend initialization