<script>
export default {
  name: 'USWDS Breadcrumb',
  props: { 
    crumbs: {
      type: Array,
      default: () => [
        { 
          name: "Home",
          link: "/" 
        },
        { 
          name: "This page",
          link: "#" 
        }
      ]
    }
  },
  computed: { 
    
  },
  methods: {
    // this component assumes the last item in the array is the current page, regardless of whether the URLs match. So determine the crumbs accordingly!
    isLast(index) {
      return index == (this.crumbs.length - 1)
    }
  }
};
</script>
<template>
  <nav v-if="crumbs.length > 0" class="usa-breadcrumb" aria-label="Breadcrumbs">
    <ol
      vocab="http://schema.org/"
      typeof="BreadcrumbList"
      class="usa-breadcrumb__list">
      <li      
        v-for="(crumb, index) in crumbs"
        :key="crumb.name"
        :class="isLast(index) ? 'usa-current' : ''"
        :aria-current="isLast(index) ? 'page' : null"
        property="itemListElement"
        typeof="ListItem"
        class="usa-breadcrumb__list-item">
        <span v-if="isLast(index)" property="name">
          {{ crumb.name }}
        </span>
        <a v-else
          property="item"
          typeof="WebPage"
          :href="crumb.link"
          class="usa-breadcrumb__link">
            <span property="name">
              {{ crumb.name }}
            </span>
          </a>
        <meta property="position" :content="(index + 1)" />
      </li>    
    </ol>
  </nav>
</template>
