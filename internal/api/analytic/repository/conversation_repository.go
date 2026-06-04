package repository

import (
	"encoding/json"
	"time"

	"github.com/Ajulll22/payment-ai-assistant/internal/api/analytic/dto/conversation"
	conversationdataset "github.com/Ajulll22/payment-ai-assistant/internal/api/analytic/dto/conversation_dataset"
	"github.com/Ajulll22/payment-ai-assistant/pkg/constant"
	"github.com/Ajulll22/payment-ai-assistant/pkg/handling"
	"gorm.io/gorm"
)

type ConversationRepository interface {
	GetConversationByRoomID(db *gorm.DB, roomID string) (conversation.ConversationEntity, error)
	GetConversationDatasetByDatasetID(db *gorm.DB, datasetID string) (conversationdataset.ConversationDatasetEntity, error)
	UpsertConversation(db *gorm.DB, roomID string, userID int, stateJSON conversation.ConversationState) error
	UpsertConversationDataset(db *gorm.DB, datasetID string, roomID string, dataJSON []map[string]any) error
}

type conversationRepository struct{}

func NewConversationRepository() ConversationRepository {
	return &conversationRepository{}
}

func (r *conversationRepository) GetConversationByRoomID(db *gorm.DB, roomID string) (conversation.ConversationEntity, error) {
	m := conversation.ConversationEntity{}
	raw := conversationSPResult{}

	query := db.Raw("exec spEP_ai_conversation_data ?", roomID).Scan(&raw)
	if query.Error != nil {
		return m, handling.NewErrorWrapper(constant.CodeDBError, "Error execute sp", nil, query.Error)
	}

	m = conversation.ConversationEntity{
		RoomID:    raw.RoomID,
		UserID:    raw.UserID,
		CreatedAt: raw.CreatedAt,
		UpdatedAt: raw.UpdatedAt,
	}

	if raw.StateJSON != "" {
		json.Unmarshal([]byte(raw.StateJSON), &m.State)
	}

	return m, nil
}

func (r *conversationRepository) GetConversationDatasetByDatasetID(db *gorm.DB, datasetID string) (conversationdataset.ConversationDatasetEntity, error) {
	m := conversationdataset.ConversationDatasetEntity{}
	raw := conversationDatasetSPResult{}

	query := db.Raw("exec spEP_ai_conversation_dataset_data ?", datasetID).Scan(&raw)
	if query.Error != nil {
		return m, handling.NewErrorWrapper(constant.CodeDBError, "Error execute sp", nil, query.Error)
	}

	m = conversationdataset.ConversationDatasetEntity{
		DatasetID: raw.DatasetID,
		RoomID:    raw.RoomID,
		CreatedAt: raw.CreatedAt,
	}

	if raw.DataJSON != "" {
		json.Unmarshal([]byte(raw.DataJSON), &m.Data)
	}

	return m, nil
}

func (r *conversationRepository) UpsertConversation(db *gorm.DB, roomID string, userID int, stateJSON conversation.ConversationState) error {
	stateString, err := json.Marshal(stateJSON)
	if err != nil {
		return handling.NewErrorWrapper(constant.CodeDBError, "Error marshal model", nil, err)
	}

	query := db.Exec("exec spEP_ai_conversation_data_upsert ?, ?, ?", roomID, userID, string(stateString))
	if query.Error != nil {
		return handling.NewErrorWrapper(constant.CodeDBError, "Error execute sp", nil, query.Error)
	}

	return nil
}

func (r *conversationRepository) UpsertConversationDataset(db *gorm.DB, datasetID string, roomID string, dataJSON []map[string]any) error {
	dataString, err := json.Marshal(dataJSON)
	if err != nil {
		return handling.NewErrorWrapper(constant.CodeDBError, "Error marshal model", nil, err)
	}

	query := db.Exec("exec spEP_ai_conversation_dataset_data_upsert ?, ?, ?", datasetID, roomID, string(dataString))
	if query.Error != nil {
		return handling.NewErrorWrapper(constant.CodeDBError, "Error execute sp", nil, query.Error)
	}

	return nil
}

type conversationSPResult struct {
	RoomID    string     `gorm:"column:room_id"`
	UserID    int        `gorm:"column:user_id"`
	StateJSON string     `gorm:"column:state_json"`
	CreatedAt *time.Time `gorm:"column:created_at"`
	UpdatedAt *time.Time `gorm:"column:updated_at"`
}

type conversationDatasetSPResult struct {
	DatasetID string     `gorm:"column:dataset_id"`
	RoomID    string     `gorm:"column:room_id"`
	DataJSON  string     `gorm:"column:data_json"`
	CreatedAt *time.Time `gorm:"column:created_at"`
}
