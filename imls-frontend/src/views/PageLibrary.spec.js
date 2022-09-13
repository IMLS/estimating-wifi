import { mount, shallowMount, flushPromises } from "@vue/test-utils";
import { expect } from "vitest";
import PageLibrary from "./PageLibrary.vue";
import { createRouter, createWebHistory } from "vue-router";
import { routes } from "../router/index.js";
import { startOfYesterday } from "date-fns";

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
          "USWDSCard",
          "FetchData",
          "Histogram",
          "Heatmap",
          "HeatmapWeeklyCalendar",
          "USWDSTable",
        ],
      },
    });
    expect(wrapper.find("h1").text()).toEqual("Library KnownGoodId");
    expect(wrapper.findAll(".usa-card").length).toBeGreaterThanOrEqual(1);
    expect(wrapper.vm.activeDate).toEqual(
      startOfYesterday().toISOString().split("T")[0]
    );
  });
  it("should render with a preset date if one is provided", () => {
    const wrapper = shallowMount(PageLibrary, {
      props: {
        id: "KnownGoodId",
        selectedDate: "2022-05-02",
      },
      global: {
        stubs: ["router-link", "router-view", "RouterView", "RouterLink"],
      },
    });
    expect(wrapper.vm.activeDate).toEqual("2022-05-02");
  });

  it("should format day labels for n days given a date and count", () => {
    expect(
      PageLibrary.methods.generateDayLabels("1999-12-31", 3)
    ).toStrictEqual(["12/31/99", "1/1/00", "1/2/00"]);
  });
  it("should return the first day of the week in ISO", () => {
    expect(
      PageLibrary.computed.startOfWeekInISO.call({ selectedDate: "1999-12-31" })
    ).toBe("1999-12-26");
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
        stubs: [
          "USWDSDatePicker",
          "USWDSCard",
          "FetchData",
          "Histogram",
          "Heatmap",
          "HeatmapWeeklyCalendar",
          "USWDSTable",
        ],
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
