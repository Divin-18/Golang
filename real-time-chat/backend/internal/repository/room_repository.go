package repository

import (
	"database/sql"
	"real-time-chat/internal/models"
)

type RoomRepository struct {
	db *sql.DB
}

func NewRoomRepository(db *sql.DB) *RoomRepository {
	return &RoomRepository{db: db}
}

func (r *RoomRepository) Create(room *models.Room) error {
	query := `
		INSERT INTO rooms (name, description, created_by, is_private)
		VALUES ($1, $2, $3, $4)
		RETURNING id, created_at
	`
	return r.db.QueryRow(query, room.Name, room.Description, room.CreatedBy, room.IsPrivate).
		Scan(&room.ID, &room.CreatedAt)
}

func (r *RoomRepository) GetByID(id int) (*models.Room, error) {
	room := &models.Room{}
	query := `
		SELECT id, name, description, created_by, is_private, created_at
		FROM rooms WHERE id = $1
	`
	err := r.db.QueryRow(query, id).Scan(
		&room.ID, &room.Name, &room.Description, &room.CreatedBy, &room.IsPrivate, &room.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return room, nil
}

func (r *RoomRepository) GetAll() ([]*models.Room, error) {
	query := `
		SELECT id, name, description, created_by, is_private, created_at
		FROM rooms WHERE is_private = false
		ORDER BY created_at DESC
	`
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var rooms []*models.Room
	for rows.Next() {
		room := &models.Room{}
		err := rows.Scan(
			&room.ID, &room.Name, &room.Description, &room.CreatedBy, &room.IsPrivate, &room.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		rooms = append(rooms, room)
	}
	return rooms, nil
}

func (r *RoomRepository) GetUserRooms(userID int) ([]*models.Room, error) {
	query := `
		SELECT r.id, r.name, r.description, r.created_by, r.is_private, r.created_at
		FROM rooms r
		INNER JOIN room_members rm ON r.id = rm.room_id
		WHERE rm.user_id = $1
		ORDER BY r.created_at DESC
	`
	rows, err := r.db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var rooms []*models.Room
	for rows.Next() {
		room := &models.Room{}
		err := rows.Scan(
			&room.ID, &room.Name, &room.Description, &room.CreatedBy, &room.IsPrivate, &room.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		rooms = append(rooms, room)
	}
	return rooms, nil
}

func (r *RoomRepository) AddMember(roomID, userID int) error {
	query := `
		INSERT INTO room_members (room_id, user_id)
		VALUES ($1, $2)
		ON CONFLICT (room_id, user_id) DO NOTHING
	`
	_, err := r.db.Exec(query, roomID, userID)
	return err
}

func (r *RoomRepository) RemoveMember(roomID, userID int) error {
	query := `DELETE FROM room_members WHERE room_id = $1 AND user_id = $2`
	_, err := r.db.Exec(query, roomID, userID)
	return err
}

func (r *RoomRepository) IsMember(roomID, userID int) (bool, error) {
	var exists bool
	query := `SELECT EXISTS(SELECT 1 FROM room_members WHERE room_id = $1 AND user_id = $2)`
	err := r.db.QueryRow(query, roomID, userID).Scan(&exists)
	return exists, err
}

func (r *RoomRepository) GetMembers(roomID int) ([]*models.User, error) {
	query := `
		SELECT u.id, u.username, u.email, u.avatar_url, u.is_online, u.created_at, u.updated_at
		FROM users u
		INNER JOIN room_members rm ON u.id = rm.user_id
		WHERE rm.room_id = $1
	`
	rows, err := r.db.Query(query, roomID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*models.User
	for rows.Next() {
		user := &models.User{}
		err := rows.Scan(
			&user.ID, &user.Username, &user.Email, &user.AvatarURL,
			&user.IsOnline, &user.CreatedAt, &user.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return users, nil
}

func (r *RoomRepository) Delete(id int) error {
	query := `DELETE FROM rooms WHERE id = $1`
	_, err := r.db.Exec(query, id)
	return err
}
