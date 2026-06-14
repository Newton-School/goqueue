package backend

import (
	"errors"
	"testing"
)

func TestBackendErrorsSupportErrorsIs(t *testing.T) {
	err := errors.Join(ErrTaskMessageNotFound, ErrConsumerGroupNotFound)

	if !errors.Is(err, ErrTaskMessageNotFound) {
		t.Fatal("errors.Is should match ErrTaskMessageNotFound")
	}
	if !errors.Is(err, ErrConsumerGroupNotFound) {
		t.Fatal("errors.Is should match ErrConsumerGroupNotFound")
	}
}
