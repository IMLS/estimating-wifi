import { fileURLToPath, URL } from "url";
import process from "process";

import { defineConfig } from "vite";
import vue from "@vitejs/plugin-vue";

// https://vitejs.dev/config/
export default defineConfig({
  plugins: [
    vue({
      template: {
        // compilerOptions: {
        //   isCustomElement: () => true
        // }
      },
    }),
  ],
  resolve: {
    alias: {
      "@": fileURLToPath(new URL("./src", import.meta.url)),
    },
  },
  test: {
    // https://github.com/vitest-dev/vitest/issues/740
    // threads: false,
    globals: true,
    environment: "jsdom",
    coverage: {
      "check-coverage": true,
      reporter: ["lcovonly", "text", "text-summary"],
      all: true,
      include: ["src/**/*.vue"],
      lines: 90,
      branches: 90,
      functions: 90,
    },
    reporters: "verbose",
    setupFiles: ["test/setup.js"],
  },
  // base: process.env.NODE_ENV === 'production'
  // // process.env.BASEURL should be '/site/[ORG_NAME]/[REPO_NAME]' on federalist
  //   ? process.env.BASEURL + '/'
  //   : '/'
});
