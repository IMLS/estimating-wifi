import { mount } from "@vue/test-utils";
import App from "./App.vue";
import { routes } from "./router/index.js";
import { createRouter, createWebHistory } from "vue-router";
let router;

beforeEach(async () => {
  router = createRouter({
    history: createWebHistory(import.meta.env.BASE_URL),
    routes: routes,
  });
  router.push("/");
  await router.isReady();
});

describe("App", () => {
  it("should render a Header and Footer", () => {
    const wrapper = mount(App, {
      global: {
        plugins: [router],
      },
    });

    expect(wrapper.findAll(".usa-header")).toHaveLength(1);
    expect(wrapper.findAll(".usa-footer")).toHaveLength(1);
  });
});
