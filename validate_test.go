package obp

import (
	"testing"
)

func Test_BasicValidation(t *testing.T) {
	t.Run("KeyMissing", func(t *testing.T) {
		t.Fail()
	})

	t.Run("KeyTypeMismatch", func(t *testing.T) {
		t.Fail()
	})
}
