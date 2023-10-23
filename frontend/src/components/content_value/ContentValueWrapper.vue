<script setup>
import { types } from '@/consts/value_view_type.js'
import { types as redisTypes } from '@/consts/support_redis_type.js'
import ContentValueString from '@/components/content_value/ContentValueString.vue'
import ContentValueHash from '@/components/content_value/ContentValueHash.vue'
import ContentValueList from '@/components/content_value/ContentValueList.vue'
import ContentValueSet from '@/components/content_value/ContentValueSet.vue'
import ContentValueZset from '@/components/content_value/ContentValueZSet.vue'
import ContentValueStream from '@/components/content_value/ContentValueStream.vue'
import { useThemeVars } from 'naive-ui'

const themeVars = useThemeVars()

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
    viewAs: {
        type: String,
        default: types.PLAIN_TEXT,
    },
})

const emit = defineEmits(['reload'])

const valueComponents = {
    [redisTypes.STRING]: ContentValueString,
    [redisTypes.HASH]: ContentValueHash,
    [redisTypes.LIST]: ContentValueList,
    [redisTypes.SET]: ContentValueSet,
    [redisTypes.ZSET]: ContentValueZset,
    [redisTypes.STREAM]: ContentValueStream,
}
</script>

<template>
    <n-empty v-if="props.blank" :description="$t('interface.nonexist_tab_content')" class="empty-content">
        <template #extra>
            <n-button :focusable="false" @click="emit('reload')">{{ $t('interface.reload') }}</n-button>
        </template>
    </n-empty>
    <component
        class="content-value-wrapper"
        :is="valueComponents[props.type]"
        v-else
        :db="props.db"
        :key-code="props.keyCode"
        :key-path="props.keyPath"
        :name="props.name"
        :size="props.size"
        :ttl="props.ttl"
        :value="props.value"
        :view-as="props.viewAs" />
</template>

<style scoped lang="scss">
.content-value-wrapper {
    background-color: v-bind('themeVars.bodyColor');
}
</style>
