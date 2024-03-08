<script setup>
import { computed, useSlots } from 'vue'
import { NIcon } from 'naive-ui'

const props = defineProps({
    tooltip: String,
    tTooltip: String,
    tooltipDelay: {
        type: Number,
        default: 800,
    },
    type: String,
    icon: [String, Object],
    size: {
        type: [Number, String],
        default: 20,
    },
    color: {
        type: String,
        default: '',
    },
    strokeWidth: {
        type: [Number, String],
        default: 3,
    },
    loading: Boolean,
    border: Boolean,
    disabled: Boolean,
    buttonStyle: [String, Object],
    buttonClass: [String, Object],
    small: Boolean,
    secondary: Boolean,
    tertiary: Boolean,
})

const emit = defineEmits(['click'])

const slots = useSlots()

const hasTooltip = computed(() => {
    return props.tooltip || props.tTooltip || slots.tooltip
})
</script>

<template>
    <n-tooltip v-if="hasTooltip" :delay="tooltipDelay" :show-arrow="false">
        <template #trigger>
            <n-button
                :class="props.buttonClass"
                :color="props.color"
                :disabled="props.disabled"
                :focusable="false"
                :loading="loading"
                :secondary="props.secondary"
                :size="props.small ? 'small' : ''"
                :style="props.buttonStyle"
                :tertiary="props.tertiary"
                :text="!props.border"
                :type="props.type"
                @click.prevent="emit('click')">
                <template #icon>
                    <slot>
                        <n-icon :color="props.color || 'currentColor'" :size="props.size">
                            <component :is="props.icon" :stroke-width="props.strokeWidth" />
                        </n-icon>
                    </slot>
                </template>
            </n-button>
        </template>
        <slot name="tooltip">
            {{ props.tTooltip ? $t(props.tTooltip) : props.tooltip }}
        </slot>
    </n-tooltip>
    <n-button
        v-else
        :class="props.buttonClass"
        :color="props.color"
        :disabled="props.disabled"
        :focusable="false"
        :loading="loading"
        :secondary="props.secondary"
        :size="props.small ? 'small' : ''"
        :style="props.buttonStyle"
        :tertiary="props.tertiary"
        :text="!props.border"
        :type="props.type"
        @click.prevent="emit('click')">
        <template #icon>
            <slot>
                <n-icon :color="props.color || 'currentColor'" :size="props.size">
                    <component :is="props.icon" :stroke-width="props.strokeWidth" />
                </n-icon>
            </slot>
        </template>
    </n-button>
</template>

<style lang="scss"></style>
