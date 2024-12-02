package moderations

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Moderation_Getters(t *testing.T) {
	m := NewFakeModeration(t).Build()

	assert.Equal(t, m.id, m.ID())
	assert.Equal(t, m.postID, m.PostID())
	assert.Equal(t, m.reason, m.Reason())
	assert.Equal(t, m.createdAt, m.CreatedAt())
	assert.Equal(t, m.createdBy, m.CreatedBy())
}
