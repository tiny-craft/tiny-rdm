import usePreferencesStore from 'stores/preferences.js'
import { createDiscreteApi, darkTheme } from 'naive-ui'
import { themeOverrides } from '@/utils/theme.js'
import { i18nGlobal } from '@/utils/i18n.js'
import { computed } from 'vue'

function setupMessage(message) {
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

function setupDialog(dialog) {
    return {
        warning: (content, onConfirm) => {
            return dialog.warning({
                title: i18nGlobal.t('warning'),
                content: content,
                closable: false,
                autoFocus: false,
                transformOrigin: 'center',
                positiveText: i18nGlobal.t('confirm'),
                negativeText: i18nGlobal.t('cancel'),
                onPositiveClick: () => {
                    onConfirm && onConfirm()
                },
            })
        },
    }
}

/**
 * setup discrete api and bind global component (like dialog, message, alert) to window
 * @return {Promise<void>}
 */
export async function setupDiscreteApi() {
    const prefStore = usePreferencesStore()
    const configProviderProps = computed(() => ({
        theme: prefStore.isDark ? darkTheme : undefined,
        themeOverrides,
    }))
    const { message, dialog } = createDiscreteApi(['message', 'dialog'], { configProviderProps })

    window.$message = setupMessage(message)
    window.$dialog = setupDialog(dialog)
}
