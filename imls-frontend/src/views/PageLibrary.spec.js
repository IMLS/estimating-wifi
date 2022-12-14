import { mount, shallowMount, flushPromises } from "@vue/test-utils";
import { expect } from "vitest";
import PageLibrary from "./PageLibrary.vue";
import { createRouter, createWebHistory } from "vue-router";
import { routes } from "../router/index.js";
import { startOfMonth } from "date-fns";

let router;


const MOCK_ERROR_MSG = "mocked error message";
// the API currently returns null instead of an empty array on no matches
const MOCK_NO_LIBS_FOUND = null;
const MOCK_ONE_LIB_FOUND = {
  "stabr":"MK",
  "fscskey":"MOCK001",
  "fscs_seq":1,
  "libname":"MOCKED PUBLIC LIBRARY",
  "address":"1234 MOCKINGBIRD ROAD",
  "city":"MOUNT MOCKINGTON",
  "zip":"00000"
};
const MOCK_ANOTHER_LIB_FOUND = {
  "stabr":"MK",
  "fscskey":"MOCK001",
  "fscs_seq":2,
  "libname":"ANOTHER MOCKED PUBLIC LIBRARY",
  "address":"5678 MOCKINGBIRD ROAD",
  "city":"MOUNT MOCKINGTON",
  "zip":"00000"
}



beforeEach(async () => {
  router = createRouter({
    history: createWebHistory(),
    routes: routes,
  });
  router.push("/");
  await router.isReady();
  await flushPromises();
  fetch.resetMocks();
});

describe("PageLibrary", () => {
  it("should render with May first data if no date is provided", () => {
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
      startOfMonth(new Date(2022, 4)).toISOString().split("T")[0]
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


  it("should request and display new library data when the id prop changes", async () => {
    const wrapper = await shallowMount(PageLibrary, {
      props: {
        id: "oneMockedLibrary",
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

    fetch.mockResponseOnce(JSON.stringify(MOCK_ONE_LIB_FOUND))
    await wrapper.vm.fetchLibraryData();
    await wrapper.vm.$nextTick();
    expect(wrapper.find("h1").text()).toEqual("MOCKED PUBLIC LIBRARY");
    expect(wrapper.vm.fetchedLibraryData).toHaveProperty('libname')

    wrapper.setProps({ id: "anotherMockedLibrary" });
    fetch.mockResponseOnce(JSON.stringify(MOCK_ANOTHER_LIB_FOUND))
    await wrapper.vm.fetchLibraryData();
    await wrapper.vm.$nextTick();
    expect(wrapper.find("h1").text()).toEqual("ANOTHER MOCKED PUBLIC LIBRARY");
    expect(wrapper.vm.fetchedLibraryData).toHaveProperty('libname')


  });
    // note that this should not be required when all REST endpoints return a usable unique library ID
  it("should format a FSCS ID and sequence into a library key", () => {
    // sequence as int
    expect(
      PageLibrary.methods.formatFSCSandSequence("AA0001", 1)
    ).toStrictEqual("AA0001-001");
    // sequence as string
    expect(
      PageLibrary.methods.formatFSCSandSequence("AA0001", "2")
    ).toStrictEqual("AA0001-002");
    // sequence as already-padded string
    expect(
      PageLibrary.methods.formatFSCSandSequence("AA0001", "003")
    ).toStrictEqual("AA0001-003");
    // sequence as already-padded string that's too long
    expect(
      PageLibrary.methods.formatFSCSandSequence("AA0001", "00004")
    ).toStrictEqual("AA0001-004");
  });

});
