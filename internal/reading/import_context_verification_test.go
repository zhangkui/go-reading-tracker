package reading

import (
	"context"
	"errors"
	"io"
	"testing"
)

type cancellationProbeReader struct{ read bool }

func (reader *cancellationProbeReader) Read(buffer []byte) (int, error) {
	reader.read = true
	copy(buffer, "[]")
	return 2, io.EOF
}

func TestImportChecksCanceledContextBeforeReading(t *testing.T) {
	service := NewService()
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	reader := &cancellationProbeReader{}
	_, err := service.Import(ctx, reader)
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("expected context cancellation, got %v", err)
	}
	if reader.read {
		t.Fatal("import read input after cancellation")
	}
}
