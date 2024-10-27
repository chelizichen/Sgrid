import 'element-plus/dist/index.css'
import './assets/tailwind.css' // 确保创建此文件
import './assets/main.css'
import { createApp } from 'vue'
import { createPinia } from 'pinia'
import App from './App.vue'
import router from './router'
import ElementPlus from 'element-plus'
import * as ElementPlusIconsVue from '@element-plus/icons-vue'
import 'highlight.js/styles/github.css'
import 'highlight.js/lib/common'
import hljsVuePlugin from '@highlightjs/vue-plugin'
import DataVVue3 from '@kjgl77/datav-vue3'

const app = createApp(App)
for (const [key, component] of Object.entries(ElementPlusIconsVue)) {
  app.component(key, component)
}
app.use(ElementPlus, { zIndex: 3000 })
app.use(hljsVuePlugin)
app.use(createPinia())
app.use(router)
app.use(DataVVue3)

app.mount('#app')
