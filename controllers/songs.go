package controllers

import (
	"api-songs-docker/db"
	"api-songs-docker/models"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"sort"
	"strings"
	"time"
	"github.com/golang-jwt/jwt"
)


var SECRET = []byte("super-secret-auth-key")
var api_key = "1234"

//Crea un json web token
func CreateJWT() (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)

	claims["exp"] = time.Now().Add(time.Hour).Unix()

	tokenStr, err := token.SignedString(SECRET)

	if err != nil {
		fmt.Println(err.Error())
		return "", err
	}

	return tokenStr, nil
}

//Valida el jwt generado en CreatJWT
func ValidateJWT(next func(w http.ResponseWriter, r* http.Request)) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		if r.Header["Token"] != nil {
			token, err := jwt.Parse(r.Header["Token"][0], func(t *jwt.Token) (interface{}, error) {
				_, ok := t.Method.(*jwt.SigningMethodHMAC)
				if !ok {
					w.WriteHeader(http.StatusUnauthorized)
					w.Write([]byte("not authorized"))
				}
				return SECRET, nil
			})

			if err != nil {
				w.WriteHeader(http.StatusUnauthorized)
				w.Write([]byte("not authorized: " + err.Error()))
			}

			if token.Valid {
				next(w, r)
			}
		} else {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("not authorized"))
		}
	})
}
//Brinda el token generado al usuario
func GetJwt(w http.ResponseWriter, r *http.Request) {
	if r.Header["Access"] != nil {
		if r.Header["Access"][0] != api_key {
			return
		} else {
			token, err := CreateJWT()
			if err != nil {
				return
			}
			fmt.Fprint(w, token)
		}
	}
}

// Busca canciones con la api de Itunes y devuelve 10 resultados en formato json
func SearchSongs(w http.ResponseWriter, r *http.Request) {
	var responseObject models.Response

	term := r.URL.Query().Get("term")
	newTerm := strings.Replace(term, "-", "+", 10)

	response, err := http.Get("https://itunes.apple.com/search?term=" + newTerm + "&limit=10")
	fmt.Println("https://itunes.apple.com/search?term=" + newTerm + "&limit=10")

	if err == nil {
		fmt.Println(response.Status)
		defer response.Body.Close()
	}
	if err != nil {
		fmt.Printf("server not responding %s", err.Error())
		return // the return statement here helps to handle the error
	}

	responseData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}

	json.Unmarshal(responseData, &responseObject)

	for i := 0; i < len(responseObject.Songs); i++ {

		jsonStruct := models.Songs{
			ArtistId:        responseObject.Songs[i].ArtistId,
			TrackName:       responseObject.Songs[i].TrackName,
			ArtistName:      responseObject.Songs[i].ArtistName,
			TrackTimeMillis: responseObject.Songs[i].TrackTimeMillis,
			CollectionName:  responseObject.Songs[i].CollectionName,
			ArtworkUrl30:    responseObject.Songs[i].ArtworkUrl30,
			TrackPrice:      responseObject.Songs[i].TrackPrice,
			Origin:          "apple",
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(jsonStruct)

	}

}

//Busca canciones desde la API de ChartLyric,devuelve un resultado de tipo json
func SearchChartlyric(w http.ResponseWriter, r *http.Request){
	term := r.URL.Query().Get("artist")
	newTerm := strings.Replace(term,"-","+",10)
	
	songTerm := r.URL.Query().Get("song")
	newSongTerm := strings.Replace(songTerm,"-","+",10)

	respon, err := http.Get("http://api.chartlyrics.com/apiv1.asmx/SearchLyricDirect?artist=" + newTerm + "&song=" + newSongTerm)
	fmt.Println("http://api.chartlyrics.com/apiv1.asmx/SearchLyric?artist=" + newTerm + "&song=" + newSongTerm)
	responData, err := ioutil.ReadAll(respon.Body)
	if err != nil {
		log.Fatal(err)
	}

	var data models.GetLyricResult
	if err = xml.Unmarshal(responData, &data); err != nil {
		log.Println(err)
	}

	
	xmlStruct := models.Songs{
		ArtistId:  0,
		TrackName: data.LyricSong,
		ArtistName: data.LyricArtist,
		TrackTimeMillis: 0,
		CollectionName: "",
		ArtworkUrl30: data.LyricCovertArtUrl,
		TrackPrice: 0,
		Origin: "ChartLyric", 
	}
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(xmlStruct)
}


//Muestra las canciones guardadas en la base de datos en postgres
func VerSongs(w http.ResponseWriter, r *http.Request) {
	
	songs := []models.Songs{}

	db.DB.Find(&songs)
	sort.Slice(songs, func(i, j int) bool {
		return songs[i].TrackName < songs[j].TrackName
	})
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(&songs)

}


//Guarda la informacion  de la canción en la base de datos, buscada desde Itunes
func GuardarSong(w http.ResponseWriter, r *http.Request) {
	
	var responseObject models.Response

	term := r.URL.Query().Get("term")
	newTerm := strings.Replace(term, "-", "+", 10)
	response, err := http.Get("https://itunes.apple.com/search?term=" + newTerm + "&limit=1")
	fmt.Println("https://itunes.apple.com/search?term=" + newTerm + "&limit=1")

	if err == nil {
		fmt.Println(response.Status)
		defer response.Body.Close()
	}
	if err != nil {
		fmt.Printf("server not responding %s", err.Error())
		return // the return statement here helps to handle the error
	}

	responseData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}

	json.Unmarshal(responseData, &responseObject)

	for i := 0; i < len(responseObject.Songs); i++ {

		jsonStruct := models.Songs{
			ArtistId:        responseObject.Songs[i].ArtistId,
			TrackName:       responseObject.Songs[i].TrackName,
			ArtistName:      responseObject.Songs[i].ArtistName,
			TrackTimeMillis: responseObject.Songs[i].TrackTimeMillis,
			CollectionName:  responseObject.Songs[i].CollectionName,
			ArtworkUrl30:    responseObject.Songs[i].ArtworkUrl30,
			TrackPrice:      responseObject.Songs[i].TrackPrice,
			Origin:          "apple",
		}
		song := &models.Songs{}
		createJson, _ := json.Marshal(jsonStruct)
		json.NewDecoder(strings.NewReader(string(createJson))).Decode(&song)
		agregarSong := db.DB.Create(&song)
		err := agregarSong.Error

		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
		}
		json.NewEncoder(w).Encode(&song)

	}
}

//Guarda una canción desde la API de ChartLyrics en la base de datos.
func SaveSong(w http.ResponseWriter, r *http.Request){
	term := r.URL.Query().Get("artist")
	newTerm := strings.Replace(term,"-","+",10)
	
	songTerm := r.URL.Query().Get("song")
	newSongTerm := strings.Replace(songTerm,"-","+",10)

	respon, err := http.Get("http://api.chartlyrics.com/apiv1.asmx/SearchLyricDirect?artist=" + newTerm + "&song=" + newSongTerm)
	fmt.Println("http://api.chartlyrics.com/apiv1.asmx/SearchLyric?artist=" + newTerm + "&song=" + newSongTerm)
	responData, err := ioutil.ReadAll(respon.Body)
	if err != nil {
		log.Fatal(err)
	}

	var data models.GetLyricResult
	if err = xml.Unmarshal(responData, &data); err != nil {
		log.Println(err)
	}

	
	xmlStruct := models.Songs{
		ArtistId:  0,
		TrackName: data.LyricSong,
		ArtistName: data.LyricArtist,
		TrackTimeMillis: 0,
		CollectionName: "",
		ArtworkUrl30: data.LyricCovertArtUrl,
		TrackPrice: 0,
		Origin: "ChartLyric", 
	}
	
	song := &models.Songs{}
	createJson, _ := json.Marshal(xmlStruct)
	json.NewDecoder(strings.NewReader(string(createJson))).Decode(&song)
	agregarSong := db.DB.Create(&song)
	error := agregarSong.Error

	if error != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
	}
	json.NewEncoder(w).Encode(&song)	
		
}


