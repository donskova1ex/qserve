package database

import (
	"context"
	"fmt"
)

func (c *ConnectionManager) ExecuteQuery(ctx context.Context, query string) ([]map[string]interface{}, error) {
	if c.db == nil {
		return nil, fmt.Errorf("database not connected")
	}

	rows, err := c.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}
	defer rows.Close()

	cols, err := rows.Columns()
	if err != nil {
		return nil, fmt.Errorf("failed to get columns: %w", err)
	}

	var result []map[string]interface{}
	values := make([]interface{}, len(cols))
	valuesPtrs := make([]interface{}, len(cols))

	for i := range values {
		valuesPtrs[i] = &values[i]
	}

	for rows.Next() {
		if err := rows.Scan(valuesPtrs...); err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}

		row := make(map[string]interface{})
		for i, col := range cols {
			val := values[i]
			if val != nil {
				switch v := val.(type) {
				case []byte:
					row[col] = string(v)
				default:
					row[col] = v
				}
			} else {
				row[col] = nil
			}
		}

		result = append(result, row)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to iterate over rows: %w", err)
	}

	return result, nil
}

func (c *ConnectionManager) ExecuteTransaction(ctx context.Context, query string) (int64, error) {
	if c.db == nil {
		return 0, fmt.Errorf("database not connected")
	}

	result, err := c.db.ExecContext(ctx, query)
	if err != nil {
		return 0, fmt.Errorf("failed to execute transaction: %w", err)
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("failed to get rows affected: %w", err)
	}

	return rowsAffected, nil
}
