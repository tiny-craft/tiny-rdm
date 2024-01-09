<script setup>
import { computed, h, ref } from 'vue'
import { get, map } from 'lodash'
import { NIcon, NText } from 'naive-ui'
import { useRender } from '@/utils/render.js'

const props = defineProps({
    value: {
        type: String,
        value: '',
    },
    options: {
        type: Object,
        value: {},
    },
    tooltip: {
        type: String,
    },
    icon: [String, Object],
    default: String,
    disabled: Boolean,
})

const emit = defineEmits(['update:value'])
const render = useRender()

const renderHeader = () => {
    return h('div', { class: 'type-selector-header' }, [h(NText, null, () => props.tooltip)])
}

const dropdownOption = computed(() => {
    const options = [
        {
            key: 'header',
            type: 'render',
            render: renderHeader,
        },
        {
            key: 'header-divider',
            type: 'divider',
        },
    ]
    return [
        ...options,
        ...map(props.options, (t) => {
            return {
                key: t,
                label: t,
            }
        }),
    ]
})

const onDropdownSelect = (key) => {
    emit('update:value', key)
}

const buttonText = computed(() => {
    return props.value || get(dropdownOption.value, [1, 'label'], props.default)
})

const showDropdown = ref(false)
const onDropdownShow = (show) => {
    showDropdown.value = show === true
}
</script>

<template>
    <n-dropdown
        :disabled="props.disabled"
        :options="dropdownOption"
        :render-label="({ label }) => render.renderLabel(label, { class: 'type-selector-item' })"
        :show-arrow="true"
        :title="props.tooltip"
        :value="props.value"
        trigger="click"
        @select="onDropdownSelect"
        @update:show="onDropdownShow">
        <n-tooltip :disabled="showDropdown" :show-arrow="false">
            {{ props.tooltip }}
            <template #trigger>
                <n-button :disabled="disabled" :focusable="false" quaternary>
                    <template #icon>
                        <n-icon>
                            <component :is="icon" />
                        </n-icon>
                    </template>
                    {{ buttonText }}
                </n-button>
            </template>
        </n-tooltip>
    </n-dropdown>
</template>

<style lang="scss">
.type-selector-header {
    height: 30px;
    line-height: 30px;
    font-size: 15px;
    font-weight: bold;
    text-align: center;
    padding: 0 10px;
}

.type-selector-item {
    min-width: 100px;
    text-align: center;
}
</style>
