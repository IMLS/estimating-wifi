import { shallowMount, flushPromises } from "@vue/test-utils";
import FetchData from "./FetchData.vue";
import { afterEach, vi } from "vitest";
import "whatwg-fetch";
import { enableAutoUnmount } from "@vue/test-utils";

enableAutoUnmount(afterEach);
beforeEach(async () => {
  vi.useFakeTimers();
});

describe("FetchData", () => {
  it("should render a loading area", async () => {
    const wrapper = await shallowMount(FetchData);
    expect(wrapper.findAll(".loading-area")).toHaveLength(1);
  });

  it("should render loaded data when the backend returns good data", async () => {
    const wrapper = await shallowMount(FetchData, {
      props: {
        path: "/rpc/bin_devices_per_hour",
        fscsId: "KnownGoodId",
        selectedDate: "2022-05-01",
      },
    });
    setTimeout(() => {
      expect(wrapper.findAll(".loaded--has-data")).toHaveLength(1);
    }, 50)
  });

  it("should render a message when the backend returns no data for the requested query (such as a missing date)", async () => {
    const wrapper = await shallowMount(FetchData, {
      props: { path: "/rpc/bin_devices_per_hour", fscsId: "KnownEmptyId" },
    });
    setTimeout(() => {
      expect(wrapper.findAll(".loaded--no-data")).toHaveLength(1);
    }, 50)

  });

  it("should update with new data when the library id prop changes", async () => {
    const wrapper = await shallowMount(FetchData, {
      props: {
        path: "/rpc/bin_devices_per_hour",
        fscsId: "KnownGoodId",
        selectedDate: "2022-05-01",
      },
    });
    
    setTimeout(() => {
      expect(wrapper.findAll(".loaded--has-data")).toHaveLength(1);
    }, 50)

    wrapper.setProps({ fscsId: "KnownEmptyId" });

    setTimeout(() => {
      expect(wrapper.findAll(".loaded--has-data")).toHaveLength(0);
      expect(wrapper.findAll(".loaded--no-data")).toHaveLength(1);
    }, 50)

    wrapper.unmount();
  });

  it("should update with new data when the selected date prop changes", async () => {
    const spyChangeDate = vi.spyOn(FetchData.methods, "fetchData");
    const wrapper = await shallowMount(FetchData, {
      props: {
        path: "/rpc/bin_devices_per_hour",
        fscsId: "KnownGoodId",
        selectedDate: "2022-05-01",
      },
    });

    setTimeout(() => {
      expect(spyChangeDate).toHaveBeenCalledTimes(1);
      expect(wrapper.vm.fetchedData).toEqual([0, 1, 2, 3, 4]);
    }, 50)

    wrapper.setProps({ selectedDate: "9999-99-99" });

    setTimeout(() => {
      expect(spyChangeDate).toHaveBeenCalledTimes(2);
      expect(wrapper.vm.fetchedData).toEqual([9, 9, 9, 9, 9]);
    }, 50)
  });

  it("should render a loading indicator before the backend resolves", async () => {
    const wrapper = await shallowMount(FetchData, {
      props: { path: "/rpc/bin_devices_per_hour", fscsId: "notARealID" },
    });
    setTimeout(() => {
      expect(wrapper.findAll(".loading-indicator")).toHaveLength(1);
    }, 10)
  
  });

  it("should not render a loading indicator after the backend resolves", async () => {
    const wrapper = await shallowMount(FetchData, {
      props: { path: "/rpc/bin_devices_per_hour", fscsId: "notARealID" },
    });
    setTimeout(() => {
      expect(wrapper.findAll(".loading-indicator")).toHaveLength(0);
    }, 50)
  });

  it("should render a loading error when the backend returns an error", async () => {
    const wrapper = await shallowMount(FetchData, {
      props: { path: "/rpc/bin_devices_per_hour", fscsId: "notARealID" },
    });
    setTimeout(() => {
      expect(wrapper.findAll(".loaded--error")).toHaveLength(1);
    }, 50)
  });

  it("should compose params in a query string from an object's key/value pairs", () => {
    expect(
      FetchData.computed.queryString.call({
        queryParams: { key1: "val1", key2: "val2" },
      })
    ).toBe("&key1=val1&key2=val2");
  });
});
