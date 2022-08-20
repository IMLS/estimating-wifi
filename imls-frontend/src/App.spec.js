import { mount } from "@vue/test-utils";
import App from "./App.vue";
// import { createRouter, createWebHistory } from 'vue-router'
// let router;

// use a real vue router, mock the router, or stub <router-view> and <router-link> built-in vue components
// beforeEach(async () => {
//   router = createRouter({
//     history: createWebHistory(import.meta.env.BASE_URL),
//     routes: [],
//   })
//   router.push('/')
//   await router.isReady()
// });

describe("App", () => {
  it("should render a Header and Footer", () => {
    const wrapper = mount(App, {
      global: {
      //  plugins: [router],
        stubs: ["router-link", "router-view", "RouterView"], 
      }
    });

    expect(wrapper.findAll(".usa-header")).toHaveLength(1);
    expect(wrapper.findAll(".usa-footer")).toHaveLength(1);
  });
});
