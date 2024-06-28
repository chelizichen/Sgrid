import 'element-plus/dist/index.css'
import './assets/main.css'
import { createApp } from 'vue'
import { createPinia } from 'pinia'
import App from './App.vue'
import router from './router'
import ElementPlus from 'element-plus'
import * as ElementPlusIconsVue from '@element-plus/icons-vue'
import 'highlight.js/styles/github.css'
import 'highlight.js/lib/common'
import hljsVuePlugin from "@highlightjs/vue-plugin";
const app = createApp(App)
for (const [key, component] of Object.entries(ElementPlusIconsVue)) {
  app.component(key, component)
}
app.use(ElementPlus, { zIndex: 3000 })
app.use(hljsVuePlugin)
app.use(createPinia())
app.use(router)
app.mount('#app')
