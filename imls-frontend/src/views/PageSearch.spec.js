import { mount, shallowMount, flushPromises } from "@vue/test-utils";
import PageSearch from "./PageSearch.vue";

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


describe("PageSearch", () => {
  it("should render", () => {
    const wrapper = mount(PageSearch, {
      props: {
        query: "search string",
      },
      global: {
        stubs: ["router-link", "router-view", "RouterView", "RouterLink"],
      },
    });
    expect(wrapper.find("h1").text()).toEqual("Libraries matching \"search string\"");
    expect(wrapper.text()).toContain("search string");
  });
});

describe("PageSearch", () => {
  it("should render a list of libraries matching the given search string", async () => {
    const wrapper = mount(PageSearch, {
      props: {
        query: "anchor point",
      },
      global: {
        stubs: [
          "router-link",
          "router-view",
          "RouterView",
          "RouterLink",
        ],
      },
    });
    await flushPromises();
    await wrapper.vm.$nextTick();
    expect(wrapper.find("h1").text()).toEqual("Libraries matching \"anchor point\"");
    expect(wrapper.findAll("ol.usa-list li").length).toBeGreaterThanOrEqual(1);
    expect(wrapper.vm.fetchedLibraries.length).toBeGreaterThanOrEqual(1);
  });

  it("should update with new results when the query changes", async () => {
    const wrapper = await mount(PageSearch, {
      props: {
        query: "foo",
      },
      global: {
        stubs: [
          "router-link",
          "router-view",
          "RouterView",
          "RouterLink",
        ],
      },
    });
    await flushPromises();
    await wrapper.vm.$nextTick();
    expect(wrapper.find("h1").text()).toEqual("Libraries matching \"foo\"");
    expect(wrapper.findAll("ol.usa-list li").length).toEqual(0);

    await wrapper.setProps({ query: "anchor point" });
    await flushPromises();
    await wrapper.vm.$nextTick();

    expect(wrapper.find("h1").text()).toEqual("Libraries matching \"anchor point\"");
    expect(wrapper.findAll("ol.usa-list li").length).toBeGreaterThanOrEqual(1);

    wrapper.unmount();
  });

});
