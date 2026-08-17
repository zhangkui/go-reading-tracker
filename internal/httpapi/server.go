package httpapi

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/zhangkui/go-reading-tracker/internal/reading"
)

type Server struct {
	service *reading.Service
	mux     *http.ServeMux
}

func NewServer(service *reading.Service) *Server {
	server := &Server{service: service, mux: http.NewServeMux()}
	server.mux.HandleFunc("/health", server.health)
	server.mux.HandleFunc("/books", server.books)
	server.mux.HandleFunc("/books/", server.bookByID)
	server.mux.HandleFunc("/goals", server.goals)
	server.mux.HandleFunc("/goals/progress", server.goalProgress)
	server.mux.HandleFunc("/imports", server.imports)
	server.mux.HandleFunc("/stats", server.statistics)
	return server
}

func (s *Server) Handler() http.Handler { return s.mux }

func writeJSON(response http.ResponseWriter, status int, value any) {
	response.Header().Set("Content-Type", "application/json")
	response.WriteHeader(status)
	_ = json.NewEncoder(response).Encode(value)
}

func writeError(response http.ResponseWriter, err error) {
	status := http.StatusBadRequest
	if errors.Is(err, reading.ErrNotFound) || errors.Is(err, reading.ErrGoalNotFound) {
		status = http.StatusNotFound
	}
	if errors.Is(err, reading.ErrDuplicateISBN) {
		status = http.StatusConflict
	}
	writeJSON(response, status, map[string]string{"error": err.Error()})
}

func decodeJSON(request *http.Request, destination any) error {
	decoder := json.NewDecoder(request.Body)
	decoder.DisallowUnknownFields()
	return decoder.Decode(destination)
}

func (s *Server) health(response http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodGet {
		response.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	writeJSON(response, http.StatusOK, map[string]string{"status": "ok"})
}

func (s *Server) books(response http.ResponseWriter, request *http.Request) {
	switch request.Method {
	case http.MethodPost:
		var input reading.CreateBookInput
		if err := decodeJSON(request, &input); err != nil {
			writeError(response, err)
			return
		}
		book, err := s.service.AddBook(input)
		if err != nil {
			writeError(response, err)
			return
		}
		writeJSON(response, http.StatusCreated, book)
	case http.MethodGet:
		query := request.URL.Query()
		page, _ := strconv.Atoi(query.Get("page"))
		pageSize, _ := strconv.Atoi(query.Get("page_size"))
		result := s.service.SearchBooks(reading.SearchOptions{Title: query.Get("title"), Author: query.Get("author"), Status: reading.Status(query.Get("status")), Tags: query["tag"], Page: page, PageSize: pageSize})
		writeJSON(response, http.StatusOK, result)
	default:
		response.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func parseBookPath(path string) (int64, bool, error) {
	parts := strings.Split(strings.TrimPrefix(path, "/books/"), "/")
	id, err := strconv.ParseInt(parts[0], 10, 64)
	return id, len(parts) == 2 && parts[1] == "progress", err
}

func (s *Server) bookByID(response http.ResponseWriter, request *http.Request) {
	id, progress, err := parseBookPath(request.URL.Path)
	if err != nil {
		writeError(response, reading.ErrInvalidInput)
		return
	}
	if progress {
		if request.Method != http.MethodPut {
			response.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		var input struct {
			CurrentPage int `json:"current_page"`
		}
		if err := decodeJSON(request, &input); err != nil {
			writeError(response, err)
			return
		}
		book, err := s.service.UpdateProgress(id, input.CurrentPage)
		if err != nil {
			writeError(response, err)
			return
		}
		writeJSON(response, http.StatusOK, book)
		return
	}
	switch request.Method {
	case http.MethodGet:
		book, err := s.service.GetBook(id)
		if err != nil {
			writeError(response, err)
			return
		}
		writeJSON(response, http.StatusOK, book)
	case http.MethodPut:
		var input reading.UpdateBookInput
		if err := decodeJSON(request, &input); err != nil {
			writeError(response, err)
			return
		}
		book, err := s.service.UpdateBook(id, input)
		if err != nil {
			writeError(response, err)
			return
		}
		writeJSON(response, http.StatusOK, book)
	case http.MethodDelete:
		if err := s.service.DeleteBook(id); err != nil {
			writeError(response, err)
			return
		}
		response.WriteHeader(http.StatusNoContent)
	default:
		response.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (s *Server) goals(response http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodPost {
		response.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	var goal reading.Goal
	if err := decodeJSON(request, &goal); err != nil {
		writeError(response, err)
		return
	}
	created, err := s.service.SetGoal(goal)
	if err != nil {
		writeError(response, err)
		return
	}
	writeJSON(response, http.StatusCreated, created)
}

func (s *Server) goalProgress(response http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodGet {
		response.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	progress, err := s.service.GetGoalProgress(request.URL.Query().Get("month"), reading.GoalType(request.URL.Query().Get("type")))
	if err != nil {
		writeError(response, err)
		return
	}
	writeJSON(response, http.StatusOK, progress)
}

func (s *Server) imports(response http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodPost {
		response.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	result, err := s.service.Import(request.Context(), request.Body)
	if err != nil {
		writeError(response, err)
		return
	}
	writeJSON(response, http.StatusOK, result)
}

func (s *Server) statistics(response http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodGet {
		response.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	result, err := s.service.Statistics(request.URL.Query().Get("month"))
	if err != nil {
		writeError(response, err)
		return
	}
	writeJSON(response, http.StatusOK, result)
}
