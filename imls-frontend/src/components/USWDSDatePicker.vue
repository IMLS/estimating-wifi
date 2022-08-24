<script>
import datePicker from "uswds/src/js/components/date-picker";
import { state } from "@/store/store.js";

// artificially limits selectable dates to the time we have 
// deterministic placeholder data for. 
// todo: refactor when we're using actual data
const MIN_DATE = "2022-05-01";
const MAX_DATE = "2022-06-01";

export default {
  name: "USWDS Date Picker",
  props: {
    label: {  
      type: String,
      default: "Date" 
    },
    initialDate: {
      type: String,
      default: () => new Date().toISOString().split("T")[0]
    },
    
  },
  data() {
    return {
      state,
      minDate: MIN_DATE,
      maxDate: MAX_DATE,
    }
  },
  mounted() {
    this.enableUSWDSFeatures();
  },
  methods: {
    async enableUSWDSFeatures(){
      await datePicker.init();
      await datePicker.enable(this.$refs.picker);
    },
    detectChange(e) {
      if (e && e.detail ) state.selectedDate = e.detail.value;
      return state.selectedDate;
    }
  }
}
</script>

<template>
  <div class="usa-form-group">
    <label class="usa-label" id="date-label" for="date"
      >{{ label }}</label
    >
    <div class="usa-hint" id="date-hint">mm/dd/yyyy</div>
    <div class="usa-date-picker maxw-card-lg" 
      ref="picker" 
      :data-default-value="initialDate"
      :data-selected-date="state.selectedDate"
      :data-min-date="minDate"
      :data-max-date="maxDate"
      >
      <input
        @change="detectChange"
        class="usa-input"
        id="date"
        name="date"
        aria-labelledby="date-label"
        aria-describedby="date-hint"
      />
    </div>
  </div>
</template>

<style>

</style>