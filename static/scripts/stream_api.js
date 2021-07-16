/*
    The stream_api lib contains all functions
    to setup a bytestream and delivers functionalities
    to encode float32, string, int and binary data into
    a websocket bytestream, which can be send to the server

    Provided are offset objects and member objects corresponding
    to the project of vb3d (video based 3D avatars), the video
    conferencing tool for the 6th semester project in webenginnering

    students are:
        jonas scholz
        katharina jüttner
        soheil rassouli
        marvin winkler
        vincent benedict winter


        author@ ben
 */


/*
    Stream related code will be wrapped in a const <i>stream_api<i> placeholder
    so the programmer or any other user, is able to easily communicate with the stream
    api using the dot syntax =>
            lines of code...
            var o = stream_api.offsets(); //object representation of all offsets within the bytestream
 */



const stream_api = {

    /**
     * Initializes the stream api, socket connections and data transport
     * @param server_url
     * @param room_id
     * @param user_id
     */
    init : function init(server_url, ws_proto, room_id, user_id) {

        stream_api.stream_session.roomID = room_id
        console.log(ws_proto + "WEBSOCKET PROTOCOL")
        //initialize the sockets
        stream_api.WSConnection = new WebSocket(ws_proto+"://"+server_url+":8080/rooms/" +room_id+ "/"+user_id+"/stream")
        stream_api.EventSocketConnection = new WebSocket(ws_proto+"://" + server_url + ":8080/rooms/" + room_id + "/event");

        //set the datatype
        stream_api.WSConnection.binaryType = "arraybuffer";
        stream_api.EventSocketConnection.binaryType = "arraybuffer";

        //recieve data from the websocket for the image user data stream
        stream_api.WSConnection.onmessage = (e) => {
            stream_api.stream_buffer = new Uint8Array(e.data)
            stream_api.recieveMediaStream(stream_api.stream_buffer)
        }

        //recieve data from the event websocket
        stream_api.EventSocketConnection.onmessage = (e) => {
            stream_api.event_buffer = new Uint8Array(e.data)
            stream_api.recieveEventStream(stream_api.event_buffer)
        }
    },

    //byte buffer 5 kbyte
    // go with 5 kByte for 30 fps => 150 kByte/2 upstream and downstream + http/tcp protocol definitions
    // still subject of further inspection, possibly even smaller buffer is possible when using 128x128 video + 1 second buffered audio as base64

    BUFFER_SIZE : 4500,

    EVENT_BUFFER_SIZE : 256,

    //defines the size of the header
    HEAD_SIZE : 128,


    WSConnection : WebSocket,
    EventSocketConnection : WebSocket,

    stream_buffer : Uint8Array,
    event_buffer : Uint8Array,

    dec : new TextDecoder("utf-8"),
    enc : new TextEncoder(),

    interval : 17,

    stream_session : {

        roomID : "",
        userID : "",
    },


    stream_recieve_meta : {

        /*
            helper object, to better gain acces over parsed data, that have been accepted by the webclient
         */
        audio_length : 0,
        audio_offset : 0,


        video_length : 0,
        video_offset : 0,

        the_video : "",
        the_audio : "",
    },


    //bytestream offsets
    offsets: {

            /*
                Defines the offset of bytes within the streambuffer, that stores the timestamp for the send package
                <br>
                <i><u>type</u>: uint32 (go)</i>
            */
                DATA_STREAM_TIMESTAMP_OFFSET : 0,
            /*
                Defines the offset of the userID for the member within the room stream
                This will contain a 16 byte long String representation of the user id, randomly serverside generated
                type: string (len = 16)
             */
                DATA_STREAM_STREAM_MEMBER_ID : 4,
                DATA_STREAM_STREAM_MEMBER_ID_LENGTH : 16,
            /*
                Defines the offset of the userID for the member within the room stream
                This will contain a 16 byte long String representation of the user id, randomly serverside generated
                type: string (len = 16)
            */
                DATA_STREAM_STREAM_ROOM_ID : 20,
                DATA_STREAM_STREAM_ROOM_ID_LENGTH : 16,
            /*
                Defines the byte length of the video payload within the datastream
                Is necessary when sending the data to the recieving client, such that the other client is able
                to disect the video stream data into slices
            */
                DATA_STREAM_VIDEO_PAYLOAD_BYTE_LENGTH_OFFSET : 36,
            /*
                Defines the offset at which the byte stream data for the video base64 data is stored and transmitted
            */
                DATA_STREAM_VIDEO_PAYLOAD_OFFSET : 128, //Buffer the video data within the stream directly after the head is
                DATA_STREAM_AUDIO_PAYLOAD_BYTE_LENGTH_OFFSET : 40,
            //This variable will be unused, the Audio Stream will immideatly after the video data be put into the videostream
            //This ensures maximum effiency, even if the buffer size is set with a prefixed size

                DATA_STREAM_AUDIO_PAYLOAD_OFFSET : 3500,


                DATA_STREAM_AUDIO_PAYLOAD_DURATION_OFFSET : 48,//TODO

                DATA_STREAM_PLAYER_POSITION_OFFSET :  52,
                DATA_STREAM_PLAYER_ROTATION_OFFSET :  64,

                DATA_STREAM_KILL_SWITCH_OFFSET : 76,


                DATA_STREAM_AUDIO_RECORDER_ID : 80,

    },
    //stream participant

    // streaming constraints
    stream_constraints : {
        BUFFER_PIXEL_WIDTH : 128,
        BUFFER_PIXEL_HEIGHT : 128,
        BUFFER_IMAGE_TYPE : "image/jpeg",
        BUFFER_IMAGE_QUALITY : 0.25,            //Can be adjusted to the network bandwidth
    },


    //packs the data of a user and the video/audio io buffers into the USER.STREAM_DATA []byte stream
    pack : function pack(USER, io_buffer){
        const ts =  USER.getTimeStamp();

        encoder_api.writeI2B(
            USER.STREAM_BUFFER,
            ts,
            stream_api.offsets.DATA_STREAM_TIMESTAMP_OFFSET);

        encoder_api.writeS2B(
            USER.STREAM_BUFFER,
            USER.USER_ID,
            stream_api.offsets.DATA_STREAM_STREAM_MEMBER_ID,
            stream_api.offsets.DATA_STREAM_STREAM_MEMBER_ID_LENGTH);

        encoder_api.writeS2B(
            USER.STREAM_BUFFER,
            stream_api.stream_session.roomID,
            stream_api.offsets.DATA_STREAM_STREAM_ROOM_ID,
            stream_api.offsets.DATA_STREAM_STREAM_ROOM_ID_LENGTH);

        const vl = io_buffer.video.length;
        encoder_api.writeI2B(USER.STREAM_BUFFER,
            vl ,
            stream_api.offsets.DATA_STREAM_VIDEO_PAYLOAD_BYTE_LENGTH_OFFSET);

        encoder_api.writeS2B(USER.STREAM_BUFFER,
            io_buffer.video,
            stream_api.offsets.DATA_STREAM_VIDEO_PAYLOAD_OFFSET,
            vl);


        if (video_api.buffer.audio_ready){

            if (video_api.audio_ready == 0) {
                rooms.USER.audio_buffer_0.src = video_api.buffer.audio;
                rooms.USER.audio_buffer_0.play();
                video_api.buffer.audio_ready = false;
            }
            else{
                rooms.USER.audio_buffer_1.src = video_api.buffer.audio;
                rooms.USER.audio_buffer_1.play();
            }
        }
        else{
            encoder_api.writeI2B(USER.STREAM_BUFFER,
                0,
                stream_api.offsets.DATA_STREAM_AUDIO_PAYLOAD_BYTE_LENGTH_OFFSET);
        }



    },


    sendMediaStream : function send(b) {
        if (stream_api.WSConnection != undefined) {
            stream_api.WSConnection.send(b);
        }
    },


    recieveEventStream : function recieveEventStream(buffer){






    },

    recieveMediaStream : function recieveMediaStream(buffer){

        var user_id_recieved = encoder_api.parseB2S(
            buffer,
            stream_api.offsets.DATA_STREAM_STREAM_MEMBER_ID,
            stream_api.offsets.DATA_STREAM_STREAM_MEMBER_ID_LENGTH);

        var vl = encoder_api.parseB2I(buffer, stream_api.offsets.DATA_STREAM_VIDEO_PAYLOAD_BYTE_LENGTH_OFFSET);
        var v = encoder_api.parseB2S(buffer, stream_api.offsets.DATA_STREAM_VIDEO_PAYLOAD_OFFSET, vl);
        stream_api.stream_recieve_meta.video_offset = stream_api.offsets.DATA_STREAM_VIDEO_PAYLOAD_OFFSET;

        stream_api.stream_recieve_meta.video_length = vl;
        stream_api.stream_recieve_meta.the_video = v;
        stream_api.stream_recieve_meta.audio_length = aL;

        var aL = encoder_api.parseB2I(buffer, stream_api.offsets.DATA_STREAM_AUDIO_PAYLOAD_BYTE_LENGTH_OFFSET);

        var ao = stream_api.offsets.DATA_STREAM_VIDEO_PAYLOAD_OFFSET + vl;
        var a = encoder_api.parseB2S(buffer, ao, aL);

        var a_id = encoder_api.parseB2I(buffer, stream_api.offsets.DATA_STREAM_AUDIO_RECORDER_ID);


        stream_api.stream_recieve_meta.audio_offset = ao;

        stream_api.stream_recieve_meta.the_audio = a;

        if (a_id == 0){
            if (!rooms.USER.audio_playing_0) {
                rooms.USER.audio_playing_0 = true;
                rooms.USER.audio_buffer_0.src = a;
                rooms.USER.audio_buffer_0.play();
                setTimeout(() => {rooms.USER.audio_playing_0=false}, 150)
            }
        }
        else {
            //a_id == 1
            if (!rooms.USER.audio_playing_1) {
                rooms.USER.audio_playing_1 = true;
                rooms.USER.audio_buffer_1.src = a;
                rooms.USER.audio_buffer_1.play();
                setTimeout(() => {rooms.USER.audio_playing_1=false}, 150)
            }
        }




        rooms.getParticipant(user_id_recieved).then((e)=> {
            console.log(e.USER_ID)




        }).catch(() => {

        });








        video_api.image_buffer.src = stream_api.stream_recieve_meta.the_video
        document.getElementById("recieve_canvas").getContext("2d").drawImage(video_api.image_buffer,0,0)

    },

    /*
          When the VideoStream by the Webcam is available,
          the websocket will start streaming the video data the server
    */
    start_stream : function (interval) {

        setInterval(stream_api.stream_data, interval);

    },

    /*
        Packs the necessary stream data into a Uint8ArrayBuffer and sends it via a websocket to the server
     */
    stream_data : function() {
        /*
            retains the current face from the video stream and packs it into a encoded data url for websocket transport
        */





        if (rooms.USER.io_constaints.video){
            video_api.resizeToImageBuffer(video_api.buffer);
        }
        /*
            checks wether the audio recorder (ms timeout see above) has finished recording and is able to be transported, then the data will be appended to the data stream
        */
        //checkForAudioBuffer();

        /*
            Packs all the data into the current datastream, that is send to the peers via the server
        */
        stream_api.pack(rooms.USER, video_api.buffer);

        video_api.clearGC();

        /*
            Clears any existing strings, dataurls and other buffered objects that are ready to be garbage collected
        */
        stream_api.sendMediaStream(rooms.USER.STREAM_BUFFER);
    },

    stream_event : () => {

    }
}