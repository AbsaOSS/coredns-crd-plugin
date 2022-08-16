// Code generated by MockGen. DO NOT EDIT.
// Source: k8s.io/client-go/tools/cache (interfaces: Indexer)

// Package k8sctrl is a generated GoMock package.
package k8sctrl

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	cache "k8s.io/client-go/tools/cache"
)

// MockIndexer is a mock of Indexer interface.
type MockIndexer struct {
	ctrl     *gomock.Controller
	recorder *MockIndexerMockRecorder
}

// MockIndexerMockRecorder is the mock recorder for MockIndexer.
type MockIndexerMockRecorder struct {
	mock *MockIndexer
}

// NewMockIndexer creates a new mock instance.
func NewMockIndexer(ctrl *gomock.Controller) *MockIndexer {
	mock := &MockIndexer{ctrl: ctrl}
	mock.recorder = &MockIndexerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockIndexer) EXPECT() *MockIndexerMockRecorder {
	return m.recorder
}

// Add mocks base method.
func (m *MockIndexer) Add(arg0 interface{}) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Add", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// Add indicates an expected call of Add.
func (mr *MockIndexerMockRecorder) Add(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Add", reflect.TypeOf((*MockIndexer)(nil).Add), arg0)
}

// AddIndexers mocks base method.
func (m *MockIndexer) AddIndexers(arg0 cache.Indexers) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddIndexers", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// AddIndexers indicates an expected call of AddIndexers.
func (mr *MockIndexerMockRecorder) AddIndexers(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddIndexers", reflect.TypeOf((*MockIndexer)(nil).AddIndexers), arg0)
}

// ByIndex mocks base method.
func (m *MockIndexer) ByIndex(arg0, arg1 string) ([]interface{}, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ByIndex", arg0, arg1)
	ret0, _ := ret[0].([]interface{})
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ByIndex indicates an expected call of ByIndex.
func (mr *MockIndexerMockRecorder) ByIndex(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ByIndex", reflect.TypeOf((*MockIndexer)(nil).ByIndex), arg0, arg1)
}

// Delete mocks base method.
func (m *MockIndexer) Delete(arg0 interface{}) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Delete", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// Delete indicates an expected call of Delete.
func (mr *MockIndexerMockRecorder) Delete(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockIndexer)(nil).Delete), arg0)
}

// Get mocks base method.
func (m *MockIndexer) Get(arg0 interface{}) (interface{}, bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Get", arg0)
	ret0, _ := ret[0].(interface{})
	ret1, _ := ret[1].(bool)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// Get indicates an expected call of Get.
func (mr *MockIndexerMockRecorder) Get(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*MockIndexer)(nil).Get), arg0)
}

// GetByKey mocks base method.
func (m *MockIndexer) GetByKey(arg0 string) (interface{}, bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetByKey", arg0)
	ret0, _ := ret[0].(interface{})
	ret1, _ := ret[1].(bool)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// GetByKey indicates an expected call of GetByKey.
func (mr *MockIndexerMockRecorder) GetByKey(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetByKey", reflect.TypeOf((*MockIndexer)(nil).GetByKey), arg0)
}

// GetIndexers mocks base method.
func (m *MockIndexer) GetIndexers() cache.Indexers {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetIndexers")
	ret0, _ := ret[0].(cache.Indexers)
	return ret0
}

// GetIndexers indicates an expected call of GetIndexers.
func (mr *MockIndexerMockRecorder) GetIndexers() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetIndexers", reflect.TypeOf((*MockIndexer)(nil).GetIndexers))
}

// Index mocks base method.
func (m *MockIndexer) Index(arg0 string, arg1 interface{}) ([]interface{}, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Index", arg0, arg1)
	ret0, _ := ret[0].([]interface{})
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Index indicates an expected call of Index.
func (mr *MockIndexerMockRecorder) Index(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Index", reflect.TypeOf((*MockIndexer)(nil).Index), arg0, arg1)
}

// IndexKeys mocks base method.
func (m *MockIndexer) IndexKeys(arg0, arg1 string) ([]string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "IndexKeys", arg0, arg1)
	ret0, _ := ret[0].([]string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// IndexKeys indicates an expected call of IndexKeys.
func (mr *MockIndexerMockRecorder) IndexKeys(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "IndexKeys", reflect.TypeOf((*MockIndexer)(nil).IndexKeys), arg0, arg1)
}

// List mocks base method.
func (m *MockIndexer) List() []interface{} {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "List")
	ret0, _ := ret[0].([]interface{})
	return ret0
}

// List indicates an expected call of List.
func (mr *MockIndexerMockRecorder) List() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "List", reflect.TypeOf((*MockIndexer)(nil).List))
}

// ListIndexFuncValues mocks base method.
func (m *MockIndexer) ListIndexFuncValues(arg0 string) []string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListIndexFuncValues", arg0)
	ret0, _ := ret[0].([]string)
	return ret0
}

// ListIndexFuncValues indicates an expected call of ListIndexFuncValues.
func (mr *MockIndexerMockRecorder) ListIndexFuncValues(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListIndexFuncValues", reflect.TypeOf((*MockIndexer)(nil).ListIndexFuncValues), arg0)
}

// ListKeys mocks base method.
func (m *MockIndexer) ListKeys() []string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListKeys")
	ret0, _ := ret[0].([]string)
	return ret0
}

// ListKeys indicates an expected call of ListKeys.
func (mr *MockIndexerMockRecorder) ListKeys() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListKeys", reflect.TypeOf((*MockIndexer)(nil).ListKeys))
}

// Replace mocks base method.
func (m *MockIndexer) Replace(arg0 []interface{}, arg1 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Replace", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// Replace indicates an expected call of Replace.
func (mr *MockIndexerMockRecorder) Replace(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Replace", reflect.TypeOf((*MockIndexer)(nil).Replace), arg0, arg1)
}

// Resync mocks base method.
func (m *MockIndexer) Resync() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Resync")
	ret0, _ := ret[0].(error)
	return ret0
}

// Resync indicates an expected call of Resync.
func (mr *MockIndexerMockRecorder) Resync() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Resync", reflect.TypeOf((*MockIndexer)(nil).Resync))
}

// Update mocks base method.
func (m *MockIndexer) Update(arg0 interface{}) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Update", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// Update indicates an expected call of Update.
func (mr *MockIndexerMockRecorder) Update(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Update", reflect.TypeOf((*MockIndexer)(nil).Update), arg0)
}