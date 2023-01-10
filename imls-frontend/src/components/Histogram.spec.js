import { mount } from "@vue/test-utils";
import Histogram from "./Histogram.vue";
import { Bar } from "vue-chartjs";

describe("Histogram", () => {
  it("should render a Bar Chart", () => {
    const wrapper = mount(Histogram, {
      shallow: true,
    });
    expect(wrapper.findAllComponents(Bar)).toHaveLength(1);
  });
  it("should format numbers and replace 0 with an en dash", () => {
    expect(Histogram.methods.formatNumbers(1)).toEqual(1);
    // note that for now, we aren't using commas
    expect(Histogram.methods.formatNumbers(1000)).toEqual(1000);
    expect(Histogram.methods.formatNumbers(0)).toEqual('–');
  });
});
