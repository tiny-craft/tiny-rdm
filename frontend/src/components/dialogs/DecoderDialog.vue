<script setup>
import useDialog from 'stores/dialog.js'
import { computed, reactive, ref, toRaw, watch } from 'vue'
import FileOpenInput from '@/components/common/FileOpenInput.vue'
import Delete from '@/components/icons/Delete.vue'
import Add from '@/components/icons/Add.vue'
import IconButton from '@/components/common/IconButton.vue'
import { cloneDeep, get, isEmpty } from 'lodash'
import usePreferencesStore from 'stores/preferences.js'
import { joinCommand } from '@/utils/decoder_cmd.js'
import Help from '@/components/icons/Help.vue'

const editName = ref('')
const decoderForm = reactive({
    name: '',
    auto: true,
    decodePath: '',
    decodeArgs: [],
    encodePath: '',
    encodeArgs: [],
})

const dialogStore = useDialog()
const prefStore = usePreferencesStore()

watch(
    () => dialogStore.decodeDialogVisible,
    (visible) => {
        if (visible) {
            const name = get(dialogStore.decodeParam, 'name', '')
            if (!isEmpty(name)) {
                editName.value = decoderForm.name = name
                decoderForm.auto = dialogStore.decodeParam.auto !== false
                decoderForm.decodePath = get(dialogStore.decodeParam, 'decodePath', '')
                decoderForm.decodeArgs = get(dialogStore.decodeParam, 'decodeArgs', [])
                decoderForm.encodePath = get(dialogStore.decodeParam, 'encodePath', '')
                decoderForm.encodeArgs = get(dialogStore.decodeParam, 'encodeArgs', [])
            } else {
                editName.value = ''
                decoderForm.decodePath = ''
                decoderForm.encodePath = ''
                decoderForm.decodeArgs = []
                decoderForm.encodeArgs = []
            }
        } else {
            editName.value = ''
        }
    },
)

const decodeCmdPreview = computed(() => {
    return joinCommand(decoderForm.decodePath, decoderForm.decodeArgs, '')
})

const encodeCmdPreview = computed(() => {
    return joinCommand(decoderForm.encodePath, decoderForm.encodeArgs, '')
})

const onAddOrUpdate = () => {
    if (isEmpty(editName.value)) {
        // add decoder
        prefStore.addCustomDecoder(toRaw(decoderForm))
    } else {
        // update decoder
        const param = cloneDeep(toRaw(decoderForm))
        param.newName = param.name
        param.name = editName.value
        prefStore.updateCustomDecoder(param)
    }
}
const onClose = () => {}
</script>

<template>
    <n-modal
        v-model:show="dialogStore.decodeDialogVisible"
        :closable="false"
        :mask-closable="false"
        :negative-button-props="{ focusable: false, size: 'medium' }"
        :negative-text="$t('common.cancel')"
        :positive-button-props="{ focusable: false, size: 'medium' }"
        :positive-text="$t('common.confirm')"
        :show-icon="false"
        :title="editName ? $t('dialogue.decoder.edit_name') : $t('dialogue.decoder.name')"
        close-on-esc
        preset="dialog"
        transform-origin="center"
        @esc="onClose"
        @positive-click="onAddOrUpdate"
        @negative-click="onClose">
        <n-form :model="decoderForm" :show-require-mark="false" label-align="left" label-placement="top">
            <n-form-item :label="$t('dialogue.decoder.decoder_name')" required show-require-mark>
                <n-input v-model:value="decoderForm.name" />
            </n-form-item>
            <n-tabs type="line">
                <!-- decode pane -->
                <n-tab-pane :tab="$t('dialogue.decoder.decoder')" name="decode">
                    <n-form-item required show-require-mark>
                        <template #label>
                            <n-space :size="5" :wrap-item="false" align="center" justify="center">
                                <span>{{ $t('dialogue.decoder.decode_path') }}</span>
                                <n-tooltip trigger="hover">
                                    <template #trigger>
                                        <n-icon :component="Help" />
                                    </template>
                                    <div class="text-block" style="max-width: 600px">
                                        {{ $t('dialogue.decoder.path_help') }}
                                    </div>
                                </n-tooltip>
                            </n-space>
                        </template>
                        <file-open-input
                            v-model:value="decoderForm.decodePath"
                            :placeholder="$t('dialogue.decoder.decode_path')" />
                    </n-form-item>
                    <n-form-item required>
                        <template #label>
                            <n-space :size="5" :wrap-item="false" align="center" justify="center">
                                <span>{{ $t('dialogue.decoder.args') }}</span>
                                <n-tooltip trigger="hover">
                                    <template #trigger>
                                        <n-icon :component="Help" />
                                    </template>
                                    <div class="text-block" style="max-width: 600px">
                                        {{ $t('dialogue.decoder.args_help').replace('[', '{').replace(']', '}') }}
                                    </div>
                                </n-tooltip>
                            </n-space>
                        </template>
                        <n-dynamic-input v-model:value="decoderForm.decodeArgs" @create="() => ''">
                            <template #action="{ index, create, remove, move }">
                                <icon-button :icon="Add" size="18" @click="() => create(index)" />
                                <icon-button :icon="Delete" size="18" @click="() => remove(index)" />
                            </template>
                        </n-dynamic-input>
                    </n-form-item>
                    <n-card
                        v-if="decodeCmdPreview"
                        content-class="cmd-line"
                        content-style="padding: 10px;"
                        embedded
                        size="small">
                        {{ decodeCmdPreview }}
                    </n-card>
                </n-tab-pane>

                <!-- encode pane -->
                <n-tab-pane :tab="$t('dialogue.decoder.encoder')" name="encode">
                    <n-form-item required show-require-mark>
                        <template #label>
                            <n-space :size="5" :wrap-item="false" align="center" justify="center">
                                <span>{{ $t('dialogue.decoder.encode_path') }}</span>
                                <n-tooltip trigger="hover">
                                    <template #trigger>
                                        <n-icon :component="Help" />
                                    </template>
                                    <div class="text-block" style="max-width: 600px">
                                        {{ $t('dialogue.decoder.path_help') }}
                                    </div>
                                </n-tooltip>
                            </n-space>
                        </template>
                        <file-open-input
                            v-model:value="decoderForm.encodePath"
                            :placeholder="$t('dialogue.decoder.encode_path')" />
                    </n-form-item>
                    <n-form-item :label="$t('dialogue.decoder.args')" required>
                        <template #label>
                            <n-space :size="5" :wrap-item="false" align="center" justify="center">
                                <span>{{ $t('dialogue.decoder.args') }}</span>
                                <n-tooltip trigger="hover">
                                    <template #trigger>
                                        <n-icon :component="Help" />
                                    </template>
                                    <div class="text-block" style="max-width: 600px">
                                        {{ $t('dialogue.decoder.args_help').replace('[', '{').replace(']', '}') }}
                                    </div>
                                </n-tooltip>
                            </n-space>
                        </template>
                        <n-dynamic-input v-model:value="decoderForm.encodeArgs" @create="() => ''">
                            <template #action="{ index, create, remove, move }">
                                <icon-button :icon="Add" size="18" @click="() => create(index)" />
                                <icon-button :icon="Delete" size="18" @click="() => remove(index)" />
                            </template>
                        </n-dynamic-input>
                    </n-form-item>
                    <n-card
                        v-if="encodeCmdPreview"
                        content-class="cmd-line"
                        content-style="padding: 10px;"
                        embedded
                        size="small">
                        {{ encodeCmdPreview }}
                    </n-card>
                </n-tab-pane>
            </n-tabs>
            <n-form-item :show-feedback="false">
                <n-checkbox v-model:checked="decoderForm.auto" :label="$t('dialogue.decoder.auto')" />
            </n-form-item>
        </n-form>
    </n-modal>
</template>

<style lang="scss" scoped>
@import '@/styles/content';
</style>
