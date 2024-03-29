// Code generated by MockGen. DO NOT EDIT.
// Source: internal/bikes/repository/repository.go
//
// Generated by this command:
//
//	mockgen -source=internal/bikes/repository/repository.go -destination=internal/bikes/repository/mocks/repository_mock.go -package=mocks
//

// Package mocks is a generated GoMock package.
package mocks

import (
	models "bikesRentalAPI/internal/bikes/models"
	reflect "reflect"

	gomock "go.uber.org/mock/gomock"
)

// MockBikeRepository is a mock of BikeRepository interface.
type MockBikeRepository struct {
	ctrl     *gomock.Controller
	recorder *MockBikeRepositoryMockRecorder
}

// MockBikeRepositoryMockRecorder is the mock recorder for MockBikeRepository.
type MockBikeRepositoryMockRecorder struct {
	mock *MockBikeRepository
}

// NewMockBikeRepository creates a new mock instance.
func NewMockBikeRepository(ctrl *gomock.Controller) *MockBikeRepository {
	mock := &MockBikeRepository{ctrl: ctrl}
	mock.recorder = &MockBikeRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockBikeRepository) EXPECT() *MockBikeRepositoryMockRecorder {
	return m.recorder
}

// CreateBike mocks base method.
func (m *MockBikeRepository) CreateBike(bike models.CreateUpdateBikeRequest) (int64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateBike", bike)
	ret0, _ := ret[0].(int64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateBike indicates an expected call of CreateBike.
func (mr *MockBikeRepositoryMockRecorder) CreateBike(bike any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateBike", reflect.TypeOf((*MockBikeRepository)(nil).CreateBike), bike)
}

// GetBikeByID mocks base method.
func (m *MockBikeRepository) GetBikeByID(bikeID int64) (*models.Bike, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetBikeByID", bikeID)
	ret0, _ := ret[0].(*models.Bike)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetBikeByID indicates an expected call of GetBikeByID.
func (mr *MockBikeRepositoryMockRecorder) GetBikeByID(bikeID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetBikeByID", reflect.TypeOf((*MockBikeRepository)(nil).GetBikeByID), bikeID)
}

// GetBikeCostPerMinute mocks base method.
func (m *MockBikeRepository) GetBikeCostPerMinute(bikeID int64) (float64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetBikeCostPerMinute", bikeID)
	ret0, _ := ret[0].(float64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetBikeCostPerMinute indicates an expected call of GetBikeCostPerMinute.
func (mr *MockBikeRepositoryMockRecorder) GetBikeCostPerMinute(bikeID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetBikeCostPerMinute", reflect.TypeOf((*MockBikeRepository)(nil).GetBikeCostPerMinute), bikeID)
}

// IsBikeAvailable mocks base method.
func (m *MockBikeRepository) IsBikeAvailable(bikeID int64) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "IsBikeAvailable", bikeID)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// IsBikeAvailable indicates an expected call of IsBikeAvailable.
func (mr *MockBikeRepositoryMockRecorder) IsBikeAvailable(bikeID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "IsBikeAvailable", reflect.TypeOf((*MockBikeRepository)(nil).IsBikeAvailable), bikeID)
}

// ListAllBikes mocks base method.
func (m *MockBikeRepository) ListAllBikes(PageID int64) (*models.BikeList, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListAllBikes", PageID)
	ret0, _ := ret[0].(*models.BikeList)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListAllBikes indicates an expected call of ListAllBikes.
func (mr *MockBikeRepositoryMockRecorder) ListAllBikes(PageID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListAllBikes", reflect.TypeOf((*MockBikeRepository)(nil).ListAllBikes), PageID)
}

// ListAvailableBikes mocks base method.
func (m *MockBikeRepository) ListAvailableBikes(PageID int64) (*models.BikeList, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListAvailableBikes", PageID)
	ret0, _ := ret[0].(*models.BikeList)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListAvailableBikes indicates an expected call of ListAvailableBikes.
func (mr *MockBikeRepositoryMockRecorder) ListAvailableBikes(PageID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListAvailableBikes", reflect.TypeOf((*MockBikeRepository)(nil).ListAvailableBikes), PageID)
}

// SetBikeAvailability mocks base method.
func (m *MockBikeRepository) SetBikeAvailability(bikeID int64, isAvailable bool) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SetBikeAvailability", bikeID, isAvailable)
	ret0, _ := ret[0].(error)
	return ret0
}

// SetBikeAvailability indicates an expected call of SetBikeAvailability.
func (mr *MockBikeRepositoryMockRecorder) SetBikeAvailability(bikeID, isAvailable any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetBikeAvailability", reflect.TypeOf((*MockBikeRepository)(nil).SetBikeAvailability), bikeID, isAvailable)
}

// UpdateBike mocks base method.
func (m *MockBikeRepository) UpdateBike(bikeID int64, fieldsToUpdate map[string]any) (int64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateBike", bikeID, fieldsToUpdate)
	ret0, _ := ret[0].(int64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UpdateBike indicates an expected call of UpdateBike.
func (mr *MockBikeRepositoryMockRecorder) UpdateBike(bikeID, fieldsToUpdate any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateBike", reflect.TypeOf((*MockBikeRepository)(nil).UpdateBike), bikeID, fieldsToUpdate)
}
