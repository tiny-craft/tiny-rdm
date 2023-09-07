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
        loading: (content, option = {}) => {
            option.duration = option.duration != null ? option.duration : 30000
            option.keepAliveOnHover = option.keepAliveOnHover !== undefined ? option.keepAliveOnHover : true
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

function setupNotification(notification) {
    return {
        error: (content, option = {}) => {
            option.content = content
            option.title = option.title || i18nGlobal.t('common.error')
            return notification.error(option)
        },
        info: (content, option = {}) => {
            option.content = content
            return notification.info(option)
        },
        success: (content, option = {}) => {
            option.content = content
            option.title = option.title || i18nGlobal.t('common.success')
            return notification.success(option)
        },
        warning: (content, option = {}) => {
            option.content = content
            option.title = option.title || i18nGlobal.t('common.warning')
            return notification.warning(option)
        },
    }
}

function setupDialog(dialog) {
    return {
        warning: (content, onConfirm) => {
            return dialog.warning({
                title: i18nGlobal.t('common.warning'),
                content: content,
                closable: false,
                autoFocus: false,
                transformOrigin: 'center',
                positiveText: i18nGlobal.t('common.confirm'),
                negativeText: i18nGlobal.t('common.cancel'),
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
    const { message, dialog, notification } = createDiscreteApi(['message', 'notification', 'dialog'], {
        configProviderProps,
        messageProviderProps: {
            placement: 'bottom-right',
        },
        notificationProviderProps: {
            max: 5,
            placement: 'top-right',
            keepAliveOnHover: true,
        },
    })

    window.$message = setupMessage(message)
    window.$notification = setupNotification(notification)
    window.$dialog = setupDialog(dialog)
}
