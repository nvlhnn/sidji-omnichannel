package postgres

import (
	"database/sql"

	"github.com/google/uuid"
	"github.com/sidji-omnichannel/internal/domain/ports/repository"
	"github.com/sidji-omnichannel/internal/models"
)

type channelRepository struct {
	db *sql.DB
}

// NewChannelRepository creates a new postgres channel repository
func NewChannelRepository(db *sql.DB) repository.ChannelRepository {
	return &channelRepository{db: db}
}

func (r *channelRepository) List(orgID uuid.UUID) ([]*models.Channel, error) {
	rows, err := r.db.Query(`
		SELECT id, organization_id, type, provider, name, config, phone_number_id, ig_user_id, facebook_page_id, tiktok_open_id, status, created_at
		FROM channels
		WHERE organization_id = $1
		ORDER BY created_at DESC
	`, orgID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var channelsList []*models.Channel
	for rows.Next() {
		channel := &models.Channel{}
		var pnID, igID, fbID, ttID sql.NullString
		err := rows.Scan(
			&channel.ID, &channel.OrganizationID, &channel.Type, &channel.Provider, &channel.Name,
			&channel.Config, &pnID, &igID, &fbID, &ttID, &channel.Status, &channel.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		if pnID.Valid { channel.PhoneNumberID = pnID.String }
		if igID.Valid { channel.IGUserID = igID.String }
		if fbID.Valid { channel.FacebookPageID = fbID.String }
		if ttID.Valid { channel.TikTokOpenID = ttID.String }
		channelsList = append(channelsList, channel)
	}

	return channelsList, nil
}

func (r *channelRepository) Get(channelID uuid.UUID) (*models.Channel, error) {
	channel := &models.Channel{}
	var pnID, igID, fbID, ttID sql.NullString
	err := r.db.QueryRow(`
		SELECT id, organization_id, type, provider, name, config, access_token, phone_number_id, ig_user_id, facebook_page_id, tiktok_open_id, status, created_at
		FROM channels
		WHERE id = $1
	`, channelID).Scan(
		&channel.ID, &channel.OrganizationID, &channel.Type, &channel.Provider, &channel.Name,
		&channel.Config, &channel.AccessToken, &pnID, &igID, &fbID, &ttID,
		&channel.Status, &channel.CreatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, repository.ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	if pnID.Valid { channel.PhoneNumberID = pnID.String }
	if igID.Valid { channel.IGUserID = igID.String }
	if fbID.Valid { channel.FacebookPageID = fbID.String }
	if ttID.Valid { channel.TikTokOpenID = ttID.String }
	return channel, nil
}

func (r *channelRepository) GetByID(orgID, channelID uuid.UUID) (*models.Channel, error) {
	channel := &models.Channel{}
	var pnID, igID, fbID, ttID sql.NullString
	err := r.db.QueryRow(`
		SELECT id, organization_id, type, provider, name, config, access_token, phone_number_id, ig_user_id, facebook_page_id, tiktok_open_id, status, created_at
		FROM channels
		WHERE id = $1 AND organization_id = $2
	`, channelID, orgID).Scan(
		&channel.ID, &channel.OrganizationID, &channel.Type, &channel.Provider, &channel.Name,
		&channel.Config, &channel.AccessToken, &pnID, &igID, &fbID, &ttID,
		&channel.Status, &channel.CreatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, repository.ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	if pnID.Valid { channel.PhoneNumberID = pnID.String }
	if igID.Valid { channel.IGUserID = igID.String }
	if fbID.Valid { channel.FacebookPageID = fbID.String }
	if ttID.Valid { channel.TikTokOpenID = ttID.String }
	return channel, nil
}

func (r *channelRepository) GetByPhoneNumberID(phoneNumberID string) (*models.Channel, error) {
	channel := &models.Channel{}
	var pnID, igID, fbID, ttID sql.NullString
	err := r.db.QueryRow(`
		SELECT id, organization_id, type, provider, name, config, access_token, phone_number_id, ig_user_id, facebook_page_id, tiktok_open_id, status
		FROM channels
		WHERE phone_number_id = $1 AND status = 'active'
	`, phoneNumberID).Scan(
		&channel.ID, &channel.OrganizationID, &channel.Type, &channel.Provider, &channel.Name,
		&channel.Config, &channel.AccessToken, &pnID, &igID, &fbID, &ttID, &channel.Status,
	)
	if err == sql.ErrNoRows {
		return nil, repository.ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	if pnID.Valid { channel.PhoneNumberID = pnID.String }
	if igID.Valid { channel.IGUserID = igID.String }
	if fbID.Valid { channel.FacebookPageID = fbID.String }
	if ttID.Valid { channel.TikTokOpenID = ttID.String }
	return channel, nil
}

func (r *channelRepository) GetByIGUserID(igUserID string) (*models.Channel, error) {
	channel := &models.Channel{}
	var pnID, igID, fbID, ttID sql.NullString
	err := r.db.QueryRow(`
		SELECT id, organization_id, type, provider, name, config, access_token, phone_number_id, ig_user_id, facebook_page_id, tiktok_open_id, status
		FROM channels
		WHERE ig_user_id = $1 AND status = 'active'
	`, igUserID).Scan(
		&channel.ID, &channel.OrganizationID, &channel.Type, &channel.Provider, &channel.Name,
		&channel.Config, &channel.AccessToken, &pnID, &igID, &fbID, &ttID, &channel.Status,
	)
	if err == sql.ErrNoRows {
		return nil, repository.ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	if pnID.Valid { channel.PhoneNumberID = pnID.String }
	if igID.Valid { channel.IGUserID = igID.String }
	if fbID.Valid { channel.FacebookPageID = fbID.String }
	if ttID.Valid { channel.TikTokOpenID = ttID.String }
	return channel, nil
}

func (r *channelRepository) GetByFacebookPageID(pageID string) (*models.Channel, error) {
	channel := &models.Channel{}
	var pnID, igID, fbID, ttID sql.NullString
	err := r.db.QueryRow(`
		SELECT id, organization_id, type, provider, name, config, access_token, phone_number_id, ig_user_id, facebook_page_id, tiktok_open_id, status
		FROM channels
		WHERE facebook_page_id = $1 AND status = 'active'
	`, pageID).Scan(
		&channel.ID, &channel.OrganizationID, &channel.Type, &channel.Provider, &channel.Name,
		&channel.Config, &channel.AccessToken, &pnID, &igID, &fbID, &ttID, &channel.Status,
	)
	if err == sql.ErrNoRows {
		return nil, repository.ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	if pnID.Valid { channel.PhoneNumberID = pnID.String }
	if igID.Valid { channel.IGUserID = igID.String }
	if fbID.Valid { channel.FacebookPageID = fbID.String }
	if ttID.Valid { channel.TikTokOpenID = ttID.String }
	return channel, nil
}

func (r *channelRepository) GetByTikTokOpenID(tiktokOpenID string) (*models.Channel, error) {
	channel := &models.Channel{}
	var pnID, igID, fbID, ttID sql.NullString
	err := r.db.QueryRow(`
		SELECT id, organization_id, type, provider, name, config, access_token, phone_number_id, ig_user_id, facebook_page_id, tiktok_open_id, status
		FROM channels
		WHERE tiktok_open_id = $1 AND status = 'active'
	`, tiktokOpenID).Scan(
		&channel.ID, &channel.OrganizationID, &channel.Type, &channel.Provider, &channel.Name,
		&channel.Config, &channel.AccessToken, &pnID, &igID, &fbID, &ttID, &channel.Status,
	)
	if err == sql.ErrNoRows {
		return nil, repository.ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	if pnID.Valid { channel.PhoneNumberID = pnID.String }
	if igID.Valid { channel.IGUserID = igID.String }
	if fbID.Valid { channel.FacebookPageID = fbID.String }
	if ttID.Valid { channel.TikTokOpenID = ttID.String }
	return channel, nil
}

func (r *channelRepository) Create(channel *models.Channel) error {
	_, err := r.db.Exec(`
		INSERT INTO channels (id, organization_id, type, provider, name, config, access_token, phone_number_id, ig_user_id, facebook_page_id, tiktok_open_id, status)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
	`, channel.ID, channel.OrganizationID, channel.Type, channel.Provider, channel.Name, channel.Config,
		channel.AccessToken, channel.PhoneNumberID, channel.IGUserID, channel.FacebookPageID, channel.TikTokOpenID, channel.Status)
	return err
}

func (r *channelRepository) Update(channel *models.Channel) error {
	_, err := r.db.Exec(`
		UPDATE channels 
		SET access_token=$1, status=$2, updated_at=NOW() 
		WHERE id=$3
	`, channel.AccessToken, channel.Status, channel.ID)
	return err
}

func (r *channelRepository) UpdateStatus(channelID uuid.UUID, status models.ChannelStatus) error {
	_, err := r.db.Exec(`
		UPDATE channels SET status = $1, updated_at = NOW() WHERE id = $2
	`, status, channelID)
	return err
}

func (r *channelRepository) Delete(orgID, channelID uuid.UUID) error {
	res, err := r.db.Exec(`
		DELETE FROM channels WHERE id = $1 AND organization_id = $2
	`, channelID, orgID)
	if err != nil {
		return err
	}
	rows, _ := res.RowsAffected()
	if rows == 0 {
		return repository.ErrNotFound
	}
	return nil
}

func (r *channelRepository) ListActiveMeta() ([]*models.Channel, error) {
	rows, err := r.db.Query(`
		SELECT id, name, type, access_token, status
		FROM channels 
		WHERE status = 'active' AND type IN ('instagram', 'whatsapp', 'facebook')
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var channels []*models.Channel
	for rows.Next() {
		ch := &models.Channel{}
		if err := rows.Scan(&ch.ID, &ch.Name, &ch.Type, &ch.AccessToken, &ch.Status); err != nil {
			return nil, err
		}
		channels = append(channels, ch)
	}
	return channels, nil
}
