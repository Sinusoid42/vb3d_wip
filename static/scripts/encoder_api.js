






const encoder_api = {
    /*
        Writes a 4 byte Integer into a predefined byte buffer
        <br>
        <b><u>param:</b></u> b : The Byte Buffer<br>
        <b><u>param:</b></u> i : The Integer to be stored in the buffer<br>
        <b><u>param:</b></u> o : The offset, the 4 bytes will be stored at with the starting index being <i>o</o><br>
     */
    writeI2B : function writeI2B(b, i, o) {
        b[ o+0 ] = (i>>24)&0xFF;
        b[ o+1 ] = (i>>16)&0xFF;
        b[ o+2 ] = (i>> 8)&0xFF;
        b[ o+3 ] = (i>> 0)&0xFF;
    },

    /*
        Writes a String into a predefined byte buffer
        <br>
        <b><u>param:</b></u> b : The Byte Buffer<br>
        <b><u>param:</b></u> s : The String to be stored within the byte array<br>
        <b><u>param:</b></u> o : The offset representing the starting index for the utf-8 string<br>
        <b><u>param:</b></u> l : The length of the string to be put into the buffer<br>
     */
    writeS2B : function writeS2B (b, s, o, l) {
        var sl = s.length
        l = sl<l?sl:l;
        var bl = b.length;
        l = o + l>bl?bl-o:l;
        var e = stream_api.enc.encode(s);
        e.forEach((n, i ) => {
            b[o + i] = n;
        });
        e = null; //help GC
    },

    /*
        Reads a Int32 from the bytestream
        <br>
        <b><u>param:</b></u> b : The Byte Buffer to read from<br>
        <b><u>param:</b></u> o : The offset representing the starting index for 4 byte representation for an integer<br>
     */
    parseB2I : function parseB2I(b, o) {
        return b[o+0]<<24 | b[o+1]<<16 | b[o+2]<<8 | b[o+3];
    },

    /*
        Reads a Float32 from the bytestream
        <br>
        <b><u>param:</b></u> b : The Byte Buffer to read from<br>
        <b><u>param:</b></u> o : The offset representing the starting index for 4 byte representation for an float<br>
     */
    parseB2F32 : function parseB2F32(b, o) {
        let e = new Uint8Array(4);
        e[0] = b[o+0];
        e[1] = b[o+1];
        e[2] = b[o+2];
        e[3] = b[o+3];
        return new Float32Array(e.buffer, 0)[0]
    },

    /*
        Writes a 4 byte Float32 into a predefined byte buffer
        <br>
        <b><u>param:</b></u> b : The Byte Buffer<br>
        <b><u>param:</b></u> o : The offset, the 4 bytes will be stored at with the starting index being <i>o</i><br>
        <b><u>param:</b></u> f : The Integer to be stored in the buffer<br>
     */
    writeF322B : function writeF322B (b, o, f){
        var e = new Float32Array(1);
        e[0] = f;
        e = new Uint8Array(e.buffer);
        b[o + 0] = e[0];
        b[o + 1] = e[1];
        b[o + 2] = e[2];
        b[o + 3] = e[3];
        e = null;
    },

    /*
        Reads a String from a byte buffer
        <br>
        <b><u>param:</b></u> b : The Byte Buffer<br>
        <b><u>param:</b></u> o : The offset representing the starting index for the utf-8 string<br>
        <b><u>param:</b></u> l : The length of the string to be read<br>
     */
    parseB2S : function parseB2S(b, o, l){
        var e = b.slice(o, o+l);
        return stream_api.dec.decode(e);
    },

    /*
        Reads in a single Boolean from a byte Buffer
        <br>
     */
    parseB2Bool : async function parseB2Bool(b, o){
        return b[o] == 1
    },

    writeBool2B : function writeBool2B(b, o, bool){
        b[o] = (bool)?1:0;
    },
}