import { createRouter, createWebHistory } from "vue-router";
import Home from "../views/PageHome.vue";

const routes = [
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
    path: "/library/:fscs_id/",
    component: () => import("../views/PageLibrary.vue"),
    props: (route) => ({
      selectedDateFromParams: route.query.date,
      id: route.params.fscs_id,
    }),
  },
  {
    path: "/state/:state_initials/",
    component: () => import("../views/PageState.vue"),
    props: (route) => ({
      stateInitials: route.params.state_initials,
    }),
  },
  // will match everything and put it under `$route.params.pathMatch`
  {
    path: "/:pathMatch(.*)*",
    name: "NotFound",
    component: () => import("../views/PageNotFound.vue"),
  },
];

let router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: routes,
  scrollBehavior(to, from, savedPosition) {
    // always scroll to top
    return { top: 0 };
  },
  linkActiveClass: "usa-current",
  linkExactActiveClass: "usa-current",
});

export { routes };
export default router;
