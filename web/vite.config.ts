import { fileURLToPath, URL } from 'node:url'

import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'
import vueDevTools from 'vite-plugin-vue-devtools'

// https://vite.dev/config/
export default defineConfig({
  plugins: [
    vue(),
    vueDevTools(),
  ],
  resolve: {
    alias: {
      '@': fileURLToPath(new URL('./src', import.meta.url))
    },
  },
  server: {
    port: 3200,
    proxy: {
      '/api': {
        target: 'http://localhost:7777',
        changeOrigin: true,
      },
    },
  },
  build: {
    rollupOptions: {
      output: {
        // 优化文件名，将下划线替换为中划线，避免特殊字符
        chunkFileNames: (chunkInfo) => {
          const name = chunkInfo.name.replace(/_/g, '-')
          return `js/${name}-[hash].js`
        },
        entryFileNames: (chunkInfo) => {
          const name = chunkInfo.name.replace(/_/g, '-')
          return `js/${name}-[hash].js`
        },
        assetFileNames: (assetInfo) => {
          const name = assetInfo.name.replace(/_/g, '-')
          const info = name.split('.')
          const ext = info[info.length - 1]
          if (/\.(png|jpe?g|gif|svg|ico)$/i.test(name)) {
            return `images/${name.replace(/\.[^.]+$/, '')}-[hash].${ext}`
          }
          if (/\.(woff2?|eot|ttf|otf)$/i.test(name)) {
            return `fonts/${name.replace(/\.[^.]+$/, '')}-[hash].${ext}`
          }
          return `assets/${name.replace(/\.[^.]+$/, '')}-[hash].${ext}`
        },
      },
    },
  },
})
