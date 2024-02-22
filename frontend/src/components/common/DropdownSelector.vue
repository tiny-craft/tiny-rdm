<script setup>
import { computed, h, ref } from 'vue'
import { get, isEmpty, some } from 'lodash'
import { NIcon, NText } from 'naive-ui'
import { useRender } from '@/utils/render.js'
import { useI18n } from 'vue-i18n'

const props = defineProps({
    value: {
        type: String,
        value: '',
    },
    options: {
        type: Array,
        value: () => [],
    },
    menuOption: {
        type: Array,
        value: () => [],
    },
    tooltip: {
        type: String,
    },
    icon: [String, Object],
    default: String,
    disabled: Boolean,
})

const emit = defineEmits(['update:value', 'menu'])
const i18n = useI18n()
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
    if (get(props.options, 0) instanceof Array) {
        // multiple group
        for (let i = 0; i < props.options.length; i++) {
            if (i !== 0 && !isEmpty(props.options[i])) {
                // add divider
                options.push({
                    key: 'header-divider' + (i + 1),
                    type: 'divider',
                })
            }
            for (const option of props.options[i]) {
                options.push({
                    key: option,
                    label: option,
                })
            }
        }
    } else {
        for (const option of props.options) {
            options.push({
                key: option,
                label: option,
            })
        }
    }

    if (!isEmpty(props.menuOption)) {
        options.push({
            key: 'header-divider',
            type: 'divider',
        })
        for (const { key, label } of props.menuOption) {
            options.push({
                key,
                label: i18n.t(label),
            })
        }
    }
    return options
})

const onDropdownSelect = (key) => {
    if (some(props.menuOption, { key })) {
        emit('menu', key)
    } else {
        emit('update:value', key)
    }
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
