package backend

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const mailingListCollection string = "mailing-list"
const subscriberCollection string = "subscriber"

type AlreadySubscribedError struct {}

func (e *AlreadySubscribedError)Error() string {
	return "Already Subscribed."
}

type Store interface {
	Insert(ctx context.Context, db *mongo.Database) error
	Delete(ctx context.Context, db *mongo.Database) (int64, error)
}

type MailingList struct {
	DisplayName   string `bson:"display_name"`
	Topic         string
	Subscriptions []Subscription 
}

func (m MailingList) Insert(ctx context.Context, db *mongo.Database) error {
	coll := db.Collection(mailingListCollection)
	_, err := coll.InsertOne(ctx, m)
	return err
}

func (m MailingList) Delete(ctx context.Context, db *mongo.Database) (int64, error) {
	//filter on display name
	filter := bson.D{{Key: "display_name", Value: bson.D{{Key: "$eq", Value: m.DisplayName}}}}
	opts := options.Delete().SetHint(bson.D{{Key: "display_name", Value: 1}})
	coll := db.Collection(mailingListCollection)

	res, err := coll.DeleteOne(ctx, filter, opts)
	if err != nil {
		return 0, err
	}

	return res.DeletedCount, nil
}

type Subscription struct {
	mailingListDisplayName string    
	SubscriberEmail        string    `bson:"subscriber_email"`
	SubscriptionDate       time.Time `bson:"subscription_date"`
}

func (m Subscription) Insert(ctx context.Context, db *mongo.Database) error {
	var subs[]Subscription 
	ml := MailingList{Subscriptions: subs}
	coll := db.Collection(mailingListCollection)
	filter := bson.D{{Key: "display_name", Value: m.mailingListDisplayName}}

	//search for mailing list with given name
	err := coll.FindOne(ctx, filter).Decode(&ml)
	if err != nil {
		//if ==ErrNoDocumnts, no mailing list with display_name has been found
		return err
	}
	//search if ml contains m, if so no need to add the subscription
	for i := 0; i < len(ml.Subscriptions); i++ {
		if ml.Subscriptions[i].SubscriberEmail == m.SubscriberEmail {
			return &AlreadySubscribedError{}
		}
	}

	//append subscription info the mailing list subscriptions array
	update := bson.D{{Key: "$push", Value: bson.D{{Key: "subscriptions", Value: bson.D{{Key: "subscriber_email", Value: m.SubscriberEmail}, {Key: "subscription_date", Value: m.SubscriptionDate}}}}}}

	_, err = coll.UpdateMany(ctx, filter, update)
	if err != nil {
		return err
	}

	return nil
}

func (m Subscription) Delete(ctx context.Context, db *mongo.Database) (int64, error) {
	//find subscription to be deleted
	coll := db.Collection(mailingListCollection)
	//filter on mailing list finding the one from which the subscription must be deleted
	filter := bson.D{{Key: "display_name", Value: m.mailingListDisplayName}}
	//remove operation from the array of subscriptions
	update := bson.D{{Key: "$pull", Value: bson.D{{Key: "subscriptions",  Value: bson.D{{Key: "subscriber_email", Value: m.SubscriberEmail}}}}}}

	updateRes, err := coll.UpdateMany(ctx, filter, update)
	if err != nil {
		return 0, err
	}

	return updateRes.ModifiedCount, nil
}

type Subscriber struct {
	Name        string 
	Surname     string 
	DisplayName string `bson:"display_name"`
	Email       string
}

func (m Subscriber) Insert(ctx context.Context, db *mongo.Database) error {
	coll := db.Collection(subscriberCollection)
	_, err := coll.InsertOne(ctx, m)
	return err
}

func (m Subscriber) Delete(ctx context.Context, db *mongo.Database) (int64, error) {
	//filter on email
	filter := bson.D{{Key: "email", Value: bson.D{{Key: "$eq", Value: m.Email}}}}
	opts := options.Delete().SetHint(bson.D{{Key: "email", Value: 1}})
	coll := db.Collection(subscriberCollection)

	deleteRes, err := coll.DeleteOne(ctx, filter, opts)
	if err != nil {
		return 0, err
	}

	//find subscription to be deleted
	coll = db.Collection(mailingListCollection)
	//filter on subscriptions in all mailing list finding those where subscriber email match the deleted one
	filter = bson.D{{Key: "subscriptions", Value: bson.D{{Key: "$elemMatch", Value: bson.D{{Key: "subscriber_email", Value: m.Email}}}}}}
	//remove operation from the array of subscriptions
	update := bson.D{{Key: "$pull", Value: bson.D{{Key: "subscriptions", Value: bson.D{{Key: "subscriber_email", Value: m.Email}}}}}}

	_, err = coll.UpdateMany(ctx, filter, update)
	if err != nil {
		return deleteRes.DeletedCount, err
	}

	return deleteRes.DeletedCount, nil
}

func initDatabaseCollections(ctx context.Context, db *mongo.Database) error {
	names, err := db.ListCollectionNames(ctx, bson.D{})

	if err != nil {
		return err
	}

	collNames, collIndex := collectionInfo()
	if !containsAll(names, collNames) {
		for i, coll := range collNames {
			err = db.CreateCollection(ctx, coll)
			if err != nil {
				return err
			}
			//create a unique index on specific columns
			collection := db.Collection(coll)
			mod := mongo.IndexModel{
				Keys: bson.M{
					collIndex[i]: 1,
				},
				Options: options.Index().SetUnique(true),
			}
			collection.Indexes().CreateOne(ctx, mod)
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

func collectionInfo() ([]string, []string) {
	return []string{mailingListCollection, subscriberCollection}, []string{"display_name", "email"}
}
