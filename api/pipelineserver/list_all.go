package pipelineserver

import (
	"encoding/json"
	"net/http"

	"github.com/concourse/atc/api/present"
	"github.com/concourse/atc/auth"
	"github.com/concourse/atc/dbng"
)

// show all public pipelines and team private pipelines if authorized
func (s *Server) ListAllPipelines(w http.ResponseWriter, r *http.Request) {
	logger := s.logger.Session("list-all-pipelines")
	authTeam, authTeamFound := auth.GetTeam(r)

	var pipelines []dbng.Pipeline
	var err error
	if authTeamFound {
		team, found, err := s.teamFactory.FindTeam(authTeam.Name())
		if err != nil {
			logger.Error("failed-to-get-team", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		if !found {
			logger.Info("team-not-found")
			w.WriteHeader(http.StatusNotFound)
			return
		}

		pipelines, err = team.VisiblePipelines()
	} else {
		pipelines, err = s.pipelineFactory.PublicPipelines()
	}

	if err != nil {
		logger.Error("failed-to-get-all-active-pipelines", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(present.Pipelines(pipelines))
}
