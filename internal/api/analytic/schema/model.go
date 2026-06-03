package schema

type Column struct {
	Name               string   `json:"name"`
	Type               string   `json:"type"`
	BusinessName       string   `json:"business_name"`
	Description        string   `json:"description"`
	Keywords           []string `json:"-"`
	SemanticHints      []string `json:"semantic_hints,omitempty"`
	CommonAggregations []string `json:"common_aggregations,omitempty"`
	IsMetric           bool     `json:"is_metric,omitempty"`
	IsDimension        bool     `json:"is_dimension,omitempty"`
}

type RelationCondition struct {
	FromColumn string `json:"from_column"`
	ToColumn   string `json:"to_column"`
}

type Relation struct {
	FromTable   string              `json:"from_table"`
	ToTable     string              `json:"to_table"`
	Conditions  []RelationCondition `json:"conditions"`
	Description string              `json:"description"`
}

type TableSchema struct {
	TableName   string   `json:"table_name"`
	Description string   `json:"description"`
	Keywords    []string `json:"-"`
	Columns     []Column `json:"columns"`
}

type DatabaseSchema struct {
	Tables    []TableSchema `json:"tables"`
	Relations []Relation    `json:"relations"`
}

type RetrievalResult struct {
	Tables []TableSchema `json:"tables"`

	Relations []Relation `json:"relations"`
}
