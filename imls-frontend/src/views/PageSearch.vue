<script>
export default {
  props: {
    query: {
      type: String,
      required: true
    }
  }
}
</script>

<template>
  <div class="search">
    <h1>This is a Search page</h1>
    <p>Your query contains "{{ query }}".</p>
  </div>
</template>

<style>
</style>
