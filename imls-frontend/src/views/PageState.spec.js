import { mount, shallowMount, flushPromises } from "@vue/test-utils";
import { expect, v1 } from "vitest";
import PageState from "./PageState.vue";
import { createRouter, createWebHistory } from "vue-router";
import { routes } from "../router/index.js";
import "whatwg-fetch";

let router;

beforeEach(async () => {
  router = createRouter({
    history: createWebHistory(),
    routes: routes,
  });
  router.push("/");
  await router.isReady();
  vi.useFakeTimers();
});

describe("PageState", () => {

  it("should render a list of libraries", async () => {
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
    setTimeout(() => {
      expect(wrapper.find("h1").text()).toEqual("Alaska Public Libraries");
      expect(wrapper.findAll("ol.usa-list li").length).toBeGreaterThanOrEqual(1);
      expect(wrapper.vm.fetchedData.length).toBeGreaterThanOrEqual(1);
    }, 50)

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

    setTimeout(() => {    
      expect(spyRedirect).toHaveBeenCalledTimes(1);
      expect(spyRedirect).toHaveBeenCalledWith({"name": "NotFound"});
    }, 50)

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
  
    wrapper.setProps({ stateInitials: "AL" });

    setTimeout(() => {
      expect(wrapper.findAll(".loaded--has-data")).toHaveLength(1);
      expect(wrapper.findAll(".loaded--no-data")).toHaveLength(0);
      expect(wrapper.find("h1").text()).toEqual("Alabama Public Libraries");
      expect(wrapper.findAll("ol.usa-list li").length).toBeGreaterThanOrEqual(1);
      expect(wrapper.findAll(".loaded--has-data")).toHaveLength(1);
    }, 50)
  });
  
});
