import { createDiscreteApi } from 'naive-ui'

const { message } = createDiscreteApi(['message'])

export function useMessage() {
    return {
        error: (content, option = null) => {
            message.error(content, option)
        },
        info: (content, option = null) => {
            message.info(content, option)
        },
        loading: (content, option = null) => {
            message.loading(content, option)
        },
        success: (content, option = null) => {
            message.success(content, option)
        },
        warning: (content, option = null) => {
            message.warning(content, option)
        },
    }
}
