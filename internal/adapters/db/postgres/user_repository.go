package postgres

import (
	"database/sql"

	"github.com/google/uuid"
	"github.com/sidji-omnichannel/internal/domain/ports/repository"
	"github.com/sidji-omnichannel/internal/models"
)

type userRepository struct {
	db *sql.DB
}

// NewUserRepository creates a new postgres user repository
func NewUserRepository(db *sql.DB) repository.UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) GetPublic(userID uuid.UUID) (*models.UserPublic, error) {
	user := &models.UserPublic{}
	var avatarURL sql.NullString
	err := r.db.QueryRow(`
		SELECT id, name, avatar_url, status FROM users WHERE id = $1
	`, userID).Scan(&user.ID, &user.Name, &avatarURL, &user.Status)

	if err != nil {
		return nil, err
	}

	if avatarURL.Valid {
		user.AvatarURL = avatarURL.String
	}

	return user, nil
}
