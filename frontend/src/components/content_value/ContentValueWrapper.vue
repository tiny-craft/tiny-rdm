<script setup>
import { decodeTypes, formatTypes } from '@/consts/value_view_type.js'
import { types as redisTypes } from '@/consts/support_redis_type.js'
import ContentValueString from '@/components/content_value/ContentValueString.vue'
import ContentValueHash from '@/components/content_value/ContentValueHash.vue'
import ContentValueList from '@/components/content_value/ContentValueList.vue'
import ContentValueSet from '@/components/content_value/ContentValueSet.vue'
import ContentValueZset from '@/components/content_value/ContentValueZSet.vue'
import ContentValueStream from '@/components/content_value/ContentValueStream.vue'
import { useThemeVars } from 'naive-ui'
import useConnectionStore from 'stores/connections.js'

const themeVars = useThemeVars()
const connectionStore = useConnectionStore()

const props = defineProps({
    blank: Boolean,
    type: String,
    name: String,
    db: Number,
    keyPath: String,
    keyCode: {
        type: Array,
        default: null,
    },
    ttl: {
        type: Number,
        default: -1,
    },
    value: [String, Object],
    size: Number,
    length: Number,
    viewAs: {
        type: String,
        default: formatTypes.PLAIN_TEXT,
    },
    decode: {
        type: String,
        default: decodeTypes.NONE,
    },
})

const valueComponents = {
    [redisTypes.STRING]: ContentValueString,
    [redisTypes.HASH]: ContentValueHash,
    [redisTypes.LIST]: ContentValueList,
    [redisTypes.SET]: ContentValueSet,
    [redisTypes.ZSET]: ContentValueZset,
    [redisTypes.STREAM]: ContentValueStream,
}

/**
 * reload current selection key
 * @returns {Promise<null>}
 */
const onReloadKey = async () => {
    await connectionStore.loadKeyValue(props.name, props.db, props.key, props.viewAs)
}
</script>

<template>
    <n-empty v-if="props.blank" :description="$t('interface.nonexist_tab_content')" class="empty-content">
        <template #extra>
            <n-button :focusable="false" @click="onReloadKey">{{ $t('interface.reload') }}</n-button>
        </template>
    </n-empty>
    <keep-alive v-else>
        <component
            :is="valueComponents[props.type]"
            :db="props.db"
            :decode="props.decode"
            :key-code="props.keyCode"
            :key-path="props.keyPath"
            :length="props.length"
            :name="props.name"
            :size="props.size"
            :ttl="props.ttl"
            :value="props.value"
            :view-as="props.viewAs" />
    </keep-alive>
</template>

<style lang="scss" scoped></style>
