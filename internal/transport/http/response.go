package response

import (
	"encoding/json"
	"log/slog"
	"net/http"
)

func ResponseJSON(log *slog.Logger, w http.ResponseWriter, status int, data any) {
	b, err := json.Marshal(data)
	if err != nil {
		log.Error("response http", slog.Any("failed to marshal", err))
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)

		w.Write([]byte(`{"error": "internal server error"}`))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(b)
}

func ResponseErr(log *slog.Logger, w http.ResponseWriter, status int, message string) {
	ResponseJSON(log, w, status, map[string]string{"error": message})
}
