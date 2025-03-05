package services

import (
	"fmt"
	"my-cucumber-backend/models"
)

// CreateChart creates a new chart configuration
func CreateChart(chart *models.Chart) error {
	result, err := DB.Exec(
		`INSERT INTO charts (name, type, config, query, user_id) 
		 VALUES (?, ?, ?, ?, ?)`,
		chart.Name, chart.Type, chart.Config, chart.Query,
		chart.UserID,
	)
	if err != nil {
		return fmt.Errorf("failed to create chart: %v", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("failed to get last insert id: %v", err)
	}
	chart.ID = int(id)
	return nil
}

// GetChartsByUser retrieves all charts for a user
func GetChartsByUser(userID int) ([]models.Chart, error) {
	rows, err := DB.Query(
		`SELECT id, name, type, config, query, user_id, 
		 created_at, updated_at FROM charts WHERE user_id = ?`,
		userID,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to query charts: %v", err)
	}
	defer rows.Close()

	var charts []models.Chart
	for rows.Next() {
		var chart models.Chart
		err := rows.Scan(
			&chart.ID, &chart.Name, &chart.Type, &chart.Config,
			&chart.Query, &chart.UserID,
			&chart.CreatedAt, &chart.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan chart: %v", err)
		}
		charts = append(charts, chart)
	}
	return charts, nil
}

// CreateDataTable creates a new data table configuration
func CreateDataTable(table *models.DataTable) error {
	result, err := DB.Exec(
		`INSERT INTO data_tables (name, columns, query, user_id) 
		 VALUES (?, ?, ?, ?)`,
		table.Name, table.Columns, table.Query, table.UserID,
	)
	if err != nil {
		return fmt.Errorf("failed to create data table: %v", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("failed to get last insert id: %v", err)
	}
	table.ID = int(id)
	return nil
}

// GetDataTablesByUser retrieves all data tables for a user
func GetDataTablesByUser(userID int) ([]models.DataTable, error) {
	rows, err := DB.Query(
		`SELECT id, name, columns, query, user_id, 
		 created_at, updated_at FROM data_tables WHERE user_id = ?`,
		userID,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to query data tables: %v", err)
	}
	defer rows.Close()

	var tables []models.DataTable
	for rows.Next() {
		var table models.DataTable
		err := rows.Scan(
			&table.ID, &table.Name, &table.Columns, &table.Query,
			&table.UserID, &table.CreatedAt, &table.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan data table: %v", err)
		}
		tables = append(tables, table)
	}
	return tables, nil
}
