package models

// ChartType represents the type of chart
type ChartType string

const (
	ChartTypePie  ChartType = "pie"
	ChartTypeBar  ChartType = "bar"
	ChartTypeLine ChartType = "line"
)

// Chart represents a chart configuration saved by a user
type Chart struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Type      ChartType `json:"type"`
	Config    string    `json:"config"` // JSON string storing chart-specific configuration
	Query     string    `json:"query"`  // Query parameters used to get the data
	UserID    int       `json:"user_id"`
	CreatedAt string    `json:"created_at"`
	UpdatedAt string    `json:"updated_at"`
}

// DataTable represents a saved data table configuration
type DataTable struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	Columns   string `json:"columns"` // JSON string storing column configurations
	Query     string `json:"query"`   // Query parameters used to get the data
	UserID    int    `json:"user_id"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}
