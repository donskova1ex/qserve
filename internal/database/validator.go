package database

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
)

type QueryValidator struct{}

func NewQueryValidator() *QueryValidator {
	return &QueryValidator{}
}

func (q *QueryValidator) ValidateQuery(query string) error {
	if strings.TrimSpace(query) == "" {
		return errors.New("query is empty")
	}

	cleanQuery := q.removeComments(strings.TrimSpace(query))

	dangerousKeywords := []string{
		"DROP", "TRUNCATE", "ALTER", "CREATE", "GRANT", "REVOKE",
		"EXEC", "EXECUTE", "DELETE", "UPDATE", "INSERT",
	}
	upperQuery := strings.ToUpper(cleanQuery)
	for _, keyword := range dangerousKeywords {
		if q.containsWholeWord(upperQuery, keyword) {
			return fmt.Errorf("dangerous keyword found: %s", keyword)
		}
	}
	if q.containsSQLServerProcedures(upperQuery) {
		return fmt.Errorf("SQL Server system procedures not allowed")
	}

	return nil
}

func (q *QueryValidator) removeComments(query string) string {
	re := regexp.MustCompile(`--.*$`)
	query = re.ReplaceAllString(query, "")

	re = regexp.MustCompile(`/\*.*?\*/`)
	query = re.ReplaceAllString(query, "")

	return query
}

func (q *QueryValidator) containsWholeWord(text, word string) bool {
	pattern := `\b` + regexp.QuoteMeta(word) + `\b`
	re := regexp.MustCompile(pattern)
	return re.MatchString(text)
}

func (q *QueryValidator) containsSQLServerProcedures(query string) bool {
	procedures := []string{"sp_", "xp_", "fn_"}
	for _, proc := range procedures {
		if strings.Contains(query, proc) {
			return true
		}
	}
	return false
}

func (q *QueryValidator) GetQueryType(query string) string {
	words := strings.Fields(strings.TrimSpace(query))
	if len(words) == 0 {
		return "UNKNOWN"
	}

	firstWord := strings.ToUpper(words[0])

	switch firstWord {
	case "SELECT":
		return "SELECT"
	case "INSERT":
		return "INSERT"
	case "UPDATE":
		return "UPDATE"
	case "DELETE":
		return "DELETE"
	case "WITH":
		// TODO: Додумать как обрабатывать CTE (Common Table Expression)
		return "SELECT"
	default:
		return "UNKNOWN"
	}
}
