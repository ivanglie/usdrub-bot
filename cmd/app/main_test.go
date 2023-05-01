package main

import (
	"context"
	"testing"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"github.com/stretchr/testify/assert"
)

func Test_getReplyMessageID(t *testing.T) {
	id := getReplyMessageID(&models.Message{ID: 1, Chat: models.Chat{ID: 11, Type: "group"}})
	assert.Equal(t, 1, id)

	id = getReplyMessageID(&models.Message{ID: 2, Chat: models.Chat{ID: 11, Type: "supergroup"}})
	assert.Equal(t, 2, id)

	id = getReplyMessageID(&models.Message{ID: 3, Chat: models.Chat{ID: 11, Type: "channel"}})
	assert.Equal(t, 3, id)

	id = getReplyMessageID(&models.Message{ID: 4, Chat: models.Chat{ID: 11, Type: "private"}})
	assert.Zero(t, id)
}

func Test_forexHandler(t *testing.T) {
	ctx := context.TODO()
	m := &models.Message{ID: 1, Chat: models.Chat{ID: 11, Type: "group"}}

	forexHandler(ctx, &bot.Bot{}, &models.Update{Message: m})

}
