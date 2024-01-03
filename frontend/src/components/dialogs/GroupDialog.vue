<script setup>
import { computed, reactive, ref, watchEffect } from 'vue'
import useDialog from 'stores/dialog'
import { useI18n } from 'vue-i18n'
import useConnectionStore from 'stores/connections.js'
import { every, get, includes, isEmpty } from 'lodash'

/**
 * Dialog for create or rename group
 */

const editGroup = ref('')
const groupForm = reactive({
    name: '',
})
const groupFormRef = ref(null)

const formRules = computed(() => {
    const requiredMsg = i18n.t('dialogue.field_required')
    const illegalChars = ['/', '\\']
    return {
        name: [
            { required: true, message: requiredMsg, trigger: 'input' },
            {
                validator: (rule, value) => {
                    return every(illegalChars, (c) => !includes(value, c))
                },
                message: i18n.t('dialogue.illegal_characters'),
                trigger: 'input',
            },
        ],
    }
})

const isRenameMode = computed(() => !isEmpty(editGroup.value))

const dialogStore = useDialog()
const connectionStore = useConnectionStore()
watchEffect(() => {
    if (dialogStore.groupDialogVisible) {
        groupForm.name = editGroup.value = dialogStore.editGroup
    }
})

const i18n = useI18n()
const onConfirm = async () => {
    try {
        await groupFormRef.value?.validate((errs) => {
            const err = get(errs, '0.0.message')
            if (err != null) {
                $message.error(err)
            }
        })

        const { name } = groupForm
        if (isRenameMode.value) {
            const { success, msg } = await connectionStore.renameGroup(editGroup.value, name)
            if (success) {
                $message.success(i18n.t('dialogue.handle_succ'))
            } else {
                $message.error(msg)
            }
        } else {
            const { success, msg } = await connectionStore.createGroup(name)
            if (success) {
                $message.success(i18n.t('dialogue.handle_succ'))
            } else {
                $message.error(msg)
            }
        }
    } catch (e) {
        const msg = get(e, 'message')
        if (!isEmpty(msg)) {
            $message.error(msg)
        }
    }
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
        :negative-text="$t('common.cancel')"
        :positive-button-props="{ size: 'medium' }"
        :positive-text="$t('common.confirm')"
        :show-icon="false"
        :title="isRenameMode ? $t('dialogue.group.rename') : $t('dialogue.group.new')"
        preset="dialog"
        transform-origin="center"
        @positive-click="onConfirm"
        @negative-click="onClose">
        <n-form
            ref="groupFormRef"
            :model="groupForm"
            :rules="formRules"
            :show-label="false"
            :show-require-mark="false"
            label-placement="top">
            <n-form-item :label="$t('dialogue.group.name')" path="name" required>
                <n-input v-model:value="groupForm.name" placeholder="" />
            </n-form-item>
        </n-form>
    </n-modal>
</template>

<style lang="scss" scoped></style>
