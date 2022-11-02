import { mount } from "@vue/test-utils";
import USWDSBreadcrumb from "./USWDSBreadcrumb.vue";

describe("USWDSBreadcrumb", () => {
  it("should render breadcrumbs", () => {
    const wrapper = mount(USWDSBreadcrumb);
    expect(wrapper.findAll(".usa-breadcrumb__list")).toHaveLength(1);
    expect(wrapper.findAll(".usa-breadcrumb__list-item")).toHaveLength(2);
    expect(wrapper.findAll(".usa-breadcrumb__link")).toHaveLength(1);
  });

});
