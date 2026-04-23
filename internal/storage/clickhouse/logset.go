package clickhouse

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/ClickHouse/clickhouse-go/v2/lib/driver"
	"github.com/google/uuid"
)

type LogsetStorage struct {
	ch driver.Conn
}

func NewLogsetStorage(ch driver.Conn) *LogsetStorage {
	return &LogsetStorage{ch: ch}
}

func (s *LogsetStorage) Get(ctx context.Context, logsetID uuid.UUID) ([]map[string]any, error) {
	rows, err := s.ch.Query(ctx, "SELECT raw FROM logs WHERE logset_id = $1", logsetID)
	if err != nil {
		return nil, fmt.Errorf("query logs: %w", err)
	}
	defer rows.Close()

	var logs []map[string]any
	for rows.Next() {
		var raw string
		if err := rows.Scan(&raw); err != nil {
			return nil, fmt.Errorf("scan row: %w", err)
		}
		var obj map[string]any
		if err := json.Unmarshal([]byte(raw), &obj); err != nil {
			return nil, fmt.Errorf("unmarshal log: %w", err)
		}
		logs = append(logs, obj)
	}

	return logs, nil
}

func (s *LogsetStorage) Insert(ctx context.Context, logsetID uuid.UUID, logs []map[string]any) error {
	batch, err := s.ch.PrepareBatch(ctx, "INSERT INTO logs (logset_id, raw)")
	if err != nil {
		return fmt.Errorf("prepare batch: %w", err)
	}

	for _, log := range logs {
		raw, err := json.Marshal(log)
		if err != nil {
			return fmt.Errorf("marshal log: %w", err)
		}
		if err := batch.Append(logsetID, string(raw)); err != nil {
			return fmt.Errorf("append to batch: %w", err)
		}
	}

	if err := batch.Send(); err != nil {
		return fmt.Errorf("send batch: %w", err)
	}

	return nil
}
