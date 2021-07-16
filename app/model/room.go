package model

import (
	"errors"
	"fmt"
	"web_projekt/v6/app/model/utils"
	"web_projekt/v6/app/scene"
)

type Room struct {
	room_id string `json:"room_id"` //The id of the room

	members map[string]*Stream_User //The members within the room

	host *Stream_User //The host of the room

	scene string //The shown scene as id

}

func (r *Room) Init() error {

	r.room_id = utils.GenerateStreamID()      //generates new random seed for the room id
	r.members = make(map[string]*Stream_User) //stores all
	r.scene = scene.DEFAULT
	return nil
}

func (r *Room) SetHost(u *Stream_User) {
	r.host = u
}

func (r *Room) GetRoomID() string {
	return r.room_id
}

func (r *Room) GetUser(user_id string) (*Stream_User, error) {
	u := r.members[user_id]
	var err error
	if u != nil {
		err = nil
	} else {
		err = errors.New("The User does not exist within this meeting")
	}
	return u, err
}

func (r *Room) Contains(user_id string) bool {
	e, _ := r.GetUser(user_id)
	return e != nil
}

func (r *Room) StoreUserByID(user_id string) *Stream_User {

	u := Stream_User{}
	u.Init()
	u.user_id = user_id
	u.room_id = r.room_id

	r.members[user_id] = &u

	fmt.Println(r.members)
	return &u
}

func (r *Room) StoreUser(user *Stream_User) *Stream_User {
	r.members[user.user_id] = user
	user.Init()
	return user
}

func (r *Room) ParticipantCount() int {
	return len(r.members)
}

func (r *Room) GetParticipants() map[string]*Stream_User {
	return r.members
}
