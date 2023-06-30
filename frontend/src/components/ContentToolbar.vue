<script setup>
import { validType } from '../consts/support_redis_type'
import useDialog from '../stores/dialog'
import Delete from './icons/Delete.vue'
import Edit from './icons/Edit.vue'
import Refresh from './icons/Refresh.vue'
import Timer from './icons/Timer.vue'
import RedisTypeTag from './common/RedisTypeTag.vue'
import { useI18n } from 'vue-i18n'
import { useMessage } from 'naive-ui'
import IconButton from './common/IconButton.vue'
import useConnectionStore from '../stores/connections.js'

const props = defineProps({
    server: String,
    db: Number,
    keyType: {
        type: String,
        validator(value) {
            return validType(value)
        },
        default: 'STRING',
    },
    keyPath: String,
    ttl: {
        type: Number,
        default: -1,
    },
})

const dialogStore = useDialog()
const connectionStore = useConnectionStore()
const message = useMessage()
const i18n = useI18n()

const onReloadKey = () => {
    connectionStore.loadKeyValue(props.server, props.db, props.keyPath)
}

const onConfirmDelete = async () => {
    const success = await connectionStore.removeKey(props.server, props.db, props.keyPath)
    if (success) {
        message.success(i18n.t('delete_key_succ', { key: props.keyPath }))
    }
}
</script>

<template>
    <div class="content-toolbar flex-box-h">
        <n-input-group>
            <redis-type-tag :type="props.keyType" size="large"></redis-type-tag>
            <n-input v-model:value="props.keyPath">
                <template #suffix>
                    <icon-button :icon="Refresh" :tooltip="$t('reload_key')" size="18" @click="onReloadKey" />
                </template>
            </n-input>
        </n-input-group>
        <n-button-group>
            <n-tooltip>
                <template #trigger>
                    <n-button @click="dialogStore.openTTLDialog(props.ttl)">
                        <template #icon>
                            <n-icon :component="Timer" size="18" />
                        </template>
                        <template v-if="ttl < 0">
                            {{ $t('forever') }}
                        </template>
                        <template v-else> {{ ttl }} {{ $t('second') }}</template>
                    </n-button>
                </template>
                TTL
            </n-tooltip>
            <n-button @click="dialogStore.openRenameKeyDialog(props.server, props.db, props.keyPath)">
                <template #icon>
                    <n-icon :component="Edit" size="18" />
                </template>
                {{ $t('rename_key') }}
            </n-button>
        </n-button-group>
        <n-tooltip>
            <template #trigger>
                <n-popconfirm
                    :negative-text="$t('cancel')"
                    :positive-text="$t('confirm')"
                    @positive-click="onConfirmDelete"
                >
                    <template #trigger>
                        <n-button>
                            <template #icon>
                                <n-icon :component="Delete" size="18" />
                            </template>
                        </n-button>
                    </template>
                    {{ $t('delete_key_tip', { key: props.keyPath }) }}
                </n-popconfirm>
            </template>
            {{ $t('delete_key') }}
        </n-tooltip>
    </div>
</template>

<style lang="scss" scoped>
.content-toolbar {
    align-items: center;
    gap: 5px;
}
</style>
