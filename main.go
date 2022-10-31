package main

import (
	"api-songs-docker/controllers"
	"api-songs-docker/db"
	"api-songs-docker/models"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	db.Connect()
	db.DB.AutoMigrate(&models.Songs{})

	router := mux.NewRouter()
	api := "/api"

	//Rutas

	router.HandleFunc(api + "/", controllers.Home).Methods("GET")

	router.HandleFunc(api + "/search", controllers.SearchSongs).Methods("GET")

	router.HandleFunc(api + "/searchlyric", controllers.SearchChartlyric).Methods("GET")

	router.HandleFunc(api + "/jwt", controllers.GetJwt).Methods("GET")

	router.Handle(api + "/mysongs", controllers.ValidateJWT(controllers.VerSongs)).Methods("GET")
	
	router.Handle(api + "/mysongs", controllers.ValidateJWT(controllers.GuardarSong)).Methods("POST")

	router.Handle(api + "/songs", controllers.ValidateJWT(controllers.SaveSong)).Methods("POST")

	http.ListenAndServe(":8080", router)
}
