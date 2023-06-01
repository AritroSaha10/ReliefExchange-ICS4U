package main

import (
    "cloud.google.com/go/firestore"
)

type FirestoreClient interface {
    Collection(string) *firestore.CollectionRef
}

type FirestoreDocIterator interface {
    Next() (*firestore.DocumentSnapshot, error)
}
