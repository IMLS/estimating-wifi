import { mount, shallowMount } from "@vue/test-utils";
import App from "./App.vue";
import { routes } from "./router/index.js";
import { createRouter, createWebHistory } from "vue-router";
import { createMetaManager, plugin as vueMetaPlugin } from 'vue-meta'
import { expect } from "vitest";

let router;
let vueMetaManager = createMetaManager(); 

beforeEach(async () => {
  router = createRouter({
    history: createWebHistory(import.meta.env.BASE_URL),
    routes: routes,
  });
  router.push("/");
  await router.isReady();
});

describe("App", () => {
  it("should render a Header and Footer", () => {
    const wrapper = mount(App, {
      global: {
        plugins: [router, vueMetaManager, vueMetaPlugin],
      },
    });
    expect(wrapper.findAll(".usa-header")).toHaveLength(1);
    expect(wrapper.findAll(".usa-footer")).toHaveLength(1);
  });

  it("should provide default site metadata", () => {
    const wrapper = shallowMount(App, {
      global: {
        plugins: [router, vueMetaManager, vueMetaPlugin],
      },
    });
    expect(App.metaInfo).toHaveProperty("title")
    expect(App.metaInfo).toHaveProperty("description")
  });
});
