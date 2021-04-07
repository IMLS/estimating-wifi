---
title: The API Token
layout: page
sidenav: false
---

<script type="text/javascript" src="{{ '/js/wordlist.js' | prepend: site.baseurl }}"></script>

You should have received an email from api.data.gov containing your authorization token. In general, you should **never share that token with anyone**. However, getting it from your email onto the Raspberry Pi could be tricky; therefore, we have provided a tool, here, to help you.

**This page does not transmit your token anywhere**. It will, entirely in your web browser, transform it into a series of common words that you can type into the Pi, and we will transform those words back into your token. We accomplish this through the magic of mathematics.

**Paste your key in and hit enter.**

<div class="grid-container">
  <div class="grid-row">
    <div class="tablet:grid-col">
        <label class="usa-label" for="api-key">Your API Key</label>
        <input class="usa-input" id="api-key" name="api-key" type="text" placeholder="library">
    </div>
    <div class="grid-col-fill" style="padding-left: 2em;">
        <table class="usa-table usa-table--borderless usa-table--striped">
            <thead>
                <tr>
                    <th scope="col">Words</th>
                    <th scope="col">Decoded</th>
                </tr>
            </thead>
            <tbody id="tablebody">
            </tbody>
        </table>
    </div>
</div>
</div>

To set up your Pi, you will need to enter each pair of words (at right) when prompted. 

For example, if your key was <b>library</b>, it would decode into three pairs of words. On the Pi, you would type each pair of words when prompted. First, you would type "government depression" and press return. Then, you would type "faith choice," press return, and so on.

This involves more typing than just entering your key, but we thought this would be easier than typing 40 random characters.

<script>
    const keyField = document.getElementById('api-key');
    keyField.addEventListener('change', update);
    
    // chunkIntoN :: string integer -> list-of strings
    function chunkIntoN(s, N) {
        chunks = [];
        console.log("s", s, "length", s.length);
        for (var ndx = 0 ; ndx < s.length; ndx += N) {
            theSlice = s.slice(ndx, ndx + N);
            chunks.push(theSlice);
            // console.log("ndx", ndx, "triple", theSlice)
        }
        return chunks;
    }

    function chunkIntoThrees (s) {
        return chunkIntoN(s, 3);
    }

    const A = "A".charCodeAt(0);
    const Z = "Z".charCodeAt(0);
    const a = "a".charCodeAt(0);
    const z = "z".charCodeAt(0);
    const zero = "0".charCodeAt(0);
    const nine = "9".charCodeAt(0);

    function stringToDec (s) {
        var result = 0;

        for (var ndx = 0 ; ndx < s.length ; ndx++) {
            var ascii = s[ndx].charCodeAt(0);
            if (ascii >= zero && ascii <= nine) {
                ascii = ascii - zero;
            } else if (ascii >= A && ascii <= Z) {
                ascii = ascii - A + 10;
            } else if (ascii >= a && ascii <= z) {
                ascii = ascii - zero + 10 + 26;
            } else {
                console.log("ERROR. Character not in range: ", s[ndx]);
            }
            // Keep only the six rightmost bits.
            // That's all we should have at this point.
            ascii = ascii & 255;
            // Shift the values
            shifted = (ascii << (6*(s.length - ndx - 1)));
            // Or with the result
            result = result | shifted;
        }
        console.log("chunk", s, "dec", result, "bin", result.toString(2));

        return result;
    } 

    // chunksToDec :: list-of string -> list-of integers
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
        console.log("wordlist length: ", wordlist.length);

        results = [];
        for (var ndx = 0 ; ndx < indexes.length ; ndx++) {
            const lookupNdx = indexes[ndx];
            const encoded = wordlist[lookupNdx];
            const decoded = chunks[ndx];
            console.log("lookup", lookupNdx, "enc", encoded, "dec", decoded);
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

            word.textContent = encoded;
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
</script>