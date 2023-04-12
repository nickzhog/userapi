package repositories

import (
	"context"

	"github.com/nickzhog/userapi/internal/server/config"
	"github.com/nickzhog/userapi/internal/server/service/user"
	userdb "github.com/nickzhog/userapi/internal/server/service/user/storage/filestorage"
	"github.com/nickzhog/userapi/pkg/logging"
)

type Repositories struct {
	User user.Repository
}

func GetRepositories(ctx context.Context, logger *logging.Logger, cfg config.Config) Repositories {

	return Repositories{
		User: userdb.NewFileStorage(cfg.Stores.UserStoreFile, logger),
	}
}
