package assistant

type ChatAssistantRequest struct {
	RoomID string `json:"" binding:"required"`
	Prompt string `json:"" binding:"required"`
}
