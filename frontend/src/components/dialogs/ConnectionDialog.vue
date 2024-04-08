<script setup>
import { every, get, includes, isEmpty, map, reject, sortBy, toNumber } from 'lodash'
import { computed, nextTick, ref, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import { ListSentinelMasters, TestConnection } from 'wailsjs/go/services/connectionService.js'
import useDialog, { ConnDialogType } from 'stores/dialog'
import Close from '@/components/icons/Close.vue'
import useConnectionStore from 'stores/connections.js'
import FileOpenInput from '@/components/common/FileOpenInput.vue'
import { KeyViewType } from '@/consts/key_view_type.js'
import { useThemeVars } from 'naive-ui'
import useBrowserStore from 'stores/browser.js'
import Delete from '@/components/icons/Delete.vue'
import Add from '@/components/icons/Add.vue'
import IconButton from '@/components/common/IconButton.vue'

/**
 * Dialog for new or edit connection
 */

const themeVars = useThemeVars()
const dialogStore = useDialog()
const connectionStore = useConnectionStore()
const browserStore = useBrowserStore()
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
    return browserStore.isConnected(editName.value)
})

const groupOptions = computed(() => {
    const options = map(connectionStore.groups, (group) => ({
        label: group,
        value: group,
    }))
    options.splice(0, 0, {
        label: 'dialogue.connection.no_group',
        value: '',
    })
    return options
})

const dbFilterList = ref([])
const onUpdateDBFilterType = (t) => {
    if (t !== 'none') {
        // set default filter index if empty
        if (isEmpty(dbFilterList.value)) {
            dbFilterList.value = ['0']
        }
    }
}

const aliasPair = ref([
    /*{ db: 0, alias: '' }*/
])
const onCreateAlias = () => {
    return {
        db: 0,
        alias: '',
    }
}
const onUpdateAlias = () => {
    const val = reject(aliasPair.value, (v) => v == null || isEmpty(v.alias))
    const result = {}
    for (const elem of val) {
        result[elem.db] = elem.alias
    }
    generalForm.value.alias = result
}

watch(
    () => dbFilterList.value,
    (list) => {
        const dbList = map(list, (item) => {
            const idx = toNumber(item)
            return isNaN(idx) ? 0 : idx
        })
        generalForm.value.dbFilterList = sortBy(dbList)
    },
    { deep: true },
)

const sshLoginType = computed(() => {
    return get(generalForm.value, 'ssh.loginType', 'pwd')
})

const loadingSentinelMaster = ref(false)
const masterNameOptions = ref([])
const onLoadSentinelMasters = async () => {
    try {
        loadingSentinelMaster.value = true
        const { success, data, msg } = await ListSentinelMasters(generalForm.value)
        if (!success || isEmpty(data)) {
            $message.error(msg || 'list sentinel master fail')
        } else {
            const options = []
            for (const m of data) {
                options.push({
                    label: m['name'],
                    value: m['name'],
                })
            }

            // select default names
            if (!isEmpty(options)) {
                generalForm.value.sentinel.master = options[0].value
            }
            masterNameOptions.value = options
        }
    } catch (e) {
        $message.error(e.message)
    } finally {
        loadingSentinelMaster.value = false
    }
}

const tab = ref('general')
const testing = ref(false)
const testResult = ref(null)
const showTestResult = computed(() => {
    return !testing.value && testResult.value != null
})
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

    // trim addr by network type
    if (get(generalForm.value, 'network', 'tcp') === 'unix') {
        generalForm.value.network = 'unix'
        generalForm.value.addr = ''
        generalForm.value.port = 0
    } else {
        generalForm.value.network = 'tcp'
        generalForm.value.sock = ''
    }

    // trim advance data
    if (get(generalForm.value, 'dbFilterType', 'none') === 'none') {
        generalForm.value.dbFilterList = []
    }

    // trim ssl data
    if (!!!generalForm.value.ssl.enable) {
        generalForm.value.ssl = {}
    }

    // trim ssh login data
    if (!!generalForm.value.ssh.enable) {
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
        generalForm.value.ssh = {}
    }

    // trim sentinel data
    if (!!!generalForm.value.sentinel.enable) {
        generalForm.value.sentinel = {}
    }

    // trim cluster data
    if (!!!generalForm.value.cluster.enable) {
        generalForm.value.cluster = {}
    }

    // trim proxy data
    if (generalForm.value.proxy.type !== 2) {
        generalForm.value.proxy.schema = ''
        generalForm.value.proxy.addr = ''
        generalForm.value.proxy.port = 0
        generalForm.value.proxy.auth = false
        generalForm.value.proxy.username = ''
        generalForm.value.proxy.password = ''
    } else if (!generalForm.value.proxy.auth) {
        generalForm.value.proxy.username = ''
        generalForm.value.proxy.password = ''
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
    testResult.value = null
    tab.value = 'general'
    loadingSentinelMaster.value = false
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
            // update alias display
            const alias = get(generalForm.value, 'alias', {})
            const pairs = []
            for (const db in alias) {
                pairs.push({ db: parseInt(db), alias: alias[db] })
            }
            aliasPair.value = pairs
            generalForm.value.proxy.auth = !isEmpty(generalForm.value.proxy.username)
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

const pasteFromClipboard = async () => {
    // url example:
    // rediss://user:password@localhost:6789/3?dial_timeout=3&db=1&read_timeout=6s&max_retries=2
    let opt = {}
    try {
        opt = await connectionStore.parseUrlFromClipboard()
    } catch (e) {
        $message.error(i18n.t('dialogue.connection.parse_fail', { reason: e.message }))
        return
    }
    generalForm.value.network = opt.network || 'tcp'
    generalForm.value.name = generalForm.value.addr = opt.addr
    generalForm.value.port = opt.port
    generalForm.value.username = opt.username
    generalForm.value.password = opt.password
    if (opt.connTimeout > 0) {
        generalForm.value.connTimeout = opt.connTimeout
    }
    if (opt.execTimeout > 0) {
        generalForm.value.execTimeout = opt.execTimeout
    }
    const { sslServerName = null } = opt
    if (sslServerName != null) {
        generalForm.value.ssl.enable = true
        if (!isEmpty(sslServerName)) {
            generalForm.value.ssl.sni = sslServerName
        }
    }
    $message.success(i18n.t('dialogue.connection.parse_pass', { url: opt.url }))
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
        style="width: 600px"
        transform-origin="center">
        <n-spin :show="closingConnection">
            <n-tabs
                v-model:value="tab"
                animated
                pane-style="min-height: 50vh;"
                placement="left"
                tab-style="justify-content: right; font-weight: 420;"
                type="line">
                <!-- General pane -->
                <n-tab-pane :tab="$t('dialogue.connection.general')" display-directive="show:lazy" name="general">
                    <n-form
                        ref="generalFormRef"
                        :model="generalForm"
                        :rules="generalFormRules()"
                        :show-require-mark="false"
                        label-placement="top">
                        <n-grid :x-gap="10">
                            <n-form-item-gi
                                :label="$t('dialogue.connection.conn_name')"
                                :span="24"
                                path="name"
                                required>
                                <n-input
                                    v-model:value="generalForm.name"
                                    :placeholder="$t('dialogue.connection.name_tip')" />
                            </n-form-item-gi>
                            <n-form-item-gi
                                v-if="!isEditMode"
                                :label="$t('dialogue.connection.group')"
                                :span="24"
                                required>
                                <n-select
                                    v-model:value="generalForm.group"
                                    :options="groupOptions"
                                    :render-label="({ label, value }) => (value === '' ? $t(label) : label)" />
                            </n-form-item-gi>
                            <n-form-item-gi :label="$t('dialogue.connection.addr')" :span="24" path="addr" required>
                                <n-input-group>
                                    <n-select
                                        v-model:value="generalForm.network"
                                        :options="[
                                            { value: 'tcp', label: 'TCP' },
                                            { value: 'unix', label: 'UNIX' },
                                        ]"
                                        style="max-width: 100px" />
                                    <template v-if="generalForm.network === 'unix'">
                                        <n-input
                                            v-model:value="generalForm.sock"
                                            :placeholder="$t('dialogue.connection.sock_tip')" />
                                    </template>
                                    <template v-else>
                                        <n-input
                                            v-model:value="generalForm.addr"
                                            :placeholder="$t('dialogue.connection.addr_tip')" />
                                        <n-text style="width: 40px; text-align: center">:</n-text>
                                        <n-input-number
                                            v-model:value="generalForm.port"
                                            :max="65535"
                                            :min="1"
                                            :show-button="false"
                                            placeholder="6379"
                                            style="width: 200px" />
                                    </template>
                                </n-input-group>
                            </n-form-item-gi>
                            <n-form-item-gi :label="$t('dialogue.connection.pwd')" :span="12" path="password">
                                <n-input
                                    v-model:value="generalForm.password"
                                    :placeholder="$t('dialogue.connection.pwd_tip')"
                                    show-password-on="click"
                                    type="password" />
                            </n-form-item-gi>
                            <n-form-item-gi :label="$t('dialogue.connection.usr')" :span="12" path="username">
                                <n-input
                                    v-model:value="generalForm.username"
                                    :placeholder="$t('dialogue.connection.usr_tip')" />
                            </n-form-item-gi>
                        </n-grid>
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
                                :label="$t('dialogue.connection.advn.filter')"
                                :span="12"
                                path="defaultFilter">
                                <n-input
                                    v-model:value="generalForm.defaultFilter"
                                    :placeholder="$t('dialogue.connection.advn.filter_tip')" />
                            </n-form-item-gi>
                            <n-form-item-gi
                                :label="$t('dialogue.connection.advn.separator')"
                                :span="12"
                                path="keySeparator">
                                <n-input
                                    v-model:value="generalForm.keySeparator"
                                    :placeholder="$t('dialogue.connection.advn.separator_tip')" />
                            </n-form-item-gi>
                            <n-form-item-gi
                                :label="$t('dialogue.connection.advn.conn_timeout')"
                                :span="12"
                                path="connTimeout">
                                <n-input-number
                                    v-model:value="generalForm.connTimeout"
                                    :max="999999"
                                    :min="1"
                                    :show-button="false"
                                    style="width: 100%">
                                    <template #suffix>
                                        {{ $t('common.second') }}
                                    </template>
                                </n-input-number>
                            </n-form-item-gi>
                            <n-form-item-gi
                                :label="$t('dialogue.connection.advn.exec_timeout')"
                                :span="12"
                                path="execTimeout">
                                <n-input-number
                                    v-model:value="generalForm.execTimeout"
                                    :max="999999"
                                    :min="1"
                                    :show-button="false"
                                    style="width: 100%">
                                    <template #suffix>
                                        {{ $t('common.second') }}
                                    </template>
                                </n-input-number>
                            </n-form-item-gi>
                            <n-form-item-gi :label="$t('dialogue.connection.advn.key_view')" :span="12">
                                <n-radio-group v-model:value="generalForm.keyView">
                                    <n-radio-button
                                        :label="$t('dialogue.connection.advn.key_view_tree')"
                                        :value="KeyViewType.Tree" />
                                    <n-radio-button
                                        :label="$t('dialogue.connection.advn.key_view_list')"
                                        :value="KeyViewType.List" />
                                </n-radio-group>
                            </n-form-item-gi>
                            <n-form-item-gi :label="$t('dialogue.connection.advn.load_size')" :span="12">
                                <n-input-number
                                    v-model:value="generalForm.loadSize"
                                    :min="0"
                                    :show-button="false"
                                    style="width: 100%" />
                            </n-form-item-gi>
                            <n-form-item-gi :label="$t('dialogue.connection.advn.dbfilter_type')" :span="24">
                                <n-radio-group
                                    v-model:value="generalForm.dbFilterType"
                                    @update:value="onUpdateDBFilterType">
                                    <n-radio-button :label="$t('dialogue.connection.advn.dbfilter_all')" value="none" />
                                    <n-radio-button
                                        :label="$t('dialogue.connection.advn.dbfilter_show')"
                                        value="show" />
                                    <n-radio-button
                                        :label="$t('dialogue.connection.advn.dbfilter_hide')"
                                        value="hide" />
                                </n-radio-group>
                            </n-form-item-gi>
                            <n-form-item-gi
                                v-show="generalForm.dbFilterType !== 'none'"
                                :label="$t('dialogue.connection.advn.dbfilter_input')"
                                :span="24">
                                <n-select
                                    v-model:value="dbFilterList"
                                    :clearable="true"
                                    :disabled="generalForm.dbFilterType === 'none'"
                                    :placeholder="$t('dialogue.connection.advn.dbfilter_input_tip')"
                                    :show="false"
                                    :show-arrow="false"
                                    filterable
                                    multiple
                                    tag />
                            </n-form-item-gi>
                            <n-form-item-gi
                                :label="$t('dialogue.connection.advn.mark_color')"
                                :span="24"
                                path="markColor">
                                <div
                                    v-for="color in predefineColors"
                                    :key="color"
                                    :style="{
                                        backgroundColor: color,
                                        borderColor:
                                            generalForm.markColor === color
                                                ? themeVars.textColorBase
                                                : themeVars.borderColor,
                                    }"
                                    class="color-preset-item"
                                    @click="generalForm.markColor = color">
                                    <n-icon v-if="isEmpty(color)" :component="Close" size="24" />
                                </div>
                            </n-form-item-gi>
                        </n-grid>
                    </n-form>
                </n-tab-pane>

                <!-- Alias pane -->
                <n-tab-pane :tab="$t('dialogue.connection.alias.title')" display-directive="show:lazy" name="alias">
                    <n-form
                        :model="generalForm.alias"
                        :show-label="false"
                        :show-require-mark="false"
                        label-placement="top">
                        <n-form-item required>
                            <n-dynamic-input
                                v-model:value="aliasPair"
                                @create="onCreateAlias"
                                @update:value="onUpdateAlias">
                                <template #default="{ value }">
                                    <n-input-number
                                        v-model:value="value.db"
                                        :min="0"
                                        :placeholder="$t('dialogue.connection.alias.db')"
                                        :show-button="false"
                                        @update:value="onUpdateAlias" />
                                    <n-text>:</n-text>
                                    <n-input
                                        v-model:value="value.alias"
                                        :placeholder="$t('dialogue.connection.alias.value')"
                                        type="text"
                                        @update:value="onUpdateAlias" />
                                </template>
                                <template #action="{ index, create, remove, move }">
                                    <icon-button :icon="Delete" size="18" @click="() => remove(index)" />
                                    <icon-button :icon="Add" size="18" @click="() => create(index)" />
                                </template>
                            </n-dynamic-input>
                        </n-form-item>
                    </n-form>
                </n-tab-pane>

                <!-- SSL pane -->
                <n-tab-pane :tab="$t('dialogue.connection.ssl.title')" display-directive="show:lazy" name="ssl">
                    <n-form-item label-placement="left">
                        <n-checkbox v-model:checked="generalForm.ssl.enable" size="medium">
                            {{ $t('dialogue.connection.ssl.enable') }}
                        </n-checkbox>
                    </n-form-item>
                    <n-form
                        :disabled="!generalForm.ssl.enable"
                        :model="generalForm.ssl"
                        :show-require-mark="false"
                        label-placement="top">
                        <n-form-item :label="$t('dialogue.connection.ssl.cert_file')">
                            <file-open-input
                                v-model:value="generalForm.ssl.certFile"
                                :disabled="!generalForm.ssl.enable"
                                :placeholder="$t('dialogue.connection.ssl.cert_file_tip')" />
                        </n-form-item>
                        <n-form-item :label="$t('dialogue.connection.ssl.key_file')">
                            <file-open-input
                                v-model:value="generalForm.ssl.keyFile"
                                :disabled="!generalForm.ssl.enable"
                                :placeholder="$t('dialogue.connection.ssl.key_file_tip')" />
                        </n-form-item>
                        <n-form-item :label="$t('dialogue.connection.ssl.ca_file')" :show-feedback="false">
                            <file-open-input
                                v-model:value="generalForm.ssl.caFile"
                                :disabled="!generalForm.ssl.enable"
                                :placeholder="$t('dialogue.connection.ssl.ca_file_tip')" />
                        </n-form-item>
                        <n-form-item>
                            <n-checkbox v-model:checked="generalForm.ssl.allowInsecure" size="medium">
                                {{ $t('dialogue.connection.ssl.allow_insecure') }}
                            </n-checkbox>
                        </n-form-item>
                        <n-form-item :label="$t('dialogue.connection.ssl.sni')">
                            <n-input
                                v-model:value="generalForm.ssl.sni"
                                :placeholder="$t('dialogue.connection.ssl.sni')" />
                        </n-form-item>
                    </n-form>
                </n-tab-pane>

                <!-- SSH pane -->
                <n-tab-pane :tab="$t('dialogue.connection.ssh.title')" display-directive="show:lazy" name="ssh">
                    <n-form-item label-placement="left">
                        <n-checkbox v-model:checked="generalForm.ssh.enable" size="medium">
                            {{ $t('dialogue.connection.ssh.enable') }}
                        </n-checkbox>
                    </n-form-item>
                    <n-form
                        :disabled="!generalForm.ssh.enable"
                        :model="generalForm.ssh"
                        :show-require-mark="false"
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
                                :show-button="false"
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
                            <file-open-input
                                v-model:value="generalForm.ssh.pkFile"
                                :disabled="!generalForm.ssh.enable"
                                :placeholder="$t('dialogue.connection.ssh.pkfile_tip')" />
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
                <n-tab-pane
                    :tab="$t('dialogue.connection.sentinel.title')"
                    display-directive="show:lazy"
                    name="sentinel">
                    <n-form-item label-placement="left">
                        <n-checkbox v-model:checked="generalForm.sentinel.enable" size="medium">
                            {{ $t('dialogue.connection.sentinel.enable') }}
                        </n-checkbox>
                    </n-form-item>
                    <n-form
                        :disabled="!generalForm.sentinel.enable"
                        :model="generalForm.sentinel"
                        :show-require-mark="false"
                        label-placement="top">
                        <n-form-item :label="$t('dialogue.connection.sentinel.master')">
                            <n-input-group>
                                <n-select
                                    v-model:value="generalForm.sentinel.master"
                                    :options="masterNameOptions"
                                    filterable
                                    tag />
                                <n-button
                                    :disabled="!generalForm.sentinel.enable"
                                    :loading="loadingSentinelMaster"
                                    @click="onLoadSentinelMasters">
                                    {{ $t('dialogue.connection.sentinel.auto_discover') }}
                                </n-button>
                            </n-input-group>
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

                <!-- Cluster pane -->
                <n-tab-pane :tab="$t('dialogue.connection.cluster.title')" display-directive="show:lazy" name="cluster">
                    <n-form-item label-placement="left">
                        <n-checkbox v-model:checked="generalForm.cluster.enable" size="medium">
                            {{ $t('dialogue.connection.cluster.enable') }}
                        </n-checkbox>
                    </n-form-item>
                    <!--                    <n-form-->
                    <!--                        :model="generalForm.cluster"-->
                    <!--                        :show-require-mark="false"-->
                    <!--                        :disabled="!generalForm.cluster.enable"-->
                    <!--                        label-placement="top">-->
                    <!--                    </n-form>-->
                </n-tab-pane>

                <!-- Proxy pane -->
                <n-tab-pane :tab="$t('dialogue.connection.proxy.title')" display-directive="show:lazy" name="proxy">
                    <n-radio-group v-model:value="generalForm.proxy.type" name="radiogroup">
                        <n-space size="large" vertical>
                            <n-radio :label="$t('dialogue.connection.proxy.type_none')" :value="0" />
                            <n-radio :label="$t('dialogue.connection.proxy.type_system')" :value="1" />
                            <n-radio :label="$t('dialogue.connection.proxy.type_custom')" :value="2" />
                            <n-form
                                :disabled="generalForm.proxy.type !== 2"
                                :model="generalForm.proxy"
                                :show-require-mark="false"
                                label-placement="top">
                                <n-grid :x-gap="10">
                                    <n-form-item-gi :show-label="false" :span="24" path="addr" required>
                                        <n-input-group>
                                            <n-select
                                                v-model:value="generalForm.proxy.schema"
                                                :consistent-menu-width="false"
                                                :options="[
                                                    { value: 'http', label: 'HTTP' },
                                                    { value: 'https', label: 'HTTPS' },
                                                    { value: 'socks5', label: 'SOCKS5' },
                                                    { value: 'socks5h', label: 'SOCKS5H' },
                                                ]"
                                                default-value="http"
                                                style="max-width: 100px" />
                                            <n-input
                                                v-model:value="generalForm.proxy.addr"
                                                :placeholder="$t('dialogue.connection.proxy.host')" />
                                            <n-text style="width: 40px; text-align: center">:</n-text>
                                            <n-input-number
                                                v-model:value="generalForm.proxy.port"
                                                :max="65535"
                                                :min="0"
                                                :show-button="false"
                                                style="width: 200px" />
                                        </n-input-group>
                                    </n-form-item-gi>
                                    <n-form-item-gi :show-label="false" :span="24" path="auth">
                                        <n-checkbox v-model:checked="generalForm.proxy.auth" size="medium">
                                            {{ $t('dialogue.connection.proxy.auth') }}
                                        </n-checkbox>
                                    </n-form-item-gi>
                                    <n-form-item-gi :label="$t('dialogue.connection.usr')" :span="12" path="username">
                                        <n-input
                                            v-model:value="generalForm.proxy.username"
                                            :disabled="!!!generalForm.proxy.auth"
                                            :placeholder="$t('dialogue.connection.proxy.usr_tip')" />
                                    </n-form-item-gi>
                                    <n-form-item-gi :label="$t('dialogue.connection.pwd')" :span="12" path="password">
                                        <n-input
                                            v-model:value="generalForm.proxy.password"
                                            :disabled="!!!generalForm.proxy.auth"
                                            :placeholder="$t('dialogue.connection.proxy.pwd_tip')"
                                            show-password-on="click"
                                            type="password" />
                                    </n-form-item-gi>
                                </n-grid>
                            </n-form>
                        </n-space>
                    </n-radio-group>
                </n-tab-pane>
            </n-tabs>

            <!-- test result alert-->
            <n-alert
                v-if="showTestResult"
                :on-close="() => (testResult = '')"
                :title="isEmpty(testResult) ? '' : $t('dialogue.connection.test_fail')"
                :type="isEmpty(testResult) ? 'success' : 'error'"
                closable>
                <template v-if="isEmpty(testResult)">{{ $t('dialogue.connection.test_succ') }}</template>
                <template v-else>{{ testResult }}</template>
            </n-alert>
        </n-spin>

        <template #action>
            <div class="flex-item-expand">
                <n-button :disabled="closingConnection" :focusable="false" :loading="testing" @click="onTestConnection">
                    {{ $t('dialogue.connection.test') }}
                </n-button>
            </div>
            <div class="flex-item n-dialog__action">
                <n-button :disabled="closingConnection" :focusable="false" @click="pasteFromClipboard">
                    {{ $t('dialogue.connection.parse_url_clipboard') }}
                </n-button>
                <n-button :disabled="closingConnection" :focusable="false" @click="onClose">
                    {{ $t('common.cancel') }}
                </n-button>
                <n-button :disabled="closingConnection" :focusable="false" type="primary" @click="onSaveConnection">
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
    border-width: 3px;
    border-style: solid;
    cursor: pointer;
    border-radius: 50%;
}
</style>
