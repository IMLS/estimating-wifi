import { reactive, computed, readonly } from "vue";

const BACKEND_BASEURL = import.meta.env.VITE_BACKEND_BASEURL;

export const store = readonly({
  fscs_ids: [
    { id: "GA0029-002" },
    { id: "GA0029-003" },
    { id: "GA0029-004" },
    { id: "GA0027-004" },
    { id: "GA0027-008" },
    { id: "GA0027-011" },
  ],
  hourlyLabels: [
    "12am",
    "1am",
    "2am",
    "3am",
    "4am",
    "5am",
    "6am",
    "7am",
    "8am",
    "9am",
    "10am",
    "11am",
    "12pm",
    "1pm",
    "2pm",
    "3pm",
    "4pm",
    "5pm",
    "6pm",
    "7pm",
    "8pm",
    "9pm",
    "10pm",
    "11pm",
  ],
  dayOfWeekLabels: ["Sun", "Mon", "Tue", "Wed", "Thu", "Fri", "Sat"],
  backendBaseUrl: BACKEND_BASEURL,
  backendPaths: {
    get24HoursBinnedByHour: "/rpc/bin_devices_per_hour",
    get24HoursBinnedByHourForNDays: "/rpc/bin_devices_over_time",
    getLibraryDetailsById: "/rpc/lib_search_fscs",
    getAllSystemsByStateInitials: "/rpc/lib_search_state"
  },
  states:{"AL":"Alabama","AK":"Alaska","AZ":"Arizona","AR":"Arkansas","CA":"California","CO":"Colorado","CT":"Connecticut","DE":"Delaware","FL":"Florida","GA":"Georgia","HI":"Hawaii","ID":"Idaho","IL":"Illinois","IN":"Indiana","IA":"Iowa","KS":"Kansas","KY":"Kentucky","LA":"Louisiana","ME":"Maine","MD":"Maryland","MA":"Massachusetts","MI":"Michigan","MN":"Minnesota","MS":"Mississippi","MO":"Missouri","MT":"Montana","NE":"Nebraska","NV":"Nevada","NH":"New Hampshire","NJ":"New Jersey","NM":"New Mexico","NY":"New York","NC":"North Carolina","ND":"North Dakota","OH":"Ohio","OK":"Oklahoma","OR":"Oregon","PA":"Pennsylvania","RI":"Rhode Island","SC":"South Carolina","SD":"South Dakota","TN":"Tennessee","TX":"Texas","UT":"Utah","VT":"Vermont","VA":"Virginia","WA":"Washington","WV":"West Virginia","WI":"Wisconsin","WY":"Wyoming"}
});
