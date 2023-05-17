package main

import (
	"context"
	"testing"

	"cloud.google.com/go/firestore"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"google.golang.org/api/iterator"

	"your_project_path/mocks"
)

func TestGetAllDonations(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockFirestoreClient := mocks.NewMockFirestoreClient(ctrl)
	mockDocIterator := mocks.NewMockFirestoreDocIterator(ctrl)

	mockFirestoreClient.EXPECT().Collection("donations").Return(&firestore.CollectionRef{})
	mockDocIterator.EXPECT().Next().Return(&firestore.DocumentSnapshot{}, iterator.Done)

	donations, err := getAllDonations(context.Background(), mockFirestoreClient)

	assert.Nil(t, err)
	assert.NotNil(t, donations)
}
