package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/maxheckel/maxs-marvelous-manuscript/internal/db"
)

type API struct {
	recordingRepo *db.RecordingRepository
	dataDir       string
}

func NewAPI(recordingRepo *db.RecordingRepository, dataDir string) *API {
	return &API{
		recordingRepo: recordingRepo,
		dataDir:       dataDir,
	}
}

// RegisterRoutes registers all API routes
func (a *API) RegisterRoutes(r *mux.Router) {
	api := r.PathPrefix("/api").Subrouter()

	// Recordings endpoints
	api.HandleFunc("/recordings", a.listRecordings).Methods("GET")
	api.HandleFunc("/recordings/{id}", a.getRecording).Methods("GET")
	api.HandleFunc("/recordings/{id}", a.deleteRecording).Methods("DELETE")
	api.HandleFunc("/recordings/{id}/audio", a.streamAudio).Methods("GET")

	// Health check
	api.HandleFunc("/health", a.healthCheck).Methods("GET")
}

// listRecordings returns all recordings
func (a *API) listRecordings(w http.ResponseWriter, r *http.Request) {
	recordings, err := a.recordingRepo.List()
	if err != nil {
		respondError(w, http.StatusInternalServerError, fmt.Sprintf("Failed to list recordings: %v", err))
		return
	}

	respondJSON(w, http.StatusOK, recordings)
}

// getRecording returns a specific recording
func (a *API) getRecording(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid recording ID")
		return
	}

	recording, err := a.recordingRepo.GetByID(id)
	if err != nil {
		respondError(w, http.StatusNotFound, "Recording not found")
		return
	}

	respondJSON(w, http.StatusOK, recording)
}

// deleteRecording deletes a recording
func (a *API) deleteRecording(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid recording ID")
		return
	}

	if err := a.recordingRepo.Delete(id); err != nil {
		respondError(w, http.StatusInternalServerError, fmt.Sprintf("Failed to delete recording: %v", err))
		return
	}

	respondJSON(w, http.StatusOK, map[string]string{"message": "Recording deleted"})
}

// streamAudio streams the audio file for a recording
func (a *API) streamAudio(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid recording ID")
		return
	}

	recording, err := a.recordingRepo.GetByID(id)
	if err != nil {
		respondError(w, http.StatusNotFound, "Recording not found")
		return
	}

	// Serve the audio file
	w.Header().Set("Content-Type", "audio/wav")
	w.Header().Set("Content-Disposition", fmt.Sprintf("inline; filename=%s", recording.Filename))
	http.ServeFile(w, r, recording.FilePath)
}

// healthCheck returns the API health status
func (a *API) healthCheck(w http.ResponseWriter, r *http.Request) {
	respondJSON(w, http.StatusOK, map[string]string{
		"status": "healthy",
	})
}

// Helper functions

func respondJSON(w http.ResponseWriter, status int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(payload)
}

func respondError(w http.ResponseWriter, status int, message string) {
	respondJSON(w, status, map[string]string{"error": message})
}
