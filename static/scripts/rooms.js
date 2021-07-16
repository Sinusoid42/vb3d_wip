/*

    Script loaded dynamically and starts the stream
    dynamically

    When loading the webpage with this script, the start

 */




const rooms = {

    USER : null,

    video_html_element : Element,

    face_video_canvas : Element,

    context : CanvasRenderingContext2D,

    Participants : [],

    _init : function init(){

        rooms.Participants = []



    },


    getParticipant : (id, cbfunc) =>{
       rooms.Participants.forEach((e) => {

           if (e.USER_ID == id){
               return new Promise((resolve, rejected) => {
                 resolve(e);
               });
           }
       });
        return new Promise((resolve, rejected) => {
            rejected();
        });
    },

    startRoomSession : function (server_url, ws_proto, room_id, user_id) {

        rooms.video_html_element = document.getElementById("video_webcam_element");
        rooms.face_video_canvas = document.getElementById("face_canvas");
        rooms.context = rooms.face_video_canvas.getContext("2d");

        rooms._init(); //Initalize the buffer, create the room participant storage

        rooms.USER = rooms.newMember(user_id)

        video_api.init(rooms.video_html_element);

        stream_api.init(server_url, ws_proto, room_id, user_id);

    },

    newMember: function newMember(user_id) {
        return new class {
            constructor() {
                this.USER_ID = user_id      //The Member ID, the stream member is identified by, see HEAD_STREAM definition, in order to view the byte stream information
                // The in/out byte stream buffer, for fast sending of dataT
                this.STREAM_BUFFER = new Uint8Array(stream_api.BUFFER_SIZE)

                this.io_constaints = new class {
                    constructor() {
                        this.video = true;
                        this.audio = false;
                    }
                }

                this.image_buffer = new Image(stream_api.stream_constraints.BUFFER_PIXEL_WIDTH, stream_api.stream_constraints.BUFFER_PIXEL_HEIGHT);

                this.audio_buffer_0 = new Audio("");
                this.audio_buffer_1 = new Audio("");

                this.audio_buffer_0.innerHTML = `<source src=""\>`;


                this.audio_playing_0 = false;
                this.audio_playing_1 = false;
                //Simple integer upcount check to prevent the video or audio, from playing in incorrect order
                //My personal belief is, that skipped audio and video is not as worse, as repeated and intersected audio and video
                this.STREAM_TIMESTAMP = 0
                this.position = new class {
                    constructor(x, y, z) {
                        this.x = x;
                        this.y = y;
                        this.z = z;
                    }
                }
                this.rotation = new class {
                    constructor(r_x, r_y, r_z) {
                        this.r_x = r_x;
                        this.r_y = r_y;
                        this.r_z = r_z;
                    }
                }
            }

            getThreePosition() {

                //parse data

                //write data to threejs vector

                //TODO implement the THREE JS Library and return a THREE JS Vector on call
                return null
            }
            getThreeRotation() {

                //TODO implement the THREE JS Library and return a THREE JS Vector on call
                return null
            }
            getTimeStamp() {
                this.STREAM_TIMESTAMP += 1;
                return this.STREAM_TIMESTAMP;
            }
            setStreamID(user_id){
                this.USER_ID = user_id
            }
        }
    },
}













