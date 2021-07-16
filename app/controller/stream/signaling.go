package stream

import (
	"errors"
	"fmt"
	"golang.org/x/net/websocket"
	"unsafe"
	"web_projekt/v7/app/controller/stream_api"
	"web_projekt/v7/app/model"
)

var stream_buffer = make([]byte, stream_api.BUFFER_SIZE)
var event_buffer = make([]byte, stream_api.EVENT_BUFFER_SIZE)

func Rooms_Event(ws *websocket.Conn) {

	for {

		if err := websocket.Message.Receive(ws, &stream_buffer); err != nil {
			continue
		} else {
			fmt.Println("Hello World")
		}

	}
}

func Rooms_Stream(ws *websocket.Conn) {
	var err error
	var r_id = make([]byte, 16)
	var m_id = make([]byte, 16)
	var i = 0
	var uri = &ws.Request().URL.Path
	var room *model.Room
	var user *model.Stream_User

	var buffer []byte

	for {
		for ; i < 16; i++ {
			if 7+i < len(*uri) {
				r_id[i] = (*uri)[7+i]
			}
			if 24+i < len(*uri) {
				m_id[i] = (*uri)[24+i]
			}
		}
		i = 0

		room, _ = Session.GetRoomByID(string(r_id))
		user, _ = room.GetUser(string(m_id))

		buffer = user.StreamBuffer()
		if err = websocket.Message.Receive(ws, &buffer); err != nil {
			continue
		}
		user.StoreBuffer(buffer)

		for k, u := range room.GetParticipants() {
			if k != user.StreamUserID() {
				/*
					Parse the data to other streaming members
				*/
				buffer = u.StreamBuffer()
				websocket.Message.Send(ws, buffer)
			}
		}

		//l , _ := parseB2I(stream_buffer, stream_api.DATA_STREAM_STREAM_MEMBER_ID_LENGTH)
		//room_id, _ := ParseB2S(stream_buffer, stream_api.Stream_api.DATA_STREAM_STREAM_ROOM_ID, 16)

		//println(room_id)

	}
}

// Float32bits returns the IEEE 754 binary representation of f.
func Uint32fromFloat32(f float32) uint32 {
	return *(*uint32)(unsafe.Pointer(&f))
}

// Float32frombits returns the floating point number corresponding
// to the IEEE 754 binary representation b.
func Float32FromUint32(b uint32) float32 {
	return *(*float32)(unsafe.Pointer(&b))
}

func ParseB2I(b []byte, offset int) (uint32, error) {
	if offset < 0 || offset >= len(b) {
		return 0, errors.New("Offset index out of Bounce")
	}
	return uint32(b[offset+0]<<24) | uint32(b[offset+1]<<16) | uint32(b[offset+2]<<8) | uint32(b[offset+3]<<0), nil
}

func ParseB2S(b []byte, offset int, l int) (string, error) {
	var s string
	var i = 0
	if offset+l >= len(b) {
		return "", errors.New("Index Out of Bounce for parsing ByteStream => String")
	}
	for ; i < l; i++ {
		s += string(b[i+offset])
	}
	return s, nil
}

func ParseB2F32(b *[]byte, offset int) (float32, error) {
	e, err := ParseB2I(*b, offset)
	return Float32FromUint32(e), err
}

func WriteI2B(b []byte, offset int, I32 uint32) error {
	if offset < 0 || offset >= len(b) {
		return errors.New("Index Out of Bounce for parsing Integer => ByteStream")
	}
	b[offset+0] = byte((I32 >> 24) & 0xFF)
	b[offset+1] = byte((I32 >> 16) & 0xFF)
	b[offset+2] = byte((I32 >> 8) & 0xFF)
	b[offset+3] = byte((I32 >> 0) & 0xFF)

	return nil
}

func WriteS2B(b []byte, offset int, s string, l int) error {
	if offset < 0 || l+offset >= len(b) {
		return errors.New("Index Out of Bounce for parsing String => ByteStream")
	}
	t := []byte(s)
	var i = 0
	for ; i < l; i++ {
		b[offset+i] = t[i]
	}
	return nil
}

func WriteF322B(b *[]byte, offset int, f32 float32) error {
	return WriteI2B(*b, offset, Uint32fromFloat32(f32))
}
