package prowlredis_test

import (
	"context"
	"testing"

	"github.com/jonesrussell/page-prowler/mocks"
)

func TestPing(t *testing.T) {
	mockClient := mocks.NewMockClient()
	err := mockClient.Ping(context.Background())
	if err != nil {
		t.Errorf("Ping() error = %v", err)
	}
}

func TestSAdd(t *testing.T) {
	mockClient := mocks.NewMockClient()
	err := mockClient.SAdd(context.Background(), "key", "member")
	if err != nil {
		t.Errorf("SAdd() error = %v", err)
	}
}

func TestDel(t *testing.T) {
	mockClient := mocks.NewMockClient()
	err := mockClient.Del(context.Background(), "key")
	if err != nil {
		t.Errorf("Del() error = %v", err)
	}
}

func TestSMembers(t *testing.T) {
	mockClient := mocks.NewMockClient()
	_, err := mockClient.SMembers(context.Background(), "key")
	if err != nil {
		t.Errorf("SMembers() error = %v", err)
	}
}

func TestSIsMember(t *testing.T) {
	mockClient := mocks.NewMockClient()
	_, err := mockClient.SIsMember(context.Background(), "key", "member")
	if err != nil {
		t.Errorf("SIsMember() error = %v", err)
	}
}

func TestOptions(t *testing.T) {
	mockClient := mocks.NewMockClient()
	options := mockClient.Options()
	if options == nil {
		t.Error("Options() returned nil")
	}
}
