package handler

import (
	"cmp"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"slices"
	"strconv"
	"strings"

	"github.com/nxvrlxv/hosts_monitor/internal/checker"
	"github.com/nxvrlxv/hosts_monitor/internal/storage"
)

var (
	OrderByAsc  = 1
	OrderByAsIs = 0
	OrderByDesc = -1

	BadNumber             = errors.New("not a number")
	NotAValidNumberStatus = errors.New("Not valid status int (must be 1 or 0)")
)

type StorageHandler struct {
	storage *storage.Storage
}

type SearchRequest struct {
	Status     bool
	OrderField string
	OrderBy    int
}

type targetResponse struct {
	ID        int     `json:"id"`
	Host      string  `json:"host"`
	Up        bool    `json:"up"`
	Status    int     `json:"status_code,omitempty"`
	LatencyMS float64 `json:"latency_ms"`
	Error     string  `json:"error,omitempty"`
}

func parseParams(paramsURL url.Values) (*SearchRequest, error) {
	if len(paramsURL) == 0 {
		return nil, nil
	}
	status, err := strconv.Atoi(paramsURL.Get("status"))
	if err != nil {
		return nil, fmt.Errorf("status: %w", BadNumber)
	}

	if status != 0 && status != 1 {
		return nil, NotAValidNumberStatus
	}
	var statusBool bool
	if status == 0 {
		statusBool = false
	} else {
		statusBool = true
	}

	orderField := strings.ToLower(paramsURL.Get("order_field"))
	if orderField == "" {
		orderField = "id"
	}

	orderBy, err := strconv.Atoi(paramsURL.Get("order_by"))
	if err != nil {
		return nil, fmt.Errorf("orderby: %w", BadNumber)
	}
	return &SearchRequest{
		Status:     statusBool,
		OrderField: orderField,
		OrderBy:    orderBy,
	}, nil

}

func selectTargetsWithStatus(targets []checker.Result, status bool) []targetResponse {
	var res []targetResponse

	for _, v := range targets {
		if v.IsUp == status {
			errField := ""
			if v.Error != nil {
				errField = v.Error.Error()
			}
			target := targetResponse{
				ID:        v.ID,
				Host:      v.Host,
				Up:        v.IsUp,
				Status:    v.Status,
				LatencyMS: float64(v.ResponseTime.Milliseconds()),
				Error:     errField,
			}
			res = append(res, target)
		}
	}

	return res
}

func OrderingTargets(targets []targetResponse, orderField string, orderBy int) []targetResponse {

	if orderBy == OrderByAsIs {
		return targets
	}

	var comparator func(a, b targetResponse) int
	switch orderField {
	case "host":
		comparator = func(a targetResponse, b targetResponse) int {
			return cmp.Compare(a.Host, b.Host)
		}
	default:
		comparator = func(a, b targetResponse) int {
			return cmp.Compare(a.ID, b.ID)
		}
	}

	var flipper func(a, b targetResponse) int
	if orderBy == OrderByDesc {
		flipper = func(a, b targetResponse) int {
			return comparator(b, a)
		}
	} else {
		flipper = comparator
	}

	slices.SortFunc(targets, flipper)
	return targets

}

func NewStorageHandler(storage *storage.Storage) *StorageHandler {
	return &StorageHandler{
		storage: storage,
	}
}

func (h *StorageHandler) GetTargets(w http.ResponseWriter, r *http.Request) {
	params, err := parseParams(r.URL.Query())
	if err != nil {
		w.Header().Set("Content-Type", "application/json")

		if errors.Is(err, BadNumber) {
			w.WriteHeader(http.StatusBadRequest)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
		}

		w.Write([]byte(fmt.Sprintf(`{error: %v}`, err.Error())))
		return
	}
	targets := h.storage.GetAllData()
	var res []targetResponse
	if params == nil {
		for _, v := range targets {
			errField := ""
			if v.Error != nil {
				errField = v.Error.Error()
			}
			target := targetResponse{
				ID:        v.ID,
				Host:      v.Host,
				Up:        v.IsUp,
				Status:    v.Status,
				LatencyMS: float64(v.ResponseTime.Milliseconds()),
				Error:     errField,
			}
			res = append(res, target)
		}
	} else {
		statusTargets := selectTargetsWithStatus(targets, params.Status)
		res = OrderingTargets(statusTargets, params.OrderField, params.OrderBy)
	}

	body, err := json.Marshal(res)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("{'error': 'error internal'}"))
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(body)

}

func (h *StorageHandler) GetTargetById(w http.ResponseWriter, r *http.Request) {
	strID := r.PathValue("id")
	id, err := strconv.Atoi(strID)
	if err != nil {
		http.Error(w, fmt.Errorf("id: %w", BadNumber).Error(), http.StatusBadRequest)
		return
	}

	target, err := h.storage.GetResultById(id)
	if err != nil {
		http.Error(w, "Not Found", http.StatusNotFound)
		return
	}

	resultTarget := targetResponse{
		ID:        target.ID,
		Host:      target.Host,
		Up:        target.IsUp,
		Status:    target.Status,
		LatencyMS: float64(target.ResponseTime.Milliseconds()),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resultTarget)

}
