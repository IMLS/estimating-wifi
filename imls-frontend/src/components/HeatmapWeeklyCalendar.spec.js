import { mount } from "@vue/test-utils";
import { expect } from "vitest";
import HeatmapWeeklyCalendar from "./HeatmapWeeklyCalendar.vue";

describe("HeatmapWeeklyCalendar", () => {
  it("should render a table", () => {
    const wrapper = mount(HeatmapWeeklyCalendar, {});
    expect(wrapper.findAll(".weekly-calendar")).toHaveLength(1);
  });
  it("should render a matrix of values in columns and rows, with percentile shading", () => {
    const wrapper = mount(HeatmapWeeklyCalendar, {
      props: {
        selectedDate: "1999-12-31", // note: this is on a Friday
        weekStartDateISO: "1999-12-26",
        dataset: [
          [5, 4],
          [1, 2],
          [2, 3],
          [3, 4],
          [4, 5],
          [5, 1],
          [5, 0],
        ],
        binLabels: ["Column A", "Column B"],
      },
    });

    function getAlphaFromRGBAColor(colorString) {
      if (colorString.startsWith("rgb(")) return 1;
      if (colorString.startsWith("rgba("))
        return parseFloat(colorString.split(", ").pop().split(")")[0]);
      return null;
    }
    // store the rendered cells for later use
    const allValuesRendered = wrapper.findAll(
      ".weekly-calendar__day .weekly-calendar__cell"
    );

    // should be the size of the matrix
    expect(allValuesRendered).toHaveLength(2 * 7);
    // the first value in the sample dataset is also the highest percentile
    expect(allValuesRendered[0].attributes("data-percentile")).toEqual("100");
    // 100th percentile color should have no alpha channel / be at 100% opacity
    expect(
      getAlphaFromRGBAColor(allValuesRendered[0].element.style.backgroundColor)
    ).toEqual(1);

    // the sixth value in the sample dataset is also at the 40th percentile
    expect(allValuesRendered[5].attributes("data-percentile")).toEqual("50");
    // 100th percentile color should have no alpha channel / be at 100% opacity
    expect(
      getAlphaFromRGBAColor(allValuesRendered[5].element.style.backgroundColor)
    ).toEqual(0.5);

    // the last value in the sample dataset is also the lowest percentile
    expect(allValuesRendered[13].attributes("data-percentile")).toEqual("0");
    // 0th percentile color should have 0% alpha channel / be at 0% opacity
    expect(
      getAlphaFromRGBAColor(allValuesRendered[13].element.style.backgroundColor)
    ).toEqual(0);

    // the selected date should have special styling, but no others
    expect(
      wrapper.findAll(".weekly-calendar__day")[5].classes("isSelectedDate")
    ).toBe(true);
    expect(
      wrapper.findAll(".weekly-calendar__day")[6].classes("isSelectedDate")
    ).toBe(false);
  });
});
