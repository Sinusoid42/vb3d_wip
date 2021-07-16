package model

import (
	"encoding/base64"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"log"
	"web_projekt/v7/app/controller/stream_api"
	"web_projekt/v7/app/model/database"
	"web_projekt/v7/app/model/utils"
)

type Stream_User struct {
	/*
		Local User identification
		For Chat/Video and Voice, the misc. user info is relevant
	*/
	user_id string
	room_id string
	/*
		The Timestamp for the datastream and ready flag
		@see Participant.parse()
	*/
	last_time_stamp uint32
	buffer_ready    bool
	/*
		ID's
	*/

	avatar_data player_controller

	/*
		Holds the buffered Bytestream which is send or recieved, when the user has enabled webcam
	*/
	stream_data   []byte
	playload_size uint32
}

/*
	Initializes a struct Stream_User

	This object is stored actively in RAM on the Server
	being the Buffer for the Stream data


*/
func (p *Stream_User) Init() (bool, error) {

	p.stream_data = make([]byte, stream_api.BUFFER_SIZE)
	p.buffer_ready = false
	return true, nil
}

/*
	Defines a User Account within the Database
	@see couchDB

	Will be a json representation of the object that
	was created on registration
*/
type User struct {
	REV string `json:"_rev"`

	ID string `json:"_id"`

	user_stream_id string `json:"user_id"`

	user_password string `json:"user_password"`

	user_name string `json:"user_name"`

	user_email_address string `json:"user_email_address"`
}

/*
	Defines the float 32 player controller infos

	=> Can be used to only send data to the subscriber for any datastream, if the clients are of a minimum distance of each other
*/
type player_controller struct {
	x, y, z    float32 //player position
	rx, ry, rz float32 //Euler Angles Rotation of the Player
}

func (u *User) M2U(m map[string]interface{}) *User {

	u.ID = m["_id"].(string)
	u.REV = m["_rev"].(string)
	u.user_stream_id = m["user_stream_id"].(string)
	u.user_password = m["user_password"].(string)
	u.user_email_address = m["user_email_address"].(string)

	return u
}

/*
	Checks wether a given Username and Password matches
	some within the database

	returns

*/
func Check(u string, p string) (bool, bool, error) {
	exist := false
	auth := false
	query := `{
		"selector":{
			"user_name":"%s"
			}
	}`
	um, err := database.Database.QueryJSON(fmt.Sprintf(query, u))
	if err != nil || len(um) == 0 {

		fmt.Println("The user does probbaly not exist")
		return auth, exist, err
	}
	exist = true
	if len(p) == 0 {

		return auth, exist, nil
	}

	us := User{}
	us = *us.M2U(um[0])
	passwordDB, _ := base64.StdEncoding.DecodeString(us.GetPassword())
	err = bcrypt.CompareHashAndPassword(passwordDB, []byte(utils.Decode(p)))
	if err != nil {
		fmt.Println("Wrong pw")
		return auth, exist, nil

	}
	auth = true
	fmt.Println("Login Success: " + us.user_name + ", s_id: " + us.user_stream_id)

	return auth, exist, nil
}

func (u *User) GetPassword() string {
	return u.user_password
}

/*

	Decrypts the base 10-shift13 encryption (rot 13 algorithm)
	Creates a hash encrypted web version of the password, send by the user
	<br>
	@returns a new User struct

*/
func NewUser(user_name string, user_password_shift_13 string, user_email_address string) User {

	pwb64 := utils.Decode(user_password_shift_13)

	pw, _ := utils.EncryptUserPassword(pwb64)

	u := User{
		REV: "", //Set by couchDB
		ID:  "", //Set by couchDB

		user_stream_id:     utils.GenerateStreamID(),
		user_password:      pw,
		user_name:          user_name,
		user_email_address: user_email_address,
	}
	return u

}

/*
	Setter for private variable user_stream_id
*/
func (u *User) SetStreamID(stream_id string) {
	u.user_stream_id = stream_id
}

func (u *User) GetStreamID() string {
	return u.user_stream_id
}

func GetUserByStreamID(id string) (*User, error) {

	query := `{
		"selector":{
			"user_stream_id":"%s"
			}
	}`
	um, err := database.Database.QueryJSON(fmt.Sprintf(query, id))
	if err != nil || len(um) == 0 {

		fmt.Println("The user does probbaly not exist")
		return nil, err
	}

	us := User{}
	us = *us.M2U(um[0])
	return &us, nil
}

func (u *User) U2M() (map[string]interface{}, error) {
	m := make(map[string]interface{})
	m["_id"] = u.ID
	m["_rev"] = u.REV
	m["user_stream_id"] = u.user_stream_id
	m["user_name"] = u.user_name
	m["user_password"] = u.user_password
	m["user_email_address"] = u.user_email_address
	return m, nil
}

func (u *User) Store() (bool, error) {
	m, _ := u.U2M()
	delete(m, "_id")
	delete(m, "_rev")
	id, rev, err := database.Database.Save(m, nil)
	if err != nil {
		log.Fatal(err)
		return false, err
	}
	u.ID = id
	u.REV = rev

	return true, nil
}

func GetUser(un string) (*User, error) {
	query := `{
		"selector":{
			"user_name":"%s"
			}
	}`
	um, err := database.Database.QueryJSON(fmt.Sprintf(query, un))
	if err != nil || len(um) == 0 {
		return nil, err
	}
	u := User{}
	u.M2U(um[0])
	return &u, nil
}

func (u *User) Save() error {
	m, _ := u.U2M()
	err := database.Database.Set(u.ID, m)
	return err
}

func (u *User) ToStreamUser(room_id string) *Stream_User {
	us := Stream_User{}
	us.Init()
	us.user_id = u.user_stream_id
	us.room_id = room_id
	return &us
}

func (u *Stream_User) StreamBuffer() []byte {
	return u.stream_data
}

func (u *Stream_User) StreamUserID() string {
	return u.user_id
}

func (u *Stream_User) StreamRoomID() string {
	return u.room_id
}

func (u *Stream_User) StoreBuffer(buffer []byte) {
	u.stream_data = buffer
}
