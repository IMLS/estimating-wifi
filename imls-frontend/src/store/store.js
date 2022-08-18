import { reactive, computed, readonly } from "vue";

function sortByAttributes(attributeKey) {
  return (a, b) => {
    const first = a.attributes[attributeKey] || "zzz"; // sort undefined values to the end
    const second = b.attributes[attributeKey] || "zzz";
    return first.localeCompare(second);
  };
}

export const store = readonly({
  // replace with real data
  sensors: [{ id: "GA0027-008-01" }],
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
});
export const state = reactive({
  selectedDate: '2022-05-01'
});
