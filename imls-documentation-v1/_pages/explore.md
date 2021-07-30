---
title: Exploring the data
layout: wide
sidenav: false
---
<script src="https://cdn.jsdelivr.net/npm/chart.js@3.4.1/dist/chart.min.js"></script>

<iframe name="dummyframe" id="dummyframe" style="display: none;"></iframe>
<form id="das-form" style="margin-bottom: 2em;" target="dummyframe">
    <div class="grid-container">
        <div class="grid-row">
            <div class="grid-col-12">
                <h2>Explore the data</h2>
                <p>Your sensor transmits data every night around midnight.
                </p> 
                <p>If you want, you can enter your key, select another library by FCFS Sequence ID, and see what their data looks like while you wait!</p>
            </div>
        </div>
    </div>
    <div class="grid-container">
        <div class="grid-row grid-gap">
            <div class="grid-col-6">
                <label class="usa-label" for="device-tag-text">Your API Key</label>
                <input class="usa-input" id="api-key-text" name="api-key-text" type="text">
            </div>
            <div class="grid-col-2" style="margin-top:3.33em">
                <button type="submit" class="usa-button" >Add key</button>    
            </div>
            <div class="grid-col-4"> &nbsp;</div>
        </div>
        <div class="grid-row grid-gap">
            <div class="grid-col-3">
                <label class="usa-label" for="options">FCFS Seq Id</label>
                <select class="usa-select" name="fcfs_seq_id_dropdown" id="fcfs_seq_id_dropdown">
                    <option value>Add key...</option>
                </select>
            </div>
            <div class="grid-col-3">
                <label class="usa-label" for="options">Session</label>
                <select class="usa-select" name="sessions_dropdown" id="sessions_dropdown">
                    <option value>Add key...</option>
                </select>
            </div>         
            <div class="grid-col-2">
                    <button  class="usa-button" 
                             style="margin-top:3.33em" 
                             id="update_button"
                             >Update</button>    
            </div>      
            <div class="grid-col-1"><span id="toggleme" style="display:none">updating...</span></div>
        </div>
                <div class="grid-row grid-gap" style="margin-top:2em">
            <div class="grid-col-12" style="text-align: center;">
            <p id="chartsummary" style="text-align: center;"></p>
            </div>
        </div>
        <div class="grid-row grid-gap" style="margin-top:0em">
            <div class="grid-col-9 grid-offset-1">
                <canvas id="das-chart"></canvas>
            </div>
        </div>
    </div>
</form>


<script>
    var DateTime = luxon.DateTime;
    var Info = luxon.Info;
    const form = document.getElementById("das-form");
    var ERROR = 0;
    fcfs_to_sessions = {}
    const MILLIS_PER_MINUTE = 1000 * 60;
    const MIN_MINUTES = 5;
    const MAX_MINUTES = 600;

</script>

<!-- handlers -->
<script>


    // This is called every time an FCFS Seq Id is selected.
    // It rebuilds the session dropdown from the options stored
    // in the map. 
    function fcfsSelectHandler() {
        try {
            var e = document.getElementById("fcfs_seq_id_dropdown");
            var session = e.options[e.selectedIndex].text;        
            var select = document.getElementById("sessions_dropdown");
            select.options.length = 0;
            console.log("attempting to look up session ", session, " in map");
            arr = fcfs_to_sessions[session];
            console.log("found ", arr);
            if (typeof arr !== 'undefined') {
                for (id of arr) {
                    var opt = document.createElement('option');
                    opt.value = id;
                    opt.innerHTML = id;
                    select.appendChild(opt);
                }
                // Update the chart with the first entry
                // console.log("UPDATING CHART FROM fcfsSelectHandler");
                // drawChartGate();
            }
      } catch (error) {
            console.error("select handler", error);
        }
    }
    
    // This is called to build the hashmap.
    function sessionHandler(fcfs) {
        return (data) => {
            // console.log("session handler", data)
            var arr = data.data.items.durations_v2
            ids = distinct(arr, "session_id")
            fcfs_to_sessions[fcfs] = ids
        }
    }

    // This is called whenever the FCFS dropdown is updated.
    // On initial button-press, that means it will be called multiple
    // times; it does a series of queries to build up both this dropdown
    // as well as the unique sessions for each FCFS Seq Id.
    function fcfsSeqResultHandler(api_key) {
    return async (data) => {
        var arr = data.data.items.durations_v2
        fcfs_ids = distinct(arr, "fcfs_seq_id")
        var select = document.getElementById("fcfs_seq_id_dropdown");
        select.options.length = 0;
        for (id of fcfs_ids) {
            var opt = document.createElement('option');
            opt.value = id;                
            opt.innerHTML = id;
            select.appendChild(opt);
            // Update session hash
            options = buildUniqSessionQuery(id);
            // console.log("options", options);
            await fetch(gqlUrl(api_key), gqlOptions(options))
                .then(res => res.json())
                .then(sessionHandler(id))
                .catch(eventFailHandler);       
            }
        }
    }

    var allow_charting = false;

    // This is called when the button is pressed.
    async function loadFCFSSeqIds() {
        var elem = document.getElementById("toggleme");
        elem.style.display = "block";
        const api_key = document.getElementById("api-key-text").value;
        if (api_key != "") {
            var fcfsQuery = buildFCFSQuery();
            fcfs_to_sessions = {}
            // Do the events query
            await fetch(gqlUrl(api_key), gqlOptions(fcfsQuery))
                .then(res => res.json())
                .then(fcfsSeqResultHandler(api_key))
                .then(fcfsSelectHandler)
                .catch(eventFailHandler);
            console.log("map", fcfs_to_sessions);
            allow_charting = true;
        }
        elem.style.display = "none";
    }

    function drawChartGate() {
        if (allow_charting) {
            console.log("UPDATING CHART");
            drawChart()
        }
    }
    
    var select = document.getElementById("fcfs_seq_id_dropdown");
    select.addEventListener(
        'change',
        function() { fcfsSelectHandler() },
        false
    );

    // var select = document.getElementById("sessions_dropdown")
    // select.addEventListener(
    //     'change',
    //     function() { drawChartGate() },
    //     false
    // );
    
</script>

<!-- chart -->
<script>

    var ctx = document.getElementById("das-chart").getContext("2d");
    
    //const labels = Utils.months({count: 7});
    var data = {
    labels: ["N", "E", "E", "D", "K", "E", "Y"],
    datasets: [{
        axis: 'y',
        label: 'Session Map',
        data: [[36, 65], 59, 80, 81, 56, 55, 40],
        fill: false,
        backgroundColor: [
        'rgba(255, 99, 132, 0.2)',
        'rgba(255, 159, 64, 0.2)',
        'rgba(255, 205, 86, 0.2)',
        'rgba(75, 192, 192, 0.2)',
        'rgba(54, 162, 235, 0.2)',
        'rgba(153, 102, 255, 0.2)',
        'rgba(201, 203, 207, 0.2)'
        ],
        borderColor: [
        'rgb(255, 99, 132)',
        'rgb(255, 159, 64)',
        'rgb(255, 205, 86)',
        'rgb(75, 192, 192)',
        'rgb(54, 162, 235)',
        'rgb(153, 102, 255)',
        'rgb(201, 203, 207)'
        ],
        borderWidth: 0
    }]
    };

    const config = {
        type: 'bar',
        data,
        options: {
            plugins: {
                legend: {
                    display: false
                },
            },
            indexAxis: 'y',
            scales: {
                x: {
                    title: {
                        text: "Time of Day (midnight to midnight)",
                        display: true
                    },
                    // The horizontal chart still uses this axis for the data.
                    stacked: true,
                    min: 0,
                    max: 24,
                },
                y: {
                    stacked: true,
                    title: {
                        text: "Patron Number",
                        display: true
                    },
                }},
    }};
    var chart = new Chart(ctx, config);

    var base_url = window.location.origin;
    
    function drawChartCallback(data) {
        //console.log("chart callback", data)
        var arr = data.data.items.durations_v2
        //console.log("chart callback", arr)
        total_minutes = 0;
        patron_devices = 0;

        sqldata = []
        sqllabels = []
        converted = []
        for (o of arr) {
            console.log(o)
            startdt = DateTime.fromISO(o.start)
            start = startdt.toSeconds();
            enddt = DateTime.fromISO(o.end)
            end = enddt.toSeconds();
            console.log(o.start, "becomes", start)
            console.log(o.end, "becomes", end)
            minutes = ((end - start) / 60)
            //console.log("start", start, "end", end);
            console.log("minutes: ", minutes)
            if ((minutes > MIN_MINUTES)
                &&
                (minutes < MAX_MINUTES)) {
            console.log("Found one")
            start_decimal = startdt.hour + (startdt.minute / 60);
            end_decimal = enddt.hour + (enddt.minute / 60);
            if (startdt.day != enddt.day) {
                end_decimal = 24;
            }
            pusharr = [o.patron_index, o.manufacturer_index, start_decimal, end_decimal];
            console.log("pushing", pusharr);
            converted.push(pusharr);
            total_minutes += minutes
            patron_devices += 1
            }
            
        }
        converted.sort((a, b) => {
            first = a[2];
            second = b[2];
            if (first < second) return -1;
            if  (first > second) return 1;
            return 0;
        });
        console.log(converted);

        for (o of converted) {
            pid = o[0]
            console.log(o[2], o[3]);
            stack = [o[2], o[3]]
            console.log(stack);
            sqldata.push(stack)
            sqllabels.push(pid)
        }
        chart.data.datasets[0].data = sqldata;
        chart.data.labels = sqllabels;
        console.log("data", sqldata)
        console.log("labels", sqllabels)
        chart.update();
        var elem = document.getElementById("chartsummary");
        elem.innerHTML = total_minutes + " minutes across " + patron_devices + " devices." 

    }
    async function drawChart() {
        //const uInt8Array = new Uint8Array(xhr.response);
        //const db = new SQL.Database(uInt8Array);
        //const contents = db.exec(`SELECT patron_index, manufacturer_index, start, end  FROM durations WHERE session_id="0"`);
         var elem = document.getElementById("toggleme");
        elem.style.display = "block";
        var e = document.getElementById("fcfs_seq_id_dropdown");
        var fcfs_seq_id = e.options[e.selectedIndex].text;     
        var e = document.getElementById("sessions_dropdown");
        var session_id = e.options[e.selectedIndex].text;     
        
        console.log("drawing " + fcfs_seq_id)
        const api_key = document.getElementById("api-key-text").value;
        var q = getChartDataQuery(fcfs_seq_id, session_id);
        console.log("query", q)
        // Do the events query
        await fetch(gqlUrl(api_key), gqlOptions(q))
            .then(res => res.json())
            .then(drawChartCallback)
            .catch(eventFailHandler);
        elem.style.display = "none";

    }
</script>

<!-- HELPERS -->
<script>

    function distinct(objarr, key) {
        d = {}
        for (o of objarr) {
            d[o[key]] = true
        }
        console.log(d)
        ls = []
        for (const [key, value] of Object.entries(d)) {
            ls.push(key)
        }
        return ls
    }

    function pad(min) {
        if (min < 10) {
            return `0${min}`;
        } else {
            return `${min}`;
        }
    }

    function eventFailHandler(e) {
        ERROR=1;
        console.log("eventHandler", e);
    }

    function wifiFailHandler(e) {
        ERROR=1;
        console.log("wifiHandler", e);
    }


    async function handleSubmit(event) {
        console.log(event)
        event.preventDefault();
        if (event.submitter.outerText.includes("Add")) {
        // RESET ERROR FLAG
        ERROR=0;
        console.log("button pressed")
        loadFCFSSeqIds()
        } else {
            console.log("update button pressed")
            drawChartGate()
        }
    } // end wifiQuery

</script>

<!-- QUERIES -->
<script>

    const SEARCH_LIMIT = 1000;
    const FCFS_SEQ_ID_SEARCH_LIMIT = 20000;

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

    function buildFCFSQuery() {
        return `
        {
            items {
                durations_v2(limit: ${FCFS_SEQ_ID_SEARCH_LIMIT}, sort: ["-id"]) {
                    fcfs_seq_id
                }
            }
        }
        `
    }

    function buildUniqSessionQuery(fcfs_seq_id) {
        return `
        {
            items {
                durations_v2(limit: ${FCFS_SEQ_ID_SEARCH_LIMIT},
                             filter: { fcfs_seq_id: {_eq: "${fcfs_seq_id}"}}) {
                    session_id
                }
            }
        }
        `
    }
    function getChartDataQuery(fcfs_seq_id, session_id) {
        return `
        {
            items {
                durations_v2(limit: ${FCFS_SEQ_ID_SEARCH_LIMIT},
                             filter: { fcfs_seq_id: {_eq: "${fcfs_seq_id}"},
                                       session_id: {_eq: "${session_id}"}, 
                             }) {
                    patron_index
                    manufacturer_index
                    start
                    end
                }
            }
        }
        `
    }
</script>


<script>
    // Do this last
    form.addEventListener("submit", handleSubmit);

</script>