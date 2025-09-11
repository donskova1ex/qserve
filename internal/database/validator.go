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

	// Удаляем комментарии для проверки
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

	// Дополнительные проверки для SQL Server
	if q.containsSQLServerProcedures(upperQuery) {
		return fmt.Errorf("SQL Server system procedures not allowed")
	}

	return nil
}

// removeComments - удаляет комментарии из SQL запроса
func (q *QueryValidator) removeComments(query string) string {
	// Удаляем однострочные комментарии
	re := regexp.MustCompile(`--.*$`)
	query = re.ReplaceAllString(query, "")

	// Удаляем многострочные комментарии
	re = regexp.MustCompile(`/\*.*?\*/`)
	query = re.ReplaceAllString(query, "")

	return query
}

// containsWholeWord - проверяет наличие целого слова
func (q *QueryValidator) containsWholeWord(text, word string) bool {
	pattern := `\b` + regexp.QuoteMeta(word) + `\b`
	re := regexp.MustCompile(pattern)
	return re.MatchString(text)
}

// containsSQLServerProcedures - проверяет системные процедуры SQL Server
func (q *QueryValidator) containsSQLServerProcedures(query string) bool {
	procedures := []string{"sp_", "xp_", "fn_"}
	for _, proc := range procedures {
		if strings.Contains(query, proc) {
			return true
		}
	}
	return false
}

// GetQueryType - определяет тип SQL запроса
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
		// CTE (Common Table Expression) - обычно SELECT
		return "SELECT"
	default:
		return "UNKNOWN"
	}
}

// IsReadOnlyQuery - проверяет, является ли запрос только для чтения
func (q *QueryValidator) IsReadOnlyQuery(query string) bool {
	queryType := q.GetQueryType(query)
	return queryType == "SELECT"
}
