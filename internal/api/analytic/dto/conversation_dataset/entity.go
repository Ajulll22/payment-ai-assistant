package conversationdataset

import "time"

type ConversationDatasetEntity struct {
	DatasetID string
	RoomID    string
	CreatedAt *time.Time
	Data      []map[string]any
}
