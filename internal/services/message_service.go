package services

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
	"github.com/sidji-omnichannel/internal/models"
)

// MessageService handles message operations
type MessageService struct {
	db *sql.DB
}

// NewMessageService creates a new message service
func NewMessageService(db *sql.DB) *MessageService {
	return &MessageService{db: db}
}

// List returns messages for a conversation
func (s *MessageService) List(convID uuid.UUID, filter *models.MessageFilter) (*models.MessageList, error) {
	query := `
		SELECT 
			m.id, m.conversation_id, m.sender_type, m.sender_id, m.content,
			m.message_type, m.media_url, m.media_mime_type, m.media_file_name,
			m.external_id, m.reply_to_id, m.status, m.created_at, m.updated_at
		FROM messages m
		WHERE m.conversation_id = $1
	`
	args := []interface{}{convID}
	argIndex := 2

	if filter.Before != nil {
		query += ` AND m.created_at < $` + string(rune('0'+argIndex))
		args = append(args, filter.Before)
		argIndex++
	}

	if filter.After != nil {
		query += ` AND m.created_at > $` + string(rune('0'+argIndex))
		args = append(args, filter.After)
		argIndex++
	}

	query += ` ORDER BY m.created_at DESC`

	limit := filter.Limit
	if limit <= 0 || limit > 100 {
		limit = 50
	}
	query += ` LIMIT $` + string(rune('0'+argIndex))
	args = append(args, limit+1) // +1 to check if there are more

	rows, err := s.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	messages := []*models.Message{}
	for rows.Next() {
		msg := &models.Message{}
		var mediaURL, mimeType, fileName, externalID, replyToID sql.NullString

		err := rows.Scan(
			&msg.ID, &msg.ConversationID, &msg.SenderType, &msg.SenderID, &msg.Content,
			&msg.MessageType, &mediaURL, &mimeType, &fileName,
			&externalID, &replyToID, &msg.Status, &msg.CreatedAt, &msg.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		if mediaURL.Valid { msg.MediaURL = mediaURL.String }
		if mimeType.Valid { msg.MediaMimeType = mimeType.String }
		if fileName.Valid { msg.MediaFileName = fileName.String }
		if externalID.Valid { msg.ExternalID = externalID.String }
		if replyToID.Valid {
			id := uuid.MustParse(replyToID.String)
			msg.ReplyToID = &id
		}

		messages = append(messages, msg)
	}

	hasMore := len(messages) > limit
	if hasMore {
		messages = messages[:limit]
	}

	// Get total count
	var total int
	err = s.db.QueryRow("SELECT COUNT(*) FROM messages WHERE conversation_id = $1", convID).Scan(&total)
	if err != nil {
		return nil, err
	}

	// Reverse to show oldest first
	for i, j := 0, len(messages)-1; i < j; i, j = i+1, j-1 {
		messages[i], messages[j] = messages[j], messages[i]
	}

	return &models.MessageList{
		Messages:   messages,
		TotalCount: total,
		HasMore:    hasMore,
	}, nil
}

// Create creates a new message
func (s *MessageService) Create(msg *models.Message) error {
	msg.ID = uuid.New()
	msg.CreatedAt = time.Now()
	msg.UpdatedAt = time.Now()

	_, err := s.db.Exec(`
		INSERT INTO messages (
			id, conversation_id, sender_type, sender_id, content,
			message_type, media_url, media_mime_type, media_file_name,
			external_id, reply_to_id, status, created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)
	`,
		msg.ID, msg.ConversationID, msg.SenderType, msg.SenderID, msg.Content,
		msg.MessageType, msg.MediaURL, msg.MediaMimeType, msg.MediaFileName,
		msg.ExternalID, msg.ReplyToID, msg.Status, msg.CreatedAt, msg.UpdatedAt,
	)

	if err != nil {
		return err
	}

	// Update conversation's last_message_at
	_, err = s.db.Exec(`
		UPDATE conversations SET last_message_at = $1, updated_at = NOW()
		WHERE id = $2
	`, msg.CreatedAt, msg.ConversationID)

	// If sender is AGENT, update last_human_reply_at
	if msg.SenderType == models.SenderAgent {
		_, err = s.db.Exec(`
			UPDATE conversations SET last_human_reply_at = $1, ai_paused_until = NULL, updated_at = NOW()
			WHERE id = $2
		`, msg.CreatedAt, msg.ConversationID)
	}

	return err
}

// UpdateStatus updates message status (sent, delivered, read) by ExternalID
func (s *MessageService) UpdateStatus(externalID string, status models.MessageStatus) error {
	_, err := s.db.Exec(`
		UPDATE messages SET status = $1, updated_at = NOW()
		WHERE external_id = $2
	`, status, externalID)
	return err
}

// UpdateStatusByID updates message status by internal UUID
func (s *MessageService) UpdateStatusByID(id uuid.UUID, status models.MessageStatus) error {
	_, err := s.db.Exec(`
		UPDATE messages SET status = $1, updated_at = NOW()
		WHERE id = $2
	`, status, id)
	return err
}

// MarkAsRead marks all contact messages in a conversation as read
func (s *MessageService) MarkAsRead(convID uuid.UUID) error {
	_, err := s.db.Exec(`
		UPDATE messages SET status = 'read', updated_at = NOW()
		WHERE conversation_id = $1 AND sender_type = 'contact' AND status != 'read'
	`, convID)
	return err
}

// CountUnread counts unread messages from contact in a conversation
func (s *MessageService) CountUnread(convID uuid.UUID) (int, error) {
	var count int
	err := s.db.QueryRow(`
		SELECT COUNT(*) FROM messages 
		WHERE conversation_id = $1 AND sender_type = 'contact' AND status != 'read'
	`, convID).Scan(&count)
	return count, err
}

// GetByExternalID finds a message by its external ID (WhatsApp/Instagram message ID)
func (s *MessageService) GetByExternalID(externalID string) (*models.Message, error) {
	msg := &models.Message{}
	var mediaURL, replyToID sql.NullString

	err := s.db.QueryRow(`
		SELECT 
			id, conversation_id, sender_type, sender_id, content,
			message_type, media_url, status, created_at
		FROM messages WHERE external_id = $1
	`, externalID).Scan(
		&msg.ID, &msg.ConversationID, &msg.SenderType, &msg.SenderID, &msg.Content,
		&msg.MessageType, &mediaURL, &msg.Status, &msg.CreatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	if mediaURL.Valid { msg.MediaURL = mediaURL.String }

	if replyToID.Valid {
		id := uuid.MustParse(replyToID.String)
		msg.ReplyToID = &id
	}

	return msg, nil
}

// GetUserPublic fetches public info for a user
func (s *MessageService) GetUserPublic(userID uuid.UUID) (*models.UserPublic, error) {
	user := &models.UserPublic{}
	var avatarURL sql.NullString
	err := s.db.QueryRow(`
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
