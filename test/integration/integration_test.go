package integration

import (
	"context"
	"github.com/stretchr/testify/assert"
	tc "github.com/testcontainers/testcontainers-go/modules/compose"
	"testing"
)

func startDockerCompose(t *testing.T) {
	compose, err := tc.NewDockerCompose("./docker-compose.yml")
	assert.NoError(t, err, "NewDockerComposeAPI()")
	t.Cleanup(func() {
		assert.NoError(t, compose.Down(context.Background(), tc.RemoveOrphans(true), tc.RemoveImagesLocal), "compose.Down()")
	})
	ctx, cancel := context.WithCancel(context.Background())
	t.Cleanup(cancel)

	assert.NoError(t, compose.Up(ctx, tc.Wait(true)), "compose.Up()")
}

func TestQuestion(t *testing.T) {
	startDockerCompose(t)

}

func TestAnswer(t *testing.T) {

}

func TestComment(t *testing.T) {

}
