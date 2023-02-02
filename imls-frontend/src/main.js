import { createApp } from "vue";
import App from "./App.vue";
import router from "./router";
import { createMetaManager } from 'vue-meta'
import { plugin as vueMetaPlugin } from 'vue-meta'

const app = createApp(App);

app.use(router);
app.use(createMetaManager())
app.use(vueMetaPlugin)

app.mount("#app");
