<script setup>
import { every, get, includes, isEmpty, map } from 'lodash'
import { computed, nextTick, ref, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import { TestConnection } from 'wailsjs/go/services/connectionService.js'
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

const tab = ref('general')
const testing = ref(false)
const showTestResult = ref(false)
const testResult = ref('')
const predefineColors = ref(['', '#F75B52', '#F7A234', '#F7CE33', '#4ECF60', '#348CF7', '#B270D3'])
const generalFormRef = ref(null)
const safeFormRef = ref(null)
const advanceFormRef = ref(null)
const fileRef = ref(null)

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

    // store new connection
    const { success, msg } = await connectionStore.saveConnection(editName.value, generalForm.value)
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
    showTestResult.value = false
    testResult.value = ''
    tab.value = 'general'
}

const choose_file = () => {
  //弹出选择本地文件
  fileRef.value.click()
}
const fileChange = (e) => {
  const file = e.target.files ? e.target.files[0] : null
  if (file) {
    const reader = new FileReader();
    reader.onload = (event) => {
      generalForm.value.sshKeyPath = event.target.result
    };
    reader.readAsText(e.target.files[0]);
  }
}

watch(
    () => dialogStore.connDialogVisible,
    (visible) => {
        if (visible) {
            editName.value = get(dialogStore.connParam, 'name', '')
            generalForm.value = dialogStore.connParam || connectionStore.newDefaultConnection()
        }
    },
)

const onTestConnection = async () => {
    testResult.value = ''
    testing.value = true
    let result = ''
    try {
        const opt = JSON.stringify(generalForm.value)
        const { success = false, msg } = await TestConnection(opt)
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
        style="width:580px"
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
                        <n-form-item :label="$t('dialogue.connection.pwd')" path="password" required>
                            <n-input
                                v-model:value="generalForm.password"
                                :placeholder="$t('dialogue.connection.pwd_tip')"
                                show-password-on="click"
                                type="password" />
                        </n-form-item>
                        <n-form-item :label="$t('dialogue.connection.usr')" path="username">
                            <n-input v-model="generalForm.username" :placeholder="$t('dialogue.connection.usr_tip')" />
                        </n-form-item>
                    </n-form>
                </n-tab-pane>

                <n-tab-pane :tab="$t('dialogue.connection.safe')" display-directive="show" name="safe">
                  <n-form
                      ref="safeFormRef"
                      :model="generalForm"
                      :rules="generalFormRules()"
                      :show-require-mark="false"
                      label-placement="top">
                    <n-form-item :label="$t('dialogue.connection.safe_link')" path="safeLink" required>
                      <n-radio-group v-model:value="generalForm.safeLink" name="safeLink">
                        <n-space>
                          <n-radio :value="1" name="no">
                            地址直连
                          </n-radio>
                          <n-radio :value="2" name="ssh">
                            SSH隧道
                          </n-radio>
                        </n-space>
                      </n-radio-group>
                    </n-form-item>
                    <n-collapse-transition :show="generalForm.safeLink === 2">
                      <n-form-item :label="$t('dialogue.connection.ssh_user')" path="ssh_user">
                        <n-input v-model:value="generalForm.sshUser" :placeholder="$t('dialogue.connection.ssh_user_tip')" />
                      </n-form-item>
                      <n-form-item :label="$t('dialogue.connection.addr')" path="ssh_addr" required>
                        <n-input
                            v-model:value="generalForm.sshAddr"
                            :placeholder="$t('dialogue.connection.ssh_addr_tip')" />
                        <n-text style="width: 40px; text-align: center">:</n-text>
                        <n-input-number
                            v-model:value="generalForm.sshPort"
                            :max="65535"
                            :min="1"
                            style="width: 200px" />
                      </n-form-item>
                      <n-form-item :label="$t('dialogue.connection.ssh_auth')" path="ssh_auth" required>
                        <n-radio-group v-model:value="generalForm.sshAuth" name="sshAuth">
                          <n-space>
                            <n-radio :value="1" name="pwd">
                              密码
                            </n-radio>
                            <n-radio :value="2" name="key_file">
                              秘钥
                            </n-radio>
                          </n-space>
                        </n-radio-group>
                      </n-form-item>
                      <n-collapse-transition :show="generalForm.sshAuth === 1">
                        <n-form-item :label="$t('dialogue.connection.pwd')" path="ssh_password">
                          <n-input
                              v-model:value="generalForm.sshPassword"
                              show-password-on="click"
                              type="password" />
                        </n-form-item>
                      </n-collapse-transition>
                      <n-collapse-transition :show="generalForm.sshAuth === 2">
                        <n-form-item :label="$t('dialogue.connection.ssh_key_path')" path="ssh_key_path">
                          <input ref="fileRef" v-show="false" type="file" @change="fileChange($event)" />
                          <n-button type="primary" ghost @click="choose_file">
                            {{ $t('dialogue.connection.choose_file') }}
                          </n-button>
                        </n-form-item>
                        <n-form-item style="margin-top: -30px">
                          <n-input
                              type="textarea"
                              size="small"
                              :placeholder="generalForm.sshKeyPath"
                              disabled
                              round
                              :rows="6"
                          />
                        </n-form-item>
                        <n-form-item :label="$t('dialogue.connection.ssh_key_pwd')" path="ssh_key_pwd">
                          <n-input
                              v-model:value="generalForm.sshKeyPwd"
                              show-password-on="click"
                              type="password"
                              :placeholder="$t('dialogue.connection.ssh_key_pwd_tip')"
                          />
                        </n-form-item>
                      </n-collapse-transition>
                    </n-collapse-transition>
                  </n-form>
                </n-tab-pane>

                <n-tab-pane :tab="$t('dialogue.connection.advanced')" display-directive="show" name="advanced">
                    <n-form
                        ref="advanceFormRef"
                        :model="generalForm"
                        :rules="generalFormRules()"
                        :show-require-mark="false"
                        label-placement="top">
                        <n-form-item :label="$t('dialogue.connection.advn_filter')" path="defaultFilter">
                            <n-input
                                v-model:value="generalForm.defaultFilter"
                                :placeholder="$t('dialogue.connection.advn_filter_tip')" />
                        </n-form-item>
                        <n-form-item :label="$t('dialogue.connection.advn_separator')" path="keySeparator">
                            <n-input
                                v-model:value="generalForm.keySeparator"
                                :placeholder="$t('dialogue.connection.advn_separator_tip')" />
                        </n-form-item>
                        <n-form-item :label="$t('dialogue.connection.advn_conn_timeout')" path="connTimeout">
                            <n-input-number v-model:value="generalForm.connTimeout" :max="999999" :min="1">
                                <template #suffix>
                                    {{ $t('common.second') }}
                                </template>
                            </n-input-number>
                        </n-form-item>
                        <n-form-item :label="$t('dialogue.connection.advn_exec_timeout')" path="execTimeout">
                            <n-input-number v-model:value="generalForm.execTimeout" :max="999999" :min="1">
                                <template #suffix>
                                    {{ $t('common.second') }}
                                </template>
                            </n-input-number>
                        </n-form-item>
                        <n-form-item :label="$t('dialogue.connection.advn_mark_color')" path="markColor">
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
                        </n-form-item>
                    </n-form>
                </n-tab-pane>
            </n-tabs>

            <!-- test result alert-->
            <n-alert
                v-if="showTestResult"
                :title="isEmpty(testResult) ? '' : $t('dialogue.connection.test_fail')"
                :type="isEmpty(testResult) ? 'success' : 'error'">
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
