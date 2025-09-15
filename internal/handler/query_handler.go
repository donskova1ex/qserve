package handler

import (
	"encoding/json"
	"net/http"
	"qserve/internal/database"
	"strings"
)

type QueryRequest struct {
	Query string `json:"query"`
}

type QueryResponse struct {
	Data         []map[string]interface{} `json:"data,omitempty"`
	RowsAffected int64                    `json:"rows_affected,omitempty"`
	Error        string                   `json:"error,omitempty"`
	Status       string                   `json:"status"`
}

type QueryHandler struct {
	dbManager *database.ConnectionManager
	validator *database.QueryValidator
}

func NewQueryHandler(dbManager *database.ConnectionManager, validator *database.QueryValidator) *QueryHandler {
	return &QueryHandler{
		dbManager: dbManager,
		validator: validator,
	}
}

func (h *QueryHandler) HandleQuery(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var req QueryRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(QueryResponse{
			Status: "error",
			Error:  "Invalid request body",
		})
		return
	}

	if strings.TrimSpace(req.Query) == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(QueryResponse{
			Status: "error",
			Error:  "Query is empty",
		})
		return
	}

	if err := h.validator.ValidateQuery(req.Query); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(QueryResponse{
			Status: "error",
			Error:  err.Error(),
		})
		return
	}

	queryType := h.validator.GetQueryType(req.Query)
	if queryType == "UNKNOWN" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(QueryResponse{
			Status: "error",
			Error:  "Unknown query type",
		})
		return
	}

	switch queryType {
	case "SELECT":
		result, err := h.dbManager.ExecuteQuery(r.Context(), req.Query)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(QueryResponse{
				Status: "error",
				Error:  err.Error(),
			})
			return
		}
		json.NewEncoder(w).Encode(QueryResponse{
			Status: "success",
			Data:   result,
		})
	case "INSERT", "UPDATE", "DELETE":
		rowsAffected, err := h.dbManager.ExecuteTransaction(r.Context(), req.Query)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(QueryResponse{
				Status: "error",
				Error:  err.Error(),
			})
			return
		}
		json.NewEncoder(w).Encode(QueryResponse{
			Status:       "success",
			RowsAffected: rowsAffected,
		})
	default:
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(QueryResponse{
			Status: "error",
			Error:  "Unsupported query type",
		})
		return
	}

}

func (h *QueryHandler) HandleHealthCheck(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if err := h.dbManager.Ping(r.Context()); err != nil {
		w.WriteHeader(http.StatusServiceUnavailable)
		json.NewEncoder(w).Encode(QueryResponse{
			Status: "error",
			Error:  err.Error(),
		})
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(QueryResponse{
		Status: "ok",
	})
}
