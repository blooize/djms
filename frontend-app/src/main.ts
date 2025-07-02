/**
 * main.ts
 *
 * Bootstraps Vuetify and other plugins then mounts the App`
 */

// Plugins
import { registerPlugins } from '@/plugins'

// Components
import App from './App.vue'

// Composables
import { createApp } from 'vue'
import { VueUmamiPlugin } from '@jaseeey/vue-umami-plugin';
import router from './router';
// Styles
import 'unfonts.css'

const app = createApp(App)

registerPlugins(app)

app.use(VueUmamiPlugin({
  websiteID: 'ad1e70c2-9959-473e-8510-be545edd26b4',
  scriptSrc: 'http://praxis.club:3000/script.js',
  router,
}))

app.use(router).mount('#app')
