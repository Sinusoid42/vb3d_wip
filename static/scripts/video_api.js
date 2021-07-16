"use scrict";




/*
    TODO


    Mute Mic, Stop video => on Button


    create room from id
    create new member stream, when joining a room

    parse player data by stream

    retrieve Participant[]





 */


const video_api = {

    the_video_stream : null,
    the_audio_stream : null,

    audio_recorder_0 : MediaRecorder,
    audio_recorder_1 : MediaRecorder,

    audio_recorder_toggle : true,
    audio_ready : 0,


    _video_buffer : class{
        constructor() {
            this.video = "";
            this.audio = "";
            this.audio_ready = false;
        }

    },

    _video_constraints : class {
        constructor() {
            this.video = true;
            this.audio = false;
        }
    },

    _audio_constraints : class {
        constructor() {
            this.video = false;
            this.audio = true;
        }
    },

    audio_recorder_constraints : {
        audioBitsPerSecond : 16000,
        mimeType : "audio/webm",
    },

    _face_roi : class {
        constructor() {
            this.x = 0,
            this.y = 0,
            this.w = 0,
            this.h = 0,

            this.x_offset = 0,
            this.y_offset = 0,
            this.w_offset = 0,
            this.h_offset = 0
        }
    },

    initialize_nets : function a() {
        Promise.all([
            faceapi.nets.tinyFaceDetector.loadFromUri('/models/'),
            faceapi.nets.faceRecognitionNet.loadFromUri('/models/'),
            faceapi.nets.faceLandmark68Net.loadFromUri('/models/')
        ]).then(video_api.start_Video_Capture);
     },

    /**
     * This function is the main init function using promisses to start the recording
     * when the initalization is done
     */
    start_Video_Capture : function (){
        navigator.mediaDevices.getUserMedia(video_api.video_constraints)
        .then((videostream) => {
            rooms.video_html_element.srcObject = videostream;
            video_api.the_video_stream = videostream;

            navigator.mediaDevices.getUserMedia(video_api.audio_contraints)
                .then((audiostream) => {
                    video_api.the_audio_stream = audiostream;

                    //start the Audio Capture
                    video_api.start_Audio_Capture(audiostream);

                    //roi fetching and audio fetching
                    setInterval(video_api.fetchMedia, 100); //approximate 20fps for detections

                    //start the streaming from the stream API
                    stream_api.start_stream(stream_api.interval);
                });
        });
    },

    /**
     * Starts to capture the audio by the user
     */
    start_Audio_Capture : async (audiostream) => {

        video_api.audio_recorder_0 = new MediaRecorder(audiostream, video_api.audio_recorder_constraints)
        video_api.audio_recorder_1 = new MediaRecorder(audiostream, video_api.audio_recorder_constraints)

        console.log("obtainedn audio stream")

        const a1 = document.createElement("audio")

        a1.innerHTML =`<source src = ""/>`
        const a2 = document.createElement("audio")

        a2.innerHTML =`<source src = ""/>`
        document.body.append(a1)
        document.body.append(a2)
        video_api.audio_recorder_0.ondataavailable =  async(e) => {

            var fr = new FileReader();
            fr.onloadend =  () => {


                video_api.audio_ready = 0;
                video_api.buffer.audio_ready = true;
                video_api.buffer.audio = fr.result;

                var al =fr.result.length;
                encoder_api.writeI2B(rooms.USER.STREAM_BUFFER,
                    video_api.audio_ready,
                    stream_api.offsets.DATA_STREAM_AUDIO_RECORDER_ID);

                encoder_api.writeI2B(rooms.USER.STREAM_BUFFER,
                    al,
                    stream_api.offsets.DATA_STREAM_AUDIO_PAYLOAD_BYTE_LENGTH_OFFSET);
                encoder_api.writeS2B(rooms.USER.STREAM_BUFFER,
                    fr.result,
                    stream_api.offsets.DATA_STREAM_AUDIO_PAYLOAD_OFFSET, al);

            }
            fr.readAsDataURL(e.data)
        }
        video_api.audio_recorder_1.ondataavailable =  async(e) => {

            var fr = new FileReader();
            fr.onloadend =  () => {

                video_api.audio_ready = 1
                video_api.buffer.audio_ready = true;
                video_api.buffer.audio = fr.result;


                var al =fr.result.length;
                encoder_api.writeI2B(rooms.USER.STREAM_BUFFER,
                    video_api.audio_ready,
                    stream_api.offsets.DATA_STREAM_AUDIO_RECORDER_ID);

                encoder_api.writeI2B(rooms.USER.STREAM_BUFFER,
                    al,
                    stream_api.offsets.DATA_STREAM_AUDIO_PAYLOAD_BYTE_LENGTH_OFFSET);
                encoder_api.writeS2B(rooms.USER.STREAM_BUFFER,
                    fr.result,
                    stream_api.offsets.DATA_STREAM_AUDIO_PAYLOAD_OFFSET, al);

            }
            fr.readAsDataURL(e.data)
        }
        video_api.audio_recorder_0.start();
    },


    /**
        Takes a snapshot in {@code stream_constraints#BUFFER_IMAGE_TYPE}
        format with given quality in order to pack the image
        <br><br>
        Stores the resized image data as <u><i>base 64</u></i> encoded
        data, that will be streamed
     */
    resizeToImageBuffer : function (video_buffer) {
        rooms.context.drawImage(rooms.video_html_element,
            video_api.FACE_ROI.x,
            video_api.FACE_ROI.y,
            video_api.FACE_ROI.w,
            video_api.FACE_ROI.h,
            0,
            0,
            stream_api.stream_constraints.BUFFER_PIXEL_WIDTH,
            stream_api.stream_constraints.BUFFER_PIXEL_HEIGHT);

        video_api.buffer.video = rooms.face_video_canvas.toDataURL(
            stream_api.stream_constraints.BUFFER_IMAGE_TYPE,
            stream_api.stream_constraints.BUFFER_IMAGE_QUALITY);

        video_api.image_buffer.src = video_api.buffer.video;
    },

    /**
     *  Clears generated data urls
     *  TODO: garbage collection to prevent data accumulation when streaming video data
     */
    clearGC : function clearGC(){
        URL.revokeObjectURL(video_api.buffer.video)
        URL.revokeObjectURL(video_api.buffer.audio)
    },

    /**
     *  Calls the face-api to detect the first found face,
     *  if no face was detected the function will be escaped
     *  <br>br>
     *  Stores the ROI information (region of interest) as AABB
     *  within the {@ FACE_ROI} encapsulated object
     */
    fetchMedia : async function(){


        if (video_api.audio_recorder_toggle){
            video_api.audio_recorder_1.start();
            try{
                setTimeout( () => {

                    video_api.audio_recorder_toggle = !video_api.audio_recorder_toggle;
                    video_api.audio_recorder_0.stop()
                }, 40);
            }
            catch (e) {
            }

        }
        else {
            video_api.audio_recorder_0.start();
            try {

                setTimeout(() => {
                    video_api.audio_recorder_toggle = !video_api.audio_recorder_toggle;
                    video_api.audio_recorder_1.stop()
                }, 40);
            }
            catch (e) {
            }
        }

        const detections = await faceapi.detectAllFaces(video_api.HTML_Video_Element,video_api.face_api_constraints
            ).withFaceLandmarks();
        if (detections[0] == undefined) {
            //Error handling, if the detections array is undefined, just escape
            return;
        }

        video_api.FACE_ROI.x = parseFloat(detections[0].alignedRect._box.x);
        video_api.FACE_ROI.y = parseFloat(detections[0].alignedRect._box.y);
        video_api.FACE_ROI.w = parseFloat(detections[0].alignedRect._box.width);
        video_api.FACE_ROI.h = parseFloat(detections[0].alignedRect._box.height);

        video_api.FACE_ROI.x += video_api.FACE_ROI.x_offset;
        video_api.FACE_ROI.y += video_api.FACE_ROI.y_offset;
        video_api.FACE_ROI.w += video_api.FACE_ROI.w_offset;
        video_api.FACE_ROI.h += video_api.FACE_ROI.h_offset;
    },

    /**
     * Initializes the video api entirely,
     *
     * starts the video stream after all media sources have successfully been loaded
     * @param the_video_element
     */
    init : (the_video_element) => {

        video_api.video_constraints = new video_api._video_constraints();
        video_api.audio_contraints = new video_api._audio_constraints();

        video_api.HTML_Video_Element = the_video_element;

        video_api.FACE_ROI = new video_api._face_roi();

        video_api.buffer = new video_api._video_buffer();

        //the Image for the current User as imagebuffer
        video_api.image_buffer = new Image(stream_api.stream_constraints.BUFFER_PIXEL_WIDTH, stream_api.stream_constraints.BUFFER_PIXEL_HEIGHT);

        //faceapi detection contraints
        video_api.face_api_constraints = new faceapi.TinyFaceDetectorOptions();

        //promise based network init
        video_api.initialize_nets();


    }
}
