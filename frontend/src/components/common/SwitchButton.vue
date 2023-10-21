<script setup>
import { NIcon } from 'naive-ui'

const props = defineProps({
    value: {
        type: Number,
        default: 0,
    },
    size: {
        type: String,
        default: 'small',
    },
    icons: Array,
    tTooltips: Array,
    iconSize: {
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
    unselectStrokeWidth: {
        type: [Number, String],
        default: 3,
    },
})

const emit = defineEmits(['update:value'])

const handleSwitch = (idx) => {
    if (idx !== props.value) {
        emit('update:value', idx)
    }
}
</script>

<template>
    <n-button-group>
        <n-tooltip
            :show-arrow="false"
            v-for="(icon, i) in props.icons"
            :key="i"
            :disabled="!(props.tTooltips && props.tTooltips[i])">
            <template #trigger>
                <n-button :tertiary="i !== props.value" :focusable="false" :size="props.size" @click="handleSwitch(i)">
                    <template #icon>
                        <n-icon :size="props.iconSize">
                            <component
                                :is="icon"
                                :stroke-width="i !== props.value ? props.unselectStrokeWidth : props.strokeWidth" />
                        </n-icon>
                    </template>
                </n-button>
            </template>
            {{ props.tTooltips ? $t(props.tTooltips[i]) : '' }}
        </n-tooltip>
    </n-button-group>
</template>

<style scoped lang="scss"></style>
