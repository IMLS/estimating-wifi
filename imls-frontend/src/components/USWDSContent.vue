<script>
export default {
  props: {
    multilineContent: {
      type: String,
      required: true,
    },
    noIntro: {
      type: Boolean,
      required: false,
    },
  },
  computed: {
    parsedContent() {
      return this.parseNewlines(this.multilineContent);
    },
  },
  methods: {
    parseNewlines(text) {
      // add period to the end of fragments
      if (".!?".indexOf(text.slice(-1)) < 0) {
        text = text + ".";
      }
      return text.split(/\r\n|\r|\n/).filter(function (item) {
        return item;
      });
    },
  },
};
</script>
<template>
  <div>
    <p
      v-for="(newline, index) in parsedContent"
      :key="newline"
      :class="index == 0 && !noIntro ? 'usa-intro' : ''"
    >
      {{ newline }}
    </p>
    <p v-if="multilineContent.length < 1">
      No description available.
    </p>
  </div>
</template>
