package services

import (
	"database/sql"
	"errors"

	"github.com/google/uuid"
	"github.com/sidji-omnichannel/internal/models"
)

var (
	ErrConversationNotFound = errors.New("conversation not found")
)

// ConversationService handles conversation operations
type ConversationService struct {
	db *sql.DB
}

// NewConversationService creates a new conversation service
func NewConversationService(db *sql.DB) *ConversationService {
	return &ConversationService{db: db}
}

// List returns conversations for an organization with filters
func (s *ConversationService) List(orgID uuid.UUID, filter *models.ConversationFilter) ([]*models.ConversationListItem, int, error) {
	// Base query
	query := `
		SELECT 
			c.id, c.status, c.last_message_at,
			ch.id, ch.type, ch.name, ch.status,
			ct.id, ct.name, ct.phone, ct.avatar_url, ct.whatsapp_id, ct.instagram_id,
			u.id, u.name, u.avatar_url, u.status,
			(SELECT COUNT(*) FROM messages m WHERE m.conversation_id = c.id AND m.sender_type = 'contact' AND m.status != 'read') as unread_count,
			lm.content, lm.sender_type, lm.message_type
		FROM conversations c
		JOIN channels ch ON ch.id = c.channel_id
		JOIN contacts ct ON ct.id = c.contact_id
		LEFT JOIN users u ON u.id = c.assigned_to
		LEFT JOIN messages lm ON lm.conversation_id = c.id AND lm.created_at = c.last_message_at
		WHERE c.organization_id = $1
	`

	args := []interface{}{orgID}
	argIndex := 2

	// ... (Filters logic remains same, adjust query concatenation if needed, but it's fine as query string is built before execution)
	
	// Copy-paste the filters part to ensure it uses the modified 'query' variable correctly if it was split. 
	// The original code appended to 'query'. So as long as I start with the full SELECT...FROM...JOIN, the appends work.

	// Apply filters
	if filter.Status != "" {
		query += ` AND c.status = $` + string(rune('0'+argIndex))
		args = append(args, filter.Status)
		argIndex++
	}

	if filter.ChannelID != uuid.Nil {
		query += ` AND c.channel_id = $` + string(rune('0'+argIndex))
		args = append(args, filter.ChannelID)
		argIndex++
	}

	if filter.AssignedTo != uuid.Nil {
		query += ` AND c.assigned_to = $` + string(rune('0'+argIndex))
		args = append(args, filter.AssignedTo)
		argIndex++
	}

	if filter.Unassigned {
		query += ` AND c.assigned_to IS NULL`
	}

	// Count total
	countQuery := `SELECT COUNT(*) FROM (` + query + `) as count_query`
	var total int
	if err := s.db.QueryRow(countQuery, args...).Scan(&total); err != nil {
		return nil, 0, err
	}

	// Add ordering and pagination
	query += ` ORDER BY c.last_message_at DESC NULLS LAST`
	
	limit := filter.Limit
	if limit <= 0 || limit > 100 {
		limit = 20
	}
	offset := (filter.Page - 1) * limit
	if offset < 0 {
		offset = 0
	}

	query += ` LIMIT $` + string(rune('0'+argIndex)) + ` OFFSET $` + string(rune('0'+argIndex+1))
	args = append(args, limit, offset)

	rows, err := s.db.Query(query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	conversations := []*models.ConversationListItem{}
	for rows.Next() {
		item := &models.ConversationListItem{
			Channel: &models.ChannelPublic{},
			Contact: &models.Contact{},
		}
		var assignedUser struct {
			ID        sql.NullString
			Name      sql.NullString
			AvatarURL sql.NullString
			Status    sql.NullString
		}

		var phone, avatarURL, whatsappID, instagramID sql.NullString
		var lastMsgContent, lastMsgSenderType, lastMsgType sql.NullString
		
		err := rows.Scan(
			&item.ID, &item.Status, &item.LastMessageAt,
			&item.Channel.ID, &item.Channel.Type, &item.Channel.Name, &item.Channel.Status,
			&item.Contact.ID, &item.Contact.Name, &phone, &avatarURL,
			&whatsappID, &instagramID,
			&assignedUser.ID, &assignedUser.Name, &assignedUser.AvatarURL, &assignedUser.Status,
			&item.UnreadCount,
			&lastMsgContent, &lastMsgSenderType, &lastMsgType,
		)
		if err != nil {
			return nil, 0, err
		}

		if phone.Valid { item.Contact.Phone = phone.String }
		if avatarURL.Valid { item.Contact.AvatarURL = avatarURL.String }
		if whatsappID.Valid { item.Contact.WhatsAppID = whatsappID.String }
		if instagramID.Valid { item.Contact.InstagramID = instagramID.String }

		if assignedUser.ID.Valid {
			item.AssignedUser = &models.UserPublic{
				ID:        uuid.MustParse(assignedUser.ID.String),
				Name:      assignedUser.Name.String,
				AvatarURL: assignedUser.AvatarURL.String,
				Status:    models.UserStatus(assignedUser.Status.String),
			}
		}

		if lastMsgContent.Valid {
			item.LastMessage = &models.Message{
				Content:     lastMsgContent.String,
				SenderType:  models.SenderType(lastMsgSenderType.String),
				MessageType: models.MessageType(lastMsgType.String),
			}
		}

		conversations = append(conversations, item)
	}

	return conversations, total, nil
}

// GetByID retrieves a single conversation
func (s *ConversationService) GetByID(orgID, convID uuid.UUID) (*models.Conversation, error) {
	conv := &models.Conversation{
		Channel: &models.ChannelPublic{},
		Contact: &models.Contact{},
	}

	var assignedTo sql.NullString
	var lastMsgAt, lastHumanReplyAt sql.NullTime
	var subject, phone, email, avatarURL, whatsappID, instagramID sql.NullString
	err := s.db.QueryRow(`
		SELECT 
			c.id, c.organization_id, c.channel_id, c.contact_id, c.assigned_to, 
			c.status, c.subject, c.last_message_at, c.last_human_reply_at, c.created_at, c.updated_at,
			ch.id, ch.type, ch.name, ch.status,
			ct.id, ct.name, ct.phone, ct.email, ct.avatar_url, ct.whatsapp_id, ct.instagram_id,
			(SELECT COUNT(*) FROM messages m WHERE m.conversation_id = c.id AND m.sender_type = 'contact' AND m.status != 'read') as unread_count
		FROM conversations c
		JOIN channels ch ON ch.id = c.channel_id
		JOIN contacts ct ON ct.id = c.contact_id
		WHERE c.id = $1 AND c.organization_id = $2
	`, convID, orgID).Scan(
		&conv.ID, &conv.OrganizationID, &conv.ChannelID, &conv.ContactID, &assignedTo,
		&conv.Status, &subject, &lastMsgAt, &lastHumanReplyAt, &conv.CreatedAt, &conv.UpdatedAt,
		&conv.Channel.ID, &conv.Channel.Type, &conv.Channel.Name, &conv.Channel.Status,
		&conv.Contact.ID, &conv.Contact.Name, &phone, &email,
		&avatarURL, &whatsappID, &instagramID,
		&conv.UnreadCount,
	)

	if err == sql.ErrNoRows {
		return nil, ErrConversationNotFound
	}
	if err != nil {
		return nil, err
	}
	
	if lastMsgAt.Valid {
		conv.LastMessageAt = &lastMsgAt.Time
	}
	if lastHumanReplyAt.Valid {
		conv.LastHumanReplyAt = &lastHumanReplyAt.Time
	}
	
	if subject.Valid {
		conv.Subject = subject.String
	}

	if phone.Valid { conv.Contact.Phone = phone.String }
	if email.Valid { conv.Contact.Email = email.String }
	if avatarURL.Valid { conv.Contact.AvatarURL = avatarURL.String }
	if whatsappID.Valid { conv.Contact.WhatsAppID = whatsappID.String }
	if instagramID.Valid { conv.Contact.InstagramID = instagramID.String }

	if assignedTo.Valid {
		id := uuid.MustParse(assignedTo.String)
		conv.AssignedTo = &id
	}

	return conv, nil
}

// Assign assigns a conversation to an agent
func (s *ConversationService) Assign(orgID, convID, userID uuid.UUID) error {
	result, err := s.db.Exec(`
		UPDATE conversations SET assigned_to = $1, updated_at = NOW()
		WHERE id = $2 AND organization_id = $3
	`, userID, convID, orgID)

	if err != nil {
		return err
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		return ErrConversationNotFound
	}

	return nil
}

// UpdateStatus updates conversation status
func (s *ConversationService) UpdateStatus(orgID, convID uuid.UUID, status models.ConversationStatus) error {
	result, err := s.db.Exec(`
		UPDATE conversations SET status = $1, updated_at = NOW()
		WHERE id = $2 AND organization_id = $3
	`, status, convID, orgID)

	if err != nil {
		return err
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		return ErrConversationNotFound
	}

	return nil
}

// FindOrCreate finds or creates a conversation for a contact on a channel
func (s *ConversationService) FindOrCreate(orgID, channelID, contactID uuid.UUID) (*models.Conversation, error) {
	// Try to find existing open conversation
	var convID uuid.UUID
	err := s.db.QueryRow(`
		SELECT id FROM conversations 
		WHERE organization_id = $1 AND channel_id = $2 AND contact_id = $3 
		AND status NOT IN ('closed')
		ORDER BY created_at DESC LIMIT 1
	`, orgID, channelID, contactID).Scan(&convID)

	if err == nil {
		return s.GetByID(orgID, convID)
	}

	if err != sql.ErrNoRows {
		return nil, err
	}

	// Create new conversation
	newID := uuid.New()
	_, err = s.db.Exec(`
		INSERT INTO conversations (id, organization_id, channel_id, contact_id, status)
		VALUES ($1, $2, $3, $4, 'open')
	`, newID, orgID, channelID, contactID)

	if err != nil {
		return nil, err
	}

	return s.GetByID(orgID, newID)
}

// Delete permanently removes a conversation
func (s *ConversationService) Delete(orgID, convID uuid.UUID) error {
	result, err := s.db.Exec(`
		DELETE FROM conversations 
		WHERE id = $1 AND organization_id = $2
	`, convID, orgID)

	if err != nil {
		return err
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		return ErrConversationNotFound
	}

	return nil
}
