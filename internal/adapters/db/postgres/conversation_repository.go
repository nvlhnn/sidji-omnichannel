package postgres

import (
	"database/sql"

	"github.com/google/uuid"
	"github.com/sidji-omnichannel/internal/domain/ports/repository"
	"github.com/sidji-omnichannel/internal/models"
)

type conversationRepository struct {
	db *sql.DB
}

// NewConversationRepository creates a new postgres conversation repository
func NewConversationRepository(db *sql.DB) repository.ConversationRepository {
	return &conversationRepository{db: db}
}

func (r *conversationRepository) List(orgID uuid.UUID, filter *models.ConversationFilter) ([]*models.ConversationListItem, int, error) {
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
	if err := r.db.QueryRow(countQuery, args...).Scan(&total); err != nil {
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

	rows, err := r.db.Query(query, args...)
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

func (r *conversationRepository) GetByID(orgID, convID uuid.UUID) (*models.Conversation, error) {
	conv := &models.Conversation{
		Channel: &models.ChannelPublic{},
		Contact: &models.Contact{},
	}

	var assignedTo sql.NullString
	var lastMsgAt, lastHumanReplyAt sql.NullTime
	var subject, phone, email, avatarURL, whatsappID, instagramID sql.NullString
	err := r.db.QueryRow(`
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
		return nil, nil // Error handling should be done in service
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

func (r *conversationRepository) Assign(orgID, convID, userID uuid.UUID) error {
	res, err := r.db.Exec(`
		UPDATE conversations SET assigned_to = $1, updated_at = NOW()
		WHERE id = $2 AND organization_id = $3
	`, userID, convID, orgID)
	if err != nil {
		return err
	}
	rows, _ := res.RowsAffected()
	if rows == 0 {
		return repository.ErrNotFound
	}
	return nil
}

func (r *conversationRepository) UpdateStatus(orgID, convID uuid.UUID, status models.ConversationStatus) error {
	res, err := r.db.Exec(`
		UPDATE conversations SET status = $1, updated_at = NOW()
		WHERE id = $2 AND organization_id = $3
	`, status, convID, orgID)
	if err != nil {
		return err
	}
	rows, _ := res.RowsAffected()
	if rows == 0 {
		return repository.ErrNotFound
	}
	return nil
}

func (r *conversationRepository) FindOpen(orgID, channelID, contactID uuid.UUID) (uuid.UUID, error) {
	var id uuid.UUID
	err := r.db.QueryRow(`
		SELECT id FROM conversations 
		WHERE organization_id = $1 AND channel_id = $2 AND contact_id = $3 
		AND status NOT IN ('closed')
		ORDER BY created_at DESC LIMIT 1
	`, orgID, channelID, contactID).Scan(&id)
	return id, err
}

func (r *conversationRepository) Create(conv *models.Conversation) error {
	_, err := r.db.Exec(`
		INSERT INTO conversations (id, organization_id, channel_id, contact_id, status)
		VALUES ($1, $2, $3, $4, $5)
	`, conv.ID, conv.OrganizationID, conv.ChannelID, conv.ContactID, conv.Status)
	return err
}

func (r *conversationRepository) Delete(orgID, convID uuid.UUID) error {
	res, err := r.db.Exec(`
		DELETE FROM conversations 
		WHERE id = $1 AND organization_id = $2
	`, convID, orgID)
	if err != nil {
		return err
	}
	rows, _ := res.RowsAffected()
	if rows == 0 {
		return repository.ErrNotFound
	}
	return nil
}

func (r *conversationRepository) UpdateLastMessage(convID uuid.UUID, lastMessageAt interface{}) error {
	_, err := r.db.Exec(`
		UPDATE conversations SET last_message_at = $1, updated_at = NOW()
		WHERE id = $2
	`, lastMessageAt, convID)
	return err
}

func (r *conversationRepository) UpdateLastHumanReply(convID uuid.UUID, lastReplyAt interface{}) error {
	_, err := r.db.Exec(`
		UPDATE conversations SET last_human_reply_at = $1, ai_paused_until = NULL, updated_at = NOW()
		WHERE id = $2
	`, lastReplyAt, convID)
	return err
}
