package backend

import (
	"context"
	"time"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/mongo"
)

type MailingList struct {
	DisplayName string
	Topic       string
}

type Subscriber struct {
	Name        string
	Surname     string
	DisplayName string
	Email       string
}

type Subscription struct {
	Id               uuid.UUID
	MailingList      *MailingList
	Subscriber       *Subscriber
	SubscriptionDate time.Time
}

// TODO: database functions
func InitCollections(ctx context.Context, db *mongo.Database) error {
	names, err := db.ListCollectionNames(ctx, nil)

	if err != nil {
		return err
	}

	if !containsAll(names, collectionNames()) {
		for _, coll := range collectionNames() {
			err = db.CreateCollection(ctx, coll)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func containsAll(s []string, values []string) bool {
	if len(values) > len(s) {
		return false
	}
	count := len(values)
	for _, str := range s {
		for _, v := range values {
			if str == v {
				count--
				break
			}
		}
		if count == 0 {
			return true
		}
	}
	return false
}

func collectionNames() []string {
	return []string{"mailing-list", "subscriber"}
}
