package services

import (
	"database/sql"
	"encoding/json"
	"errors"

	"github.com/google/uuid"
	"github.com/sidji-omnichannel/internal/models"
)

var (
	ErrContactNotFound = errors.New("contact not found")
)

// ContactService handles contact operations
type ContactService struct {
	db *sql.DB
}

// NewContactService creates a new contact service
func NewContactService(db *sql.DB) *ContactService {
	return &ContactService{db: db}
}

// List returns all contacts for an organization with pagination
func (s *ContactService) List(orgID uuid.UUID, page, limit int, search string) ([]*models.Contact, int, error) {
	if limit <= 0 || limit > 100 {
		limit = 20
	}
	if page <= 0 {
		page = 1
	}
	offset := (page - 1) * limit

	// Base query
	query := `
		SELECT id, organization_id, name, phone, email, avatar_url, whatsapp_id, instagram_id, facebook_id, created_at
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

	// Get total count
	var total int
	err := s.db.QueryRow(countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	// Add ordering and pagination
	query += ` ORDER BY created_at DESC LIMIT $` + string(rune('0'+len(args)+1)) + ` OFFSET $` + string(rune('0'+len(args)+2))
	args = append(args, limit, offset)

	rows, err := s.db.Query(query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var contacts []*models.Contact

	for rows.Next() {
		contact := &models.Contact{}
		var phone, email, avatarURL, whatsappID, instagramID, facebookID sql.NullString
		
		err := rows.Scan(
			&contact.ID, &contact.OrganizationID, &contact.Name, &phone,
			&email, &avatarURL, &whatsappID, &instagramID, &facebookID, &contact.CreatedAt,
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
		
		contacts = append(contacts, contact)
	}

	return contacts, total, nil
}

// GetByID retrieves a contact by ID
func (s *ContactService) GetByID(orgID, contactID uuid.UUID) (*models.Contact, error) {
	contact := &models.Contact{}
	var phone, email, avatarURL, whatsappID, instagramID, facebookID sql.NullString
	var metadata []byte
	
	err := s.db.QueryRow(`
		SELECT id, organization_id, name, phone, email, avatar_url, metadata, whatsapp_id, instagram_id, facebook_id, created_at, updated_at
		FROM contacts
		WHERE id = $1 AND organization_id = $2
	`, contactID, orgID).Scan(
		&contact.ID, &contact.OrganizationID, &contact.Name, &phone,
		&email, &avatarURL, &metadata, &whatsappID, &instagramID, &facebookID,
		&contact.CreatedAt, &contact.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, ErrContactNotFound
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
	
	return contact, nil
}

// Create creates a new contact
func (s *ContactService) Create(orgID uuid.UUID, input *models.CreateContactInput) (*models.Contact, error) {
	contact := &models.Contact{
		ID:             uuid.New(),
		OrganizationID: orgID,
		Name:           input.Name,
		Phone:          input.Phone,
		Email:          input.Email,
		WhatsAppID:     input.WhatsAppID,
		InstagramID:    input.InstagramID,
		FacebookID:     input.FacebookID,
	}

	_, err := s.db.Exec(`
		INSERT INTO contacts (id, organization_id, name, phone, email, whatsapp_id, instagram_id, facebook_id)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`, contact.ID, contact.OrganizationID, contact.Name, contact.Phone, contact.Email,
		contact.WhatsAppID, contact.InstagramID, contact.FacebookID)

	if err != nil {
		return nil, err
	}

	return contact, nil
}

// Update updates a contact
func (s *ContactService) Update(orgID, contactID uuid.UUID, input *models.UpdateContactInput) (*models.Contact, error) {
	// Build dynamic update query
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

	result, err := s.db.Exec(query, args...)
	if err != nil {
		return nil, err
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		return nil, ErrContactNotFound
	}

	return s.GetByID(orgID, contactID)
}

// Delete deletes a contact
func (s *ContactService) Delete(orgID, contactID uuid.UUID) error {
	result, err := s.db.Exec(`
		DELETE FROM contacts WHERE id = $1 AND organization_id = $2
	`, contactID, orgID)
	if err != nil {
		return err
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		return ErrContactNotFound
	}

	return nil
}

// FindOrCreateByWhatsAppID finds or creates a contact by WhatsApp ID
func (s *ContactService) FindOrCreateByWhatsAppID(orgID uuid.UUID, whatsAppID, name string) (*models.Contact, error) {
	// Try to find existing contact
	contact := &models.Contact{}
	var phoneNS, emailNS, avatarURLNS, whatsappIDNS, instagramIDNS sql.NullString
	
	err := s.db.QueryRow(`
		SELECT id, organization_id, name, phone, email, avatar_url, whatsapp_id, instagram_id, created_at
		FROM contacts
		WHERE organization_id = $1 AND whatsapp_id = $2
	`, orgID, whatsAppID).Scan(
		&contact.ID, &contact.OrganizationID, &contact.Name, &phoneNS,
		&emailNS, &avatarURLNS, &whatsappIDNS, &instagramIDNS, &contact.CreatedAt,
	)

	if err == nil {
		if phoneNS.Valid { contact.Phone = phoneNS.String }
		if emailNS.Valid { contact.Email = emailNS.String }
		if avatarURLNS.Valid { contact.AvatarURL = avatarURLNS.String }
		if whatsappIDNS.Valid { contact.WhatsAppID = whatsappIDNS.String }
		if instagramIDNS.Valid { contact.InstagramID = instagramIDNS.String }
		
		// Contact found, update name if changed
		if name != "" && contact.Name != name {
			s.db.Exec("UPDATE contacts SET name = $1, updated_at = NOW() WHERE id = $2", name, contact.ID)
			contact.Name = name
		}
		return contact, nil
	}

	if err != sql.ErrNoRows {
		return nil, err
	}

	// Create new contact
	contact = &models.Contact{
		ID:             uuid.New(),
		OrganizationID: orgID,
		Name:           name,
		Phone:          whatsAppID, // WhatsApp ID is usually the phone number
		WhatsAppID:     whatsAppID,
	}

	_, err = s.db.Exec(`
		INSERT INTO contacts (id, organization_id, name, phone, whatsapp_id)
		VALUES ($1, $2, $3, $4, $5)
	`, contact.ID, contact.OrganizationID, contact.Name, contact.Phone, contact.WhatsAppID)

	if err != nil {
		return nil, err
	}

	return contact, nil
}

// FindOrCreateByInstagramID finds or creates a contact by Instagram ID
func (s *ContactService) FindOrCreateByInstagramID(orgID uuid.UUID, instagramID, name, avatarURL string) (*models.Contact, error) {
	// Try to find existing contact
	contact := &models.Contact{}
	var phoneNS, emailNS, avatarURLNS, whatsappIDNS, instagramIDNS sql.NullString

	err := s.db.QueryRow(`
		SELECT id, organization_id, name, phone, email, avatar_url, whatsapp_id, instagram_id, created_at
		FROM contacts
		WHERE organization_id = $1 AND instagram_id = $2
	`, orgID, instagramID).Scan(
		&contact.ID, &contact.OrganizationID, &contact.Name, &phoneNS,
		&emailNS, &avatarURLNS, &whatsappIDNS, &instagramIDNS, &contact.CreatedAt,
	)

	if err == nil {
		if phoneNS.Valid { contact.Phone = phoneNS.String }
		if emailNS.Valid { contact.Email = emailNS.String }
		if avatarURLNS.Valid { contact.AvatarURL = avatarURLNS.String }
		if whatsappIDNS.Valid { contact.WhatsAppID = whatsappIDNS.String }
		if instagramIDNS.Valid { contact.InstagramID = instagramIDNS.String }

		// Contact found, update name/avatar if changed
		updates := false
		query := "UPDATE contacts SET updated_at = NOW()"
		args := []interface{}{}
		argIndex := 1
		
		if name != "" && contact.Name != name {
			query += ", name = $" + string(rune('0'+argIndex))
			args = append(args, name)
			argIndex++
			contact.Name = name
			updates = true
		}
		
		if avatarURL != "" && contact.AvatarURL != avatarURL {
			query += ", avatar_url = $" + string(rune('0'+argIndex))
			args = append(args, avatarURL)
			argIndex++
			contact.AvatarURL = avatarURL
			updates = true
		}

		if updates {
			query += " WHERE id = $" + string(rune('0'+argIndex))
			args = append(args, contact.ID)
			s.db.Exec(query, args...)
		}
		
		return contact, nil
	}

	if err != sql.ErrNoRows {
		return nil, err
	}

	// Create new contact
	contact = &models.Contact{
		ID:             uuid.New(),
		OrganizationID: orgID,
		Name:           name,
		InstagramID:    instagramID,
		AvatarURL:      avatarURL,
	}

	_, err = s.db.Exec(`
		INSERT INTO contacts (id, organization_id, name, instagram_id, avatar_url)
		VALUES ($1, $2, $3, $4, $5)
	`, contact.ID, contact.OrganizationID, contact.Name, contact.InstagramID, contact.AvatarURL)

	if err != nil {
		return nil, err
	}

	return contact, nil
}

// FindOrCreateByFacebookID finds or creates a contact by Facebook ID
func (s *ContactService) FindOrCreateByFacebookID(orgID uuid.UUID, facebookID, name, avatarURL string) (*models.Contact, error) {
	// Try to find existing contact
	contact := &models.Contact{}
	var phoneNS, emailNS, avatarURLNS, whatsappIDNS, instagramIDNS, facebookIDNS sql.NullString

	err := s.db.QueryRow(`
		SELECT id, organization_id, name, phone, email, avatar_url, whatsapp_id, instagram_id, facebook_id, created_at
		FROM contacts
		WHERE organization_id = $1 AND facebook_id = $2
	`, orgID, facebookID).Scan(
		&contact.ID, &contact.OrganizationID, &contact.Name, &phoneNS,
		&emailNS, &avatarURLNS, &whatsappIDNS, &instagramIDNS, &facebookIDNS, &contact.CreatedAt,
	)

	if err == nil {
		if phoneNS.Valid { contact.Phone = phoneNS.String }
		if emailNS.Valid { contact.Email = emailNS.String }
		if avatarURLNS.Valid { contact.AvatarURL = avatarURLNS.String }
		if whatsappIDNS.Valid { contact.WhatsAppID = whatsappIDNS.String }
		if instagramIDNS.Valid { contact.InstagramID = instagramIDNS.String }
		if facebookIDNS.Valid { contact.FacebookID = facebookIDNS.String }

		// Contact found, update name/avatar if changed
		updates := false
		query := "UPDATE contacts SET updated_at = NOW()"
		args := []interface{}{}
		argIndex := 1
		
		if name != "" && contact.Name != name {
			query += ", name = $" + string(rune('0'+argIndex))
			args = append(args, name)
			argIndex++
			contact.Name = name
			updates = true
		}
		
		if avatarURL != "" && contact.AvatarURL != avatarURL {
			query += ", avatar_url = $" + string(rune('0'+argIndex))
			args = append(args, avatarURL)
			argIndex++
			contact.AvatarURL = avatarURL
			updates = true
		}

		if updates {
			query += " WHERE id = $" + string(rune('0'+argIndex))
			args = append(args, contact.ID)
			s.db.Exec(query, args...)
		}
		
		return contact, nil
	}

	if err != sql.ErrNoRows {
		return nil, err
	}

	// Create new contact
	contact = &models.Contact{
		ID:             uuid.New(),
		OrganizationID: orgID,
		Name:           name,
		FacebookID:     facebookID,
		AvatarURL:      avatarURL,
	}

	_, err = s.db.Exec(`
		INSERT INTO contacts (id, organization_id, name, facebook_id, avatar_url)
		VALUES ($1, $2, $3, $4, $5)
	`, contact.ID, contact.OrganizationID, contact.Name, contact.FacebookID, contact.AvatarURL)

	if err != nil {
		return nil, err
	}

	return contact, nil
}

// GetConversations returns all conversations for a contact
func (s *ContactService) GetConversations(orgID, contactID uuid.UUID) ([]*models.Conversation, error) {
	rows, err := s.db.Query(`
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
		
		err := rows.Scan(&conv.ID, &conv.Status, &lastMessageAt, &conv.CreatedAt)
		if err != nil {
			return nil, err
		}
		
		if lastMessageAt.Valid {
			conv.LastMessageAt = &lastMessageAt.Time
		}
		
		conversations = append(conversations, conv)
	}

	return conversations, nil
}
