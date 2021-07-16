package stream

import "web_projekt/v7/app/model"

func NewRoom(Host *model.User) (*model.Room, error) {

	room := model.Room{}
	room.Init()
	su := Host.ToStreamUser(room.GetRoomID())
	room.SetHost(su)
	//room.StoreUser(su)
	return &room, nil
}
