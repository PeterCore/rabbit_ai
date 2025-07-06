package model

import (
	"database/sql"
	"errors"
	"time"
)

// ErrConversationNotFound 对话未找到错误
var ErrConversationNotFound = errors.New("conversation not found")

// ErrMessageNotFound 消息未找到错误
var ErrMessageNotFound = errors.New("message not found")

// Conversation 对话会话模型
type Conversation struct {
	ID            int64     `json:"id" db:"id"`
	UserID        int64     `json:"user_id" db:"user_id"`
	Title         string    `json:"title" db:"title"`                     // 对话标题
	Status        int       `json:"status" db:"status"`                   // 1: 活跃, 0: 已删除
	MessageCount  int       `json:"message_count" db:"message_count"`     // 消息数量
	LastMessageAt time.Time `json:"last_message_at" db:"last_message_at"` // 最后消息时间
	CreatedAt     time.Time `json:"created_at" db:"created_at"`
	UpdatedAt     time.Time `json:"updated_at" db:"updated_at"`
}

// Message 消息模型
type Message struct {
	ID             int64     `json:"id" db:"id"`
	ConversationID int64     `json:"conversation_id" db:"conversation_id"`
	Role           string    `json:"role" db:"role"`                   // user/assistant
	Content        string    `json:"content" db:"content"`             // 消息内容
	Tokens         int       `json:"tokens" db:"tokens"`               // 消耗的token数量
	Model          string    `json:"model" db:"model"`                 // 使用的模型
	FinishReason   string    `json:"finish_reason" db:"finish_reason"` // 结束原因
	CreatedAt      time.Time `json:"created_at" db:"created_at"`
}

// ConversationRepository 对话数据访问接口
type ConversationRepository interface {
	Create(conversation *Conversation) error
	GetByID(id int64) (*Conversation, error)
	GetByUserID(userID int64, limit, offset int) ([]*Conversation, error)
	Update(conversation *Conversation) error
	Delete(id int64) error
	GetUserConversationCount(userID int64) (int, error)
}

// MessageRepository 消息数据访问接口
type MessageRepository interface {
	Create(message *Message) error
	GetByID(id int64) (*Message, error)
	GetByConversationID(conversationID int64, limit, offset int) ([]*Message, error)
	GetConversationMessages(conversationID int64) ([]*Message, error)
	Update(message *Message) error
	Delete(id int64) error
	GetConversationMessageCount(conversationID int64) (int, error)
}

// ConversationRepositoryImpl 对话数据访问实现
type ConversationRepositoryImpl struct {
	db *sql.DB
}

// MessageRepositoryImpl 消息数据访问实现
type MessageRepositoryImpl struct {
	db *sql.DB
}

// NewConversationRepository 创建对话仓库实例
func NewConversationRepository(db *sql.DB) ConversationRepository {
	return &ConversationRepositoryImpl{db: db}
}

// NewMessageRepository 创建消息仓库实例
func NewMessageRepository(db *sql.DB) MessageRepository {
	return &MessageRepositoryImpl{db: db}
}

// Create 创建对话
func (r *ConversationRepositoryImpl) Create(conversation *Conversation) error {
	query := `
		INSERT INTO conversations (user_id, title, status, message_count, last_message_at, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id`

	now := time.Now()
	conversation.CreatedAt = now
	conversation.UpdatedAt = now
	conversation.LastMessageAt = now

	return r.db.QueryRow(
		query,
		conversation.UserID,
		conversation.Title,
		conversation.Status,
		conversation.MessageCount,
		conversation.LastMessageAt,
		conversation.CreatedAt,
		conversation.UpdatedAt,
	).Scan(&conversation.ID)
}

// GetByID 根据ID获取对话
func (r *ConversationRepositoryImpl) GetByID(id int64) (*Conversation, error) {
	conversation := &Conversation{}
	query := `
		SELECT id, user_id, title, status, message_count, last_message_at, created_at, updated_at
		FROM conversations WHERE id = $1 AND status = 1`

	err := r.db.QueryRow(query, id).Scan(
		&conversation.ID,
		&conversation.UserID,
		&conversation.Title,
		&conversation.Status,
		&conversation.MessageCount,
		&conversation.LastMessageAt,
		&conversation.CreatedAt,
		&conversation.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrConversationNotFound
		}
		return nil, err
	}

	return conversation, nil
}

// GetByUserID 根据用户ID获取对话列表
func (r *ConversationRepositoryImpl) GetByUserID(userID int64, limit, offset int) ([]*Conversation, error) {
	query := `
		SELECT id, user_id, title, status, message_count, last_message_at, created_at, updated_at
		FROM conversations 
		WHERE user_id = $1 AND status = 1
		ORDER BY last_message_at DESC
		LIMIT $2 OFFSET $3`

	rows, err := r.db.Query(query, userID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var conversations []*Conversation
	for rows.Next() {
		conversation := &Conversation{}
		err := rows.Scan(
			&conversation.ID,
			&conversation.UserID,
			&conversation.Title,
			&conversation.Status,
			&conversation.MessageCount,
			&conversation.LastMessageAt,
			&conversation.CreatedAt,
			&conversation.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		conversations = append(conversations, conversation)
	}

	return conversations, nil
}

// Update 更新对话
func (r *ConversationRepositoryImpl) Update(conversation *Conversation) error {
	query := `
		UPDATE conversations 
		SET title = $1, status = $2, message_count = $3, last_message_at = $4, updated_at = $5
		WHERE id = $6`

	conversation.UpdatedAt = time.Now()

	result, err := r.db.Exec(
		query,
		conversation.Title,
		conversation.Status,
		conversation.MessageCount,
		conversation.LastMessageAt,
		conversation.UpdatedAt,
		conversation.ID,
	)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrConversationNotFound
	}

	return nil
}

// Delete 删除对话（软删除）
func (r *ConversationRepositoryImpl) Delete(id int64) error {
	query := `UPDATE conversations SET status = 0, updated_at = $1 WHERE id = $2`

	result, err := r.db.Exec(query, time.Now(), id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrConversationNotFound
	}

	return nil
}

// GetUserConversationCount 获取用户对话数量
func (r *ConversationRepositoryImpl) GetUserConversationCount(userID int64) (int, error) {
	query := `SELECT COUNT(*) FROM conversations WHERE user_id = $1 AND status = 1`

	var count int
	err := r.db.QueryRow(query, userID).Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
}

// Create 创建消息
func (r *MessageRepositoryImpl) Create(message *Message) error {
	query := `
		INSERT INTO messages (conversation_id, role, content, tokens, model, finish_reason, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id`

	message.CreatedAt = time.Now()

	return r.db.QueryRow(
		query,
		message.ConversationID,
		message.Role,
		message.Content,
		message.Tokens,
		message.Model,
		message.FinishReason,
		message.CreatedAt,
	).Scan(&message.ID)
}

// GetByID 根据ID获取消息
func (r *MessageRepositoryImpl) GetByID(id int64) (*Message, error) {
	message := &Message{}
	query := `
		SELECT id, conversation_id, role, content, tokens, model, finish_reason, created_at
		FROM messages WHERE id = $1`

	err := r.db.QueryRow(query, id).Scan(
		&message.ID,
		&message.ConversationID,
		&message.Role,
		&message.Content,
		&message.Tokens,
		&message.Model,
		&message.FinishReason,
		&message.CreatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrMessageNotFound
		}
		return nil, err
	}

	return message, nil
}

// GetByConversationID 根据对话ID获取消息列表
func (r *MessageRepositoryImpl) GetByConversationID(conversationID int64, limit, offset int) ([]*Message, error) {
	query := `
		SELECT id, conversation_id, role, content, tokens, model, finish_reason, created_at
		FROM messages 
		WHERE conversation_id = $1
		ORDER BY created_at ASC
		LIMIT $2 OFFSET $3`

	rows, err := r.db.Query(query, conversationID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messages []*Message
	for rows.Next() {
		message := &Message{}
		err := rows.Scan(
			&message.ID,
			&message.ConversationID,
			&message.Role,
			&message.Content,
			&message.Tokens,
			&message.Model,
			&message.FinishReason,
			&message.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		messages = append(messages, message)
	}

	return messages, nil
}

// GetConversationMessages 获取对话的所有消息
func (r *MessageRepositoryImpl) GetConversationMessages(conversationID int64) ([]*Message, error) {
	query := `
		SELECT id, conversation_id, role, content, tokens, model, finish_reason, created_at
		FROM messages 
		WHERE conversation_id = $1
		ORDER BY created_at ASC`

	rows, err := r.db.Query(query, conversationID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messages []*Message
	for rows.Next() {
		message := &Message{}
		err := rows.Scan(
			&message.ID,
			&message.ConversationID,
			&message.Role,
			&message.Content,
			&message.Tokens,
			&message.Model,
			&message.FinishReason,
			&message.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		messages = append(messages, message)
	}

	return messages, nil
}

// Update 更新消息
func (r *MessageRepositoryImpl) Update(message *Message) error {
	query := `
		UPDATE messages 
		SET role = $1, content = $2, tokens = $3, model = $4, finish_reason = $5
		WHERE id = $6`

	result, err := r.db.Exec(
		query,
		message.Role,
		message.Content,
		message.Tokens,
		message.Model,
		message.FinishReason,
		message.ID,
	)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrMessageNotFound
	}

	return nil
}

// Delete 删除消息
func (r *MessageRepositoryImpl) Delete(id int64) error {
	query := `DELETE FROM messages WHERE id = $1`

	result, err := r.db.Exec(query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrMessageNotFound
	}

	return nil
}

// GetConversationMessageCount 获取对话消息数量
func (r *MessageRepositoryImpl) GetConversationMessageCount(conversationID int64) (int, error) {
	query := `SELECT COUNT(*) FROM messages WHERE conversation_id = $1`

	var count int
	err := r.db.QueryRow(query, conversationID).Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
}
