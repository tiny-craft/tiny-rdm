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
import { computed, onMounted, ref, watch } from 'vue'
import { isEmpty } from 'lodash'
import useDialogStore from 'stores/dialog.js'
import { useI18n } from 'vue-i18n'
import ContentToolbar from '@/components/content_value/ContentToolbar.vue'
import ContentValueJson from '@/components/content_value/ContentValueJson.vue'

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
const i18n = useI18n()

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
 *      format: String,
 *      decode: String,
 *      end: Boolean
 *      loading: Boolean
 * }>}
 */
const data = computed(() => {
    return props.content
})
const initializing = ref(false)

const loading = computed(() => {
    return data.value.loading === true || initializing.value
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
    [redisTypes.JSON]: ContentValueJson,
}

const keyName = computed(() => {
    return !isEmpty(data.value.keyCode) ? data.value.keyCode : data.value.keyPath
})

/**
 *
 * @param {boolean} reset
 * @param {boolean} [full]
 * @param {string} [selMatch]
 * @return {Promise<void>}
 */
const loadData = async (reset, full, selMatch) => {
    try {
        if (!!props.blank) {
            return
        }
        const { name, db, matchPattern } = data.value
        reset = reset === true
        await browserStore.loadKeyDetail({
            server: name,
            db: db,
            key: keyName.value,
            matchPattern: selMatch === undefined ? matchPattern : selMatch,
            decode: '',
            format: '',
            reset,
            full: full === true,
        })
    } finally {
    }
}

/**
 * reload current key
 * @param {string} [selDecode]
 * @param {string} [selFormat]
 * @return {Promise<void>}
 */
const onReload = async (selDecode, selFormat) => {
    try {
        const { name, db, keyCode, keyPath, decode, format, matchPattern } = data.value
        const targetFormat = selFormat || format
        const targetDecode = selDecode || decode
        browserStore.setSelectedFormat(name, keyPath, db, targetFormat, targetDecode)
        await browserStore.reloadKey({
            server: name,
            db,
            key: keyCode || keyPath,
            decode: targetDecode,
            format: targetFormat,
            matchPattern,
            showLoading: false,
        })
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
    $dialog.warning(i18n.t('dialogue.remove_tip', { name: data.value.keyPath }), () => {
        const { name, db } = data.value
        browserStore.deleteKey(name, db, keyName.value).then((success) => {
            if (success) {
                $message.success(i18n.t('dialogue.delete.success', { key: data.value.keyPath }))
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

const onMatch = (match) => {
    loadData(true, false, match || '')
}

const contentRef = ref(null)
const initContent = async () => {
    // onReload()
    try {
        initializing.value = true
        if (contentRef.value?.reset != null) {
            contentRef.value?.reset()
        }
        await loadData(true, false, '')
    } finally {
        initializing.value = false
    }
}

onMounted(() => {
    // onReload()
    initContent()
})

watch(() => data.value?.keyPath, initContent)
</script>

<template>
    <n-empty v-if="props.blank" :description="$t('interface.nonexist_tab_content')" class="empty-content">
        <template #extra>
            <n-button :focusable="false" @click="onReload">{{ $t('interface.reload') }}</n-button>
        </template>
    </n-empty>
    <!-- FIXME: keep alive may cause virtual list null value error. -->
    <!-- <keep-alive v-else> -->
    <component
        :is="valueComponents[data.type]"
        v-else
        ref="contentRef"
        :db="data.db"
        :decode="data.decode"
        :end="data.end"
        :format="data.format"
        :key-code="data.keyCode"
        :key-path="data.keyPath"
        :length="data.length"
        :loading="loading"
        :name="data.name"
        :size="data.size"
        :ttl="data.ttl"
        :value="data.value"
        @delete="onDelete"
        @loadall="onLoadAll"
        @loadmore="onLoadMore"
        @match="onMatch"
        @reload="onReload">
        <template #toolbar>
            <content-toolbar
                :db="data.db"
                :key-code="data.keyCode"
                :key-path="data.keyPath"
                :key-type="data.type"
                :loading="loading"
                :server="data.name"
                :ttl="data.ttl"
                class="value-item-part"
                @delete="onDelete"
                @reload="onReload"
                @rename="onRename" />
        </template>
    </component>
    <!--    </keep-alive>-->
</template>

<style lang="scss" scoped></style>
