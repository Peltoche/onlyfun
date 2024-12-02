package posts

import (
	"bytes"
	"testing"

	"github.com/Peltoche/onlyfun/internal/services/users"
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_Post_Getters(t *testing.T) {
	p := NewFakePost(t).Build()

	assert.Equal(t, p.id, p.ID())
	assert.Equal(t, p.status, p.Status())
	assert.Equal(t, p.title, p.Title())
	assert.Equal(t, p.fileID, p.FileID())
	assert.Equal(t, p.createdAt, p.CreatedAt())
	assert.Equal(t, p.createdBy, p.CreatedBy())
}

func Test_CreateCmd_is_validatable(t *testing.T) {
	assert.Implements(t, (*validation.Validatable)(nil), new(CreateCmd))
}

func Test_CreateCmd_Validate_success(t *testing.T) {
	user := users.NewFakeUser(t).Build()
	fileContent := []byte("some-content")

	err := CreateCmd{
		Title:     "This is a title",
		Media:     bytes.NewReader(fileContent),
		CreatedBy: user,
	}.Validate()

	require.NoError(t, err)
}
