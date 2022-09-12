import { mount } from "@vue/test-utils";
import { expect } from "vitest";
import Heatmap from "./Heatmap.vue";

describe("Heatmap", () => {
  it("should render a table", () => {
    const wrapper = mount(Heatmap, {
    });
    expect(wrapper.findAll('.data-grid')).toHaveLength(1);
  });
  // it("should sort an array to calculate percentile", () => {
  //   const arrayToSort = [3,4,5,1,2];
  //     expect(Heatmap.methods.sortArrayAscending.call(arrayToSort)).toBe("1, 2, 3, 4, 5")
  // });
  it("should render a matrix of values in columns and rows, with percentile shading", () => {
   const wrapper = mount(Heatmap, { props: {
        dataset: [ [5,4,3,1,2], [1,1,1,1,2], [4,4,4,4,0] ],
        datasetLabels: ["FirstRow", "SecondRow", "ThirdRow"],
        binLabels: ["Column A", "Column B", "Column C", "Column D", "Column E"]
      }
    });

    function getAlphaFromRGBAColor(colorString) {
      if (colorString.startsWith("rgb(")) return 1;
      if (colorString.startsWith("rgba(")) return parseFloat(colorString.split(', ').pop().split(')')[0]);
      return null;
    }
    
    // store the rendered cells for later use
    const allValuesRendered = wrapper.findAll('.data-grid__cell');

    // should be the size of the matrix
    expect(allValuesRendered).toHaveLength(5 * 3);
    // the first value in the sample dataset is also the highest value
    expect(allValuesRendered[0].attributes("data-percentile")).toEqual("100");
    // 100th percentile color should have no alpha channel / be at 100% opacity
    expect(getAlphaFromRGBAColor(allValuesRendered[0].element.style.backgroundColor)).toEqual(1);

    // the fourth value in the sample dataset is also at the 40th value
    expect(allValuesRendered[3].attributes("data-percentile")).toEqual("40");
    // 100th percentile color should have no alpha channel / be at 100% opacity
    expect(getAlphaFromRGBAColor(allValuesRendered[3].element.style.backgroundColor)).toEqual(0.4);

    // the last value in the sample dataset is also the lowest value
    expect(allValuesRendered[14].attributes("data-percentile")).toEqual("0");
    // 0th percentile color should have 0% alpha channel / be at 0% opacity
    expect(getAlphaFromRGBAColor(allValuesRendered[14].element.style.backgroundColor)).toEqual(0);

  });
});
