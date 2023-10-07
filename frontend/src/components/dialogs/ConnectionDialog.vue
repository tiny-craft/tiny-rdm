<script setup>
import { every, get, includes, isEmpty, map, sortBy, toNumber } from 'lodash'
import { computed, nextTick, ref, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import { SelectKeyFile, TestConnection } from 'wailsjs/go/services/connectionService.js'
import useDialog, { ConnDialogType } from 'stores/dialog'
import Close from '@/components/icons/Close.vue'
import useConnectionStore from 'stores/connections.js'

/**
 * Dialog for new or edit connection
 */

const dialogStore = useDialog()
const connectionStore = useConnectionStore()
const i18n = useI18n()

const editName = ref('')
const generalForm = ref(null)
const generalFormRules = () => {
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
        addr: { required: true, message: requiredMsg, trigger: 'input' },
        defaultFilter: { required: true, message: requiredMsg, trigger: 'input' },
        keySeparator: { required: true, message: requiredMsg, trigger: 'input' },
    }
}
const isEditMode = computed(() => dialogStore.connType === ConnDialogType.EDIT)
const closingConnection = computed(() => {
    if (isEmpty(editName.value)) {
        return false
    }
    return connectionStore.isConnected(editName.value)
})

const groupOptions = computed(() => {
    const options = map(connectionStore.groups, (group) => ({
        label: group,
        value: group,
    }))
    options.splice(0, 0, {
        label: i18n.t('dialogue.connection.no_group'),
        value: '',
    })
    return options
})

const dbFilterList = ref([])
const onUpdateDBFilterList = (list) => {
    const dbList = []
    for (const item of list) {
        const idx = toNumber(item)
        if (!isNaN(idx)) {
            dbList.push(idx)
        }
    }
    generalForm.value.dbFilterList = sortBy(dbList)
}

const sshLoginType = computed(() => {
    return get(generalForm.value, 'ssh.loginType', 'pwd')
})

const onChoosePKFile = async () => {
    const { success, data } = await SelectKeyFile(i18n.t('dialogue.connection.ssh.pkfile_selection_title'))
    if (!success) {
        generalForm.value.ssh.pkFile = ''
    } else {
        generalForm.value.ssh.pkFile = get(data, 'path', '')
    }
}

const tab = ref('general')
const testing = ref(false)
const showTestResult = ref(false)
const testResult = ref('')
const predefineColors = ref(['', '#F75B52', '#F7A234', '#F7CE33', '#4ECF60', '#348CF7', '#B270D3'])
const generalFormRef = ref(null)
const advanceFormRef = ref(null)

const onSaveConnection = async () => {
    // validate general form
    await generalFormRef.value?.validate((err) => {
        if (err) {
            nextTick(() => (tab.value = 'general'))
        }
    })

    // validate advance form
    await advanceFormRef.value?.validate((err) => {
        if (err) {
            nextTick(() => (tab.value = 'advanced'))
        }
    })

    // trim ssh login data
    if (generalForm.value.ssh.enable) {
        switch (generalForm.value.ssh.loginType) {
            case 'pkfile':
                generalForm.value.ssh.password = ''
                break
            default:
                generalForm.value.ssh.pkFile = ''
                generalForm.value.ssh.passphrase = ''
                break
        }
    } else {
        // ssh disabled, reset to default value
        const { ssh } = connectionStore.newDefaultConnection()
        generalForm.value.ssh = ssh
    }

    // store new connection
    const { success, msg } = await connectionStore.saveConnection(
        isEditMode.value ? editName.value : null,
        generalForm.value,
    )
    if (!success) {
        $message.error(msg)
        return
    }

    $message.success(i18n.t('dialogue.handle_succ'))
    onClose()
}

const resetForm = () => {
    generalForm.value = connectionStore.newDefaultConnection()
    generalFormRef.value?.restoreValidation()
    testing.value = false
    showTestResult.value = false
    testResult.value = ''
    tab.value = 'general'
}

watch(
    () => dialogStore.connDialogVisible,
    (visible) => {
        if (visible) {
            resetForm()
            editName.value = get(dialogStore.connParam, 'name', '')
            generalForm.value = dialogStore.connParam || connectionStore.newDefaultConnection()
            dbFilterList.value = map(generalForm.value.dbFilterList, (item) => item + '')
            generalForm.value.ssh.loginType = generalForm.value.ssh.loginType || 'pwd'
        }
    },
)

const onTestConnection = async () => {
    testResult.value = ''
    testing.value = true
    let result = ''
    try {
        const { success = false, msg } = await TestConnection(generalForm.value)
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
    dialogStore.closeConnDialog()
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
        :title="isEditMode ? $t('dialogue.connection.edit_title') : $t('dialogue.connection.new_title')"
        preset="dialog"
        transform-origin="center">
        <n-spin :show="closingConnection">
            <n-tabs v-model:value="tab" animated type="line">
                <!-- General pane -->
                <n-tab-pane :tab="$t('dialogue.connection.general')" display-directive="show" name="general">
                    <n-form
                        ref="generalFormRef"
                        :model="generalForm"
                        :rules="generalFormRules()"
                        :show-require-mark="false"
                        label-placement="top">
                        <n-form-item :label="$t('dialogue.connection.conn_name')" path="name" required>
                            <n-input
                                v-model:value="generalForm.name"
                                :placeholder="$t('dialogue.connection.name_tip')" />
                        </n-form-item>
                        <n-form-item v-if="!isEditMode" :label="$t('dialogue.connection.group')" required>
                            <n-select v-model:value="generalForm.group" :options="groupOptions" />
                        </n-form-item>
                        <n-form-item :label="$t('dialogue.connection.addr')" path="addr" required>
                            <n-input
                                v-model:value="generalForm.addr"
                                :placeholder="$t('dialogue.connection.addr_tip')" />
                            <n-text style="width: 40px; text-align: center">:</n-text>
                            <n-input-number
                                v-model:value="generalForm.port"
                                :max="65535"
                                :min="1"
                                style="width: 200px" />
                        </n-form-item>
                        <n-form-item :label="$t('dialogue.connection.pwd')" path="password">
                            <n-input
                                v-model:value="generalForm.password"
                                :placeholder="$t('dialogue.connection.pwd_tip')"
                                show-password-on="click"
                                type="password" />
                        </n-form-item>
                        <n-form-item :label="$t('dialogue.connection.usr')" path="username">
                            <n-input
                                v-model:value="generalForm.username"
                                :placeholder="$t('dialogue.connection.usr_tip')" />
                        </n-form-item>
                    </n-form>
                </n-tab-pane>

                <!-- Advance pane -->
                <n-tab-pane :tab="$t('dialogue.connection.advn.title')" display-directive="show" name="advanced">
                    <n-form
                        ref="advanceFormRef"
                        :model="generalForm"
                        :rules="generalFormRules()"
                        :show-require-mark="false"
                        label-placement="top">
                        <n-grid :x-gap="10">
                            <n-form-item-gi
                                :span="12"
                                :label="$t('dialogue.connection.advn.filter')"
                                path="defaultFilter">
                                <n-input
                                    v-model:value="generalForm.defaultFilter"
                                    :placeholder="$t('dialogue.connection.advn.filter_tip')" />
                            </n-form-item-gi>
                            <n-form-item-gi
                                :span="12"
                                :label="$t('dialogue.connection.advn.separator')"
                                path="keySeparator">
                                <n-input
                                    v-model:value="generalForm.keySeparator"
                                    :placeholder="$t('dialogue.connection.advn.separator_tip')" />
                            </n-form-item-gi>
                            <n-form-item-gi
                                :span="12"
                                :label="$t('dialogue.connection.advn.conn_timeout')"
                                path="connTimeout">
                                <n-input-number v-model:value="generalForm.connTimeout" :max="999999" :min="1">
                                    <template #suffix>
                                        {{ $t('common.second') }}
                                    </template>
                                </n-input-number>
                            </n-form-item-gi>
                            <n-form-item-gi
                                :span="12"
                                :label="$t('dialogue.connection.advn.exec_timeout')"
                                path="execTimeout">
                                <n-input-number v-model:value="generalForm.execTimeout" :max="999999" :min="1">
                                    <template #suffix>
                                        {{ $t('common.second') }}
                                    </template>
                                </n-input-number>
                            </n-form-item-gi>
                            <n-form-item-gi :span="24" :label="$t('dialogue.connection.advn.dbfilter_type')">
                                <n-radio-group v-model:value="generalForm.dbFilterType">
                                    <n-radio-button :label="$t('dialogue.connection.advn.dbfilter_all')" value="none" />
                                    <n-radio-button
                                        :label="$t('dialogue.connection.advn.dbfilter_show')"
                                        value="show" />
                                    <n-radio-button
                                        :label="$t('dialogue.connection.advn.dbfilter_hide')"
                                        value="hide" />
                                </n-radio-group>
                            </n-form-item-gi>
                            <n-form-item-gi :span="24" :label="$t('dialogue.connection.advn.dbfilter_input')">
                                <n-select
                                    v-model:value="dbFilterList"
                                    :disabled="generalForm.dbFilterType === 'none'"
                                    filterable
                                    multiple
                                    tag
                                    :placeholder="$t('dialogue.connection.advn.dbfilter_input_tip')"
                                    :show-arrow="false"
                                    :show="false"
                                    :clearable="true"
                                    @update:value="onUpdateDBFilterList" />
                            </n-form-item-gi>
                            <n-form-item-gi
                                :span="24"
                                :label="$t('dialogue.connection.advn.mark_color')"
                                path="markColor">
                                <div
                                    v-for="color in predefineColors"
                                    :key="color"
                                    :class="{
                                        'color-preset-item_selected': generalForm.markColor === color,
                                    }"
                                    :style="{ backgroundColor: color }"
                                    class="color-preset-item"
                                    @click="generalForm.markColor = color">
                                    <n-icon v-if="isEmpty(color)" :component="Close" size="24" />
                                </div>
                            </n-form-item-gi>
                        </n-grid>
                    </n-form>
                </n-tab-pane>

                <!-- SSH pane -->
                <n-tab-pane :tab="$t('dialogue.connection.ssh.title')" display-directive="show" name="ssh">
                    <n-form-item label-placement="left">
                        <n-checkbox v-model:checked="generalForm.ssh.enable" size="medium">
                            {{ $t('dialogue.connection.ssh.enable') }}
                        </n-checkbox>
                    </n-form-item>
                    <n-form
                        :model="generalForm.ssh"
                        :show-require-mark="false"
                        :disabled="!generalForm.ssh.enable"
                        label-placement="top">
                        <n-form-item :label="$t('dialogue.connection.addr')" required>
                            <n-input
                                v-model:value="generalForm.ssh.addr"
                                :placeholder="$t('dialogue.connection.ssh.addr_tip')" />
                            <n-text style="width: 40px; text-align: center">:</n-text>
                            <n-input-number
                                v-model:value="generalForm.ssh.port"
                                :max="65535"
                                :min="1"
                                style="width: 200px" />
                        </n-form-item>
                        <n-form-item :label="$t('dialogue.connection.ssh.login_type')">
                            <n-radio-group v-model:value="generalForm.ssh.loginType">
                                <n-radio-button :label="$t('dialogue.connection.pwd')" value="pwd" />
                                <n-radio-button :label="$t('dialogue.connection.ssh.pkfile')" value="pkfile" />
                            </n-radio-group>
                        </n-form-item>
                        <n-form-item
                            v-if="sshLoginType === 'pwd' || sshLoginType === 'pkfile'"
                            :label="$t('dialogue.connection.usr')">
                            <n-input
                                v-model:value="generalForm.ssh.username"
                                :placeholder="$t('dialogue.connection.ssh.usr_tip')" />
                        </n-form-item>
                        <n-form-item v-if="sshLoginType === 'pwd'" :label="$t('dialogue.connection.pwd')">
                            <n-input
                                v-model:value="generalForm.ssh.password"
                                :placeholder="$t('dialogue.connection.ssh.pwd_tip')"
                                show-password-on="click"
                                type="password" />
                        </n-form-item>
                        <n-form-item v-if="sshLoginType === 'pkfile'" :label="$t('dialogue.connection.ssh.pkfile')">
                            <n-input-group>
                                <n-input
                                    v-model:value="generalForm.ssh.pkFile"
                                    :placeholder="$t('dialogue.connection.ssh.pkfile_tip')" />
                                <n-button :focusable="false" @click="onChoosePKFile">...</n-button>
                            </n-input-group>
                        </n-form-item>
                        <n-form-item v-if="sshLoginType === 'pkfile'" :label="$t('dialogue.connection.ssh.passphrase')">
                            <n-input
                                v-model:value="generalForm.ssh.passphrase"
                                :placeholder="$t('dialogue.connection.ssh.passphrase_tip')"
                                show-password-on="click"
                                type="password" />
                        </n-form-item>
                    </n-form>
                </n-tab-pane>

                <!-- Sentinel pane -->
                <n-tab-pane :tab="$t('dialogue.connection.sentinel.title')" display-directive="show" name="sentinel">
                    <n-form-item label-placement="left">
                        <n-checkbox v-model:checked="generalForm.sentinel.enable" size="medium">
                            {{ $t('dialogue.connection.sentinel.enable') }}
                        </n-checkbox>
                    </n-form-item>
                    <n-form
                        :model="generalForm.sentinel"
                        :show-require-mark="false"
                        :disabled="!generalForm.sentinel.enable"
                        label-placement="top">
                        <n-form-item :label="$t('dialogue.connection.sentinel.master')">
                            <n-input
                                v-model:value="generalForm.sentinel.master"
                                :placeholder="$t('dialogue.connection.sentinel.master')" />
                        </n-form-item>
                        <n-form-item :label="$t('dialogue.connection.sentinel.password')">
                            <n-input
                                v-model:value="generalForm.sentinel.password"
                                :placeholder="$t('dialogue.connection.sentinel.pwd_tip')"
                                show-password-on="click"
                                type="password" />
                        </n-form-item>
                        <n-form-item :label="$t('dialogue.connection.sentinel.username')">
                            <n-input
                                v-model:value="generalForm.sentinel.username"
                                :placeholder="$t('dialogue.connection.sentinel.usr_tip')" />
                        </n-form-item>
                    </n-form>
                </n-tab-pane>

                <!-- TODO: SSL tab pane -->
                <!-- TODO: Sentinel tab pane -->
                <!-- TODO: Cluster tab pane -->
            </n-tabs>

            <!-- test result alert-->
            <n-alert
                v-if="showTestResult"
                :title="isEmpty(testResult) ? '' : $t('dialogue.connection.test_fail')"
                :type="isEmpty(testResult) ? 'success' : 'error'"
                closable
                :on-close="() => (showTestResult = false)">
                <template v-if="isEmpty(testResult)">{{ $t('dialogue.connection.test_succ') }}</template>
                <template v-else>{{ testResult }}</template>
            </n-alert>
        </n-spin>

        <template #action>
            <div class="flex-item-expand">
                <n-button :focusable="false" :disabled="closingConnection" :loading="testing" @click="onTestConnection">
                    {{ $t('dialogue.connection.test') }}
                </n-button>
            </div>
            <div class="flex-item n-dialog__action">
                <n-button :focusable="false" :disabled="closingConnection" @click="onClose">
                    {{ $t('common.cancel') }}
                </n-button>
                <n-button :focusable="false" :disabled="closingConnection" type="primary" @click="onSaveConnection">
                    {{ isEditMode ? $t('preferences.general.update') : $t('common.confirm') }}
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
    border-radius: 50%;

    &_selected,
    &:hover {
        border-color: #cdd0d6;
    }
}
</style>
