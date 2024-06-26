// -*- coding: utf-8 -*-

package web

import (
	. "Projet_GO_Reservation/pkg/const"
	. "Projet_GO_Reservation/pkg/json"
	. "Projet_GO_Reservation/pkg/models"
	. "Projet_GO_Reservation/pkg/reservation"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func EnableHandlers() {

	// Create the directory with static file like CSS and JS
	staticDir := http.Dir("templates/src")
	staticHandler := http.FileServer(staticDir)
	http.Handle("/static/", http.StripPrefix("/static/", staticHandler))

	// Create all the handle to "listen" the right path using const in webConst.go
	http.HandleFunc(RouteIndex, IndexHandler)

	// Reservation Handlers
	http.HandleFunc(RouteIndexReservation, ReservationHandler)
	http.HandleFunc(RouteListReservation, ListByRoomDateIdReservationHandler)
	http.HandleFunc(RouteCreateReservation, CreateReservationHandler)
	http.HandleFunc(RouteCancelReservation, CancelReservationHandler)
	http.HandleFunc(RouteUpdateReservation, UpdateReservationHandler)

	// Rooms Handlers
	http.HandleFunc(RouteGetAllRoolAvailable, GetAllRoomAvailHandler)
	http.HandleFunc(RouteGetAllRooms, GetAllRoomsHandler)
	http.HandleFunc(RouteCreateRoom, CreateRoomHandler)
	http.HandleFunc(RouteCancelSalle, CancelSalleHandler)

	// Json Handlers
	http.HandleFunc(RouteExportJson, ExportJsonHandler)
	http.HandleFunc(RouteDownloadJson, DownloadJsonHandler)
	//http.HandleFunc(RouteImportJson, UploadJsonHandler)

	Log.Infos("Handlers Enabled")

	Log.Infos("Starting server on port " + PORT)
	err := http.ListenAndServe(":"+PORT, nil)
	if err != nil {
		Log.Error("Failed to start server: ", err)
		return
	}
	Log.Infos("Server stopped on port " + PORT)

}

//
// ------------------------------------------------------------------------------------------------ //
//

func IndexHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method == http.MethodGet {
		currentTime := time.Now().Format("2006-01-02 15:04")
		templates.ExecuteTemplate(w, "menu.html", map[string]interface{}{
			"now": currentTime,
		})
	}
}

//
// ------------------------------------------------------------------------------------------------ //
//

func ReservationHandler(w http.ResponseWriter, r *http.Request) {
	//w.Header().Set("Content-Type", "text/html; charset=utf-8")

	if r.Method == http.MethodGet {
		result := ListReservations(nil)
		if result == nil {
			Log.Error("Data are null for unknown reason :/")
			var msg = "Impossible to retrieve data"
			// Exécuter le template avec l'URL et le message
			// Fait un "truc bizarre" car impossible d'utiliser la méthode avec le message dans l'URL car sinon redirections infinie
			templates.ExecuteTemplate(w, "reservations.html", map[string]interface{}{
				"message": msg,
				"result":  nil,
			})
			return
		}
		templates.ExecuteTemplate(w, "reservations.html", map[string]interface{}{
			"result":  result,
			"message": nil,
		})

	}
}

//
// ------------------------------------------------------------------------------------------------ //
//

func ListByRoomDateIdReservationHandler(w http.ResponseWriter, r *http.Request) {

	roomStr := r.URL.Query().Get("idRoom")

	dateStr := r.URL.Query().Get("idDate")

	idStr := r.URL.Query().Get("idReserv")

	if r.Method == http.MethodGet {

		// Si aucun ID n'est fourni, redirigez vers la page de liste des réservations
		if roomStr == NullString && dateStr == NullString && idStr == NullString {
			var msg = "Vous ne pouvez pas acceder à cette page sans spécifier un truc :/"
			http.Redirect(w, r, "/reservation?message="+msg, http.StatusSeeOther)
			return
		}

		var result []Reservation

		var resultType string

		// If it's with room param
		if roomStr != NullString {
			resultType = "room"
			idRoom, err := strconv.Atoi(roomStr)
			if err != nil {
				http.Error(w, "ID de salle invalide", http.StatusBadRequest)
				return
			}

			Log.Infos("Listing des réservations par Salles")
			result = ListReservationsByRoom(&idRoom)
		}

		// If it's with date param
		if dateStr != NullString {
			resultType = "date"
			Log.Infos("Listing des réservations par Date")
			dateStr = strings.Replace(dateStr, "T", " ", 1)
			dateStr = dateStr + ":00"
			result = ListReservationsByDate(&dateStr)
		}

		// If it's with id param
		if idStr != NullString {
			resultType = "ID reservation"
			Log.Infos("Listing des réservations par ID (reservation)")
			var tmp = "id_reservation=" + idStr
			result = ListReservations(&tmp)
			// It have a special pages yes
			if result != nil {
				templates.ExecuteTemplate(w, "soloReservation.html", result)
			} else {
				var msg = "Impossible de trouver cette réservation"
				http.Redirect(w, r, "/reservation?message="+msg, http.StatusSeeOther)
			}
			return
		}

		if result == nil {
			Log.Error("No result")
			var msg = "Impossible to retrieve data, this " + resultType + " don't exist"
			http.Redirect(w, r, "/reservation?message="+msg, http.StatusSeeOther)
			return
		}

		templates.ExecuteTemplate(w, "reservations.html", map[string]interface{}{
			"message": nil,
			"result":  result,
		})

	}
}

//
// ------------------------------------------------------------------------------------------------ //
//

func CreateReservationHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method == http.MethodGet {
		templates.ExecuteTemplate(w, "creerReservations.html", map[string]interface{}{
			"message": nil,
			"result":  nil,
		})
	} else if r.Method == http.MethodPost {
		horaireStartDate := r.FormValue("horaire_start_date")
		horaireStartTime := r.FormValue("horaire_start_time") + ":00"
		horaireEndDate := r.FormValue("horaire_end_date")
		horaireEndTime := r.FormValue("horaire_end_time") + ":00"
		salle := r.FormValue("id_salle")

		salleInt64, err := strconv.ParseInt(salle, 10, 64)

		resultSalle := GetAllSalle()
		leBool := false
		for _, m := range resultSalle {
			if m.IdSalle == salleInt64 {
				leBool = true
				break
			}
		}

		if leBool == false {
			var msg = "Cette ID de salle n'existe pas :/"
			Log.Error(msg)
			http.Redirect(w, r, "/reservation/create?message="+msg, http.StatusSeeOther)
			return
		}

		if err != nil {
			var msg = "Erreur dans le format de la date/heure de début"
			Log.Error(msg)
			http.Redirect(w, r, "/reservation/create?message="+msg, http.StatusSeeOther)
			return
		}

		horaireStartDateTime, err := time.Parse("2006-01-02 15:04:05", horaireStartDate+" "+horaireStartTime)
		if err != nil {
			var msg = "Erreur dans le format de la date/heure de début"
			Log.Error(msg)
			http.Redirect(w, r, "/reservation/create?message="+msg, http.StatusSeeOther)
			return
		}

		horaireEndDateTime, err := time.Parse("2006-01-02 15:04:05", horaireEndDate+" "+horaireEndTime)
		if err != nil {
			var msg = "Erreur dans le format de la date/heure de début"
			Log.Error(msg)
			http.Redirect(w, r, "/reservation/create?message="+msg, http.StatusSeeOther)
			return
		}

		today := time.Now()

		if horaireStartDateTime.Before(today) || horaireStartDateTime.Equal(today) {
			var msg = "Impossible de créer une réservation avant aujourd'hui ou aujourd'hui"
			Log.Error(msg)
			http.Redirect(w, r, "/reservation/create?message="+msg, http.StatusSeeOther)
			return
		}

		if horaireStartDateTime.After(horaireEndDateTime) {
			var msg = "La fin de la réservation doit être après le début de celle-ci"
			Log.Error(msg)
			http.Redirect(w, r, "/reservation/create?message="+msg, http.StatusSeeOther)
			return
			//La fin est avant le début ??
		}

		if horaireStartDateTime.Equal(horaireEndDateTime) {
			var msg = "Vous ne pouvez pas faire des réservation de moins de 1 minute"
			Log.Error(msg)
			http.Redirect(w, r, "/reservation/create?message="+msg, http.StatusSeeOther)
			return
			// La fin doit être différente du début
		}

		horaireStartSeconds := horaireStartDate + " " + horaireStartTime
		horaireEndSeconds := horaireEndDate + " " + horaireEndTime

		if salleInt64 == 13 {
			var msg = "Salle déjà réservé pour l'éternité par M. Sananes qui a traumatisé des générations d'élèves avec ses pointeurs."
			Log.Error(msg)
			http.Redirect(w, r, "/reservation?message="+msg, http.StatusSeeOther)
			return
		}

		result := CreateReservation(&salleInt64, &horaireStartSeconds, &horaireEndSeconds)
		if result == false {
			var msg = "An error occured"
			Log.Error(msg)
			http.Redirect(w, r, "/reservation?message="+msg, http.StatusSeeOther)
			return
		}

		http.Redirect(w, r, RouteIndexReservation, http.StatusSeeOther)
	}
}

//
// ------------------------------------------------------------------------------------------------ //
//

func CancelReservationHandler(w http.ResponseWriter, r *http.Request) {

	idReserv := r.URL.Query().Get("idReserv")
	var tmp = "id_reservation=" + idReserv

	result := ListReservations(&tmp)
	if len(result) == 0 {
		Log.Error("Aucune reservation contient cet ID")
		http.Error(w, "Aucune réservation trouvée pour l'ID de salle "+idReserv, http.StatusBadRequest)
		return
	}

	newChoix, err := strconv.Atoi(idReserv)

	if err != nil {
		Log.Error("Impossible to convert the id from string to int")
		http.Error(w, "Impossible to convert the id from string to int : "+idReserv, http.StatusBadRequest)
		return
	}

	CancelReservation(newChoix)

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Réservation annulée avec succès"))
	Log.Infos("Reservation annulée avec succès !")
	return
}

//
// ------------------------------------------------------------------------------------------------ //
//

func UpdateReservationHandler(w http.ResponseWriter, r *http.Request) {
	leString := r.URL.Query().Get("idReserv")

	Log.Debug(leString)

	idReserv := strings.Split(leString, "?")[0]
	newEtat := strings.Split(leString, "?etat=")[1]

	Log.Debug(idReserv)
	Log.Debug(newEtat)
	//return

	if idReserv == "nil" || newEtat == "nil" {
		var msg = "Vous ne pouvez pas mettre à jour une réservation avec le même statut :/"
		Log.Error(msg)
		http.Error(w, msg, http.StatusBadRequest)
		return
	}

	var err error
	var newReserv, newState int

	newReserv, err = strconv.Atoi(idReserv)
	if err != nil {
		var msg = "Impossible de transformer l'ID reserv string en int : " + idReserv
		Log.Error(msg, err)
		http.Error(w, msg, http.StatusBadRequest)
		return
	}
	newState, err = strconv.Atoi(newEtat)
	if err != nil {
		var msg = "Impossible de transformer l'ID etat string en int : " + newEtat
		Log.Error(msg, err)
		http.Error(w, msg, http.StatusBadRequest)
		return
	}

	UpdateReservation(&newState, &newReserv)
	//return
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Réservation annulée avec succès"))
	Log.Infos("Reservation annulée avec succès !")
	//http.Redirect(w, r, "", http.StatusSeeOther)
	return
}

//
// ------------------------------------------------------------------------------------------------ //
//

// ExportJsonHandler Export the BDD to the data.json
func ExportJsonHandler(w http.ResponseWriter, r *http.Request) {

	leBool := DataToJson(ListReservations(nil))

	if leBool == false {
		var msg = "L'export en JSON n'a pas réussi :/"
		Log.Error(msg)
		http.Error(w, msg, http.StatusBadRequest)
		return
	}

	var msg = "L'export en JSON à réussi, vous pouvez désormais le télécharger !"
	Log.Infos(msg)
	http.Error(w, msg, http.StatusOK)

}

//
// ------------------------------------------------------------------------------------------------ //
//

// DownloadJsonHandler Send the data.json to the client
func DownloadJsonHandler(w http.ResponseWriter, r *http.Request) {
	Log.Debug("Download begin")
	http.ServeFile(w, r, "./data.json")
}

//
// ------------------------------------------------------------------------------------------------ //
//

/*
// UploadJsonHandler Send the data to the server
func UploadJsonHandler(w http.ResponseWriter, r *http.Request) {
	Log.Debug("Upload begin")

	if r.Method == http.MethodPost {

		// Lit le fichier envoyé
		file, _, err := r.FormFile("reservation")
		if err != nil {
			var msg = "Erreur lors de l'upload du fichier : " + err.Error()
			Log.Error(msg)
			http.Error(w, msg, http.StatusBadRequest)
			return
		}

		defer file.Close()

		// Le transforme en fichier json
		jsonData, err := io.ReadAll(file)
		if err != nil {
			var msg = "Erreur lors de la lecture du fichier : " + err.Error()
			Log.Error(msg)
			http.Error(w, msg, http.StatusBadRequest)
			return
		}

		// Transforme le fichier json en map[string]interface{}
		var data []map[string]interface{}
		err = json.Unmarshal(jsonData, &data)

		// Save it inside the BDD
		if !JsonToData(data) {
			var msg = "Erreur lors de l'enregistrement des données : " + err.Error()
			Log.Error(msg)
			http.Error(w, msg, http.StatusBadRequest)
			return
		}

		var msg = "L'import et la sauvegarde dans la base de donnée ont réussi !"
		Log.Infos(msg)
		http.Error(w, msg, http.StatusOK)

	}

}
*/

//
// ------------------------------------------------------------------------------------------------ //
//

func GetAllRoomAvailHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method == http.MethodGet {
		// Affichage de toutes les salles disponible selon un horaire, etc... (OU PAS, on peut le faire aussi dans la page des salles de bases)
	} else if r.Method == http.MethodPost {

		var params struct {
			StartDateTime string `json:"startDateTime"`
			EndDateTime   string `json:"endDateTime"`
		}
		err := json.NewDecoder(r.Body).Decode(&params)
		if err != nil {
			var msg = "Erreur lors de la lecture des paramètres"
			Log.Error(msg)
			http.Error(w, msg, http.StatusBadRequest)
			return
		}

		sallesAvail := GetAllSalleDispo(&params.StartDateTime, &params.EndDateTime)
		if sallesAvail == nil {
			var msg = "Impossible de récupérer les salles disponibles"
			Log.Error(msg)
			http.Error(w, msg, http.StatusBadRequest)
			return
		}

		jsonData, err := json.Marshal(sallesAvail)
		if err != nil {
			var msg = "Erreur lors de la conversion en JSON, impossible d'avoir les salles disponibles"
			http.Error(w, msg, http.StatusInternalServerError)
			Log.Error(msg, err)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(jsonData)

		Log.Infos("Envoie des salles disponibles avec succès !")
	}

}

func GetAllRoomsHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodGet {
		var msg = "Vous ne pouvez pas faire de requête autre que GET pour ça :/"
		Log.Error(msg)
		http.Redirect(w, r, RouteGetAllRooms+"?message="+msg, http.StatusSeeOther)
		return
	}

	result := GetAllSalle()

	if len(result) == 0 {
		var msg = "Pas de salles à lister, veuillez en créer une."
		Log.Error(msg)
		http.Redirect(w, r, RouteCreateRoom+"?message="+msg, http.StatusSeeOther)
		return
	}

	templates.ExecuteTemplate(w, "salles.html", map[string]interface{}{
		"result":  result,
		"message": nil,
	})
}

func CreateRoomHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		templates.ExecuteTemplate(w, "creerSalles.html", map[string]interface{}{
			"message": nil,
			"result":  nil,
		})
	} else if r.Method == http.MethodPost {
		nomSalle := r.FormValue("nom")
		placeSalle := r.FormValue("place")

		if placeSalle == NullString || nomSalle == NullString {

			var msg = "Erreur mauvais saisie de données"
			Log.Error(msg)
			http.Redirect(w, r, RouteCreateRoom+"?message="+msg, http.StatusSeeOther)
			return
		}

		placeSalleInt, err := strconv.Atoi(placeSalle)
		if err != nil {
			var msg = "ID de salle invalide"
			Log.Error(msg)
			http.Redirect(w, r, RouteCreateRoom+"?message="+msg, http.StatusSeeOther)
			return
		}

		leBool := CreateRoom(&nomSalle, &placeSalleInt)
		if leBool == false {
			var msg = "Erreur lors de la création"
			Log.Error(msg)
			http.Redirect(w, r, RouteCreateRoom+"?message="+msg, http.StatusSeeOther)
			return
		}

		http.Redirect(w, r, RouteGetAllRooms, http.StatusSeeOther)
		return

	}
}

func CancelSalleHandler(w http.ResponseWriter, r *http.Request) {

	idSalle := r.URL.Query().Get("idSalle")

	newChoix, err := strconv.Atoi(idSalle)

	if err != nil {
		Log.Error("Impossible to convert the id from string to int")
		http.Error(w, "Impossible to convert the id from string to int : "+idSalle, http.StatusBadRequest)
		return
	}

	result := GetSalleById(&newChoix)
	if len(result) == 0 {
		Log.Error("Aucune salle contient cet ID")
		http.Error(w, "Aucune salle trouvée pour l'ID "+idSalle, http.StatusBadRequest)
		return
	}

	DeleteRoomByID(&newChoix)

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Salle supprimer avec succès"))
	Log.Infos("Salle supprimer avec succès !")
	return
}
