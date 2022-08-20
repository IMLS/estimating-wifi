import { mount } from "@vue/test-utils";
import USWDSCard from "./USWDSCard.vue";

describe("USWDSCard", () => {
  it("should render a card", () => {
    const wrapper = mount(USWDSCard);
    expect(wrapper.findAll(".usa-card__container")).toHaveLength(1);
    expect(wrapper.findAll(".usa-card__header")).toHaveLength(1);
    expect(wrapper.findAll(".usa-card__body")).toHaveLength(1);
  });

  it("should render a card title regardless of whether a title is provided", async () => {
    const wrapper = mount(USWDSCard);
    expect(wrapper.find(".usa-card__heading").text()).toEqual("No title content provided");
    await wrapper.setProps({ title: 'New title' })
    expect(wrapper.find(".usa-card__heading").text()).toEqual("New title");
  });
});
