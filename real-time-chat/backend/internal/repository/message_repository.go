package repository

import (
	"database/sql"
	"real-time-chat/internal/models"
)

type MessageRepository struct {
	db *sql.DB
}

func NewMessageRepository(db *sql.DB) *MessageRepository {
	return &MessageRepository{db: db}
}

func (r *MessageRepository) Create(message *models.Message) error {
	query := `
		INSERT INTO messages (room_id, user_id, content, message_type)
		VALUES ($1, $2, $3, $4)
		RETURNING id, created_at
	`
	return r.db.QueryRow(query, message.RoomID, message.UserID, message.Content, message.MessageType).
		Scan(&message.ID, &message.CreatedAt)
}

func (r *MessageRepository) GetByRoomID(roomID, limit, offset int) ([]*models.Message, error) {
	query := `
		SELECT m.id, m.room_id, m.user_id, COALESCE(u.username, 'Deleted User') as username, 
		       m.content, m.message_type, m.created_at
		FROM messages m
		LEFT JOIN users u ON m.user_id = u.id
		WHERE m.room_id = $1
		ORDER BY m.created_at DESC
		LIMIT $2 OFFSET $3
	`
	rows, err := r.db.Query(query, roomID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messages []*models.Message
	for rows.Next() {
		message := &models.Message{}
		err := rows.Scan(
			&message.ID, &message.RoomID, &message.UserID, &message.Username,
			&message.Content, &message.MessageType, &message.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		messages = append(messages, message)
	}

	// Reverse to get chronological order
	for i, j := 0, len(messages)-1; i < j; i, j = i+1, j-1 {
		messages[i], messages[j] = messages[j], messages[i]
	}

	return messages, nil
}

func (r *MessageRepository) GetLatestByRoomID(roomID, limit int) ([]*models.Message, error) {
	return r.GetByRoomID(roomID, limit, 0)
}

func (r *MessageRepository) Delete(id int) error {
	query := `DELETE FROM messages WHERE id = $1`
	_, err := r.db.Exec(query, id)
	return err
}
