import { mount } from "@vue/test-utils";
import USWDSDatePicker from "./USWDSDatePicker.vue";
import { vi } from "vitest";
import { format, startOfYesterday } from "date-fns";

// don't test the external UWSDS code
vi.mock("uswds/src/js/components/date-picker");

describe("USWDSDatePicker", () => {
  it("should render a date picker", async () => {
    const wrapper = mount(USWDSDatePicker);
    expect(wrapper.findAll(".usa-date-picker")).toHaveLength(1);
  });

  it("should limit date selections to established range", async () => {
    const wrapper = mount(USWDSDatePicker);
    const input = wrapper.find('.usa-input[name="date"]');
    await input.setValue(null); // USWDS code returns null on the input for all date values outside of min/max
    await wrapper.vm.$nextTick();

    let yesterdayDateReadable = format(startOfYesterday(), 'MM/dd/yyyy')
    expect(wrapper.vm.minDateReadable).toStrictEqual("01/01/2022");
    expect(wrapper.vm.maxDateReadable).toStrictEqual(yesterdayDateReadable);

    expect(wrapper.find("#date-outside-range")).toBeTruthy();
    expect(USWDSDatePicker.methods.formatReadableDate("2022-02-01")).toStrictEqual("02/01/2022");
  });

  it("should use today's date by default", async () => {
    const wrapper = mount(USWDSDatePicker);
    expect(
      wrapper.find(".usa-date-picker").attributes("data-default-value")
    ).toEqual(new Date().toISOString().split("T")[0]);
  });

  it("uses a specific date if provided", async () => {
    const wrapper = mount(USWDSDatePicker, {
      props: {
        initialDate: "2022-01-01",
      },
    });
    expect(
      wrapper.find(".usa-date-picker").attributes("data-default-value")
    ).toEqual("2022-01-01");
  });

  it("should update state if the selected date changes", async () => {
    const spyChange = vi.spyOn(USWDSDatePicker.methods, "detectChange");
    const wrapper = mount(USWDSDatePicker);
    expect(spyChange).toHaveBeenCalledTimes(0);
    const input = wrapper.find('.usa-input[name="date"]');
    await input.setValue("1999-12-31");
    expect(wrapper.emitted('date_changed')).toHaveLength(1)
    expect(wrapper.emitted('date_changed')[0]).toEqual(["1999-12-31"])
    expect(spyChange).toHaveBeenCalledTimes(1);
    await wrapper.vm.$nextTick();
  });
});
