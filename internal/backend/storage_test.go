package backend

import (
	"context"
	"testing"
	"time"
)

func insertExpectOk(t *testing.T, s Store, backend Backend) {
	err := s.Insert(context.TODO(), backend.mongoDatabase)
	if err != nil {
		t.Log(err)
		t.Log("Struct: ", s)
		t.FailNow()
	}
}

func deleteExpectOk(t *testing.T, s Store, backend Backend) {
	del, err := s.Delete(context.TODO(), backend.mongoDatabase)
	if err != nil {
		t.Log(err)
		t.Log("Struct: ", s)
		t.FailNow()
	}
	if del != 1 {
		t.Log("Deleted more than 1 entry: ", del)
		t.Log("Struct: ", s)
		t.FailNow()
	}
}

func insertExpectErr(t *testing.T, s Store, backend Backend) {
	err := s.Insert(context.TODO(), backend.mongoDatabase)
	if err == nil {
		t.Log("Should have failed.")
		t.Log("Struct: ", s)
		t.FailNow()
	}
}

func deleteExpectErr(t *testing.T, s Store, backend Backend) {
	_, err := s.Delete(context.TODO(), backend.mongoDatabase)
	if err == nil {
		t.Log("Should have failed.")
		t.Log("Struct: ", s)
		t.FailNow()
	}
}

func deleteNothing(t *testing.T, s Store, backend Backend) {
	del, err := s.Delete(context.TODO(), backend.mongoDatabase)
	if err != nil {
		t.Log(err)
		t.FailNow()
	}
	if del != 0 {
		t.Log("Deleted more than 0 entry: ", del)
		t.Log("Struct: ", s)
		t.FailNow()
	}
}

func TestInsertDeleteMailingList(t *testing.T) {
	initEnvVar(t)
	backend, err := InitBackend(context.TODO())
	if err != nil {
		t.Log(err)
		t.FailNow()
	}

	ml := MailingList{
		DisplayName:   "ML1",
		Topic:         "TOPIC1",
		Subscriptions: nil,
	}
	ml.Delete(context.TODO(), backend.mongoDatabase)
	insertExpectOk(t, ml, backend)
	insertExpectErr(t, ml, backend)
	deleteExpectOk(t, ml, backend)
	deleteNothing(t, ml, backend)
}

func TestInsertDeleteSubscriber(t *testing.T) {
	initEnvVar(t)
	backend, err := InitBackend(context.TODO())
	if err != nil {
		t.Log(err)
		t.FailNow()
	}
	ml := MailingList{
		DisplayName:   "ML2",
		Topic:         "TOPIC2",
		Subscriptions: nil,
	}
	ml.Delete(context.TODO(), backend.mongoDatabase)
	insertExpectOk(t, ml, backend)

	s := Subscriber{
		Email: "test@gmail.com",
	}
	s.Delete(context.TODO(), backend.mongoDatabase)

	insertExpectOk(t, s, backend)
	insertExpectErr(t, s, backend)

	s2 := Subscriber{
		Email: "test2@gmail.com",
	}
	s2.Delete(context.TODO(), backend.mongoDatabase)

	insertExpectOk(t, s2, backend)
	insertExpectErr(t, s2, backend)
	deleteExpectOk(t, s2, backend)
	deleteExpectOk(t, s, backend)
	deleteNothing(t, s2, backend)
	deleteNothing(t, s, backend)
	deleteExpectOk(t, ml, backend)
}

func TestInsertDeleteSubscription(t *testing.T) {
	initEnvVar(t)
	backend, err := InitBackend(context.TODO())
	if err != nil {
		t.Log(err)
		t.FailNow()
	}
	ml := MailingList{
		DisplayName:   "ML3",
		Topic:         "TOPIC3",
		Subscriptions: []Subscription{},
	}
	ml.Delete(context.TODO(), backend.mongoDatabase)
	insertExpectOk(t, ml, backend)

	ml2 := MailingList{
		DisplayName:   "ML4",
		Topic:         "TOPIC4",
		Subscriptions: []Subscription{},
	}
	ml2.Delete(context.TODO(), backend.mongoDatabase)
	insertExpectOk(t, ml2, backend)

	s := Subscriber{
		Email: "test3@gmail.com",
	}
	s.Delete(context.TODO(), backend.mongoDatabase)

	insertExpectOk(t, s, backend)
	insertExpectErr(t, s, backend)

	subs1 := Subscription{
		mailingListDisplayName: "ML3",
		SubscriberEmail:        "test3@gmail.com",
		SubscriptionDate:       time.Now(),
	}
	subs1.Delete(context.TODO(), backend.mongoDatabase)

	insertExpectOk(t, subs1, backend)
	insertExpectErr(t, subs1, backend)

	subs2 := Subscription{
		mailingListDisplayName: "ML4",
		SubscriberEmail:        "test3@gmail.com",
		SubscriptionDate:       time.Now(),
	}
	subs2.Delete(context.TODO(), backend.mongoDatabase)

	insertExpectOk(t, subs2, backend)
	insertExpectErr(t, subs2, backend)

	deleteExpectOk(t, subs1, backend)
	deleteNothing(t, subs1, backend)

	deleteExpectOk(t, s, backend)
	//Deleting a subscriber should delete all its subscriptions
	deleteNothing(t, subs2, backend)

	deleteNothing(t, s, backend)
	deleteExpectOk(t, ml, backend)
	deleteExpectOk(t, ml2, backend)
}

func TestGetUserSubscription(t *testing.T) {
	initEnvVar(t)
	backend, err := InitBackend(context.TODO())
	if err != nil {
		t.Log(err)
		t.FailNow()
	}
	ml := MailingList{
		DisplayName:   "ML5",
		Topic:         "TOPIC5",
		Subscriptions: []Subscription{},
	}
	ml.Delete(context.TODO(), backend.mongoDatabase)
	insertExpectOk(t, ml, backend)

	s := Subscriber{
		Email: "test5@gmail.com",
	}
	s.Delete(context.TODO(), backend.mongoDatabase)

	insertExpectOk(t, s, backend)
	insertExpectErr(t, s, backend)

	timeSubs1 := time.Now()

	subs1 := Subscription{
		mailingListDisplayName: "ML5",
		SubscriberEmail:        "test5@gmail.com",
		SubscriptionDate:       timeSubs1,
	}
	subs1.Delete(context.TODO(), backend.mongoDatabase)

	insertExpectOk(t, subs1, backend)
	insertExpectErr(t, subs1, backend)

	s2 := Subscriber{
		Email: "test4@gmail.com",
	}
	s2.Delete(context.TODO(), backend.mongoDatabase)

	insertExpectOk(t, s2, backend)
	insertExpectErr(t, s2, backend)

	subs2 := Subscription{
		mailingListDisplayName: "ML5",
		SubscriberEmail:        "test4@gmail.com",
		SubscriptionDate:       time.Now(),
	}
	subs2.Delete(context.TODO(), backend.mongoDatabase)

	insertExpectOk(t, subs2, backend)
	insertExpectErr(t, subs2, backend)

	defer deleteExpectOk(t, s, backend)
	defer deleteExpectOk(t, s2, backend)
	defer deleteExpectOk(t, ml, backend)

	expected := MailingList{
		DisplayName: "ML5",
		Topic:       "TOPIC5",
		Subscriptions: []Subscription{
			{
				SubscriberEmail:  "test5@gmail.com",
				SubscriptionDate: timeSubs1,
			},
		},
	}

	result, err := s.GetUserSubscriptions(context.TODO(), backend.mongoDatabase)
	if err != nil {
		t.Log(err)
		t.FailNow()
	}
	if len(result) != 1 || len(result[0].Subscriptions) != len(expected.Subscriptions) || result[0].DisplayName != expected.DisplayName || result[0].Topic != expected.Topic || result[0].Subscriptions[0].SubscriberEmail != expected.Subscriptions[0].SubscriberEmail || result[0].Subscriptions[0].SubscriptionDate.UTC().Format("2006-01-02T15:04:05Z07:00") != expected.Subscriptions[0].SubscriptionDate.UTC().Format("2006-01-02T15:04:05Z07:00") {
		t.Log("Result: ", result)
		t.Log("Expected: ", expected)
		t.FailNow()
	}
}

func TestGetMailingList(t *testing.T) {
	initEnvVar(t)
	backend, err := InitBackend(context.TODO())
	if err != nil {
		t.Log(err)
		t.FailNow()
	}
	ml := MailingList{
		DisplayName:   "ML6",
		Topic:         "TOPIC6",
		Subscriptions: []Subscription{},
	}
	ml.Delete(context.TODO(), backend.mongoDatabase)
	insertExpectOk(t, ml, backend)

	ml2 := MailingList{
		DisplayName:   "ML7",
		Topic:         "TOPIC7",
		Subscriptions: []Subscription{},
	}
	ml2.Delete(context.TODO(), backend.mongoDatabase)
	insertExpectOk(t, ml2, backend)

	timeSubs1:= time.Now()
	subs1 := Subscription{
		mailingListDisplayName: "ML6",
		SubscriberEmail:        "test50@gmail.com",
		SubscriptionDate:       timeSubs1,
	}
	subs1.Delete(context.TODO(), backend.mongoDatabase)

	insertExpectOk(t, subs1, backend)
	insertExpectErr(t, subs1, backend)

	subs2 := Subscription{
		mailingListDisplayName: "ML7",
		SubscriberEmail:        "test40@gmail.com",
		SubscriptionDate:       time.Now(),
	}
	subs2.Delete(context.TODO(), backend.mongoDatabase)

	insertExpectOk(t, subs2, backend)
	insertExpectErr(t, subs2, backend)

	defer deleteExpectOk(t, ml, backend)
	defer deleteExpectOk(t, ml2, backend)

	expected := MailingList{
		DisplayName: "ML6",
		Topic:       "TOPIC6",
		Subscriptions: []Subscription{
			{
				SubscriberEmail:  "test50@gmail.com",
				SubscriptionDate: timeSubs1,
			},
		},
	}

	result, err := GetMailingList(context.TODO(), backend.mongoDatabase, "ML6")
	if err != nil {
		t.Log(err)
		t.FailNow()
	}
	if len(result.Subscriptions) != len(expected.Subscriptions) || result.DisplayName != expected.DisplayName || result.Topic != expected.Topic || result.Subscriptions[0].SubscriberEmail != expected.Subscriptions[0].SubscriberEmail || result.Subscriptions[0].SubscriptionDate.UTC().Format("2006-01-02T15:04:05Z07:00") != expected.Subscriptions[0].SubscriptionDate.UTC().Format("2006-01-02T15:04:05Z07:00") {
		t.Log("Result: ", result)
		t.Log("Expected: ", expected)
		t.FailNow()
	}

	result2, err := GetAllMailingList(context.TODO(), backend.mongoDatabase)
	if err != nil {
		t.Log(err)
		t.FailNow()
	}
	if len(result2) < 2 {
		t.Log("Result length: ", len(result2))
		t.Log("Expected lenght >= 2")
		t.FailNow()
	}
}