// Grab the element that contains the API key.
const keyField = document.getElementById('api-key');
// Register the function that encodes everything.
keyField.addEventListener('change', update);

// chunkIntoN :: string integer -> list-of strings
// PURPOSE
// Takes a string and breaks it into a list of strings.
// Each list is of length N. The last string will be shorter.
function chunkIntoN(s, N) {
    chunks = [];
    // console.log("s", s, "length", s.length);
    for (var ndx = 0 ; ndx < s.length; ndx += N) {
        theSlice = s.slice(ndx, ndx + N);
        chunks.push(theSlice);
        // console.log("ndx", ndx, "triple", theSlice)
    }
    return chunks;
}

// chunkIntoThrees :: string -> list-of string
// PURPOSE
// A trivial helper for chunkIntoN.
function chunkIntoThrees (s) {
    return chunkIntoN(s, 3);
}

// CONSTANTS
// For the ASCII manipulations below.
const A = "A".charCodeAt(0);
const Z = "Z".charCodeAt(0);
const a = "a".charCodeAt(0);
const z = "z".charCodeAt(0);
const zero = "0".charCodeAt(0);
const nine = "9".charCodeAt(0);

// stringToDec :: string -> number
// PURPOSE
// Does a funny encoding of a string into a number.
// Takes 0-9 and maps them to the values 0-9.
// Takes A-Z and maps them to the values 10 - 36.
// Takes a-z and maps them to 37-...
// This gives us a range that is less than 64, and therefore
// we can represent each character with 6 bits.
function stringToDec (s) {
    var result = 0;
    // console.log(s)
    for (var ndx = 0 ; ndx < 3 ; ndx++) {
        var ascii = 63
        if (s[ndx]) {
            ascii = s[ndx].charCodeAt(0);
            if (ascii >= zero && ascii <= nine) {
                ascii = ascii - zero;
            } else if (ascii >= A && ascii <= Z) {
                ascii = ascii - A + 10;
            } else if (ascii >= a && ascii <= z) {
                ascii = ascii - a + 10 + 26;
            } else {
                console.log("ERROR. Character not in range: ", s[ndx]);
            }
        } 
        // Keep only the six rightmost bits.
        // That's all we should have at this point.
        ascii = ascii & (Math.pow(2, 6) - 1);
        // console.log("result in", result.toString(2))
        // console.log("ascii", ascii, ascii.toString(2));
        // Shift the values
        shifted = (ascii << (6*(3 - ndx - 1)));
        // console.log("shifted", shifted.toString(2));
        // Or with the result
        result = result | shifted;
        // console.log("result", result.toString(2));
    }
    
    // console.log("chunk", s, "dec", result, "bin", result.toString(2));

    return result;
} 

// chunksToDec :: list-of string -> list-of integers
// PURPOSE
// 
function chunksToDec (cs) {
    indexes = [];
    for (var ndx = 0; ndx < cs.length; ndx++) {
        indexes.push(stringToDec(cs[ndx]));
    }
    return indexes;
}

function updateHelper (key) {
    const table = document.getElementById('tablebody');
    chunks = chunkIntoThrees(key);
    indexes = chunksToDec(chunks);
    // console.log("wordlist length: ", wordlist.length);

    results = [];
    for (var ndx = 0 ; ndx < indexes.length ; ndx++) {
        const lookupNdx = indexes[ndx];
        const encoded = wordlist[lookupNdx];
        const decoded = chunks[ndx];
        // console.log("lookup", lookupNdx, "enc", encoded, "dec", decoded);
        results.push([encoded, decoded]);
    }

    // Reverse the list, because of the .push()
    // results = results.reverse();

    // Clear the table of current values.
    while (table.firstChild) {
        table.removeChild(table.firstChild);
    }


    for (var ndx = 0 ; ndx < results.length ; ndx++) {
        let row = document.createElement("tr");
        let word = document.createElement("td");
        let triple = document.createElement("td");

        const encoded = results[ndx][0];
        const decoded = results[ndx][1];
        
        word.innerHTML = (ndx + 1) + ". <b>" + encoded + "</b>";
        triple.textContent = decoded;
        console.log("enc", encoded, "dec", decoded);

        row.appendChild(word);
        row.appendChild(triple);
        table.appendChild(row);
    }

}

function update (e) {
    // Remove all of the table's children.
    const key = `${e.target.value}`
    updateHelper(key);

}

// Initialize the table
window.addEventListener('DOMContentLoaded', (event) => {
console.log('DOM fully loaded and parsed');
    updateHelper("library");
});



///// TESTING
// Grabbed from https://stackoverflow.com/questions/7837456/how-to-compare-arrays-in-javascript
// Warn if overriding existing method
if(Array.prototype.equals)
    console.warn("Overriding existing Array.prototype.equals. Possible causes: New API defines the method, there's a framework conflict or you've got double inclusions in your code.");
// attach the .equals method to Array's prototype to call it on any array
Array.prototype.equals = function (array) {
    // if the other array is a falsy value, return
    if (!array)
        return false;

    // compare lengths - can save a lot of time 
    if (this.length != array.length)
        return false;

    for (var i = 0, l=this.length; i < l; i++) {
        // Check if we have nested arrays
        if (this[i] instanceof Array && array[i] instanceof Array) {
            // recurse into the nested arrays
            if (!this[i].equals(array[i]))
                return false;       
        }           
        else if (this[i] != array[i]) { 
            // Warning - two different object instances will never be equal: {x:20} != {x:20}
            return false;   
        }           
    }       
    return true;
}
// Hide method from for-in loops
Object.defineProperty(Array.prototype, "equals", {enumerable: false});


function tests () {
    var keys = ["2LVtzHrVMC4u0lRPDpWg", "svHDmjmFLCUxJQxlP3qy", "YylHLkeoR1HT3uctu4Jc"];
    var valid = [
        ["state term", "native harmony", "forward metallic", "water case", "measure return", "reason spiritual", "external call"],
        ["chamber follow", "double question", "enter exhibit", "distance attack",
        "surface regular", "that intimate", "backward attend"],
        ["shoulder joint", "bearing uniform", "country weather", "form nature",
        "power language", "instrument northern", "surface belief"],
    ];

    for (var ndx = 0 ; ndx < keys.length ; ndx++) {
        key = keys[ndx];
        computed = [];
        
        chunks = chunkIntoThrees(key);
        indexes = chunksToDec(chunks);
        for (var inner = 0 ; inner < indexes.length ; inner++) {
            computed.push(wordlist[indexes[inner]]);
        }

        if (computed.equals(valid[ndx])) {
            console.log("Test passed: ", key);
        } else {
            console.log("FAIL: ", key);
            console.log("Expected: ", valid[ndx]);
            console.log("Computed: ", computed);
        }
    }
}

tests();