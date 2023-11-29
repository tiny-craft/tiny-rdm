<script setup>
import { computed } from 'vue'
import { NIcon } from 'naive-ui'

const emit = defineEmits(['click'])

const props = defineProps({
    tooltip: String,
    tTooltip: String,
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
})

const hasTooltip = computed(() => {
    return props.tooltip || props.tTooltip
})
</script>

<template>
    <n-tooltip v-if="hasTooltip" :show-arrow="false">
        <template #trigger>
            <n-button
                :class="props.buttonClass"
                :color="props.color"
                :disabled="disabled"
                :focusable="false"
                :loading="loading"
                :style="props.buttonStyle"
                :text="!border"
                :type="type"
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
        {{ props.tTooltip ? $t(props.tTooltip) : props.tooltip }}
    </n-tooltip>
    <n-button
        v-else
        :class="props.buttonClass"
        :color="props.color"
        :disabled="disabled"
        :focusable="false"
        :loading="loading"
        :style="props.buttonStyle"
        :text="!border"
        :type="type"
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
