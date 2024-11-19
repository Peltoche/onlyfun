package posts

import (
	"context"

	"github.com/Peltoche/onlyfun/internal/services/medias"
	"github.com/Peltoche/onlyfun/internal/services/users"
	"github.com/Peltoche/onlyfun/internal/tools"
	"github.com/Peltoche/onlyfun/internal/tools/sqlstorage"
)

type Service interface {
	Create(ctx context.Context, cmd *CreateCmd) (*Post, error)
	GetLatestPost(ctx context.Context) (*Post, error)
	GetPosts(ctx context.Context, start uint64, nbPosts uint64) ([]Post, error)
	GetNextPostToModerate(ctx context.Context) (*Post, error)
	CountPostsWaitingModeration(ctx context.Context) (int, error)
	GetUserStats(ctx context.Context, user *users.User) (map[Status]int, error)
	SuscribeToNewPost() <-chan Post
}

func Init(
	tools tools.Tools,
	db sqlstorage.Querier,
	medias medias.Service,
) (Service, error) {

	posts := newSqlStorage(db)

	return newService(tools, posts, medias), nil
}
