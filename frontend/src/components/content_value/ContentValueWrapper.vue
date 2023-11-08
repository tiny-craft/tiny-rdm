<script setup>
import { types as redisTypes } from '@/consts/support_redis_type.js'
import ContentValueString from '@/components/content_value/ContentValueString.vue'
import ContentValueHash from '@/components/content_value/ContentValueHash.vue'
import ContentValueList from '@/components/content_value/ContentValueList.vue'
import ContentValueSet from '@/components/content_value/ContentValueSet.vue'
import ContentValueZset from '@/components/content_value/ContentValueZSet.vue'
import ContentValueStream from '@/components/content_value/ContentValueStream.vue'
import { useThemeVars } from 'naive-ui'
import useBrowserStore from 'stores/browser.js'
import { computed, onMounted, watch } from 'vue'
import { isEmpty } from 'lodash'
import { decodeTypes, formatTypes } from '@/consts/value_view_type.js'
import useDialogStore from 'stores/dialog.js'

const themeVars = useThemeVars()
const browserStore = useBrowserStore()
const dialogStore = useDialogStore()

const props = defineProps({
    blank: Boolean,
    content: {
        type: Object,
        default: {},
    },
})

/**
 *
 * @type {ComputedRef<{
 *      type:
 *      String, name: String,
 *      db: Number,
 *      keyPath: String,
 *      keyCode: Array,
 *      ttl: Number,
 *      value: [String, Object],
 *      size: Number,
 *      length: Number,
 *      viewAs: String,
 *      decode: String,
 *      end: Boolean
 * }>}
 */
const data = computed(() => {
    return props.content
})

const binaryKey = computed(() => {
    return !!data.value.keyCode
})

const valueComponents = {
    [redisTypes.STRING]: ContentValueString,
    [redisTypes.HASH]: ContentValueHash,
    [redisTypes.LIST]: ContentValueList,
    [redisTypes.SET]: ContentValueSet,
    [redisTypes.ZSET]: ContentValueZset,
    [redisTypes.STREAM]: ContentValueStream,
}

const keyName = computed(() => {
    return !isEmpty(data.value.keyCode) ? data.value.keyCode : data.value.keyPath
})

const loadData = async (reset, full) => {
    try {
        const { name, db, view, decodeType, matchPattern } = data.value
        await browserStore.loadKeyDetail({
            server: name,
            db: db,
            key: keyName.value,
            viewType: view,
            decodeType: decodeType,
            matchPattern: matchPattern,
            reset: reset === true,
            full: full === true,
        })
    } finally {
    }
}

const onReload = async () => {
    try {
        const { name, db, keyCode, keyPath } = data.value
        await browserStore.reloadKey({ server: name, db, key: keyCode || keyPath })
    } finally {
    }
}

const onRename = () => {
    const { name, db, keyPath } = data.value
    if (binaryKey.value) {
        $message.error(i18n.t('dialogue.rename_binary_key_fail'))
    } else {
        dialogStore.openRenameKeyDialog(name, db, keyPath)
    }
}

const onDelete = () => {
    $dialog.warning(i18n.t('dialogue.remove_tip', { name: props.keyPath }), () => {
        const { name, db } = data.value
        browserStore.deleteKey(name, db, keyName.value).then((success) => {
            if (success) {
                $message.success(i18n.t('dialogue.delete_key_succ', { key: props.keyPath }))
            }
        })
    })
}

const onLoadMore = () => {
    loadData(false, false)
}

const onLoadAll = () => {
    loadData(false, true)
}

onMounted(() => {
    // onReload()
    loadData(false, false)
})

watch(
    () => data.value?.keyPath,
    () => {
        // onReload()
        loadData(false, false)
    },
)
</script>

<template>
    <n-empty v-if="props.blank" :description="$t('interface.nonexist_tab_content')" class="empty-content">
        <template #extra>
            <n-button :focusable="false" @click="onReload">{{ $t('interface.reload') }}</n-button>
        </template>
    </n-empty>
    <keep-alive v-else>
        <component
            :is="valueComponents[data.type]"
            :db="data.db"
            :decode="data.decode || decodeTypes.NONE"
            :end="data.end"
            :key-code="data.keyCode"
            :key-path="data.keyPath"
            :length="data.length"
            :loading="data.loading === true"
            :name="data.name"
            :size="data.size"
            :ttl="data.ttl"
            :value="data.value"
            :view-as="data.viewAs || formatTypes.PLAIN_TEXT"
            @delete="onDelete"
            @loadall="onLoadAll"
            @loadmore="onLoadMore"
            @reload="onReload"
            @rename="onRename" />
    </keep-alive>
</template>

<style lang="scss" scoped></style>
