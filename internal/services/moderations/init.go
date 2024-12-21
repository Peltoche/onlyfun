package moderations

import (
	"context"
	"database/sql"

	"github.com/Peltoche/onlyfun/internal/services/perms"
	"github.com/Peltoche/onlyfun/internal/tools"
)

type Service interface {
	ModeratePost(ctx context.Context, cmd *PostModerationCmd) (*Moderation, error)
}

func Init(tools tools.Tools, db *sql.DB, permsSvc perms.Service) Service {
	storage := newSqlStorage(db)

	return newService(tools, storage, permsSvc)
}
