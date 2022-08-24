<script>
export default {
  props: {
    title: {  type: String,
      default: "No title content provided" 
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
      <slot />
    </div>
  </div>
</template>

<style>

</style>