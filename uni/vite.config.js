import { defineConfig } from 'vite'
import uni from '@dcloudio/vite-plugin-uni'
import path from 'path'

const createAutoInjectPopupModalPlugin = () => {
  const skipMarker = 'no-popup-modal'
  const popupTag = '    <popup-modal client-type="uni" />'
  const selfClosingTags = new Set(['area', 'base', 'br', 'col', 'embed', 'hr', 'img', 'input', 'link', 'meta', 'param', 'source', 'track', 'wbr'])

  const findRootCloseIndex = (templateContent) => {
    const tagPattern = /<\/?([a-zA-Z][\w-]*)\b[^>]*>/g
    const stack = []
    let rootStarted = false
    let match

    while ((match = tagPattern.exec(templateContent))) {
      const token = match[0]
      const tag = match[1]
      const isClosing = token.startsWith('</')
      const isSelfClosing = token.endsWith('/>') || selfClosingTags.has(tag)

      if (!rootStarted) {
        if (isClosing || isSelfClosing) {
          continue
        }
        rootStarted = true
        stack.push(tag)
        continue
      }

      if (isClosing) {
        if (stack.length > 0) {
          stack.pop()
        }
        if (stack.length === 0) {
          return match.index
        }
        continue
      }

      if (!isSelfClosing) {
        stack.push(tag)
      }
    }

    return -1
  }

  return {
    name: 'auto-inject-popup-modal',
    enforce: 'pre',
    transform(code, id) {
      const filePath = id.split('?')[0].replace(/\\/g, '/')
      if (!filePath.endsWith('.vue') || !filePath.includes('/src/pages/')) {
        return null
      }
      if (code.includes('<popup-modal') || code.includes(skipMarker)) {
        return null
      }

      const templateStart = code.indexOf('<template')
      if (templateStart === -1) {
        return null
      }
      const templateOpenEnd = code.indexOf('>', templateStart)
      const templateEnd = code.indexOf('</template>', templateOpenEnd)
      if (templateOpenEnd === -1 || templateEnd === -1) {
        return null
      }

      const templateContent = code.slice(templateOpenEnd + 1, templateEnd)
      const rootCloseIndex = findRootCloseIndex(templateContent)
      if (rootCloseIndex === -1) {
        return null
      }

      const injectedTemplate =
        templateContent.slice(0, rootCloseIndex) +
        `\n${popupTag}\n` +
        templateContent.slice(rootCloseIndex)

      return {
        code:
          code.slice(0, templateOpenEnd + 1) +
          injectedTemplate +
          code.slice(templateEnd),
        map: null,
      }
    },
  }
}
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
    createAutoInjectPopupModalPlugin(),
    uni(),
  ],
})