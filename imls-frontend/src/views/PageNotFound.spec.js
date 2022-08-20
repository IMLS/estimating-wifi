import { mount } from "@vue/test-utils";
import PageNotFound from "./PageNotFound.vue";

describe("PageNotFound", () => {
  it("should render ", () => {
    const wrapper = mount(PageNotFound, {
      global: {
        stubs: ["router-link", "router-view", "RouterView", "RouterLink"], 
      }
    });
    expect(wrapper.find("h1").text()).toEqual("This is a 404 page");
    expect(wrapper.text()).toContain("Sorry, we couldn't find anything.");

  });
});
