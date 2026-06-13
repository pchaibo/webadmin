import { fileURLToPath, URL } from 'node:url'

import { defineConfig, loadEnv } from 'vite'
import vue from '@vitejs/plugin-vue'
import vueJsx from '@vitejs/plugin-vue-jsx'
import vueDevTools from 'vite-plugin-vue-devtools'

// https://vite.dev/config/
export default defineConfig(({ mode }) => {
  const env = loadEnv(mode, process.cwd(), '')
  const isDev = mode === 'development'

  return {
    plugins: [vue(), vueJsx(), vueDevTools()],
    base: isDev ? '/' : '/ul/',
    build: {
      outDir: '../Backend/frontend',
      emptyOutDir: true,
    },
    server: {
      proxy: {
        '/api': {
          target: env.VITE_API_TARGET_URL || 'http://127.0.0.1:4000',
          changeOrigin: true,
         // rewrite: (path) => path.replace(/^\/api/, ''),
        },
      },
    },
    resolve: {
      alias: {
        '@': fileURLToPath(new URL('./src', import.meta.url)),
      },
    },
  }
})
