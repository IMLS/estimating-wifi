import { mount, shallowMount, flushPromises } from "@vue/test-utils";
import { expect } from "vitest";
import PageLibrary from "./PageLibrary.vue";
import { createRouter, createWebHistory } from "vue-router";
import { routes } from "../router/index.js";
import USWDSDatePicker from "../components/USWDSDatePicker.vue";

let router;

beforeEach(async () => {
  router = createRouter({
    history: createWebHistory(),
    routes: routes,
  });
  router.push("/");
  await router.isReady();
});

describe("PageLibrary", () => {
  it("should render with yesterday's date if no date is provided", () => {
    const wrapper = mount(PageLibrary, {
      props: {
        id: "KnownGoodId",
      },
      global: {
        stubs: [
          "router-link",
          "router-view",
          "RouterView",
          "RouterLink",
          "USWDSDatePicker",
        ],
      },
    });
    expect(wrapper.find("h1").text()).toEqual("Library KnownGoodId");
    expect(wrapper.findAll(".usa-card").length).toBeGreaterThanOrEqual(1);
    expect(wrapper.vm.activeDate).toEqual(
      new Date(Date.now() - 86400 * 1000).toISOString().split("T")[0]
    );
  });
  it("should render with a preset date if one is provided", () => {
    const wrapper = mount(PageLibrary, {
      props: {
        id: "KnownGoodId",
        selectedDate: "2022-05-02",
      },
      global: {
        stubs: [
          "router-link",
          "router-view",
          "RouterView",
          "RouterLink",
          "USWDSDatePicker",
        ],
      },
    });
    expect(wrapper.vm.activeDate).toEqual("2022-05-02");
  });
  it("should respond by navigating to a new route query param when the selected date changes", async () => {
    const spyChangeDate = vi.spyOn(
      PageLibrary.methods,
      "navigateToSelectedDate"
    );
    const wrapper = shallowMount(PageLibrary, {
      props: {
        id: "KnownGoodId",
      },
      global: {
        plugins: [router],
        stubs: ["USWDSDatePicker"],
      },
    });
    expect(spyChangeDate).toHaveBeenCalledTimes(0);
    const childWrapper = wrapper.findComponent({ name: "USWDSDatePicker" });
    childWrapper.vm.$emit("date_changed", "2022-05-02");
    await wrapper.vm.$nextTick();
    await flushPromises();
    expect(spyChangeDate).toHaveBeenCalledTimes(1);
  });
});
