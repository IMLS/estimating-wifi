import { fileURLToPath, URL } from 'url'
import process from 'process';

import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'

// https://vitejs.dev/config/
export default defineConfig({
  plugins: [vue()],
  resolve: {
    alias: {
      '@': fileURLToPath(new URL('./src', import.meta.url))
    }
  },
  // base: process.env.NODE_ENV === 'production'
  // // process.env.BASEURL should be '/site/[ORG_NAME]/[REPO_NAME]' on federalist
  //   ? process.env.BASEURL + '/'
  //   : '/'
})
