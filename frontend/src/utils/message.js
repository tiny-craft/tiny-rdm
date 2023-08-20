import { createDiscreteApi } from 'naive-ui'
import { themeOverrides } from '@/utils/theme.js'

const { message } = createDiscreteApi(['message'], {
    configProviderProps: {
        themeOverrides,
    },
})

export function useMessage() {
    return {
        error: (content, option = null) => {
            return message.error(content, option)
        },
        info: (content, option = null) => {
            return message.info(content, option)
        },
        loading: (content, option = null) => {
            return message.loading(content, option)
        },
        success: (content, option = null) => {
            return message.success(content, option)
        },
        warning: (content, option = null) => {
            return message.warning(content, option)
        },
    }
}
