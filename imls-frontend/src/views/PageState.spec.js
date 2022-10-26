import { mount, shallowMount, flushPromises } from "@vue/test-utils";
import { expect } from "vitest";
import PageState from "./PageState.vue";
import { createRouter, createWebHistory } from "vue-router";
import { routes } from "../router/index.js";

let router;

beforeEach(async () => {
  router = createRouter({
    history: createWebHistory(),
    routes: routes,
  });
  router.push("/");
  await router.isReady();
});

describe("PageState", () => {
  it("should render a list of libraries", async () => {
    const wrapper = mount(PageState, {
      props: {
        stateInitials: "AK",
      },
      global: {
        stubs: [
          "router-link",
          "router-view",
          "RouterView",
          "RouterLink",,
        ],
      },
    });
    await flushPromises();
    await wrapper.vm.$nextTick();

    expect(wrapper.find("h1").text()).toEqual("Alaska Public Libraries");
    expect(wrapper.findAll("ol.usa-list li").length).toBeGreaterThanOrEqual(1);
    expect(wrapper.vm.fetchedData.length).toBeGreaterThanOrEqual(1);
  });
  
});
