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
            v-for="(icon, i) in props.icons"
            :key="i"
            :disabled="!(props.tTooltips && props.tTooltips[i])"
            :show-arrow="false">
            <template #trigger>
                <n-button :focusable="false" :size="props.size" :tertiary="i !== props.value" @click="handleSwitch(i)">
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

<style lang="scss" scoped></style>
