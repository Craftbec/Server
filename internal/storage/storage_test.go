package storage

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	// "github.com/stretchr/testify/mock"
	"github.com/golang/mock/gomock"
)

func TestGetReport(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockStorage := NewMockStorage(ctrl)
	report := "T"
	mockStorage.EXPECT().GetReport(context.Background(), 1).Return(report, nil)
	result, err := mockStorage.GetReport(context.Background(), 1)
	assert.Nil(t, err)
	assert.Equal(t, report, result)
}

func TestPostReport(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockStorage := NewMockStorage(ctrl)
	mockStorage.EXPECT().PostReport(context.Background(), "T").Return(nil)
	err := mockStorage.PostReport(context.Background(), "T")
	assert.Nil(t, err)
}

func TestGetObservationTime(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockStorage := NewMockStorage(ctrl)
	modelID := 1
	days := 10
	mockStorage.EXPECT().GetObservationTime(context.Background(), modelID).Return(days, nil)
	result, err := mockStorage.GetObservationTime(context.Background(), modelID)
	assert.Nil(t, err)
	assert.Equal(t, days, result)
}
