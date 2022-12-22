<script>
import datePicker from "uswds/src/js/components/date-picker";
import { startOfYesterday } from "date-fns";

const MIN_DATE = "2022-01-01";
const MAX_DATE = startOfYesterday().toISOString().split("T")[0];

export default {
  name: "USWDSDatePicker",
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
  emits: ['date_changed'],
  data() {
    return {
      selectedDate: null,
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
      this.selectedDate = this.$refs.date.value;
      this.$emit('date_changed', encodeURIComponent(e.target.value) )
    }
  }
}
</script>

<template>
  <div class="usa-form-group">
    <label id="date-label" class="usa-label" for="date"
      >{{ label }}</label
    >
    <div id="date-hint" class="usa-hint">mm/dd/yyyy</div>
    <div
ref="picker" 
      class="usa-date-picker maxw-date-picker" 
      :data-default-value="initialDate"
      :data-selected-date="!!selectedDate ? selectedDate : initialDate"
      :data-min-date="minDate"
      :data-max-date="maxDate"
      >
      <input
        id="date"
        ref="date"
        class="usa-input"
        name="date"
        aria-labelledby="date-label"
        aria-describedby="date-hint"
        @change="detectChange"
      />
    </div>
  </div>
</template>

<style>

</style>