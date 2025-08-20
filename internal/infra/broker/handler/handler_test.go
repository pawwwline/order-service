package handler

import (
	"context"
	"errors"
	"order-service/internal/domain"
	"order-service/internal/lib/logger"
	"testing"
)

type MockUseCase struct {
	called bool
	error  error
}

func (m *MockUseCase) CreateOrder(ctx context.Context, params domain.OrderParams) error {
	m.called = true
	return m.error
}

func TestProcessOrderMessage(t *testing.T) {
	logger, err := logger.InitLogger("test")
	if err != nil {
		t.Fatalf("expected logger not nil, got %v", err)
	}
	if logger == nil {
		t.Fatalf("expected logger not nil, got nil")
	}

	tests := []struct {
		name       string
		input      []byte
		mockErr    error
		wantResult Result
		wantCalled bool
	}{
		{
			name:       "valid message",
			input:      []byte(`{"order_uid": "12345"}`),
			mockErr:    nil,
			wantResult: Success,
			wantCalled: true,
		},
		{
			name:       "invalid JSON",
			input:      []byte(`{"invalid_json":`),
			mockErr:    nil,
			wantResult: DLQ,
			wantCalled: false,
		},
		{
			name:       "domain error - non-retryable",
			input:      []byte(`{"order_uid": "12345"}`),
			mockErr:    domain.ErrInvalidState,
			wantResult: DLQ,
			wantCalled: true,
		},
		{
			name:       "retryable error",
			input:      []byte(`{"order_uid": "12345"}`),
			mockErr:    errors.New("network error"),
			wantResult: Retry,
			wantCalled: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			uc := &MockUseCase{
				called: false,
				error:  tt.mockErr,
			}
			processor := NewMessageProcessor(uc, logger)

			res := processor.ProcessOrderMessage(context.Background(), tt.input)

			if res != tt.wantResult {
				t.Fatalf("expected result %v, got %v", tt.wantResult, res)
			}
			if uc.called != tt.wantCalled {
				t.Fatalf("expected use case called %v, got %v", tt.wantCalled, uc.called)
			}
		})
	}
}
