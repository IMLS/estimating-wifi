import { mount } from "@vue/test-utils";
import { expect } from "vitest";
import PageLibrary from "./PageLibrary.vue";

describe("PageLibrary", () => {
  it("should render", () => {
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
        ],
      },
    });
    expect(wrapper.find("h1").text()).toEqual("Library KnownGoodId");
    expect(wrapper.findAll(".usa-card")).toHaveLength(1);
  });
});
