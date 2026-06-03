package conversation

import (
	"time"

	"github.com/Ajulll22/payment-ai-assistant/internal/api/analytic/schema"
)

type ConversationEntity struct {
	RoomID string
	UserID int
	State  ConversationState
}

type ConversationState struct {
	CurrentTopic         string
	CurrentIntent        string
	LastUserMessage      string
	LastAssistantMessage map[string]any
	LastActivityAt       time.Time
}

type AnalyticalContext struct {
	LastSQL         string
	LastLimit       int
	LastGroupBy     []string
	LastOrderBy     []string
	LastDataSet     []map[string]any
	LastTableSchema []schema.TableSchema
}

type VisualizationContext struct {
	LastChartType string
}
