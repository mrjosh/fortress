package helpers

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_GitGetLocalConfig(t *testing.T) {

	cfg, err := GetLocalGitConfig()
	assert.NoError(t, err)

	assert.Equal(t, cfg.BaseURL, "github.com")
	assert.Equal(t, cfg.Repository, "mrjosh/fortress")

}
