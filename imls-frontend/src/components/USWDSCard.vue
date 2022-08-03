<script>
import { RouterLink } from 'vue-router'

export default {
  props: {
    title: {  type: String,
      default: "No title content provided" 
    },
    ctaText: {
      type: String,
      default: "Find out more"
    },
    ctaLink: {
      type: String,
      default: "#"
    }
  },
  computed: {
    isButtonDisabled() {
      return this.ctaLink == '' ? true : undefined
    }
  }
}
</script>

<template>
  <div class="usa-card__container">
    <header class="usa-card__header">
      <h2 class="usa-card__heading">
        {{ title }}
      </h2>
    </header>
    <div class="usa-card__body">
      <p>
        <slot />
      </p>
    </div>
    <div class="usa-card__footer">
      <router-link
        class="usa-button"
        :to="ctaLink"
      >
        {{ ctaText }}
      </router-link>
    </div>
  </div>
</template>

<style>

</style>