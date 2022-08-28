package internalhttp

import (
	"encoding/json"
	"fmt"
	"github.com/EkaterinaShamanaeva/otus-go/hw12_13_14_15_calendar/internal/app"
	"github.com/gofrs/uuid"
	"net/http"
	"time"
)

func (s *Server) createEvent(w http.ResponseWriter, r *http.Request) {
	ev := app.Event{}
	err := json.NewDecoder(r.Body).Decode(&ev)
	fmt.Println("decoded: ", ev, ev.TimeStart, ev.Duration)
	if err != nil {
		s.logg.Error(fmt.Sprintf("error to get request body: %v", err))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	err = s.app.CreateEvent(r.Context(), &ev)
	if err != nil {
		s.logg.Error(fmt.Sprintf("error to create event: %v", err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (s *Server) updateEvent(w http.ResponseWriter, r *http.Request) {
	ev := app.Event{}
	err := json.NewDecoder(r.Body).Decode(&ev)
	if err != nil {
		s.logg.Error(fmt.Sprintf("error to get request body: %v", err))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	err = s.app.UpdateEvent(r.Context(), &ev)
	if err != nil {
		s.logg.Error(fmt.Sprintf("error to update event: %v", err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (s *Server) deleteEvent(w http.ResponseWriter, r *http.Request) {
	type id struct {
		id uuid.UUID `json:"id"`
	}
	var idEv id
	err := json.NewDecoder(r.Body).Decode(&idEv)
	if err != nil {
		s.logg.Error(fmt.Sprintf("error to get request body: %v", err))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	err = s.app.DeleteEvent(r.Context(), idEv.id)
	if err != nil {
		s.logg.Error(fmt.Sprintf("error to update event: %v", err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (s *Server) getEventsPerDay(w http.ResponseWriter, r *http.Request) {
	var day time.Time
	err := json.NewDecoder(r.Body).Decode(&day)
	if err != nil {
		s.logg.Error(fmt.Sprintf("error to get request body: %v", err))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	events, err := s.app.GetEventsPerDay(r.Context(), day)
	if err != nil {
		s.logg.Error(fmt.Sprintf("error to get list of events: %v", err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = json.NewEncoder(w).Encode(events)
	if err != nil {
		s.logg.Error(fmt.Sprintf("error to get list of events: %v", err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (s *Server) getEventsPerWeek(w http.ResponseWriter, r *http.Request) {
	var day time.Time
	err := json.NewDecoder(r.Body).Decode(&day)
	if err != nil {
		s.logg.Error(fmt.Sprintf("error to get request body: %v", err))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	events, err := s.app.GetEventsPerWeek(r.Context(), day)
	if err != nil {
		s.logg.Error(fmt.Sprintf("error to get list of events: %v", err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = json.NewEncoder(w).Encode(events)
	if err != nil {
		s.logg.Error(fmt.Sprintf("error to get list of events: %v", err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (s *Server) getEventsPerMonth(w http.ResponseWriter, r *http.Request) {
	var day time.Time
	err := json.NewDecoder(r.Body).Decode(&day)
	if err != nil {
		s.logg.Error(fmt.Sprintf("error to get request body: %v", err))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	events, err := s.app.GetEventsPerMonth(r.Context(), day)
	if err != nil {
		s.logg.Error(fmt.Sprintf("error to get list of events: %v", err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = json.NewEncoder(w).Encode(events)
	if err != nil {
		s.logg.Error(fmt.Sprintf("error to get list of events: %v", err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}
