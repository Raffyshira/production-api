package service

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"

	"github.com/google/uuid"
	"github.com/raffyshira/project-rest-api/internal/mailer"
	"github.com/raffyshira/project-rest-api/internal/store"
	"go.uber.org/zap"
)

type AuthService interface {
	RegisterUser(ctx context.Context, username, email, plainPassword string) (*store.User, string, error)
}

type authService struct {
	store  store.Storage
	mailer mailer.Client
	logger *zap.SugaredLogger
	config Config
}

func NewAuthService(store store.Storage, mailer mailer.Client, logger *zap.SugaredLogger, config Config) AuthService {
	return &authService{
		store:  store,
		mailer: mailer,
		logger: logger,
		config: config,
	}
}

func (s *authService) RegisterUser(ctx context.Context, username, email, plainPassword string) (*store.User, string, error) {
	// 1. Dapatkan Role ID
	role, err := s.store.Roles.GetByName(ctx, "user")
	if err != nil {
		return nil, "", err
	}

	// 2. Buat objek User
	user := &store.User{
		Username: username,
		Email:    email,
		RoleID:   role.ID,
	}

	// 3. Hash Password
	if err := user.Password.Set(plainPassword); err != nil {
		return nil, "", err
	}

	// 4. Generate Token Aktivasi
	plainToken := uuid.New().String()
	hash := sha256.Sum256([]byte(plainToken))
	hashToken := hex.EncodeToString(hash[:])

	// 5. Simpan User dan Token ke Database
	err = s.store.Users.CreateAndInvite(ctx, user, hashToken, s.config.MailExp)
	if err != nil {
		return nil, "", err
	}

	// 6. Siapkan data email
	activationURL := fmt.Sprintf("%s/confirm/%s", s.config.FrontendURL, plainToken)
	isProdEnv := s.config.Env == "production"
	vars := struct {
		Username      string
		ActivationURL string
	}{
		Username:      user.Username,
		ActivationURL: activationURL,
	}

	// 7. Kirim Email
	status, err := s.mailer.Send(mailer.UserWelcomeTemplate, user.Username, user.Email, vars, !isProdEnv)
	if err != nil {
		s.logger.Errorw("error sending welcome email", "error", err)

		// Rollback jika gagal kirim email (SAGA pattern)
		if errDel := s.store.Users.Delete(ctx, user.ID); errDel != nil {
			s.logger.Errorw("error deleting user after email failed", "error", errDel)
		}
		return nil, "", err
	}

	s.logger.Infow("Email sent", "status code", status)

	return user, plainToken, nil
}
