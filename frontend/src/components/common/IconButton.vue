<script setup>
import { computed } from 'vue'
import { NIcon } from 'naive-ui'

const emit = defineEmits(['click'])

const props = defineProps({
    tooltip: String,
    tTooltip: String,
    icon: [String, Object],
    size: {
        type: [Number, String],
        default: 20,
    },
    color: {
        type: String,
        default: 'currentColor',
    },
    strokeWidth: {
        type: [Number, String],
        default: 3,
    },
    border: Boolean,
    disabled: Boolean,
})

const hasTooltip = computed(() => {
    return props.tooltip || props.tTooltip
})
</script>

<template>
    <n-tooltip v-if="hasTooltip">
        <template #trigger>
            <n-button :disabled="disabled" :text="!border" :focusable="false" @click.prevent="emit('click')">
                <n-icon :color="props.color" :size="props.size">
                    <component :is="props.icon" :stroke-width="props.strokeWidth" />
                </n-icon>
            </n-button>
        </template>
        {{ props.tTooltip ? $t(props.tTooltip) : props.tooltip }}
    </n-tooltip>
    <n-button v-else :disabled="disabled" :text="!border" :focusable="false" @click.prevent="emit('click')">
        <n-icon :color="props.color" :size="props.size">
            <component :is="props.icon" :stroke-width="props.strokeWidth" />
        </n-icon>
    </n-button>
</template>

<style lang="scss"></style>
