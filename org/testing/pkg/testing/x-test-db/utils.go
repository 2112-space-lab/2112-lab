package xtestdb

import (
	"bytes"
	"fmt"
	"log/slog"
	"os"
	"regexp"
	"strings"
	"text/template"

	db_models "github.com/org/2112-space-lab/org/testing/pkg/testing/x-test-db/models"
)

var regexpDatabaseNameChars = regexp.MustCompile(`[^a-zA-Z0-9_]+`)

func normalizeDatabaseName(
	logger *slog.Logger,
	dbName db_models.PotentialDatabaseName,
) db_models.DatabaseName {
	truncated := dbName
	if len(truncated) > dbNameMaxLength {
		truncated = dbName[:dbNameMaxLength]
		logger.Warn("dbName is too long - will be truncated to dbNameMaxLength",
			slog.Group("databaseName",
				slog.String("originalDbName", string(dbName)),
				slog.String("truncated", string(truncated)),
				slog.Int("dbNameMaxLength", dbNameMaxLength),
			),
		)
	}

	lowered := strings.ToLower(string(dbName))
	normalized := regexpDatabaseNameChars.ReplaceAllString(lowered, "")
	logger.Info("database name has been normalized",
		slog.Group("databaseName",
			slog.String("originalDbName", string(dbName)),
			slog.String("normalized", string(truncated)),
			slog.Int("dbNameMaxLength", dbNameMaxLength),
		),
	)
	return db_models.DatabaseName(normalized)
}

func getSQLFileFromTemplate(filePath string, placeholders map[string]string) ([]byte, error) {
	var buf bytes.Buffer
	f, err := os.ReadFile(filePath)
	if err != nil {
		return buf.Bytes(), fmt.Errorf("cannot read file [%s] - [%w]", filePath, err)
	}
	scr := string(f)
	templ, err := template.New("scr").Parse(scr)
	if err != nil {
		return buf.Bytes(), fmt.Errorf("failed to parse template from file [%s] - [%w]", filePath, err)
	}
	err = templ.Execute(&buf, placeholders)
	if err != nil {
		return buf.Bytes(), fmt.Errorf("start/stop time could not be computed")
	}
	return buf.Bytes(), nil

}
