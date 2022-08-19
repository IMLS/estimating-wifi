<script>
import datePicker from "uswds/src/js/components/date-picker";
import { state } from "@/store/store.js";

export default {
  name: "USWDS Date Picker",
  props: {
    label: {  
      type: String,
      default: "Date" 
    },
    initialDate: {
      type: String,
      default: () => new Date().toString()
    },
    
  },
  data() {
    return {
      state
    }
  },
  mounted() {
    try {
      datePicker.init();
      datePicker.enable(this.$refs.picker);
    } catch(e) {
      throw new Error(e)
    }
  },
  methods: {
    detectChange(e) {
      state.selectedDate = e.detail.value;
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
      data-min-date="2022-05-01"
      data-max-date="2022-06-01"
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