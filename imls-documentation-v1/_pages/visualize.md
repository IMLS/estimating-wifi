---
title: Visualizing the data
layout: wide
sidenav: false
---

## The Sketch

<script>
    // Refreshed everytime we hit the button.
    // Need the `let` for scoping.
    let gql = null;
</script>

<iframe name="dummyframe" id="dummyframe" style="display: none;"></iframe>
<form id="das-form" style="margin-bottom: 2em;" target="dummyframe">
    <div class="grid-container">
        <div class="grid-row">
            <div class="grid-col-4" style="margin-right:10px;">
                <label class="usa-label" for="device-tag-text">Your API Key</label>
                <input class="usa-input" id="api-key-text" name="api-key-text" type="text">
            </div>
            <div class="grid-col-3" style="margin-right:10px;">
                <label class="usa-label" for="device-tag-text">FCFS Seq Id</label>
                <input class="usa-input" id="fcfs-text" name="fcfs-text" type="text">
            </div>
            <div class="grid-col-3" style="margin-right:10px;">
                <label class="usa-label" for="device-tag-text">Device tag</label>
                <input class="usa-input" id="device-tag-text" name="device-tag-text" type="text">
            </div>
            <div class="grid-col-1">
                <button type="submit" class="usa-button" style="margin-top: 3.1em;">Go!</button>    
            </div>
        </div>
        <div class="grid-row" style="margin-top: 2em;">
            <div class="usa-alert usa-alert--error" role="alert" id="errormsg" style="display:none">
                <div class="usa-alert__body">
                    <h4 class="usa-alert__heading">OH NOES!</h4>
                    <p class="usa-alert__text">Either you entered something incorrectly, or something is broken elsewhere.</p>
                    <p class="usa-alert__text">It is beyond the abilities of this simple webpage to tell which is true.</p>
                    <p class="usa-alert__text">Try again; if problems persist, reach out to the team for support.</p>
                </div>
            </div>
        </div>
    </div>
</form>

<main>
</main>


<script>

    // Takes the data array, sorts it on event_ids,
    // and stuffs it into gql.data
    function sortData (data) {
        var arr = data.data.items.wifi_v1
        console.log("reading count", arr.length);
        // sort this by id
        arr.sort((a, b) => { 
            if (parseInt(a.id) > parseInt(b.id)) { return -1; } 
            if (parseInt(a.id) < parseInt(b.id)) { return  1; }
            return 0; 
            });

        gql.events = new Set();
        for (let ndx = 0 ; ndx < arr.length ; ndx++) {
            gql.events.add(parseInt(arr[ndx].event_id));
        }
        console.log("number of minutes", gql.events.size);

        gql.data = arr;
    }

    function successStartup (data) {
        // Load these into the GQL object
        getSessions(data);
    }

    function successWifi (data) {
        sortData(data)
        gql.unique = {};        
        // Create a set of unique mfg ids.
        mfgs = new Set();
        gql.data.map(o => mfgs.add(parseInt(o.manufacturer_index)));
        gql.unique = {};
        gql.unique.manufacturers = mfgs;
        // And a unique set of device ids.
        devices = new Set();
        gql.data.map(o => devices.add(parseInt(o.patron_index)));
        gql.unique.devices = devices;
        // Map a mfg to all the devices in the set.
        console.log("unique", gql.unique);

        mfgToDevice = {}
        for (let m of gql.unique.manufacturers) {
            mfgToDevice[m] = new Set();
            for (let o of gql.data) {
                if (o.manufacturer_index == m) {
                    mfgToDevice[m].add(parseInt(o.patron_index));
                }
            }
        }
        gql.mfgToDevice = mfgToDevice;

        console.log("mfgToDevice", mfgToDevice);
    }
    function failure (err) {
        console.log("failure", err);
    }

    async function mashTheButton (event) {
        event.preventDefault();
        console.log("BUTTON! WOO!");
        const device_tag = document.getElementById("device-tag-text").value;
        const fcfs_seq_id = document.getElementById("fcfs-text").value;
        const api_key = document.getElementById("api-key-text").value;

        gql = new GQL(api_key, fcfs_seq_id, device_tag);
        gql.setSearchLimit(-1);
        await gql.runStartupQuery((d) => {console.log("startup events grabbed");}, failure);
        // If that was successful, we can now grab the wifi events
        // for each of those sessions.
        await gql.runWifiPerSession((d) => {console.log("wifi sessions grabbed");}, failure)
        // Now, we have wifi events per session.
        // FIXME: We are not resetting every day.

    }

    const form = document.getElementById("das-form");
    form.addEventListener("submit", mashTheButton);
</script>

{% asset "gqllib.js" %}
{% comment %}
{% asset "visualize.js" %}
{% endcomment %}
