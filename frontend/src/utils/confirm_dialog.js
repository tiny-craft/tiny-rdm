import { createDiscreteApi } from 'naive-ui'
import { i18nGlobal } from '@/utils/i18n.js'

const { dialog } = createDiscreteApi(['dialog'])

export function useConfirmDialog() {
    return {
        warning: (content, onConfirm) => {
            dialog.warning({
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
