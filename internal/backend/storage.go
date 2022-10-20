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

type MailingList struct {
	DisplayName   string `bson:"display_name"`
	Topic         string
	Subscriptions []Subscription `bson:"omitempty"`
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
	SubscriberEmail  string    `bson:"subscriber_email"`
	SubscriptionDate time.Time `bson:"subscription_date"`
}

func (m Subscription) Insert(ctx context.Context, db *mongo.Database, mailingListDisplayName string) error {
	var ml MailingList
	coll := db.Collection(mailingListCollection)
	filter := bson.D{{Key: "display_name", Value: mailingListDisplayName}}

	//search for mailing list with given name
	err := coll.FindOne(ctx, filter).Decode(&ml)
	if err != nil {
		//if ==ErrNoDocumnts, no mailing list with display_name has been found
		return err
	}

	//search if ml contains m, if so no need to add the subscription
	for i := 0; i < len(ml.Subscriptions); i++ {
		if ml.Subscriptions[i].SubscriberEmail == m.SubscriberEmail {
			return nil
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

type Subscriber struct {
	Name        string `bson:"omitempty"`
	Surname     string `bson:"omitempty"`
	DisplayName string `bson:"display_name,omitempty"`
	Email       string
}

func (m Subscriber) Insert(ctx context.Context, db *mongo.Database) error {
	coll := db.Collection(subscriberCollection)
	_, err := coll.InsertOne(ctx, m)
	return err
}

func (m Subscriber) Delete(ctx context.Context, db *mongo.Database) (int64, int64, int64, error) {
	//filter on email
	filter := bson.D{{Key: "email", Value: bson.D{{Key: "$eq", Value: m.Email}}}}
	opts := options.Delete().SetHint(bson.D{{Key: "email", Value: 1}})
	coll := db.Collection(subscriberCollection)

	deleteRes, err := coll.DeleteOne(ctx, filter, opts)
	if err != nil {
		return 0, 0, 0, err
	}

	//find subscription to be deleted
	coll = db.Collection(mailingListCollection)
	//filter on subscriptions in all mailing list finding those where subscriber email match the deleted one
	filter = bson.D{{Key: "subscriptions", Value: bson.D{{Key: "$elemMatch", Value: bson.D{{Key: "subscriber_email", Value: m.Email}}}}}}
	//remove operation from the array of subscriptions
	update := bson.D{{Key: "$pull", Value: bson.D{{Key: "subscriptions", Value: bson.D{{Key: "$elemMatch", Value: bson.D{{Key: "subscriber_email", Value: m.Email}}}}}}}}

	updateRes, err := coll.UpdateMany(ctx, filter, update)
	if err != nil {
		return deleteRes.DeletedCount, 0, 0, err
	}

	return deleteRes.DeletedCount, updateRes.MatchedCount, updateRes.ModifiedCount, nil
}

func InitCollections(ctx context.Context, db *mongo.Database) error {
	names, err := db.ListCollectionNames(ctx, nil)

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
