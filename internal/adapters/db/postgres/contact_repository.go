package postgres

import (
	"database/sql"
	"encoding/json"

	"github.com/google/uuid"
	"github.com/sidji-omnichannel/internal/domain/ports/repository"
	"github.com/sidji-omnichannel/internal/models"
)

type contactRepository struct {
	db *sql.DB
}

// NewContactRepository creates a new postgres contact repository
func NewContactRepository(db *sql.DB) repository.ContactRepository {
	return &contactRepository{db: db}
}

func (r *contactRepository) List(orgID uuid.UUID, page, limit int, search string) ([]*models.Contact, int, error) {
	offset := (page - 1) * limit
	query := `
		SELECT id, organization_id, name, phone, email, avatar_url, whatsapp_id, instagram_id, facebook_id, tiktok_id, created_at
		FROM contacts
		WHERE organization_id = $1
	`
	countQuery := `SELECT COUNT(*) FROM contacts WHERE organization_id = $1`
	args := []interface{}{orgID}

	if search != "" {
		searchFilter := ` AND (name ILIKE $2 OR phone ILIKE $2 OR email ILIKE $2)`
		query += searchFilter
		countQuery += searchFilter
		args = append(args, "%"+search+"%")
	}

	var total int
	if err := r.db.QueryRow(countQuery, args...).Scan(&total); err != nil {
		return nil, 0, err
	}

	query += ` ORDER BY created_at DESC LIMIT $` + string(rune('0'+len(args)+1)) + ` OFFSET $` + string(rune('0'+len(args)+2))
	args = append(args, limit, offset)

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var contacts []*models.Contact
	for rows.Next() {
		contact := &models.Contact{}
		var phone, email, avatarURL, whatsappID, instagramID, facebookID, tiktokID sql.NullString
		err := rows.Scan(
			&contact.ID, &contact.OrganizationID, &contact.Name, &phone,
			&email, &avatarURL, &whatsappID, &instagramID, &facebookID, &tiktokID, &contact.CreatedAt,
		)
		if err != nil {
			return nil, 0, err
		}
		if phone.Valid { contact.Phone = phone.String }
		if email.Valid { contact.Email = email.String }
		if avatarURL.Valid { contact.AvatarURL = avatarURL.String }
		if whatsappID.Valid { contact.WhatsAppID = whatsappID.String }
		if instagramID.Valid { contact.InstagramID = instagramID.String }
		if facebookID.Valid { contact.FacebookID = facebookID.String }
		if tiktokID.Valid { contact.TikTokID = tiktokID.String }
		contacts = append(contacts, contact)
	}
	return contacts, total, nil
}

func (r *contactRepository) GetByID(orgID, contactID uuid.UUID) (*models.Contact, error) {
	contact := &models.Contact{}
	var phone, email, avatarURL, whatsappID, instagramID, facebookID, tiktokID sql.NullString
	var metadata []byte
	err := r.db.QueryRow(`
		SELECT id, organization_id, name, phone, email, avatar_url, metadata, whatsapp_id, instagram_id, facebook_id, tiktok_id, created_at, updated_at
		FROM contacts
		WHERE id = $1 AND organization_id = $2
	`, contactID, orgID).Scan(
		&contact.ID, &contact.OrganizationID, &contact.Name, &phone,
		&email, &avatarURL, &metadata, &whatsappID, &instagramID, &facebookID, &tiktokID,
		&contact.CreatedAt, &contact.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, repository.ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	if metadata != nil {
		contact.Metadata = json.RawMessage(metadata)
	}
	if phone.Valid { contact.Phone = phone.String }
	if email.Valid { contact.Email = email.String }
	if avatarURL.Valid { contact.AvatarURL = avatarURL.String }
	if whatsappID.Valid { contact.WhatsAppID = whatsappID.String }
	if instagramID.Valid { contact.InstagramID = instagramID.String }
	if facebookID.Valid { contact.FacebookID = facebookID.String }
	if tiktokID.Valid { contact.TikTokID = tiktokID.String }
	return contact, nil
}

func (r *contactRepository) GetByWhatsAppID(orgID uuid.UUID, whatsAppID string) (*models.Contact, error) {
	contact := &models.Contact{}
	var phone, email, avatarURL, whatsappID, instagramID sql.NullString
	err := r.db.QueryRow(`
		SELECT id, organization_id, name, phone, email, avatar_url, whatsapp_id, instagram_id, created_at
		FROM contacts
		WHERE organization_id = $1 AND whatsapp_id = $2
	`, orgID, whatsAppID).Scan(
		&contact.ID, &contact.OrganizationID, &contact.Name, &phone,
		&email, &avatarURL, &whatsappID, &instagramID, &contact.CreatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, repository.ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	if phone.Valid { contact.Phone = phone.String }
	if email.Valid { contact.Email = email.String }
	if avatarURL.Valid { contact.AvatarURL = avatarURL.String }
	if whatsappID.Valid { contact.WhatsAppID = whatsappID.String }
	if instagramID.Valid { contact.InstagramID = instagramID.String }
	return contact, nil
}

func (r *contactRepository) GetByInstagramID(orgID uuid.UUID, instagramID string) (*models.Contact, error) {
	contact := &models.Contact{}
	var phone, email, avatarURL, whatsappID, instagramIDNS sql.NullString
	err := r.db.QueryRow(`
		SELECT id, organization_id, name, phone, email, avatar_url, whatsapp_id, instagram_id, created_at
		FROM contacts
		WHERE organization_id = $1 AND instagram_id = $2
	`, orgID, instagramID).Scan(
		&contact.ID, &contact.OrganizationID, &contact.Name, &phone,
		&email, &avatarURL, &whatsappID, &instagramIDNS, &contact.CreatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, repository.ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	if phone.Valid { contact.Phone = phone.String }
	if email.Valid { contact.Email = email.String }
	if avatarURL.Valid { contact.AvatarURL = avatarURL.String }
	if whatsappID.Valid { contact.WhatsAppID = whatsappID.String }
	if instagramIDNS.Valid { contact.InstagramID = instagramIDNS.String }
	return contact, nil
}

func (r *contactRepository) GetByFacebookID(orgID uuid.UUID, facebookID string) (*models.Contact, error) {
	contact := &models.Contact{}
	var phone, email, avatarURL, whatsappID, instagramID, facebookIDNS sql.NullString
	err := r.db.QueryRow(`
		SELECT id, organization_id, name, phone, email, avatar_url, whatsapp_id, instagram_id, facebook_id, created_at
		FROM contacts
		WHERE organization_id = $1 AND facebook_id = $2
	`, orgID, facebookID).Scan(
		&contact.ID, &contact.OrganizationID, &contact.Name, &phone,
		&email, &avatarURL, &whatsappID, &instagramID, &facebookIDNS, &contact.CreatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, repository.ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	if phone.Valid { contact.Phone = phone.String }
	if email.Valid { contact.Email = email.String }
	if avatarURL.Valid { contact.AvatarURL = avatarURL.String }
	if whatsappID.Valid { contact.WhatsAppID = whatsappID.String }
	if instagramID.Valid { contact.InstagramID = instagramID.String }
	if facebookIDNS.Valid { contact.FacebookID = facebookIDNS.String }
	return contact, nil
}

func (r *contactRepository) GetByTikTokID(orgID uuid.UUID, tiktokID string) (*models.Contact, error) {
	contact := &models.Contact{}
	var phone, email, avatarURL, whatsappID, instagramID, facebookID, tiktokIDNS sql.NullString
	err := r.db.QueryRow(`
		SELECT id, organization_id, name, phone, email, avatar_url, whatsapp_id, instagram_id, facebook_id, tiktok_id, created_at
		FROM contacts
		WHERE organization_id = $1 AND tiktok_id = $2
	`, orgID, tiktokID).Scan(
		&contact.ID, &contact.OrganizationID, &contact.Name, &phone,
		&email, &avatarURL, &whatsappID, &instagramID, &facebookID, &tiktokIDNS, &contact.CreatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, repository.ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	if phone.Valid { contact.Phone = phone.String }
	if email.Valid { contact.Email = email.String }
	if avatarURL.Valid { contact.AvatarURL = avatarURL.String }
	if whatsappID.Valid { contact.WhatsAppID = whatsappID.String }
	if instagramID.Valid { contact.InstagramID = instagramID.String }
	if facebookID.Valid { contact.FacebookID = facebookID.String }
	if tiktokIDNS.Valid { contact.TikTokID = tiktokIDNS.String }
	return contact, nil
}

func (r *contactRepository) Create(contact *models.Contact) error {
	_, err := r.db.Exec(`
		INSERT INTO contacts (id, organization_id, name, phone, email, avatar_url, whatsapp_id, instagram_id, facebook_id, tiktok_id)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
	`, contact.ID, contact.OrganizationID, contact.Name, contact.Phone, contact.Email, contact.AvatarURL,
		contact.WhatsAppID, contact.InstagramID, contact.FacebookID, contact.TikTokID)
	return err
}

func (r *contactRepository) Update(orgID, contactID uuid.UUID, input *models.UpdateContactInput) (*models.Contact, error) {
	query := "UPDATE contacts SET updated_at = NOW()"
	args := []interface{}{}
	argIndex := 1
	if input.Name != "" {
		query += ", name = $" + string(rune('0'+argIndex))
		args = append(args, input.Name)
		argIndex++
	}
	if input.Phone != "" {
		query += ", phone = $" + string(rune('0'+argIndex))
		args = append(args, input.Phone)
		argIndex++
	}
	if input.Email != "" {
		query += ", email = $" + string(rune('0'+argIndex))
		args = append(args, input.Email)
		argIndex++
	}
	query += " WHERE id = $" + string(rune('0'+argIndex)) + " AND organization_id = $" + string(rune('0'+argIndex+1))
	args = append(args, contactID, orgID)
	res, err := r.db.Exec(query, args...)
	if err != nil {
		return nil, err
	}
	rows, _ := res.RowsAffected()
	if rows == 0 {
		return nil, repository.ErrNotFound
	}
	return r.GetByID(orgID, contactID)
}

func (r *contactRepository) UpdateNameAndAvatar(contactID uuid.UUID, name, avatarURL string) error {
	updates := false
	query := "UPDATE contacts SET updated_at = NOW()"
	args := []interface{}{}
	argIndex := 1
	if name != "" {
		query += ", name = $" + string(rune('0'+argIndex))
		args = append(args, name)
		argIndex++
		updates = true
	}
	if avatarURL != "" {
		query += ", avatar_url = $" + string(rune('0'+argIndex))
		args = append(args, avatarURL)
		argIndex++
		updates = true
	}
	if !updates {
		return nil
	}
	query += " WHERE id = $" + string(rune('0'+argIndex))
	args = append(args, contactID)
	_, err := r.db.Exec(query, args...)
	return err
}

func (r *contactRepository) Delete(orgID, contactID uuid.UUID) error {
	res, err := r.db.Exec(`DELETE FROM contacts WHERE id = $1 AND organization_id = $2`, contactID, orgID)
	if err != nil {
		return err
	}
	rows, _ := res.RowsAffected()
	if rows == 0 {
		return repository.ErrNotFound
	}
	return nil
}

func (r *contactRepository) GetConversations(orgID, contactID uuid.UUID) ([]*models.Conversation, error) {
	rows, err := r.db.Query(`
		SELECT c.id, c.status, c.last_message_at, c.created_at
		FROM conversations c
		WHERE c.organization_id = $1 AND c.contact_id = $2
		ORDER BY c.last_message_at DESC
	`, orgID, contactID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var conversations []*models.Conversation
	for rows.Next() {
		conv := &models.Conversation{}
		var lastMessageAt sql.NullTime
		if err := rows.Scan(&conv.ID, &conv.Status, &lastMessageAt, &conv.CreatedAt); err != nil {
			return nil, err
		}
		if lastMessageAt.Valid {
			conv.LastMessageAt = &lastMessageAt.Time
		}
		conversations = append(conversations, conv)
	}
	return conversations, nil
}
