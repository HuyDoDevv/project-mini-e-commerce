package pgx

import (
	"context"
	"fmt"
	"project-mini-e-commerce/pkg/logger"
	"reflect"
	"regexp"
	"strings"
	"time"

	"github.com/jackc/pgx/v5/tracelog"
	"github.com/rs/zerolog"
)

type ZerlogTrace struct {
	Logger         zerolog.Logger
	SlowQueryLimit time.Duration
}

type QueryInfo struct {
	QueryName     string
	OperationType string
	CleanQuery    string
	OriginalSQL   string
}

var (
	sqlcNameRegex = regexp.MustCompile(`-- name:\s*(\w+)\s*:(\w+)`)
	spaceRegex    = regexp.MustCompile(`\s+`)
	commentRegex  = regexp.MustCompile(`-- [^\r\n]*`)
)

func parseSQL(sql string) QueryInfo {
	info := QueryInfo{
		OriginalSQL: sql,
	}
	if matches := sqlcNameRegex.FindStringSubmatch(sql); len(matches) == 3 {
		info.QueryName = matches[1]
		info.OperationType = strings.ToUpper(matches[2])
	}
	cleanSQL := commentRegex.ReplaceAllString(sql, "")
	cleanSQL = strings.TrimSpace(cleanSQL)
	cleanSQL = spaceRegex.ReplaceAllString(cleanSQL, " ")
	info.CleanQuery = cleanSQL

	return info
}

func formatArgs(args any) string {
	val := reflect.ValueOf(args)

	if args == nil || val.Kind() == reflect.Ptr && val.IsNil() {
		return "NULL"
	}

	if val.Kind() == reflect.Ptr {
		args = val.Elem().Interface()
	}

	switch v := args.(type) {
	case string:
		return fmt.Sprintf("'%s'", strings.ReplaceAll(v, "'", "''"))
	case bool:
		return fmt.Sprintf("%t", v)
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
		return fmt.Sprintf("%d", v)
	case float32, float64:
		return fmt.Sprintf("%f", v)
	case time.Time:
		return fmt.Sprintf("'%s'", v.Format("2006-01-02 15:04:05"))
	case nil:
		return "NULL"
	default:
		return fmt.Sprintf("'%s'", strings.ReplaceAll(fmt.Sprintf("%v", v), "'", "''"))
	}
}

func replacePlaceholders(sql string, agrs []any) string {
	for i, arg := range agrs {
		replacer := fmt.Sprintf("$%d", i+1)
		sql = strings.ReplaceAll(sql, replacer, formatArgs(arg))
	}
	return sql
}

func (t *ZerlogTrace) Log(ctx context.Context, level tracelog.LogLevel, msg string, data map[string]any) {
	sql, _ := data["sql"].(string)
	agrs, _ := data["args"].([]any)
	duration, _ := data["duration"].(time.Duration)

	queryInfo := parseSQL(sql)

	var finalSQL string
	if len(agrs) > 0 {
		finalSQL = replacePlaceholders(queryInfo.CleanQuery, agrs)
	} else {
		finalSQL = queryInfo.CleanQuery
	}

	baseLoger := t.Logger.With().
		Str("trace_id", logger.GetTraceId(ctx)).
		Dur("duration", duration).
		Str("sql_original", queryInfo.OriginalSQL).
		Str("sql", finalSQL).
		Str("query_name", queryInfo.QueryName).
		Str("clean_query", queryInfo.CleanQuery).
		Str("operation", queryInfo.OperationType).
		Interface("args", agrs)

	finalLogger := baseLoger.Logger()
	if msg == "Query" && duration > t.SlowQueryLimit {
		finalLogger.Warn().Str("event", "Slow query").Msg("Slow SQL query")
		return
	}

	if msg == "Query" {
		finalLogger.Info().Str("event", "Query").Msg("Executed SQL query")
		return
	}
}
