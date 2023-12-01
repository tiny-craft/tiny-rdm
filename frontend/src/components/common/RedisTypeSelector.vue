<script setup>
import { computed, h } from 'vue'
import { useThemeVars } from 'naive-ui'
import { types, typesBgColor, typesColor } from '@/consts/support_redis_type.js'
import { get, map, toUpper } from 'lodash'

const props = defineProps({
    value: {
        type: String,
        default: 'ALL',
    },
})

const emit = defineEmits(['update:value', 'select'])

const options = computed(() => {
    const opts = map(types, (v) => ({
        label: v,
        key: v,
    }))
    return [{ label: 'ALL', key: 'ALL' }, ...opts]
})

const themeVars = useThemeVars()
const renderIcon = (option) => {
    if (option.key === props.value) {
        const backgroundColor = get(typesColor, option.key, themeVars.value.textColor3)
        return h('div', { style: { borderRadius: '999px', width: '10px', height: '10px', backgroundColor } }, '')
    }
}

const renderLabel = (option) => {
    const color = get(typesColor, option.key, '')
    return h('div', { style: { color, fontWeight: '450' } }, option.label)
}

const fontColor = computed(() => {
    return get(typesColor, props.value, '')
})

const backgroundColor = computed(() => {
    return get(typesBgColor, props.value, '')
})

const displayValue = computed(() => {
    return get(types, toUpper(props.value), 'ALL')
})

const handleSelect = (select) => {
    if (props.value !== select) {
        emit('update:value', select)
        emit('select', select)
    }
}
</script>

<template>
    <n-dropdown
        :options="options"
        :render-icon="renderIcon"
        :render-label="renderLabel"
        show-arrow
        @select="handleSelect">
        <n-tag
            :bordered="true"
            :color="{ color: backgroundColor, textColor: fontColor }"
            class="redis-tag"
            size="medium"
            strong>
            {{ displayValue }}
        </n-tag>
    </n-dropdown>
</template>

<style lang="scss" scoped>
.redis-tag {
    padding: 0 10px;
}

:deep(.dropdown-type-item) {
    padding: 10px;
}
</style>
