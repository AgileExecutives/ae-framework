import { fileURLToPath, URL } from 'node:url'
import path from 'node:path'

import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'
import vueDevTools from 'vite-plugin-vue-devtools'
import tailwindcss from '@tailwindcss/vite'

// https://vite.dev/config/
export default defineConfig({
  plugins: [
    vue(),
    tailwindcss(),
    vueDevTools(),
  ],
  server: {
    port: 3003,
    host: '0.0.0.0',
    allowedHosts:  [process.env.ALLOWED_HOST || 'localhost'],
    proxy: {
      '/api': {
        target: 'http://localhost:8080',
        changeOrigin: true,
        rewrite: (path) => path.replace(/^\/api/, '/api')
      }
    }
  },
  build: {
    lib: {
      entry: fileURLToPath(new URL('./src/index.ts', import.meta.url)),
      name: 'UnburdyBaseApp',
      fileName: (format) => `unburdy-base-app.${format}.js`,
      formats: ['es', 'cjs']
    },
    rollupOptions: {
      // Externalize deps that shouldn't be bundled
      external: [
        'vue',
        'vue-router',
        'pinia',
        '@vueuse/core',
        'vue-i18n',
        '@iconify/vue',
        'lucide-vue-next',
        'reka-ui',
        'class-variance-authority',
        'clsx',
        'tailwind-merge',
        'vee-validate',
        '@vee-validate/zod',
        'zod',
        'axios'
      ],
      output: {
        // Preserve CSS
        assetFileNames: (assetInfo) => {
          if (assetInfo.name === 'style.css') return 'unburdy-base-app.css'
          return assetInfo.name || 'asset'
        },
        globals: {
          vue: 'Vue',
          'vue-router': 'VueRouter',
          pinia: 'Pinia',
          '@vueuse/core': 'VueUse',
          'vue-i18n': 'VueI18n'
        }
      }
    },
    // Extract CSS
    cssCodeSplit: false,
    sourcemap: true,
    outDir: 'dist'
  },
  resolve: {
    alias: {
      '@': fileURLToPath(new URL('./src', import.meta.url)),
      '@@': path.resolve(__dirname, './src')
    },
  },
})
