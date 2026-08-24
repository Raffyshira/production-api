package service

import (
	"time"

	"github.com/raffyshira/project-rest-api/internal/mailer"
	"github.com/raffyshira/project-rest-api/internal/store"
	"go.uber.org/zap"
)

type Services struct {
	Auth  AuthService
	Users UserService
	Posts PostService
}

type Config struct {
	Env         string
	FrontendURL string
	MailExp     time.Duration
}

type ServiceDependencies struct {
	Store  store.Storage
	Mailer mailer.Client
	Logger *zap.SugaredLogger
	Config Config
}

func NewServices(deps ServiceDependencies) Services {
	return Services{
		Auth:  NewAuthService(deps.Store, deps.Mailer, deps.Logger, deps.Config),
		Users: NewUserService(deps.Store),
		Posts: NewPostsService(deps.Store),
	}
}
