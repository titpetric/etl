package query

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/titpetric/etl/server/internal/handler/query/model"
)

func TestQueryLoadConfig(t *testing.T) {
	cfg, err := model.Load("testdata/user.GetByID.yml")
	assert.NoError(t, err)
	assert.NotNil(t, cfg)

	b, err := model.Encode(cfg)
	fmt.Println(string(b))
}
