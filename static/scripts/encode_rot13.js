function parseB2SRot13(b){
    var s = ""
    var i=0
    while(i<b.length){
        s += String.fromCharCode(b[i], 2).charAt(0);
        i+=1;
    }
    i = null; //GC
    return s
}

function parseS2BRot13(s){
    var b = new Uint8Array(s.length)
    var i = 0
    while(i<b.length){
        b[i] = s.charCodeAt(i)
        i+=1;
    }
    return b
}

/*
    Rot 13 encryption algorithm

    base 10 shift 13 in binary

    is decoded serverside again
 */
function encodeRot13(p){

    b = parseS2BRot13(p);
    var i = 0;
    while (i<b.length){
        var j = b[i]
        j = j+13;
        if (j>255){
            j = j-255;
        }
        b[i] = j&0xff;
        i+=1;

    }
    p = parseB2SRot13(b)
    return p
    //could do here a base 16shift or sth like that => security ? +
    //probably serverside randomly generated shifts, for each new login site
}
