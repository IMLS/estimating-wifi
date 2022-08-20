import { mount } from "@vue/test-utils";
import USWDSTable from "./USWDSTable.vue";

describe("USWDSUSWDSTable", () => {
  it("should render a table", () => {
    const wrapper = mount(USWDSTable);
    expect(wrapper.findAll(".usa-table")).toHaveLength(1);
  });

  it("should render a caption when provided", async () => {
    const wrapper = mount(USWDSTable);
    expect(wrapper.find("caption").text()).toEqual("");
    await wrapper.setProps({ caption: "New caption" });
    expect(wrapper.find("caption").text()).toEqual("New caption");
  });

  it("should render the specificed number of table rows and columns", () => {
    const wrapper = mount(USWDSTable, {
      props: {
        rows: [
          ["row 1, column A", "row 1, column B"],
          ["row 2, column A", "row 2, column B"],
        ],
        headers: ["Header for column A", "Header for column B"],
      },
    });

    expect(wrapper.find("th").text()).toContain("Header for column");
  });
});
