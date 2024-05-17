package main

import (
	"context"
	"log/slog"

	"github.com/go-mysql-org/go-mysql/canal"
	"github.com/go-mysql-org/go-mysql/mysql"
	"github.com/go-mysql-org/go-mysql/replication"
)

var _ canal.EventHandler = (*RowOnlyEventHandler)(nil)

type RowOnlyEventHandler struct {
	logger *slog.Logger
}

func (e *RowOnlyEventHandler) OnDDL(header *replication.EventHeader, nextPos mysql.Position, queryEvent *replication.QueryEvent) error {
	return nil
}

func (e *RowOnlyEventHandler) OnGTID(header *replication.EventHeader, gtid mysql.BinlogGTIDEvent) error {
	return nil
}

func (e *RowOnlyEventHandler) OnPosSynced(header *replication.EventHeader, pos mysql.Position, set mysql.GTIDSet, force bool) error {
	return nil
}

func (e *RowOnlyEventHandler) OnRawEvent(event *replication.BinlogEvent) error {
	return nil
}

func (e *RowOnlyEventHandler) OnRotate(header *replication.EventHeader, r *replication.RotateEvent) error {
	return nil
}

func (e *RowOnlyEventHandler) OnRow(evt *canal.RowsEvent) error {
	ctx := context.Background()
	attr := make([]slog.Attr, 0)
	attr = append(attr,
		slog.String("event_type", evt.Header.EventType.String()),
		slog.Any("rows", evt.Rows),
		slog.String("action", evt.Action))

	if evt.Table == nil {
		e.logger.LogAttrs(ctx, slog.LevelInfo, "OnRow", attr...)
		return nil
	}

	columns := make([]slog.Attr, 0, len(evt.Table.Columns)*2)
	for _, v := range evt.Table.Columns {
		columns = append(columns,
			slog.String("name", v.Name),
			slog.String("raw_type", v.RawType))
	}

	attr = append(attr,
		slog.String("table", evt.Table.String()),
		slog.Group("columns", slog.Attr{Value: slog.GroupValue(columns...)}))

	e.logger.LogAttrs(ctx, slog.LevelInfo, "OnRow", attr...)

	return nil
}

func (e *RowOnlyEventHandler) OnTableChanged(header *replication.EventHeader, schema string, table string) error {
	e.logger.Info("OnTableChanged",
		slog.String("event_type", header.EventType.String()),
		slog.String("table", table),
		slog.String("schema", schema))

	return nil
}

func (e *RowOnlyEventHandler) OnUnmarshal(data []byte) (interface{}, error) {
	e.logger.Info("OnUnmarshal",
		slog.String("data", string(data)))

	return nil, nil
}

func (e *RowOnlyEventHandler) OnXID(header *replication.EventHeader, pos mysql.Position) error {
	return nil
}

func (e *RowOnlyEventHandler) OnRowsQueryEvent(evt *replication.RowsQueryEvent) error {
	return nil
}

func (e *RowOnlyEventHandler) String() string {
	return "binL.RowOnlyEventHandler"
}
