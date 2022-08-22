import { reactive, computed, readonly } from "vue";

// todo: update when the backend has a real host
const BACKEND_BASEURL = `${window.location.protocol}//127.0.0.1:3000`;

export const store = readonly({
  fscs_ids: [
    { id: "AA0001-001" },
    { id: "AA0002-001" },
    { id: "AA0003-001" },
    { id: "AA0004-001" },
    { id: "AA0005-001" },
    { id: "AA0006-001" },
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
  backendPaths: {
    get24HoursBinnedByHour: "/rpc/bin_devices_per_hour",
  },
});

export const state = reactive({
  // todo: Update selectedDate when we have real data
  selectedDate: "2022-05-01",
  fetchCount: null,
  fetchError: {},
  fetchedData: {},
  isLoading: false,

  async fetchData(path, queryString) {
    this.isLoading = true;
    try {
      // todo: let queryString be an array of params instead of a string
      const response = await fetch(`${BACKEND_BASEURL}${path}${queryString}`);
      if (await !response.ok) {
        throw new Error(response.status);
      }
      this.fetchedData = await response.json();
    } catch (error) {
      this.fetchError = error;
    }
    this.isLoading = false;
  },
});
