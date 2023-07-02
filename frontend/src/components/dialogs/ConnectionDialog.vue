<script setup>
import { get, isEmpty, map } from 'lodash'
import { computed, nextTick, ref, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import { TestConnection } from '../../../wailsjs/go/services/connectionService.js'
import useDialog from '../../stores/dialog'
import { useMessage } from 'naive-ui'
import Close from '../icons/Close.vue'
import useConnectionStore from '../../stores/connections.js'

/**
 * Dialog for create or edit connection
 */

const dialogStore = useDialog()
const connectionStore = useConnectionStore()
const message = useMessage()
const i18n = useI18n()

const editName = ref('')
const generalForm = ref(null)
const generalFormRules = () => {
    const requiredMsg = i18n.t('field_required')
    return {
        name: { required: true, message: requiredMsg, trigger: 'input' },
        addr: { required: true, message: requiredMsg, trigger: 'input' },
        defaultFilter: { required: true, message: requiredMsg, trigger: 'input' },
        keySeparator: { required: true, message: requiredMsg, trigger: 'input' },
    }
}
const isEditMode = computed(() => !isEmpty(editName.value))

const groupOptions = computed(() => {
    const options = map(connectionStore.groups, (group) => ({
        label: group,
        value: group,
    }))
    options.splice(0, 0, {
        label: i18n.t('no_group'),
        value: '',
    })
    return options
})

const tab = ref('general')
const testing = ref(false)
const showTestResult = ref(false)
const testResult = ref('')
const showTestConnSuccResult = computed(() => {
    return isEmpty(testResult.value) && showTestResult.value === true
})
const showTestConnFailResult = computed(() => {
    return !isEmpty(testResult.value) && showTestResult.value === true
})
const formLabelWidth = computed(() => {
    // Compatible with long english word
    if (tab.value === 'advanced' && i18n.locale.value === 'en') {
        return '140px'
    }
    return '80px'
})
const predefineColors = ref(['', '#FE5959', '#FEC230', '#FEF27F', '#6CEFAF', '#46C3FC', '#B388FC', '#B0BEC5'])
const generalFormRef = ref(null)
const advanceFormRef = ref(null)

const onSaveConnection = async () => {
    // Validate general form
    await generalFormRef.value?.validate((err) => {
        if (err) {
            nextTick(() => (tab.value = 'general'))
        }
    })

    // Validate advance form
    await advanceFormRef.value?.validate((err) => {
        if (err) {
            nextTick(() => (tab.value = 'advanced'))
        }
    })

    // Store new connection
    const { success, msg } = await connectionStore.saveConnection(editName.value, generalForm.value)
    if (!success) {
        message.error(msg)
        return
    }

    message.success(i18n.t('handle_succ'))
    onClose()
}

const resetForm = () => {
    generalForm.value = connectionStore.newDefaultConnection()
    generalFormRef.value?.restoreValidation()
    showTestResult.value = false
    testResult.value = ''
    tab.value = 'general'
}

watch(
    () => dialogStore.connDialogVisible,
    (visible) => {
        if (visible) {
            editName.value = get(dialogStore.connParam, 'name', '')
            generalForm.value = dialogStore.connParam || connectionStore.newDefaultConnection()
        }
    }
)

const onTestConnection = async () => {
    testResult.value = ''
    testing.value = true
    let result = ''
    try {
        const { addr, port, username, password } = generalForm.value
        const { success = false, msg } = await TestConnection(addr, port, username, password)
        if (!success) {
            result = msg
        }
    } catch (e) {
        result = e.message
    } finally {
        testing.value = false
        showTestResult.value = true
    }

    if (!isEmpty(result)) {
        testResult.value = result
    } else {
        testResult.value = ''
    }
}

const onClose = () => {
    if (isEditMode.value) {
        dialogStore.closeEditDialog()
    } else {
        dialogStore.closeNewDialog()
    }
}
</script>

<template>
    <n-modal
        v-model:show="dialogStore.connDialogVisible"
        :closable="false"
        :close-on-esc="false"
        :mask-closable="false"
        :on-after-leave="resetForm"
        :show-icon="false"
        :title="isEditMode ? $t('edit_conn_title') : $t('new_conn_title')"
        preset="dialog"
        transform-origin="center"
    >
        <n-tabs v-model:value="tab">
            <n-tab-pane :tab="$t('general')" display-directive="show" name="general">
                <n-form
                    ref="generalFormRef"
                    :label-width="formLabelWidth"
                    :model="generalForm"
                    :rules="generalFormRules()"
                    :show-require-mark="false"
                    label-align="right"
                    label-placement="left"
                >
                    <n-form-item :label="$t('conn_name')" path="name" required>
                        <n-input v-model:value="generalForm.name" :placeholder="$t('conn_name_tip')" />
                    </n-form-item>
                    <n-form-item v-if="!isEditMode" :label="$t('conn_group')" required>
                        <n-select v-model:value="generalForm.group" :options="groupOptions"></n-select>
                    </n-form-item>
                    <n-form-item :label="$t('conn_addr')" path="addr" required>
                        <n-input v-model:value="generalForm.addr" :placeholder="$t('conn_addr_tip')" />
                        <n-text style="width: 40px; text-align: center">:</n-text>
                        <n-input-number v-model:value="generalForm.port" :max="65535" :min="1" style="width: 200px" />
                    </n-form-item>
                    <n-form-item :label="$t('conn_pwd')" path="password">
                        <n-input
                            v-model:value="generalForm.password"
                            :placeholder="$t('conn_pwd_tip')"
                            show-password-on="click"
                            type="password"
                        />
                    </n-form-item>
                    <n-form-item :label="$t('conn_usr')" path="username">
                        <n-input v-model="generalForm.username" :placeholder="$t('conn_usr_tip')" />
                    </n-form-item>
                </n-form>
            </n-tab-pane>

            <n-tab-pane :tab="$t('advanced')" display-directive="show" name="advanced">
                <n-form
                    ref="advanceFormRef"
                    :label-width="formLabelWidth"
                    :model="generalForm"
                    :rules="generalFormRules()"
                    :show-require-mark="false"
                    label-align="right"
                    label-placement="left"
                >
                    <n-form-item :label="$t('conn_advn_filter')" path="defaultFilter">
                        <n-input v-model:value="generalForm.defaultFilter" :placeholder="$t('conn_advn_filter_tip')" />
                    </n-form-item>
                    <n-form-item :label="$t('conn_advn_separator')" path="keySeparator">
                        <n-input
                            v-model:value="generalForm.keySeparator"
                            :placeholder="$t('conn_advn_separator_tip')"
                        />
                    </n-form-item>
                    <n-form-item :label="$t('conn_advn_conn_timeout')" path="connTimeout">
                        <n-input-number v-model:value="generalForm.connTimeout" :max="999999" :min="1">
                            <template #suffix>
                                {{ $t('second') }}
                            </template>
                        </n-input-number>
                    </n-form-item>
                    <n-form-item :label="$t('conn_advn_exec_timeout')" path="execTimeout">
                        <n-input-number v-model:value="generalForm.execTimeout" :max="999999" :min="1">
                            <template #suffix>
                                {{ $t('second') }}
                            </template>
                        </n-input-number>
                    </n-form-item>
                    <n-form-item :label="$t('conn_advn_mark_color')" path="markColor">
                        <div
                            v-for="color in predefineColors"
                            :key="color"
                            :class="{
                                'color-preset-item_selected': generalForm.markColor === color,
                            }"
                            :style="{ backgroundColor: color }"
                            class="color-preset-item"
                            @click="generalForm.markColor = color"
                        >
                            <n-icon v-if="color === ''" :component="Close" size="24" />
                        </div>
                    </n-form-item>
                </n-form>
            </n-tab-pane>
        </n-tabs>

        <!-- test result alert-->
        <n-alert v-if="showTestConnSuccResult" title="" type="success">
            {{ $t('conn_test_succ') }}
        </n-alert>
        <n-alert v-if="showTestConnFailResult" title="" type="error">
            {{ $t('conn_test_fail') }}: {{ testResult }}
        </n-alert>

        <template #action>
            <div class="flex-item-expand">
                <n-button :loading="testing" @click="onTestConnection">{{ $t('conn_test') }}</n-button>
            </div>
            <div class="flex-item n-dialog__action">
                <n-button @click="onClose">{{ $t('cancel') }}</n-button>
                <n-button type="primary" @click="onSaveConnection">
                    {{ isEditMode ? $t('update') : $t('confirm') }}
                </n-button>
            </div>
        </template>
    </n-modal>
</template>

<style lang="scss" scoped>
.color-preset-item {
    width: 24px;
    height: 24px;
    margin-right: 2px;
    border: white 3px solid;
    cursor: pointer;

    &_selected,
    &:hover {
        border-color: #cdd0d6;
    }
}
</style>
