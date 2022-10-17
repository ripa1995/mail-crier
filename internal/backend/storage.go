package backend

import (
	"time"
	"github.com/google/uuid"
)

type MailingList struct {
	DisplayName string
	Topic string
}

type Subscriber struct {
	Name    string
	Surname string
	DisplayName string
	Email string
}

type Subscription struct {
	Id uuid.UUID
	MailingList *MailingList
	Subscriber *Subscriber
	SubscriptionDate time.Time
}

//TODO: database functions
