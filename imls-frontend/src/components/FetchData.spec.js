import { mount, flushPromises } from "@vue/test-utils";
import FetchData from "./FetchData.vue";
import { afterEach, vi } from 'vitest';
import fetch from 'node-fetch';
import { enableAutoUnmount } from '@vue/test-utils'
import { state } from "@/store/store.js";

enableAutoUnmount(afterEach)

describe("FetchData", () => {
  it("should render a loading area", () => {
    const wrapper = mount(FetchData);
    expect(wrapper.findAll(".loading-area")).toHaveLength(1);
    
  });

  it("should render loaded data when the backend returns good data", async () => {
    const wrapper = await mount(FetchData, { props: { path: "/rpc/bin_devices_per_hour", fscsId: "KnownGoodId"} });
    await flushPromises();
    expect(await wrapper.findAll(".loaded--has-data")).toHaveLength(1);
  });

  it("should render a message when the backend returns no data for the requested query", async () => {
    const wrapper = await mount(FetchData, { props: { path: "/rpc/bin_devices_per_hour", fscsId: "KnownEmptyId"} });
    await flushPromises();
    expect(await wrapper.findAll(".loaded--no-data")).toHaveLength(1);
  });

  it("should update with new data when the prop for library id changes", async () => {
    const wrapper = await mount(FetchData, { props: { path: "/rpc/bin_devices_per_hour", fscsId: "KnownGoodId"} });
      await flushPromises();
      expect(await wrapper.findAll(".loaded--has-data")).toHaveLength(1);
      await wrapper.setProps({fscsId: "KnownEmptyId" });
      await flushPromises();
      expect(await wrapper.findAll(".loaded--has-data")).toHaveLength(0);
      expect(await wrapper.findAll(".loaded--no-data")).toHaveLength(1);
      wrapper.unmount();
  });

  it("should update with new data when the selected date changes in state", async () => {
    const spyChangeDate = vi.spyOn(FetchData.methods, 'fetchDataFromState') 
    const wrapper = await mount(FetchData, { props: { path: "/rpc/bin_devices_per_hour", fscsId: "KnownGoodId"} });
      await flushPromises();
      expect(spyChangeDate).toHaveBeenCalledTimes(1)
      expect(wrapper.vm.state.fetchedData).toEqual([0,1,2,3,4])
      state.selectedDate =  "9999-99-99";
      await flushPromises();
      await wrapper.vm.$nextTick();
      expect(spyChangeDate).toHaveBeenCalledTimes(2);
      expect(wrapper.vm.state.fetchedData).toEqual([9,9,9,9,9])
  });

  it("should render a loading indicator before the backend resolves", async () => {
    const wrapper = mount(FetchData, { props: { path: "/rpc/bin_devices_per_hour", fscsId: "notARealID"} });
    expect(await wrapper.findAll(".loading-indicator")).toHaveLength(1);
    await flushPromises();
  });

    it("should not render a loading indicator after the backend resolves", async () => {
    const wrapper = mount(FetchData, { props: { path: "/rpc/bin_devices_per_hour", fscsId: "notARealID"} });
    await flushPromises();
    expect(await wrapper.findAll(".loading-indicator")).toHaveLength(0);
  });

  it("should render a loading error when the backend returns an error", async () => {
    const wrapper = await mount(FetchData, { props: { path: "/rpc/bin_devices_per_hour", fscsId: "notARealID"} });
    await flushPromises();
    expect(await wrapper.findAll(".loaded--error")).toHaveLength(1);
  });

});
