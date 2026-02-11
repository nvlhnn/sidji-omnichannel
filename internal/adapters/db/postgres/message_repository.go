package postgres

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
	"github.com/sidji-omnichannel/internal/domain/ports/repository"
	"github.com/sidji-omnichannel/internal/models"
)

type messageRepository struct {
	db *sql.DB
}

// NewMessageRepository creates a new postgres message repository
func NewMessageRepository(db *sql.DB) repository.MessageRepository {
	return &messageRepository{db: db}
}

func (r *messageRepository) List(convID uuid.UUID, filter *models.MessageFilter) (*models.MessageList, error) {
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

	rows, err := r.db.Query(query, args...)
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
	err = r.db.QueryRow("SELECT COUNT(*) FROM messages WHERE conversation_id = $1", convID).Scan(&total)
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

func (r *messageRepository) Create(msg *models.Message) error {
	if msg.ID == uuid.Nil {
		msg.ID = uuid.New()
	}
	if msg.CreatedAt.IsZero() {
		msg.CreatedAt = time.Now()
	}
	if msg.UpdatedAt.IsZero() {
		msg.UpdatedAt = time.Now()
	}

	_, err := r.db.Exec(`
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

	return err
}

func (r *messageRepository) UpdateStatus(externalID string, status models.MessageStatus) error {
	_, err := r.db.Exec(`
		UPDATE messages SET status = $1, updated_at = NOW()
		WHERE external_id = $2
	`, status, externalID)
	return err
}

func (r *messageRepository) UpdateStatusByID(id uuid.UUID, status models.MessageStatus) error {
	_, err := r.db.Exec(`
		UPDATE messages SET status = $1, updated_at = NOW()
		WHERE id = $2
	`, status, id)
	return err
}

func (r *messageRepository) MarkAsRead(convID uuid.UUID) error {
	_, err := r.db.Exec(`
		UPDATE messages SET status = 'read', updated_at = NOW()
		WHERE conversation_id = $1 AND sender_type = 'contact' AND status != 'read'
	`, convID)
	return err
}

func (r *messageRepository) CountUnread(convID uuid.UUID) (int, error) {
	var count int
	err := r.db.QueryRow(`
		SELECT COUNT(*) FROM messages 
		WHERE conversation_id = $1 AND sender_type = 'contact' AND status != 'read'
	`, convID).Scan(&count)
	return count, err
}

func (r *messageRepository) GetByExternalID(externalID string) (*models.Message, error) {
	msg := &models.Message{}
	var mediaURL, replyToID sql.NullString

	err := r.db.QueryRow(`
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
