package websessions

import (
	"testing"
	"time"

	"github.com/Peltoche/onlyfun/internal/tools/secret"
	"github.com/Peltoche/onlyfun/internal/tools/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSessionTypes(t *testing.T) {
	now := time.Now()
	session := Session{
		token:     secret.NewText("some-token"),
		userID:    uuid.UUID("3a708fc5-dc10-4655-8fc2-33b08a4b33a5"),
		ip:        "192.168.1.1",
		device:    "Android - Chrome",
		createdAt: now,
	}

	assert.Equal(t, "some-token", session.Token().Raw())
	assert.Equal(t, uuid.UUID("3a708fc5-dc10-4655-8fc2-33b08a4b33a5"), session.UserID())
	assert.Equal(t, "192.168.1.1", session.IP())
	assert.Equal(t, "Android - Chrome", session.Device())
	assert.Equal(t, now, session.CreatedAt())
}

func Test_CreateCmd_Validate(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		cmd := CreateCmd{
			UserID:     "3a708fc5-dc10-4655-8fc2-33b08a4b33a5",
			UserAgent:  "firefox 4.4.4",
			RemoteAddr: "192.168.1.1:3927",
		}

		require.NoError(t, cmd.Validate())
	})

	t.Run("with an error", func(t *testing.T) {
		cmd := CreateCmd{
			UserID:     "some-invalid-id",
			UserAgent:  "firefox 4.4.4",
			RemoteAddr: "192.168.1.1:3927",
		}

		require.EqualError(t, cmd.Validate(), "UserID: must be a valid UUID v4.")
	})
}
