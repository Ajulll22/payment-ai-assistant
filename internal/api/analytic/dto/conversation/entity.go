package conversation

import (
	"time"

	"github.com/Ajulll22/payment-ai-assistant/internal/api/analytic/schema"
)

type ConversationEntity struct {
	RoomID    string
	UserID    int
	State     ConversationState
	CreatedAt *time.Time
	UpdatedAt *time.Time
}

type ConversationState struct {
	CurrentTopic         string     `json:"current_topic"`
	CurrentIntent        string     `json:"current_intent"`
	LastUserMessage      string     `json:"last_user_message"`
	LastAssistantMessage string     `json:"last_assistant_message"`
	LastActivityAt       *time.Time `json:"last_activity_at"`

	AnalyticalContext    AnalyticalContext    `json:"analytical_context"`
	VisualizationContext VisualizationContext `json:"visualization_context"`
	ExecuteContext       ExecuteContext       `json:"execute_context"`
}

type AnalyticalContext struct {
	LastSQL                    string               `json:"last_sql"`
	LastLimit                  int                  `json:"last_limit"`
	LastMetric                 string               `json:"last_metric"`
	LastDimension              string               `json:"last_dimension"`
	LastTimeRange              string               `json:"last_time_range"`
	LastGroupBy                []string             `json:"last_group_by"`
	LastOrderBy                []string             `json:"last_order_by"`
	LastRelevantTableSchema    []schema.TableSchema `json:"last_relevant_table_schema"`
	LastRelevantRelationScehma []schema.Relation    `json:"last_relevant_relation_schema"`
	LastDataSet                DatasetAnalytical    `json:"last_dataset"`
}

type VisualizationContext struct {
	LastTable *TableVisualization `json:"last_table"`
	LastChart *ChartVisualization `json:"last_chart"`
}

type ExecuteContext struct {
	LastAction              string `json:"last_action"`
	LastToolUsed            string `json:"last_tool_used"`
	LastExecutionDurationMs int64  `json:"last_execution_duration_ms"`
	LastError               string `json:"last_error"`
	SQLExecuted             bool   `json:"sql_executed"`
	UsedPreviousDataset     bool   `json:"used_previous_dataset"`
}

type ChartVisualization struct {
	ChartType string `json:"chart_type"`
	Title     string `json:"title"`
	XAxis     string `json:"x_axis"`
	YAxis     string `json:"y_axis"`
}

type TableVisualization struct {
	Columns []ColumnTable `json:"columns"`
}

type ColumnTable struct {
	Field string `json:"field"`
	Label string `json:"label"`
}

type DatasetAnalytical struct {
	DatasetID int              `json:"dataset_id"`
	RowCount  int              `json:"row_count"`
	Preview   []map[string]any `json:"preview"`
	Columns   []string         `json:"columns"`
}
