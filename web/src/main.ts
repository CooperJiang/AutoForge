import './styles/main.css'

import { createApp } from 'vue'
import { createPinia } from 'pinia'

import App from './App.vue'
import router from './router'

// 导入安全存储工具
import SecureStorage from './utils/storage'

// 初始化安全存储
SecureStorage.migrateFromOldStorage() // 清理旧数据
SecureStorage.cleanExpiredItems() // 清理过期数据

const app = createApp(App)

app.use(createPinia())
app.use(router)

app.mount('#app')
