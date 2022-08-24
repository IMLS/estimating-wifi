import { mount } from "@vue/test-utils";
import PageHome from "./PageHome.vue";

describe("PageHome", () => {
  it("should render ", () => {
    const wrapper = mount(PageHome, {
      global: {
        stubs: ["router-link", "router-view", "RouterView", "RouterLink"],
      },
    });
    expect(wrapper.find("h1").text()).toEqual("Home");
    expect(wrapper.text()).toContain("I am a homepage");
  });
});
