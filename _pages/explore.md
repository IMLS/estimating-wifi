---
title: Exploring the data
layout: wide
sidenav: false
---

<iframe name="dummyframe" id="dummyframe" style="display: none;"></iframe>
<form id="das-form" style="margin-bottom: 2em;" target="dummyframe">
    <div class="grid-container">
        <div class="grid-row">
            <h2>Is it working?</h2>
        </div>
    </div>
    <div class="grid-container">
        <div class="grid-row">
            <div class="grid-col-12">
                <label class="usa-label" for="device-tag-text">Your API Key</label>
                <input class="usa-input" id="api-key-text" name="api-key-text" type="text">
                <label class="usa-label" for="device-tag-text">FCFS Seq Id</label>
                <input class="usa-input" id="fcfs-text" name="fcfs-text" type="text">
                <label class="usa-label" for="device-tag-text">Device tag</label>
                <input class="usa-input" id="device-tag-text" name="device-tag-text" type="text">
            </div>
        </div>
        <div class="grid-row" style="margin-top: 2em;">
            <div class="grid-col-3">
                <button type="submit" class="usa-button">Check what's up</button>    
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

<div class="grid-container" id="toggleme" style="display:none">
    <div class="grid-row">
        <div class="grid-col-9">
            <p>The device <span id="device_tag"></span> last started up on <span id="last_reboot_date"></span> at <span id="last_reboot_time"></span>.</p>
            <p>The last wifi observation was <span id="last_wifi_obs"></span>.</p>
            <p>How many devices have been seen recently?</p>
        </div>
    </div>
    <div class="grid-row">
        <div class="grid-col-9">
            <button type="download" id="downloadbutton" class="usa-button">Grab this data as a CSV file</button>
        </div>
    </div>
</div>

<div class="grid-container" >
    <div class="grid-row">
        <div class="grid-col-9">
            <div class="ct-chart ct-chart-1" style="padding-bottom: 2em;" ></div>
        </div>
    </div>
</div>

<script>
    var DateTime = luxon.DateTime;
    var Info = luxon.Info;

    const form = document.getElementById("das-form");

    const SEARCH_LIMIT = 1000;

    var TheWifiData = [];

    function gqlUrl (key) {
        return `https://api.data.gov/TEST/10x-imls/v1/graphql/?api_key=${key}`;
    }

    function gqlOptions(query) {
        const options = {
            method: "POST",
            headers: {
                "Content-Type": "application/json",
            },
            body: JSON.stringify({
                query: query
            })
        };
        return options;
    }

    function setResultText(arr) {
        // Data is sorted by the query, so first event is newest
        if (arr != null) {
            console.log("startup events", arr);
            var lastSeen = arr[0];
            var tagElem = document.getElementById("device_tag");
            var dateElem = document.getElementById("last_reboot_date");
            var timeElem = document.getElementById("last_reboot_time");
                    
            var localtime = lastSeen["localtime"];
            var dt = DateTime.fromISO(localtime);

            tagElem.innerHTML = "<b>" + lastSeen["device_tag"] + "</b>";
            dateElem.innerHTML = "<b>" + dt.weekdayLong + ", " + Info.months()[dt.month - 1] +  " " + dt.day + "</b>";
            timeElem.innerHTML = "<b>" + dt.hour + ":" + pad(dt.minute) + "</b>";
        }
    }

    function setLastWifiObsText(arr) {
        if (arr != null) {
            var lastSeen = arr[0];
            var tagElem = document.getElementById("last_wifi_obs");
            var localtime = lastSeen["localtime"];
            var dt = DateTime.fromISO(localtime);
            tagElem.innerHTML = ("<b>" + dt.weekdayLong + ", " + Info.months()[dt.month - 1] +  " " + dt.day + 
                                 " at " + dt.hour + ":" + pad(dt.minute) + "</b>"
                                 );
        }
    }

    function drawResultChart(arr) {
        event_ids = arr.map(o => o.event_id);
        

        // Store this for CSV generation
        TheWifiData = arr;

        // The array is in reverse order. This means the most
        // recent events are first.

        current_eid = -1;
        count = 0;
        counts = [];

        // Walk the list of event IDs.
        // Count the number of objects with each event ID.
        // Keep the list of counts. Each event is essentially
        // one minute.
        for (var ndx = 0; ndx < event_ids.length; ndx++) {
            if (ndx == 0) {
                current_eid = event_ids[ndx];
                count = 1;
            } else if (event_ids[ndx] != current_eid) {
                counts.push(count);
                current_eid = event_ids[ndx];
                count = 1;
            } else {
                count += 1;
            }
        }


        // Create some cute labels.
        labels = []
        for (var ndx = 0; ndx < counts.length - 1; ndx++) {
            if (ndx == 0) {
                labels.push("just now");
            } else if ((ndx % 5) == 0) {
                labels.push(`-${counts.length - (ndx + 1)}`);
            } else {
                labels.push(" ");
            }
        }
        labels.push(`-${counts.length - 1} mins ago`);

        chartData = {
            // A labels array that can contain any sort of values
            labels: labels,
            // Our series array that contains series objects or in this case series data arrays
            series: [ counts ]
        };
        chartOptions = {
            fullWidth: true,
            height: "300px",
            chartPadding: {
                right: 40
            },
            axisX: {
                offset: 70 
            },
        };
        
        new Chartist.Bar('.ct-chart-1', chartData, chartOptions)
    }

    function eventsResult(data) {
        // What comes back, if successful, looks like:
        // {data : { items : { events_v1 : [ obj ... ]}}}
        // where objects are keyed with the fields requested in the GraphQL query.
        var arr = data.data.items.events_v1
        setResultText(arr);
    }

    function wifiResult(data) {
        console.log(data);
        var arr = data.data.items.wifi_v1
        console.log("first event", arr[0]);
        console.log("last event", arr[arr.length - 1]);

        setLastWifiObsText(arr);
        drawResultChart(arr);
    }

    function pad(min) {
        if (min < 10) {
            return `0${min}`;
        } else {
            return `${min}`;
        }
    }

    var ERROR = 0;
    function eventFailHandler(e) {
        ERROR=1;
        console.log("eventHandler", e);
    }

    function wifiFailHandler(e) {
        ERROR=1;
        console.log("wifiHandler", e);
    }

    function buildEventQuery(fcfs_seq_id, device_tag) {
        return `
        {
            items {
                events_v1(limit: ${SEARCH_LIMIT},
                            filter: { fcfs_seq_id: {_eq: "${fcfs_seq_id}"}, 
                                    device_tag: {_eq: "${device_tag}"}, 
                                    tag:{_eq:"startup"}},
                            sort: ["-id"]) {
                    id
                    servertime
                    localtime
                    session_id
                    device_tag
                    tag
                }
            }
        }`;
    }

    function buildWifiQuery(fcfs_seq_id, device_tag) {
        return `
        {
            items {
                wifi_v1(limit: ${SEARCH_LIMIT}, 
                        filter: { fcfs_seq_id: {_eq:"${fcfs_seq_id}"}, 
                                  device_tag: {_eq: "${device_tag}"}
                                },
                        sort: ["-id"] 
                        ) {
                    id
                    device_tag
                    fcfs_seq_id
                    session_id
                    event_id
                    manufacturer_index
                    patron_index
                    servertime
                    localtime
                }
            }
        }`
    }

    async function handleSubmit(event) {
        event.preventDefault();
        // RESET ERROR FLAG
        ERROR=0;
        var errelem = document.getElementById("errormsg");
        errelem.style.display = "none";

        const key = 1;
        const device_tag = document.getElementById("device-tag-text").value;
        const fcfs_seq_id = document.getElementById("fcfs-text").value;
        const api_key = document.getElementById("api-key-text").value;

        var eventQuery = buildEventQuery(fcfs_seq_id, device_tag);

        var wifiQuery = buildWifiQuery(fcfs_seq_id, device_tag);

        // Do the events query
        await fetch(gqlUrl(api_key), gqlOptions(eventQuery))
            .then(res => res.json())
            .then(eventsResult)
            .catch(eventFailHandler);

        // Now the wifi query
        await fetch(gqlUrl(api_key), gqlOptions(wifiQuery))
            .then(res => res.json())
            .then(wifiResult)
            .catch(wifiFailHandler);

        // If we navigated HTTPS without error...
        if (ERROR == 0) {
            // Toggle visibility now, so that the chart draws
            var elem = document.getElementById("toggleme");
            elem.style.display = "block";
            document.querySelector('.ct-chart-1').__chartist__.update()
        } else {
            var errelem = document.getElementById("errormsg");
            errelem.style.display = "block";
        }

    } // end wifiQuery

    // See
    // https://stackoverflow.com/questions/14964035/how-to-export-javascript-array-info-to-csv-on-client-side
    async function handleDownload(event) {
        event.preventDefault();
        let csvContent = "data:text/csv;charset=utf-8,";
        var seq_id = "";
        var tag = "";

        // Add a header row. Make sure it aligns with the fields below.
        headerRow = "fcfs_seq_id,device_tag,session_id,event_id,manufacturer_index,patron_index,localtime,servertime\r\n";
        csvContent += headerRow;
        
        // The wifi objs contain the fields from the GQL
        TheWifiData.forEach(function(obj) {
            rowArray = [ obj.fcfs_seq_id, obj.device_tag, obj.session_id, obj.event_id, obj.manufacturer_index, obj.patron_index, obj.localtime, obj.servertime ]
            let row = rowArray.join(",");
            csvContent += row + "\r\n";
            // This happens many times, but it should all be the same,
            // so multiple assignments won't matter...
            seq_id = obj.fcfs_seq_id;
            tag = obj.device_tag
        });

        console.log(csvContent);
        var encodedUri = encodeURI(csvContent);
        var link = document.createElement("a");
        link.setAttribute("href", encodedUri);
        link.setAttribute("download", seq_id + "_" + tag + ".csv");
        document.body.appendChild(link); // Required for FF
        link.click()
    }

    form.addEventListener("submit", handleSubmit);
    document.getElementById("downloadbutton").onclick = handleDownload;

</script>