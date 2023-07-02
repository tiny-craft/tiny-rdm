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
    disabled: Boolean,
})

const disableColor = computed(() => {
    const baseColor = props.color
    const grayScale = Math.round(
        0.299 * parseInt(baseColor.substring(1, 2), 16) +
            0.587 * parseInt(baseColor.substring(3, 2), 16) +
            0.114 * parseInt(baseColor.substring(5, 2), 16)
    )
    const color = `#${grayScale.toString(16).repeat(3)}`
    return color
})

const hasTooltip = computed(() => {
    return props.tooltip || props.tTooltip
})
</script>

<template>
    <n-tooltip v-if="hasTooltip">
        <template #trigger>
            <n-button text :disabled="disabled" @click="emit('click')">
                <n-icon :size="props.size" :color="props.color">
                    <component :is="props.icon" :stroke-width="props.strokeWidth" />
                </n-icon>
            </n-button>
        </template>
        {{ props.tTooltip ? $t(props.tTooltip) : props.tooltip }}
    </n-tooltip>
    <n-button v-else text :disabled="disabled" @click="emit('click')">
        <n-icon :size="props.size" :color="props.color">
            <component :is="props.icon" :stroke-width="props.strokeWidth" />
        </n-icon>
    </n-button>
</template>

<style lang="scss"></style>
