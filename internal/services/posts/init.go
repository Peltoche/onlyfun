package posts

import (
	"context"

	"github.com/Peltoche/onlyfun/internal/services/medias"
	"github.com/Peltoche/onlyfun/internal/services/perms"
	"github.com/Peltoche/onlyfun/internal/services/users"
	"github.com/Peltoche/onlyfun/internal/tools"
	"github.com/Peltoche/onlyfun/internal/tools/sqlstorage"
)

type Service interface {
	Create(ctx context.Context, cmd *CreateCmd) (*Post, error)
	GetLatestPost(ctx context.Context) (*Post, error)
	GetByID(ctx context.Context, postID uint) (*Post, error)
	GetPosts(ctx context.Context, start uint, nbPosts uint) ([]Post, error)
	GetNextPostToModerate(ctx context.Context) (*Post, error)
	CountPostsWaitingModeration(ctx context.Context) (int, error)
	GetUserStats(ctx context.Context, user *users.User) (map[Status]int, error)
	SuscribeToNewPost() <-chan Post
	ValidatePost(ctx context.Context, cmd *ValidatePostcmd) error
}

func Init(
	tools tools.Tools,
	db sqlstorage.Querier,
	mediasSvc medias.Service,
	permsSvc perms.Service,
) Service {
	storage := newSqlStorage(db)

	return newService(tools, storage, mediasSvc, permsSvc)
}
