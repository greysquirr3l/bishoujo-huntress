package huntress

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/greysquirr3l/bishoujo-huntress/internal/domain/auditlog"
	"github.com/greysquirr3l/bishoujo-huntress/internal/domain/common"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// mockAuditLogRepo implements repository.AuditLogRepository using testify/mock
type mockAuditLogRepo struct {
	mock.Mock
}

func (m *mockAuditLogRepo) Get(ctx context.Context, id string) (*auditlog.AuditLog, error) {
	args := m.Called(ctx, id)
	if v := args.Get(0); v != nil {
		return v.(*auditlog.AuditLog), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *mockAuditLogRepo) List(ctx context.Context, params *auditlog.ListParams) ([]*auditlog.AuditLog, *common.Pagination, error) {
	args := m.Called(ctx, params)
	if v := args.Get(0); v != nil {
		return v.([]*auditlog.AuditLog), args.Get(1).(*common.Pagination), args.Error(2)
	}
	return nil, nil, args.Error(2)
}

func TestAuditLogService_List(t *testing.T) {
	repo := new(mockAuditLogRepo)
	svc := NewAuditLogService(repo)
	tm := time.Now()
	params := &AuditLogListParams{StartTime: &tm, Page: 1, Limit: 10}

	internalParams := &auditlog.ListParams{
		StartTime:    params.StartTime,
		EndTime:      params.EndTime,
		Actor:        params.Actor,
		Action:       params.Action,
		ResourceType: params.ResourceType,
		ResourceID:   params.ResourceID,
		Page:         params.Page,
		Limit:        params.Limit,
	}

	internalLogs := []*auditlog.AuditLog{{ID: "1", Actor: "user", Action: "login", Timestamp: tm}}
	internalPag := &common.Pagination{Page: 1, PerPage: 10, TotalPages: 1, TotalItems: 1}

	repo.On("List", mock.Anything, internalParams).Return(internalLogs, internalPag, nil)

	result, pagination, err := svc.List(context.Background(), params)
	assert.NoError(t, err)
	// Convert internalLogs to public logs for comparison
	publicLogs := []*AuditLog{{ID: "1", Actor: "user", Action: "login", Timestamp: tm}}
	expectedPag := &Pagination{CurrentPage: 1, PerPage: 10, TotalPages: 1, TotalItems: 1}
	assert.Equal(t, publicLogs, result)
	assert.Equal(t, expectedPag, pagination)
}

func TestAuditLogService_Get(t *testing.T) {
	repo := new(mockAuditLogRepo)
	svc := NewAuditLogService(repo)
	tm := time.Now()
	internalLog := &auditlog.AuditLog{ID: "1", Actor: "user", Action: "login", Timestamp: tm}
	repo.On("Get", mock.Anything, "1").Return(internalLog, nil)
	result, err := svc.Get(context.Background(), "1")
	assert.NoError(t, err)
	publicLog := &AuditLog{ID: "1", Actor: "user", Action: "login", Timestamp: tm}
	assert.Equal(t, publicLog, result)
}

func TestAuditLogService_Get_NotFound(t *testing.T) {
	repo := new(mockAuditLogRepo)
	svc := NewAuditLogService(repo)
	repo.On("Get", mock.Anything, "notfound").Return(nil, errors.New("not found"))
	result, err := svc.Get(context.Background(), "notfound")
	assert.Error(t, err)
	assert.Nil(t, result)
}
