<script setup>
import { computed, reactive, ref, watch } from 'vue'
import useDialog from 'stores/dialog'
import { useI18n } from 'vue-i18n'
import useConnectionStore from 'stores/connections.js'
import { isEmpty } from 'lodash'

/**
 * Dialog for create or rename group
 */

const editGroup = ref('')
const groupForm = reactive({
    name: '',
})

const isRenameMode = computed(() => !isEmpty(editGroup.value))

const dialogStore = useDialog()
const connectionStore = useConnectionStore()
watch(
    () => dialogStore.groupDialogVisible,
    (visible) => {
        if (visible) {
            groupForm.name = editGroup.value = dialogStore.editGroup
        }
    },
)

const i18n = useI18n()
const onConfirm = async () => {
    try {
        const { name } = groupForm
        if (isRenameMode.value) {
            const { success, msg } = await connectionStore.renameGroup(editGroup.value, name)
            if (success) {
                $message.success(i18n.t('handle_succ'))
            } else {
                $message.error(msg)
            }
        } else {
            const { success, msg } = await connectionStore.createGroup(name)
            if (success) {
                $message.success(i18n.t('handle_succ'))
            } else {
                $message.error(msg)
            }
        }
    } catch (e) {
        $message.error(e.message)
    }
    onClose()
}

const onClose = () => {
    if (isRenameMode.value) {
        dialogStore.closeNewGroupDialog()
    } else {
        dialogStore.closeRenameGroupDialog()
    }
}
</script>

<template>
    <n-modal
        v-model:show="dialogStore.groupDialogVisible"
        :closable="false"
        :close-on-esc="false"
        :mask-closable="false"
        :negative-button-props="{ size: 'medium' }"
        :negative-text="$t('cancel')"
        :positive-button-props="{ size: 'medium' }"
        :positive-text="$t('confirm')"
        :show-icon="false"
        :title="isRenameMode ? $t('rename_group') : $t('new_group')"
        preset="dialog"
        transform-origin="center"
        @positive-click="onConfirm"
        @negative-click="onClose">
        <n-form :model="groupForm" :show-label="false" :show-require-mark="false" label-placement="top">
            <n-form-item :label="$t('group_name')" required>
                <n-input v-model:value="groupForm.name" placeholder="" />
            </n-form-item>
        </n-form>
    </n-modal>
</template>

<style lang="scss" scoped></style>
