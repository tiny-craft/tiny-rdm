<script setup>
import { computed, h, ref } from 'vue'
import { map } from 'lodash'
import { NIcon, NText } from 'naive-ui'

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
})

const emit = defineEmits(['update:value'])

const renderHeader = () => {
    return h('div', { class: 'type-selector-header' }, [h(NText, null, () => props.tooltip)])
}

const renderLabel = (option) => {
    return h('div', { class: 'type-selector-item' }, option.label)
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
    return props.value || get(dropdownOption.value, [1, 'label'], 'None')
})

const showDropdown = ref(false)
const onDropdownShow = (show) => {
    showDropdown.value = show === true
}
</script>

<template>
    <n-dropdown
        :options="dropdownOption"
        :render-label="renderLabel"
        :show-arrow="true"
        :title="props.tooltip"
        :value="props.value"
        trigger="click"
        @select="onDropdownSelect"
        @update:show="onDropdownShow">
        <n-tooltip :disabled="showDropdown" :show-arrow="false">
            {{ props.tooltip }}
            <template #trigger>
                <n-button :focusable="false" quaternary>
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

<style lang="scss" scoped></style>
