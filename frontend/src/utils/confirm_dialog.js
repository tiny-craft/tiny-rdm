import { createDiscreteApi } from 'naive-ui'
import { useI18n } from 'vue-i18n'
import { onMounted } from 'vue'

const { dialog } = createDiscreteApi(['dialog'])

export function useConfirmDialog() {
    let i18n
    onMounted(() => {
        i18n = useI18n()
    })
    return {
        warning: (content, onConfirm) => {
            dialog.warning({
                title: i18n.t('warning'),
                content: content,
                closable: false,
                autoFocus: false,
                transformOrigin: 'center',
                positiveText: i18n.t('confirm'),
                negativeText: i18n.t('cancel'),
                onPositiveClick: () => {
                    onConfirm && onConfirm()
                },
            })
        },
    }
}
