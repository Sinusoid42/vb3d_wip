package stream

import (
	"errors"
	"web_projekt/v6/app/model"
)

var Session _session

func init() {
	Session = _session{}
	Session.Init()
}

type _session struct {
	active_rooms map[string]*model.Room

	count uint32
}

func (s *_session) Init() {
	s.active_rooms = make(map[string]*model.Room)
	s.count = 0
}

func (s *_session) Store(room *model.Room) error {

	//TODO
	//check if the id is already existing
	//No user may connect twice
	//user stream setup

	s.active_rooms[room.GetRoomID()] = room
	s.count += 1

	return nil

}

func (s *_session) IsActive(room_id string) bool {
	r, _ := s.GetRoomByID(room_id)
	return r != nil
}

func (s *_session) GetRoomByID(room_id string) (*model.Room, error) {
	r := s.active_rooms[room_id]
	if r != nil {
		return r, nil
	}
	return nil, errors.New("The Room does not exist")
}
