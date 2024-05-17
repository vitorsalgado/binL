package main

import (
	"bytes"
	"context"
	"log/slog"

	"github.com/go-mysql-org/go-mysql/canal"
	"github.com/go-mysql-org/go-mysql/mysql"
	"github.com/go-mysql-org/go-mysql/replication"
)

var _ canal.EventHandler = (*EventHandler)(nil)

type EventHandler struct {
	logger *slog.Logger
}

func (e *EventHandler) OnDDL(header *replication.EventHeader, nextPos mysql.Position, queryEvent *replication.QueryEvent) error {
	e.logger.Info("OnDDL",
		slog.String("event_type", header.EventType.String()),
		slog.String("statement", string(queryEvent.Query)),
		slog.String("pos", nextPos.String()))

	return nil
}

func (e *EventHandler) OnGTID(header *replication.EventHeader, gtid mysql.GTIDSet) error {
	e.logger.Info("OnGTID",
		slog.String("event_type", header.EventType.String()),
		slog.String("gtid", gtid.String()))

	return nil
}

func (e *EventHandler) OnPosSynced(header *replication.EventHeader, pos mysql.Position, set mysql.GTIDSet, force bool) error {
	gtid := ""
	if set != nil {
		gtid = set.String()
	}

	e.logger.Info("OnPosSynced",
		slog.String("event_type", header.EventType.String()),
		slog.String("pos", pos.String()),
		slog.String("set", gtid),
		slog.Bool("force", force))

	return nil
}

func (e *EventHandler) OnRawEvent(event *replication.BinlogEvent) error {
	buf := bytes.NewBuffer(make([]byte, 0, len(event.RawData)))
	event.Dump(buf)

	e.logger.Info("OnRawEvent",
		slog.String("event_type", event.Header.EventType.String()),
		slog.String("dump", buf.String()))

	return nil
}

func (e *EventHandler) OnRotate(header *replication.EventHeader, r *replication.RotateEvent) error {
	e.logger.Info("OnRotate",
		slog.String("event_type", header.EventType.String()),
		slog.Uint64("pos", r.Position))

	return nil
}

func (e *EventHandler) OnRow(evt *canal.RowsEvent) error {
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

func (e *EventHandler) OnTableChanged(header *replication.EventHeader, schema string, table string) error {
	e.logger.Info("OnTableChanged",
		slog.String("event_type", header.EventType.String()),
		slog.String("table", table),
		slog.String("schema", schema))

	return nil
}

func (e *EventHandler) OnUnmarshal(data []byte) (interface{}, error) {
	e.logger.Info("OnUnmarshal",
		slog.String("data", string(data)))

	return nil, nil
}

func (e *EventHandler) OnXID(header *replication.EventHeader, pos mysql.Position) error {
	e.logger.Info("OnXID",
		slog.String("event_type", header.EventType.String()),
		slog.String("pos", pos.String()))

	return nil
}

func (e *EventHandler) String() string {
	return "binL.EventHandler"
}
