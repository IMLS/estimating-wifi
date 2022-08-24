import { mount } from "@vue/test-utils";
import PageSearch from "./PageSearch.vue";

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
    expect(wrapper.find("h1").text()).toEqual("This is a Search page");
    expect(wrapper.text()).toContain("search string");
  });
});
