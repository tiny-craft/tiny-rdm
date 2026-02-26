import vue from '@vitejs/plugin-vue'
import AutoImport from 'unplugin-auto-import/vite'
import Icons from 'unplugin-icons/vite'
import { NaiveUiResolver } from 'unplugin-vue-components/resolvers'
import Components from 'unplugin-vue-components/vite'
import { defineConfig } from 'vite'

const rootPath = new URL('.', import.meta.url).pathname
const isWeb = process.env.VITE_WEB === 'true'

// https://vitejs.dev/config/
export default defineConfig({
    plugins: [
        vue(),
        AutoImport({
            imports: [
                {
                    'naive-ui': ['useDialog', 'useMessage', 'useNotification', 'useLoadingBar'],
                },
            ],
        }),
        Components({
            resolvers: [NaiveUiResolver()],
        }),
        Icons(),
    ],
    resolve: {
        alias: {
            '@': rootPath + 'src',
            stores: rootPath + 'src/stores',
            // Web mode: redirect wailsjs imports to HTTP/WebSocket adapters
            // Desktop mode (wails build): use real Wails RPC bindings
            ...(isWeb
                ? {
                      'wailsjs/runtime/runtime.js': rootPath + 'src/utils/wails_runtime.js',
                      'wailsjs/go/services/connectionService.js': rootPath + 'src/utils/api.js',
                      'wailsjs/go/services/browserService.js': rootPath + 'src/utils/api.js',
                      'wailsjs/go/services/cliService.js': rootPath + 'src/utils/api.js',
                      'wailsjs/go/services/monitorService.js': rootPath + 'src/utils/api.js',
                      'wailsjs/go/services/pubsubService.js': rootPath + 'src/utils/api.js',
                      'wailsjs/go/services/preferencesService.js': rootPath + 'src/utils/api.js',
                      'wailsjs/go/services/systemService.js': rootPath + 'src/utils/api.js',
                  }
                : {}),
            wailsjs: rootPath + 'wailsjs',
        },
    },
    css: {
        preprocessorOptions: {
            scss: {
                api: 'modern-compiler',
            },
        },
    },
    ...(isWeb
        ? {
              server: {
                  proxy: {
                      '/api': {
                          target: 'http://localhost:8088',
                          changeOrigin: true,
                      },
                      '/ws': {
                          target: 'ws://localhost:8088',
                          ws: true,
                      },
                  },
              },
          }
        : {}),
})
