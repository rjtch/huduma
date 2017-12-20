package api

import (
	"net/http"
)

//Info gives some informations about the given api
func (api *API) Info(w http.ResponseWriter, r *http.Request) {
	sendJSON(w, http.StatusOK, map[string]string{
		"version": "api.Vesion",
		"App":     "huduma",
	})
}

//GetAllAPIs lists all availables APIs
func GetAllAPIs() {
}

//GetOneAPI returns one specific API
func GetOneAPI() {
}

//PutOneAPI updates one specific API
func PutOneAPI() {

}

//PostOneAPI creates a new API
func PostOneAPI() {

}

//DeleteOneAPI deletes a specific API
func DeleteOneAPI() {

}

//DeleteAllAPIs deletes all available APIs in a specific context
func DeleteAllAPIs() {

}
