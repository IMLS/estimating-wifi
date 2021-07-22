class GQL {
    constructor (api_key, fcfs_seq_id, device_tag) {
        GQL.DateTime = luxon.DateTime;
        GQL.Info = luxon.Info;
        this.form = document.getElementById("das-form");
        this.SEARCH_LIMIT = 1000;
        this.data = [];    
        this.api_key = api_key;
        this.device_tag = device_tag;
        this.fcfs_seq_id = fcfs_seq_id;
        this.WIFISTATUS = false;
        this.STARTUPSTATUS = false;
    }

    setSearchLimit (lim) {
        this.SEARCH_LIMIT = lim;
    }


    getUniqueSessionsQuery () {
        return `
        {
            items {
                wifi_v1(limit: ${this.SEARCH_LIMIT},
                            filter: { fcfs_seq_id: {_eq: "${this.fcfs_seq_id}"}, 
                                    device_tag: {_eq: "${this.device_tag}"},
                                    },
                            sort: ["id"]) {
                    id
                    session_id
                    fcfs_seq_id
                    device_tag
                }
            }
        }`;
    }

    getWifiSessionQuery(session_id) {
        return `
        {
            items {
                wifi_v1(limit: ${this.SEARCH_LIMIT}, 
                        filter: { fcfs_seq_id: {_eq:"${this.fcfs_seq_id}"}, 
                                  device_tag: {_eq: "${this.device_tag}"},
                                  session_id: {_eq: "${session_id}"},
                                },
                        sort: ["id"] 
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

    getGqlUrl () {
        return `https://api.data.gov/TEST/10x-imls/v1/graphql/?api_key=${this.api_key}`;
    }
    
    gqlOptions(query) {
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

    processStartupEvents (data) {
        console.log("startup", data)
        // Working backwards from wifi data.
        // The startup events in the events_v1 table are not reliable.
        var arr = data.data.items.wifi_v1;
        // Sorted by `id`, so the sessions should be "in order"
        var sessions = new Set();
        for (let o of arr) {
            sessions.add(o.session_id);
        }

        this.sessions = {};
        for (let s of sessions) {
            this.sessions[s] = [];
        }
        this.STARTUPSTATUS = true;
    }

    processWifiEvents (data) {
        var arr = data.data.items.wifi_v1;
        console.log("wifi", arr);
        for (let o of arr) {
            this.sessions[o.session_id].push(o);
        }
    }

    async runStartupQuery (success, failure) {
        this.STATUS = false;
        var req = async () => {
            await fetch(this.getGqlUrl(), this.gqlOptions(this.getUniqueSessionsQuery()))
            .then(res => res.json())
            .then(res => { this.processStartupEvents(res); success(res); })
            .catch(failure);
        }
        await req();
    }

    async runWifiPerSession (success, failure) {
        for (let k of Object.keys(this.sessions)) {
            console.log("session_id", k)
            var req = async () => {
                await fetch(this.getGqlUrl(), this.gqlOptions(this.getWifiSessionQuery(k)))
                .then(res => res.json())
                .then(res => { this.processWifiEvents(res); success(res); })
                .catch(failure);
            } 
            await req();
        }
    }
}
