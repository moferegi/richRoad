import { defineConfig } from 'vite'
import uni from '@dcloudio/vite-plugin-uni'
import path from 'path'
// https://vitejs.dev/config/
export default defineConfig({
  resolve: {
    alias:{
      '@': path.resolve(__dirname, './src')
    }
  },
  server: {
    proxy: {
      // 把key的路径代理到target位置
      // detail: https://cli.vuejs.org/config/#devserver-proxy
      "/api": {
        // 需要代理的路径   例如 '/api'
        target: `http://localhost:8888/`, // 代理到 目标路径
        changeOrigin: true,
        rewrite: (path) =>
          path.replace(new RegExp('^/api'), '')
      }
    },
  },
  plugins: [
    uni(),
  ],
})