package medias

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_FileMedia_Getters(t *testing.T) {
	p := NewFakeFileMeta(t).Build()

	assert.Equal(t, p.id, p.ID())
	assert.Equal(t, p.mimetype, p.Mimetype())
	assert.Equal(t, p.checksum, p.Checksum())
	assert.Equal(t, p.mediaType, p.Type())
	assert.Equal(t, p.size, p.Size())
}
