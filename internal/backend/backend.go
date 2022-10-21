package backend

import (
	"context"
	"crypto/tls"
	"log"
	"net"
	"net/smtp"
	"os"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type SMTPConfig struct {
	Username string
	Password string
	Hostname string
	TLSPort  string
}

type MongoConfig struct {
	ConnectionURI string
	DatabaseName  string
}

type Backend struct {
	smtpAuth      *smtp.Auth
	mongoDatabase *mongo.Database
}

func InitBackend(ctx context.Context) (Backend, error) {
	smtpConfig := loadSMTPConfig()
	mongoConfig := loadMongoConfig()

	database, err := initMongoDatabaseConnection(ctx, mongoConfig)
	if err != nil {
		return Backend{}, err
	}

	auth, err := initSMTPAuth(smtpConfig)
	if err != nil {
		return Backend{}, err
	}

	err = InitDatabaseCollections(ctx, database)
	if err != nil {
		return Backend{}, err
	}

	return Backend{mongoDatabase: database, smtpAuth: auth}, nil
}

func (backend *Backend) Terminate(ctx context.Context) {
	if err := backend.mongoDatabase.Client().Disconnect(ctx); err != nil {
		panic(err)
	}
}

func initSMTPAuth(smtpConfig SMTPConfig) (*smtp.Auth, error) {
	auth := smtp.PlainAuth("", smtpConfig.Username, smtpConfig.Password, smtpConfig.Hostname)

	address := smtpConfig.Hostname + ":" + smtpConfig.TLSPort
	client, err := smtp.Dial(address)
	if err != nil {
		return nil, err
	}

	host, _, _ := net.SplitHostPort(address)
	_ = client.StartTLS(&tls.Config{ServerName: host})

	err = client.Auth(auth)
	if err != nil {
		return nil, err
	}

	defer client.Close()

	return &auth, nil
}

func initMongoDatabaseConnection(ctx context.Context, mongoConfig MongoConfig) (*mongo.Database, error) {
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoConfig.ConnectionURI))
	if err != nil {
		return nil, err
	}

	database := client.Database(mongoConfig.DatabaseName)

	return database, nil
}

func loadSMTPConfig() SMTPConfig {
	username := os.Getenv("SMTP_USER")
	if username == "" {
		log.Fatal("You must set your 'SMTP_USER' environmental variable.")
	}

	password := os.Getenv("SMTP_PW")
	if password == "" {
		log.Fatal("You must set your 'SMTP_PW' environmental variable.")
	}

	hostname := os.Getenv("SMTP_HOSTNAME")
	if hostname == "" {
		log.Fatal("You must set your 'SMTP_HOSTNAME' environmental variable.")
	}

	port := os.Getenv("SMTP_TLS_PORT")
	if hostname == "" {
		log.Fatal("You must set your 'SMTP_TLS_PORT' environmental variable.")
	}

	return SMTPConfig{Username: username, Password: password, Hostname: hostname, TLSPort: port}
}

func loadMongoConfig() MongoConfig {
	uri := os.Getenv("MONGODB_URI")
	if uri == "" {
		log.Fatal("You must set your 'MONGODB_URI' environmental variable. See\n\t https://www.mongodb.com/docs/drivers/go/current/usage-examples/#environment-variable")
	}

	database := os.Getenv("MONGODB_DATABASE_NAME")
	if database == "" {
		log.Fatal("You must set your 'MONGODB_DATABASE_NAME' environmental variable.")
	}

	return MongoConfig{ConnectionURI: uri, DatabaseName: database}
}
