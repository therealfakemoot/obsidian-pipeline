package obp_test

import (
	"testing"
)

func Test_BasicValidation(t *testing.T) {
	t.Parallel()
	t.Run("KeyMissing", func(t *testing.T) {
		t.Parallel()
		t.Fail()
	})

	t.Run("KeyTypeMismatch", func(t *testing.T) {
		t.Parallel()
		t.Fail()
	})
}
