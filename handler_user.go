package main
import (
	"encoding/json"
	"time"
	"log"
	"github.com/google/uuid"
	"fmt"
	"net/http"
	"github.com/sahildhargave/rss_scraper/internal/database"
)

func (apiCfg *apiConfig) handlerUsersCreate(w http.ResponseWriter, r *http.Request){
	type parameters struct{
		Name string `json:"name"`
	}
	decoder := json.NewDecoder(r.Body)

	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Error parsing JSON:", err))
		return
	}
	user, err := apiCfg.DB.CreateUser(r.Context(), database.CreateUserParams{
		ID: uuid.New(),
		CreatedAt : time.Now().UTC(),
		UpdatedAt : time.Now().UTC(),
		Name: params.Name,
	})

	if err != nil {
		log.Println(err)
		respondWithError(w, http.StatusInternalServerError, "Couldn't create user")
		return
	}

	respondWithJSON(w, 200, user)

}


// func (apiCfg *apiConfig) handlerUsersGet(w http.ResponseWrite, r *http.Request, user database.User){
//        respondWithJOSN(w , http.StatusOK, databaseUserToUser(user))
// } 