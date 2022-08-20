import { mount } from "@vue/test-utils";
import PageAbout from "./PageAbout.vue";

describe("PageAbout", () => {
  it("should render", () => {
    const wrapper = mount(PageAbout, {
      global: {
        stubs: ["router-link", "router-view", "RouterView", "RouterLink"],
      },
    });
    expect(wrapper.find("h1").text()).toEqual("This is an about page");
    expect(wrapper.text()).toContain(
      "I'm the first paragraph. I get special styling in USWDS Content."
    );
  });
});
