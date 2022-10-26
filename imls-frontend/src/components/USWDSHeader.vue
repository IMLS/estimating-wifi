<script>
import { RouterLink } from "vue-router";
import { store } from "@/store/store.js";

export default {
  name: "USWDS Header",
  components: { RouterLink },
  data() {
    return {
      searchString: "",
      store,
    };
  },
  methods: {
    submitForm: function (e) {
      e.preventDefault();
      this.$router.push({
        path: "/search",
        query: { query: this.searchString },
      });
    },
  },
};
</script>
<template>
  <div>
    <a
      class="usa-skipnav"
      href="#main-content"
    >Skip to main content</a>
    <div class="usa-overlay" />
    <header class="usa-header usa-header--extended">
      <div class="usa-navbar">
        <div
          id="extended-logo"
          class="usa-logo"
        >
          <em class="usa-logo__text">
            <RouterLink
              to="/"
              title="Public Library Wifi Estimator"
            > 
            Public Library Wifi Access Estimator </RouterLink>
          </em>
        </div>
        <button class="usa-menu-btn">
          Menu
        </button>
      </div>
      <nav
        aria-label="Primary navigation"
        class="usa-nav"
      >
        <div class="usa-nav__inner">
          <button class="usa-nav__close">
            <img
              src="~uswds/img/usa-icons/close.svg"
              role="img"
              alt="Close"
            >
          </button>
          <ul class="usa-nav__primary usa-accordion">
            <li class="usa-nav__primary-item">
              <button
                class="usa-accordion__button usa-nav__link"
                aria-expanded="false"
                aria-controls="extended-nav-section-libraries"
              >
                <span>Example Libraries</span>
              </button>

              <ul
                id="extended-nav-section-libraries"
                class="usa-nav__submenu"
              >
                <li v-for="library in store.fscs_ids" :key="library.id" class="usa-nav__submenu-item">
                  <RouterLink :to="{ path: '/library/' + library.id + '/' , query: $route.query}">
                    Example library {{ library.id }}
                  </RouterLink>
                </li>
               
              </ul>
            </li>
            <li class="usa-nav__primary-item">
              <button
                class="usa-accordion__button usa-nav__link"
                aria-expanded="false"
                aria-controls="extended-nav-section-states"
              >
                <span>All States</span>
              </button>

              <ul
                id="extended-nav-section-states"
                class="usa-nav__submenu"
              >
                <li v-for="(stateName, stateAbbr) in store.states" :key="stateName" class="usa-nav__submenu-item">
                  <RouterLink :to="{ path: '/state/' + stateAbbr + '/' , query: $route.query}">
                    {{ stateName }}
                  </RouterLink>
                </li>
               
              </ul>
            </li>
            <li class="usa-nav__primary-item">
              <RouterLink to="/about">
                About
              </RouterLink>
            </li>
          </ul>
          <div class="usa-nav__secondary">
            <ul class="usa-nav__secondary-links">
              <li class="usa-nav__secondary-item">
                <a target="_blank" href="//github.com/IMLS/estimating-wifi">Github Repo</a>
              </li>
            </ul>
            <form
              class="usa-search usa-search--small"
              role="search"
              @submit.stop="submitForm"
            >
              <label
                class="usa-sr-only"
                for="extended-search-field-en-small"
              > Search </label>
              <input
                id="extended-search-field-en-small"
                v-model="searchString"
                class="usa-input"
                type="search"
                name="search"
              >
              <button
                class="usa-button"
                type="submit"
              >
                <img
                  src="~uswds/img/usa-icons-bg/search--white.svg"
                  class="usa-search__submit-icon"
                  alt="Search"
                >
              </button>
            </form>
          </div>
        </div>
      </nav>
    </header>
  </div>
</template>
