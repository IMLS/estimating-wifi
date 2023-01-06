import { mount, flushPromises } from "@vue/test-utils";
import PageSearch from "./PageSearch.vue";
import { expect } from "vitest";

import { createRouter, createWebHistory } from "vue-router";
import { routes } from "../router/index.js";

let router;

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
]

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

  it("should render a message when the request fails", async () => {
    const wrapper = mount(PageSearch, {
      props: {
        query: "mocked error response",
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
    await wrapper.vm.searchLibraryNames();
    await wrapper.vm.$nextTick();

    expect(wrapper.find("h1").text()).toEqual("Libraries matching \"mocked error response\"");
    expect(wrapper.find("p").text()).toEqual("Sorry, no matching libraries found.");
    expect(wrapper.find("p + span").text()).to.contain(MOCK_ERROR_MSG);
    expect(wrapper.findAll("ol.usa-list li").length).toEqual(0);
    expect(wrapper.vm.fetchedLibraries).toBeNull();
    expect(wrapper.vm.fetchError).to.contain(MOCK_ERROR_MSG);
  });

  it("should render a message when no matching libraries are found", async () => {
    const wrapper = mount(PageSearch, {
      props: {
        query: "mocked empty response",
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

    fetch.mockResponseOnce(JSON.stringify(MOCK_NO_LIBS_FOUND))
    await wrapper.vm.searchLibraryNames();
    await wrapper.vm.$nextTick();

    expect(wrapper.find("h1").text()).toEqual("Libraries matching \"mocked empty response\"");
    expect(wrapper.find("p").text()).toEqual("Sorry, no matching libraries found.");
    expect(wrapper.findAll("ol.usa-list li").length).toEqual(0);
    expect(wrapper.vm.fetchedLibraries.length).toEqual(0);
    
  });

  it("should render a list of libraries matching the given search string", async () => {
    const wrapper = await mount(PageSearch, {
      props: {
        query: "Mocked library",
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
    await wrapper.vm.searchLibraryNames();
    await wrapper.vm.$nextTick();

    expect(wrapper.vm.isLoading).toBeFalsy();
    expect(wrapper.find("h1").text()).toEqual("Libraries matching \"Mocked library\"");
    expect(wrapper.findAll("ol.usa-list li").length).toBeGreaterThanOrEqual(1);
    expect(wrapper.vm.fetchedLibraries.length).toBeGreaterThanOrEqual(1);
  });

  it("should update with new results when the query changes", async () => {
    const wrapper = await mount(PageSearch, {
      props: {
        query: "",
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
    await wrapper.setProps({ query: "Mocked library" });
    await wrapper.vm.searchLibraryNames();
    await wrapper.vm.$nextTick();

    expect(wrapper.vm.isLoading).toBeFalsy();
    expect(wrapper.find("h1").text()).toEqual("Libraries matching \"Mocked library\"");
    expect(wrapper.findAll("ol.usa-list li").length).toBeGreaterThanOrEqual(1);
    expect(wrapper.vm.fetchedLibraries.length).toBeGreaterThanOrEqual(1);
  });

  // note that this should not be required when all REST endpoints return a usable unique library ID
  it("should format a FSCS ID and sequence into a library key", () => {
    // sequence as int
    expect(
      PageSearch.methods.formatFSCSandSequence("AA0001", 1)
    ).toStrictEqual("AA0001-001");
    // sequence as string
    expect(
      PageSearch.methods.formatFSCSandSequence("AA0001", "2")
    ).toStrictEqual("AA0001-002");
    // sequence as already-padded string
    expect(
      PageSearch.methods.formatFSCSandSequence("AA0001", "003")
    ).toStrictEqual("AA0001-003");
    // sequence as already-padded string that's too long
    expect(
      PageSearch.methods.formatFSCSandSequence("AA0001", "00004")
    ).toStrictEqual("AA0001-004");
  });

  it("should format a page title using the current query", () => {
    let pageMeta = PageSearch.metaInfo('searched query');
    expect(pageMeta).toHaveProperty("title");
    expect(pageMeta.title).toStrictEqual('Library search results for "searched query"');
  });
});
