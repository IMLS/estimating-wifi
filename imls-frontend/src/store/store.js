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
    { id: "AA0003-001" }
  ],
});
export const state = reactive({});
