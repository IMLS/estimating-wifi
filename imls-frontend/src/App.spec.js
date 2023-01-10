import { mount, shallowMount, flushPromises } from "@vue/test-utils";
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

  it("should skip focus to the main element on route change", async () => {
    const spyChangeFocus = vi.spyOn(
      App.methods,
      "setRouteWrapperFocus"
    );

    const wrapper = mount(App, {
      global: {
        plugins: [router, vueMetaManager, vueMetaPlugin],
      },
      attachTo: document.body
    });
    // hasn't skipped focus to main yet
    expect(wrapper.vm.$refs.focus).not.toBe(document.activeElement);
    expect(spyChangeFocus).toHaveBeenCalledTimes(0);
    router.push("/");
    // we're going to force this method to run because the watcher on $route change isn't running in the test env
    wrapper.vm.setRouteWrapperFocus();
    await wrapper.vm.$nextTick();
    expect(wrapper.vm.$refs.focus).toBe(document.activeElement);
    expect(spyChangeFocus).toHaveBeenCalledTimes(1);
  });
});
