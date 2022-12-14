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
  await flushPromises();
  fetch.resetMocks();
});

const MOCK_ERROR_MSG = "mocked error message";
// the API currently returns null instead of an empty array on no matches
const MOCK_NO_LIBS_FOUND = null;
const MOCK_ONE_LIB_FOUND = [
  {
    "stabr":"MK",
    "fscskey":"MOCK001",
    "fscs_seq":1,
    "libname":"MOCKED PUBLIC LIBRARY",
    "address":"1234 MOCKINGBIRD ROAD",
    "city":"MOUNT MOCKINGTON",
    "zip":"00000"
  }
];


describe("PageState", () => {
  it("should render", () => {
    const wrapper = mount(PageState, {
      props: {
        stateInitials: "AK",
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
    expect(wrapper.find("h1").text()).toEqual("Alaska Public Libraries");

  });

  it("should render a message when the request fails", async () => {
    const wrapper = mount(PageState, {
      props: {
        stateInitials: "AK",
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

    fetch.mockRejectOnce(new Error(MOCK_ERROR_MSG))
    await wrapper.vm.fetchLibraries();
    await wrapper.vm.$nextTick();

    expect(wrapper.find("p").text()).toEqual("Sorry, no matching libraries found.");
    expect(wrapper.find("p + span").text()).to.contain(MOCK_ERROR_MSG);
    expect(wrapper.findAll("ol.usa-list li").length).toEqual(0);
    expect(wrapper.vm.fetchedLibraries.length).toEqual(0);
    expect(wrapper.vm.fetchError).to.contain(MOCK_ERROR_MSG);
  });


  it("should render a list of libraries when the request succeeds", async () => {
    const wrapper = await mount(PageState, {
      props: {
        stateInitials: "AK",
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

    fetch.mockResponseOnce(JSON.stringify(MOCK_ONE_LIB_FOUND))
    await wrapper.vm.fetchLibraries();
    await flushPromises();
    await wrapper.vm.$nextTick();

    expect(wrapper.vm.isLoading).toBeFalsy();
    expect(wrapper.find("h1").text()).toEqual("Alaska Public Libraries");
    expect(wrapper.findAll("ol.usa-list li").length).toBeGreaterThanOrEqual(1);
    expect(wrapper.vm.fetchedLibraries.length).toBeGreaterThanOrEqual(1);

  });
  

  it("should navigate to 404 if a bad state abbr is provided", async () => {
    router.push('/');

    // After this line, router is ready
    await router.isReady();

    const spyRedirect = vi.spyOn(
      router,
      "push"
    );
    const wrapper = shallowMount(PageState, {
      props: {
        stateInitials: "ZZ",
      },
      global: {
        plugins: [router],
      },
    });

    
    PageState.beforeRouteEnter.call(wrapper.vm, undefined, undefined, (c) => c(wrapper.vm));

    expect(spyRedirect).toHaveBeenCalledTimes(1);
    expect(spyRedirect).toHaveBeenCalledWith({"name": "NotFound"});

  });

  it("should update with new libraries when the state abbr prop changes", async () => {
    const wrapper = await shallowMount(PageState, {
      props: {
        stateInitials: "AK",
      },
      global: {
        plugins: [router],
      },
    });
    fetch.mockResponseOnce(JSON.stringify(MOCK_NO_LIBS_FOUND))
    await wrapper.vm.fetchLibraries();
    await flushPromises();
    await wrapper.vm.$nextTick();

    expect(wrapper.findAll(".loaded--has-data")).toHaveLength(0);
    expect(wrapper.findAll(".loaded--no-data")).toHaveLength(1);
    expect(wrapper.find("h1").text()).toEqual("Alaska Public Libraries");
        expect(wrapper.findAll("ol.usa-list li").length).toBe(0);

    await wrapper.setProps({ stateInitials: "AL" });
    fetch.mockResponseOnce(JSON.stringify(MOCK_ONE_LIB_FOUND))
    await wrapper.vm.fetchLibraries();
    await flushPromises();
    await wrapper.vm.$nextTick();

    expect(wrapper.findAll(".loaded--has-data")).toHaveLength(1);
    expect(wrapper.findAll(".loaded--no-data")).toHaveLength(0);
    expect(wrapper.find("h1").text()).toEqual("Alabama Public Libraries");
    expect(wrapper.findAll("ol.usa-list li").length).toBeGreaterThanOrEqual(1);
  });
  
    // note that this should not be required when all REST endpoints return a usable unique library ID
  it("should format a FSCS ID and sequence into a library key", () => {
    // sequence as int
    expect(
      PageState.methods.formatFSCSandSequence("AA0001", 1)
    ).toStrictEqual("AA0001-001");
    // sequence as string
    expect(
      PageState.methods.formatFSCSandSequence("AA0001", "2")
    ).toStrictEqual("AA0001-002");
    // sequence as already-padded string
    expect(
      PageState.methods.formatFSCSandSequence("AA0001", "003")
    ).toStrictEqual("AA0001-003");
    // sequence as already-padded string that's too long
    expect(
      PageState.methods.formatFSCSandSequence("AA0001", "00004")
    ).toStrictEqual("AA0001-004");
  });
});
