import App from './App'

// #ifndef VUE3
import Vue from 'vue'
import './uni.promisify.adaptor'
Vue.config.productionTip = false
App.mpType = 'app'
const app = new Vue({
  ...App
})
// #endif

// #ifdef VUE3
import * as Pinia from 'pinia';
import { createSSRApp } from 'vue'
import uviewPlus from 'uview-plus'
export function createApp() {
  const app = createSSRApp(App)
  app.use(Pinia.createPinia())
  app.use(uviewPlus)
  return {
    app,
	Pinia
  }
}
// #endif
