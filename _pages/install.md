---
title: Install and Config
layout: page
sidenav: false
---
{% assign counter = 0 %}

# {{page.title}}

Here we go! You have an RPi, it is set up, and it can talk to the internet! This is very, *very* exciting.
## Step {% assign counter = counter | plus: 1  %}{{counter}}: Open the terminal

In the upper-left corner of the screen is a little black box with a blue bar at the top. This icon is intended to represent a command terminal. It should open up a window that looks like this:

{% asset "pi-command-prompt.png" alt="location of pi command prompt icon" %}

## Step {% assign counter = counter | plus: 1  %}{{counter}}: Copy the line below

```
bash <(curl -s https://raw.githubusercontent.com/jadudm/imls-client-pi-playbook/main/bootstrap.sh)
```

Highlight that line, right-click, and say **Copy**. 

## Step {% assign counter = counter | plus: 1  %}{{counter}}: Paste the line into the terminal

Over in the terminal, right-click, and select **Paste**. The command should now be in the terminal window.

Press Enter. Or, Return. Or whatever that key is called on your keyboard.

## Step {% assign counter = counter | plus: 1  %}{{counter}}: Enter your FCFS_Seq_Id

You are approaching the final stretch.

You are setting up your device so that it will live in a particular library. Or, perhaps, a bookmobile. Either way, there should be an FSFS_Seq_Id for your location.

You can [look this up your FCFS Seq Id for your building on the IMLS website](https://www.imls.gov/search-compare/). For example, the Berea branch in Madison County, KY is **KY0069-003**.

Type this in when prompted. You will need to use caps for the state abbreviation (e.g. you cannot say `ky`, but you can say `KY`).

## Step {% assign counter = counter | plus: 1  %}{{counter}}: Enter a tag

This is a free-form text entry that will accept up to 255 characters. It *will* become part of the dataset, so please do not put any personally identifiable information into this field. Email addresses are considered PII.

We intend this field either as a "hardware tag" (so you might enter `RPi 001`), or it could be a reminder where you put it. In the latter case, you might enter `reference desk`, or `networking closet`, or (if it is in the bookmobile)... er, perhaps `glovebox`? Point being, this is mostly for you. This will become part of the dataset, so something simple but descriptive is best.

## Step {% assign counter = counter | plus: 1  %}{{counter}}: Enter your api.data.gov key

This is a tedious part.

First, note: you should never, ever give anyone your API key. You should not paste it into websites. You should not mail it to your friends. (They keys are free... they can go get their own.) But, we're going to ask you to violate that rule right now.

On your other computer (not the Pi), you'll want to find the email where you received your API key from api.data.gov. Then, you'll want to visit this page. (Yes, this one, right here.) Then, you'll want to paste your key into the text box below:

<script type="text/javascript" src="{{ '/js/wordlist.js' | prepend: site.baseurl }}"></script>

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

Your key will not be sent anywhere. The transformation that takes place happens *entirely within this webpage*. That is the reason it is safe for you to paste your API key here (and nowhere else!). 

Your API key has been transformed into a list of "word pairs."

You'll read each word-pair (perhaps off your phone, or a laptop off to the side) and type them into the Raspberry Pi one pair at a time. The Pi will make sure that you're typing "correct" word pairs. (There's a limited set; if you make a typo, we'll ask you to try again.) You should end up entering 14 word pairs.

As you type your word pairs, the Raspberry Pi will "translate" those words back into pieces of your API key. There's some fun mathematics happening behind the scenes here, but we don't want to tell you too much... you might get so excited, you forget what you're doing!

No. Probably not.

So, type the fourteen word pairs. 

When you are done, you can type `DONE`. Or, `done`. 

This will now save the API key on the Raspberry Pi.

## Step {% assign counter = counter | plus: 1  %}{{counter}}: Installation continues

At this point, the Pi is going to do a whole bunch of things on its own. This process is *slow*. But, you don't have to do anything more; the rest is fully automated.

What is happening? This is where we take the friendly Raspberry Pi that you were sent, and we turn it into a secure device. Amongst other things, we are:

1. Installing software to do the WIFISESS data collection. (This is open source software that was developed in collaboration with IMLS specifically for this pilot.) 
2. Making sure the RPi gets security and software updates. 
3. Disabling the ability for anyone to connect to the RPi via the network.
4. Disabling the ability for anyone to log into the RPi with a keyboard and mouse.
5. Removing the graphical interface.

And a few more things. When we're done, you really have no idea if it worked, save for the fact that it will reboot, and then quietly do its job. 

How do you know if it is working? We'll have to check and see if it is logging data! (FIXME. To be written.)




<script>
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
</script>

<!-- Tests -->
<script>

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
</script>