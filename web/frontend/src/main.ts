import { createApp } from 'vue'
import { createRouter, createWebHistory } from 'vue-router'
import App from './App.vue'
import RecordingsList from './views/RecordingsList.vue'
import RecordingDetail from './views/RecordingDetail.vue'

const router = createRouter({
  history: createWebHistory(),
  routes: [
    { path: '/', component: RecordingsList },
    { path: '/recordings/:id', component: RecordingDetail },
  ]
})

const app = createApp(App)
app.use(router)
app.mount('#app')
