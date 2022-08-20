import { mount } from "@vue/test-utils";
import USWDSHeader from "./USWDSHeader.vue";
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

describe("USWDSHeader", () => {
  it("should render a header", () => {
    const wrapper = mount(USWDSHeader, {
      global: {
        plugins: [router],
      },
    });
    expect(wrapper.findAll(".usa-banner")).toHaveLength(1);
    expect(wrapper.findAll(".usa-header")).toHaveLength(1);
  });

  it("should submit a search query", async () => {
    const spySubmit = vi.spyOn(USWDSHeader.methods, "submitForm");
    const wrapper = mount(USWDSHeader, {
      global: {
        plugins: [router],
      },
    });
    const input = wrapper.find('.usa-input[type="search"]');
    await input.setValue("text to search for");
    await wrapper.find('button.usa-button[type="submit"]').trigger("submit");
    await wrapper.vm.$nextTick();
    expect(spySubmit).toHaveBeenCalledTimes(1);
  });
});
