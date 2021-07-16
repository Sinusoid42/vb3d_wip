package controller

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"golang.org/x/net/websocket"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"web_projekt/v7/app/auth"
	"web_projekt/v7/app/controller/stream"
	"web_projekt/v7/app/model"
	idutils "web_projekt/v7/app/model/utils"
	"web_projekt/v7/app/utils"
)

var templates *template.Template

const LOGIN_FORM_USER_NAME = "user_name"
const LOGIN_FORM_USER_PASSWORD = "user_password"
const REGISTER_FORM_USER_NAME = "user_name"
const REGISTER_FORM_USER_PASSWORD = "user_password"
const REGISTER_FORM_USER_EMAIL_ADDRESS = "user_email_address"

const SERVER_ADRESS = "localhost"
const SERVER_PORT = "8080"
const SERVER_PROTOCOL_HTTP = "http"
const SERVER_PROTOCOL_WEBSOCKET = "ws"

const INDEX_ADDRESS = "/"
const INDEX_ADDRESS_1 = "/index"
const ROOMS_ADDRESS = "/rooms"
const ROOMS_ID_ADDRESS = "/rooms/{room_id}"
const ROOMS_EVENT_ADDRESS = "/rooms/{room_id}/event"
const LOGIN_ADDRESS = "/login"
const REGISTER_ADDRESS = "/register"
const USER_ADDRESS = "/user"
const LOGOUT_ADDRESS = "/logout"

type Sdp struct {
	session_description string
}

func init() {
	fmt.Println(utils.GetLocalEnv())
	templates = template.Must(template.ParseGlob(utils.GetLocalEnv() + "static/template/*.tmpl"))
}

func Index(w http.ResponseWriter, r *http.Request) {
	c, err := auth.GetCookieStore().Get(r, auth.GetSessionCookie())
	if err != nil {
		b := struct {
			Auth bool `json:"Auth"`
		}{
			Auth: false,
		}
		templates.ExecuteTemplate(w, "index.tmpl", b)
		return
	}
	p := c.Values[auth.Authenticated]

	if p != nil {
		p = p.(bool)
	} else {
		p = false
	}
	m := make(map[string]interface{})

	m["auth"] = p.(bool)
	m["user_id"] = ""
	m["room_id"] = "roomId"
	m["exist"] = false
	m["server_protocol"] = SERVER_PROTOCOL_HTTP
	m["server_address"] = SERVER_ADRESS
	m["server_port"] = SERVER_PORT
	m["server_websocket_protocol"] = SERVER_PROTOCOL_WEBSOCKET

	templates.ExecuteTemplate(w, "index.tmpl", m)

}

type session struct {
	Room_id string `json:"room_id"`
}

func RoomsSession(w http.ResponseWriter, r *http.Request) {
	fmt.Println("\n\nNew RoomSession Request")
	fmt.Println(r.URL.Path)

	vars := mux.Vars(r)
	fmt.Println(vars)
	c, err := auth.GetCookieStore().Get(r, auth.GetSessionCookie())
	if err != nil {
		http.Redirect(w, r, INDEX_ADDRESS, http.StatusFound)
		return
	}
	p := c.Values[auth.Authenticated]
	user_id := c.Values[auth.UserID]
	roomId := vars["room_id"]
	if p != nil {
		p = p.(bool)
	} else {
		p = false
	}

	//TODO
	//IF Room was not created => make http redirect to Index
	//If not logged in => http open was /rooms/{room_id} => execute Template with promt, to create account or just join
	//If join worked => start stream via ajax call to this room

	//restful websocket communication via /rooms/{room_id}/stream
	if r.Method == http.MethodGet {
		//Just serve the streaming website for the room
		//create
		println("GET/rooms/{room_id}")
		if !stream.Session.IsActive(roomId) {
			println("Session not active")
			http.Redirect(w, r, INDEX_ADDRESS, http.StatusFound) //The Room Session does not exist, go back to index
			return
		}

		if !p.(bool) {

			println("Not logged in => clientside selection for login or ->connect any")

			if j, _ := strconv.ParseBool(r.FormValue("join")); j {

				m := make(map[string]interface{})
				m["auth"] = p
				m["user_id"] = ""
				m["room_id"] = roomId
				m["exist"] = true
				m["server_url"] = SERVER_ADRESS
				m["server_protocol"] = SERVER_PROTOCOL_HTTP
				m["server_address"] = SERVER_ADRESS
				m["server_port"] = SERVER_PORT
				m["server_websocket_protocol"] = SERVER_PROTOCOL_WEBSOCKET

				templates.ExecuteTemplate(w, "rooms.tmpl", m)
				return

			}

			m := make(map[string]interface{})
			m["auth"] = p
			m["user_id"] = ""
			m["room_id"] = roomId
			m["exist"] = true
			m["server_url"] = SERVER_ADRESS
			m["server_protocol"] = SERVER_PROTOCOL_HTTP
			m["server_address"] = SERVER_ADRESS
			m["server_port"] = SERVER_PORT
			m["server_websocket_protocol"] = SERVER_PROTOCOL_WEBSOCKET
			templates.ExecuteTemplate(w, "rooms.tmpl", m)
			return

		} else {

			println("Is logged in, stream starts auto")

			fmt.Println(roomId)
			fmt.Println(user_id)

			if r, _ := stream.Session.GetRoomByID(roomId); r.Contains(user_id.(string)) {

				m := make(map[string]interface{})
				m["auth"] = p
				m["user_id"] = user_id.(string)
				m["room_id"] = roomId
				m["exist"] = true
				m["server_url"] = SERVER_ADRESS
				m["server_websocket_protocol"] = SERVER_PROTOCOL_WEBSOCKET
				templates.ExecuteTemplate(w, "rooms.tmpl", m)
				println("The user already exists within the room")
				return
			}

			fmt.Println("The user is authenticated and joined the room")
			//TODO

			peerConnectionConfig := webrtc.Configuration{
				ICEServers: []webrtc.ICEServer{
					{
						URLs: []string{"stun:stun.l.google.com:19302"},
					},
				},
			}

			webrtc.PeerConnection{}

			mediaEngine := webrtc.MediaEngine{}
			mediaEngine.RegisterCodec(webrtc.NewRTPCodecType(webrtc.PayloadType(), 90000))
			api := webrtc.NewAPI(webrtc.WithMediaEngine(&mediaEngine))

			var session Sdp
			fmt.Println(session)
			offer := webrtc.SessionDescription{}
			peerConnection, _ := api.NewPeerConnection(peerConnectionConfig)

			peerConnection.SetRemoteDescription(offer)

			answer, _ := peerConnection.CreateAnswer(nil)

			err := peerConnection.SetLocalDescription(answer)

			fmt.Println(err)

			room, _ := stream.Session.GetRoomByID(roomId)
			room.StoreUserByID(user_id.(string))

			http.Handle("/rooms/"+roomId+"/"+user_id.(string)+"/stream", websocket.Handler(stream.Rooms_Stream))
		}
	}
	if r.Method == http.MethodPost {
		println("POST/rooms/{room_id}")
		if !p.(bool) { //user is not logged in
			println("Not logged in ")

			if !stream.Session.IsActive(roomId) {
				println("Session not active")

				//The Room Session does not exist, go back to index
				m := make(map[string]interface{})
				m["auth"] = p
				m["user_id"] = idutils.GenerateStreamID()
				m["room_id"] = roomId
				m["exist"] = stream.Session.IsActive(roomId)
				b, _ := json.MarshalIndent(m, "", "  ")

				w.Header().Set("Content-Type", "application/json")
				w.Write(b)
				return
			}

			if j, _ := strconv.ParseBool(r.FormValue("join")); j {
				fmt.Println("Still join the room")

				m := make(map[string]interface{})
				m["auth"] = p
				m["user_id"] = idutils.GenerateStreamID()
				m["room_id"] = roomId
				m["exist"] = stream.Session.IsActive(roomId)
				m["server_url"] = SERVER_ADRESS
				m["server_websocket_protocol"] = SERVER_PROTOCOL_WEBSOCKET
				room, _ := stream.Session.GetRoomByID(roomId)
				room.StoreUserByID(m["user_id"].(string))

				http.Handle("/rooms/"+roomId+"/"+m["user_id"].(string)+"/stream", websocket.Handler(stream.Rooms_Stream))

				b, _ := json.MarshalIndent(m, "", "  ")
				w.Header().Set("Content-Type", "application/json")
				w.Write(b)

				return
			}

			m := make(map[string]interface{})
			m["auth"] = p
			m["user_id"] = ""
			m["room_id"] = roomId
			m["exist"] = stream.Session.IsActive(roomId)
			m["server_url"] = SERVER_ADRESS
			m["server_websocket_protocol"] = SERVER_PROTOCOL_WEBSOCKET
			b, _ := json.MarshalIndent(m, "", "  ")
			w.Header().Set("Content-Type", "application/json")
			w.Write(b)
			return
		} else {

			if !stream.Session.IsActive(roomId) {
				println("Session not active")

				//The Room Session does not exist, go back to index
				m := make(map[string]interface{})
				m["auth"] = p
				m["user_id"] = idutils.GenerateStreamID()
				m["room_id"] = roomId
				m["exist"] = stream.Session.IsActive(roomId)
				b, _ := json.MarshalIndent(m, "", "  ")

				w.Header().Set("Content-Type", "application/json")
				w.Write(b)
				return
			}

		}

	}

	m := make(map[string]interface{})
	m["auth"] = true
	m["user_id"] = user_id.(string)
	m["room_id"] = roomId
	m["server_url"] = SERVER_ADRESS
	m["server_websocket_protocol"] = SERVER_PROTOCOL_WEBSOCKET
	fmt.Println(m)

	templates.ExecuteTemplate(w, "rooms.tmpl", m)

}

func Rooms(w http.ResponseWriter, r *http.Request) {

	c, _ := auth.GetCookieStore().Get(r, auth.GetSessionCookie())
	p := c.Values[auth.Authenticated]

	if p == nil || !p.(bool) {
		http.Redirect(w, r, INDEX_ADDRESS, http.StatusFound)
		return
	}
	u := c.Values[auth.UserID]

	fmt.Println(u)
	user, _ := model.GetUserByStreamID(u.(string))
	newroom, _ := stream.NewRoom(user)

	//Open Websockets for the datastream
	//???
	//http.redirect to the room

	if !stream.Session.IsActive(newroom.GetRoomID()) {
		fmt.Println("Creating the Server Event Broadcast socket")

		http.Handle("/rooms/"+newroom.GetRoomID()+"/event", websocket.Handler(stream.Rooms_Event))
	}

	stream.Session.Store(newroom)
	fmt.Println("New Meeting Room has been created")

	fmt.Println(stream.Session)
	fmt.Println(newroom)
	fmt.Println(p)

	http.Redirect(w, r, ROOMS_ADDRESS+"/"+newroom.GetRoomID(), http.StatusFound)
	return
	//templates.ExecuteTemplate(w, "rooms.tmpl", nil)
}

/*
	Handle /login request for the server
	when creating a new user
*/
func Login(w http.ResponseWriter, r *http.Request) {
	c, _ := auth.GetCookieStore().Get(r, auth.GetSessionCookie())
	p := c.Values[auth.Authenticated]
	if p != nil {
		if p.(bool) {
			http.Redirect(w, r, INDEX_ADDRESS, http.StatusFound)
			return
		}
	}
	if r.Method == http.MethodPost {
		u := r.FormValue(LOGIN_FORM_USER_NAME)
		p := r.FormValue(LOGIN_FORM_USER_PASSWORD)
		ok, exist, err := model.Check(u, p)
		if err != nil {
			log.Fatal(err)
		}
		s := r.FormValue("submit")
		if submit, _ := strconv.ParseBool(s); ok && submit { //ok
			a := map[string]interface{}{}
			a["auth"] = true
			a["exist"] = true

			us, _ := model.GetUser(u)

			c.Values[auth.Authenticated] = true
			c.Values[auth.UserID] = us.GetStreamID()
			c.Save(r, w)

			b, _ := json.MarshalIndent(a, "", "  ")
			w.Header().Set("Content-Type", "application/json")
			w.Write(b)
			return
		}
		if exist {
			a := map[string]interface{}{}
			a["auth"] = ok
			a["exist"] = exist

			b, _ := json.MarshalIndent(a, "", "  ")
			w.Header().Set("Content-Type", "application/json")
			w.Write(b)
			return
		}
		a := map[string]interface{}{}
		a["auth"] = ok
		a["exist"] = exist

		b, _ := json.MarshalIndent(a, "", "  ")
		w.Header().Set("Content-Type", "application/json")
		w.Write(b)
	}
	if r.Method == http.MethodGet {
		m := make(map[string]interface{})

		m["auth"] = false
		m["user_id"] = ""
		m["room_id"] = ""
		m["exist"] = false
		m["server_protocol"] = SERVER_PROTOCOL_HTTP
		m["server_address"] = SERVER_ADRESS
		m["server_port"] = SERVER_PORT
		m["server_websocket_protocol"] = SERVER_PROTOCOL_WEBSOCKET

		templates.ExecuteTemplate(w, "login.tmpl", m)
		return
	}
}

/*
	Registers a new User and stores its data user_data encrypted within
	the couchDB database
*/
func Register(w http.ResponseWriter, r *http.Request) {
	c, _ := auth.GetCookieStore().Get(r, auth.GetSessionCookie())
	p := c.Values[auth.Authenticated]
	if p != nil && p.(bool) {
		http.Redirect(w, r, INDEX_ADDRESS, http.StatusFound)
		return
	}
	if r.Method == http.MethodGet {
		m := make(map[string]interface{})

		m["auth"] = false
		m["user_id"] = ""
		m["room_id"] = ""
		m["exist"] = false
		m["server_protocol"] = SERVER_PROTOCOL_HTTP
		m["server_address"] = SERVER_ADRESS
		m["server_port"] = SERVER_PORT
		m["server_websocket_protocol"] = SERVER_PROTOCOL_WEBSOCKET
		templates.ExecuteTemplate(w, "register.tmpl", m)
		return
	}
	if r.Method == http.MethodPost {
		if ok, _ := strconv.ParseBool(r.FormValue("submit")); ok {
			u := model.NewUser(
				r.FormValue(REGISTER_FORM_USER_NAME),
				r.FormValue(REGISTER_FORM_USER_PASSWORD),
				r.FormValue(REGISTER_FORM_USER_EMAIL_ADDRESS))
			u.Store()
		}
		usnm := r.FormValue(REGISTER_FORM_USER_NAME)
		uspw := r.FormValue(REGISTER_FORM_USER_PASSWORD)
		auth := false
		exist := false
		if usnm != "" || uspw != "" {
			auth, exist, _ = model.Check(usnm, uspw)
		}
		a := map[string]interface{}{}
		a["auth"] = auth
		a["exist"] = exist

		b, _ := json.MarshalIndent(a, "", "  ")
		w.Header().Set("Content-Type", "application/json")

		w.Write(b)
		return
	}
	m := make(map[string]interface{})

	m["auth"] = false
	m["user_id"] = ""
	m["room_id"] = ""
	m["exist"] = false
	m["server_protocol"] = SERVER_PROTOCOL_HTTP
	m["server_address"] = SERVER_ADRESS
	m["server_port"] = SERVER_PORT
	m["server_websocket_protocol"] = SERVER_PROTOCOL_WEBSOCKET
	templates.ExecuteTemplate(w, "register.tmpl", m)
}

/*
	Logs a User out of the server and removes it from the current session users
*/
func Logout(w http.ResponseWriter, r *http.Request) {

	c, _ := auth.GetCookieStore().Get(r, auth.GetSessionCookie())
	p := c.Values[auth.Authenticated]
	if p != nil {
		if p.(bool) {
			c.Values[auth.Authenticated] = false
			c.Values[auth.UserID] = ""
			c.Save(r, w)

			http.Redirect(w, r, LOGIN_ADDRESS, http.StatusFound)
			return
		}
	}
	http.Redirect(w, r, LOGIN_ADDRESS, http.StatusFound)
}
