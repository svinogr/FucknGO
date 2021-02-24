package FucknGO

import (
	"FucknGO/config"
	"testing"
)

func TestReadConfig(t *testing.T) {
	_, err := config.GetConfig()

	if err != nil {
		t.Error(err)
	}
}
