import { mount } from "@vue/test-utils";
import Histogram from "./Histogram.vue";
import { Bar } from 'vue-chartjs'


describe("Histogram", () => {
it("should render a Bar Chart", () => {
    const wrapper = mount(Histogram, {
      shallow: true
    });
    expect(wrapper.findAllComponents(Bar)).toHaveLength(1);
  });
});
