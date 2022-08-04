import { createRouter, createWebHistory } from "vue-router";
import Home from "../views/PageHome.vue";

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: "/",
      alias: "/home/",
      name: "home",
      component: Home,
    },
    {
      path: "/about/",
      alias: "/about-us/",
      name: "about",
      // route level code-splitting
      // this generates a separate chunk (About.[hash].js) for this route
      // which is lazy-loaded when the route is visited.
      component: () => import("../views/PageAbout.vue"),
    },
    {
      path: "/search/",
      name: "search",
      component: () => import("../views/PageSearch.vue"),
      props: (route) => ({ query: route.query.query }),
    },
    {
      path: "/sensors/:id/",
      component: () => import("../views/PageSingleSensor.vue"),
      props: true,
    },
    // will match everything and put it under `$route.params.pathMatch`
    {
      path: "/:pathMatch(.*)*",
      name: "NotFound",
      component: () => import("../views/PageNotFound.vue"),
    },
  ],
  scrollBehavior(to, from, savedPosition) {
    // always scroll to top
    return { top: 0 };
  },
  linkActiveClass: "usa-current",
  linkExactActiveClass: "usa-current",
});

export default router;
