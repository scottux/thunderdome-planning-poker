package api

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/StevenWeathers/thunderdome-planning-poker/pkg/database"
	"github.com/gorilla/mux"
)

// handleGetTeamByUser gets an team with user role
func (a *api) handleGetTeamByUser() http.HandlerFunc {
	type TeamResponse struct {
		Team     *database.Team `json:"team"`
		TeamRole string         `json:"teamRole"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		TeamRole := r.Context().Value(contextKeyTeamRole).(string)
		TeamID := vars["teamId"]

		Team, err := a.db.TeamGet(TeamID)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		a.respondWithJSON(w, http.StatusOK, &TeamResponse{
			Team:     Team,
			TeamRole: TeamRole,
		})
	}
}

// handleGetTeamsByUser gets a list of teams the user is apart of
func (a *api) handleGetTeamsByUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		UserID := r.Context().Value(contextKeyUserID).(string)
		vars := mux.Vars(r)
		Limit, _ := strconv.Atoi(vars["limit"])
		Offset, _ := strconv.Atoi(vars["offset"])

		Organizations := a.db.TeamListByUser(UserID, Limit, Offset)

		a.respondWithJSON(w, http.StatusOK, Organizations)
	}
}

// handleGetTeamUsers gets a list of users associated to the team
func (a *api) handleGetTeamUsers() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		TeamID := vars["teamId"]
		Limit, _ := strconv.Atoi(vars["limit"])
		Offset, _ := strconv.Atoi(vars["offset"])

		Teams := a.db.TeamUserList(TeamID, Limit, Offset)

		a.respondWithJSON(w, http.StatusOK, Teams)
	}
}

// handleCreateTeam handles creating an team with current user as admin
func (a *api) handleCreateTeam() http.HandlerFunc {
	type CreateTeamResponse struct {
		TeamID string `json:"id"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		UserID := r.Context().Value(contextKeyUserID).(string)
		keyVal := a.getJSONRequestBody(r, w)

		TeamName := keyVal["name"].(string)
		TeamID, err := a.db.TeamCreate(UserID, TeamName)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		var NewTeam = &CreateTeamResponse{
			TeamID: TeamID,
		}

		a.respondWithJSON(w, http.StatusOK, NewTeam)
	}
}

// handleTeamAddUser handles adding user to a team
func (a *api) handleTeamAddUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		keyVal := a.getJSONRequestBody(r, w)

		vars := mux.Vars(r)
		TeamID := vars["teamId"]
		UserEmail := strings.ToLower(keyVal["email"].(string))
		Role := keyVal["role"].(string)

		User, UserErr := a.db.GetUserByEmail(UserEmail)
		if UserErr != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		_, err := a.db.TeamAddUser(TeamID, User.UserID, Role)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		return
	}
}

// handleTeamRemoveUser handles removing user from a team
func (a *api) handleTeamRemoveUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		keyVal := a.getJSONRequestBody(r, w)

		vars := mux.Vars(r)
		TeamID := vars["teamId"]
		UserID := keyVal["id"].(string)

		err := a.db.TeamRemoveUser(TeamID, UserID)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		return
	}
}

// handleGetTeamBattles gets a list of battles associated to the team
func (a *api) handleGetTeamBattles() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		TeamID := vars["teamId"]
		Limit, _ := strconv.Atoi(vars["limit"])
		Offset, _ := strconv.Atoi(vars["offset"])

		Battles := a.db.TeamBattleList(TeamID, Limit, Offset)

		a.respondWithJSON(w, http.StatusOK, Battles)
	}
}

// handleTeamRemoveBattle handles removing battle from a team
func (a *api) handleTeamRemoveBattle() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		keyVal := a.getJSONRequestBody(r, w)

		vars := mux.Vars(r)
		TeamID := vars["teamId"]
		BattleID := keyVal["id"].(string)

		err := a.db.TeamRemoveBattle(TeamID, BattleID)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		return
	}
}

// handleDeleteTeam handles deleting a team
func (a *api) handleDeleteTeam() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		keyVal := a.getJSONRequestBody(r, w)
		TeamID := keyVal["id"].(string)

		err := a.db.TeamDelete(TeamID)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		return
	}
}