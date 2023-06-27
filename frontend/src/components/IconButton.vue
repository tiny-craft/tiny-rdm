<script setup>
import { computed } from 'vue'
import { NIcon } from 'naive-ui'

const emit = defineEmits(['click'])

const props = defineProps({
    tooltip: {
        type: String,
    },
    tTooltip: {
        type: String,
    },
    icon: {
        type: [String, Object],
    },
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
})

const hasTooltip = computed(() => {
    return props.tooltip || props.tTooltip
})
</script>

<template>
    <n-tooltip v-if="hasTooltip">
        <template #trigger>
            <n-icon :color="props.color" :size="props.size" class="icon-btn" @click="emit('click')">
                <component :is="props.icon" :stroke-width="props.strokeWidth"></component>
            </n-icon>
        </template>
        {{ props.tTooltip ? $t(props.tTooltip) : props.tooltip }}
    </n-tooltip>
    <n-icon v-else :color="props.color" :size="props.size" class="icon-btn" @click="emit('click')">
        <component :is="props.icon" :stroke-width="props.strokeWidth"></component>
    </n-icon>
</template>

<style lang="scss"></style>
