package repository

import "gorm.io/gorm"

type ConversationRepository interface {
}

type conversationRepository struct{}

func NewConversationRepository() ConversationRepository {
	return &conversationRepository{}
}

func GetConversationByRoomID(db *gorm.DB, roomID string) {

}
