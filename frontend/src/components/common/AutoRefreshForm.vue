<script setup>
import { isNumber } from 'lodash'

const props = defineProps({
    loading: {
        type: Boolean,
        default: false,
    },
    on: {
        type: Boolean,
        default: false,
    },
    defaultValue: {
        type: Number,
        default: 2,
    },
    interval: {
        type: Number,
        default: 2,
    },
    onRefresh: {
        type: Function,
        default: () => {},
    },
})

const emit = defineEmits(['toggle', 'update:on', 'update:interval'])

const onToggle = (on) => {
    emit('update:on', on === true)
    if (on) {
        let interval = props.interval
        if (!isNumber(interval)) {
            interval = props.defaultValue
        }
        interval = Math.max(1, interval)
        emit('update:interval', interval)
        emit('toggle', true)
    } else {
        emit('toggle', false)
    }
}
</script>

<template>
    <n-form :show-feedback="false" label-align="right" label-placement="left" label-width="auto" size="small">
        <n-form-item :label="$t('interface.auto_refresh')">
            <n-switch :loading="props.loading" :value="props.on" @update:value="onToggle" />
        </n-form-item>
        <n-form-item :label="$t('interface.refresh_interval')">
            <n-input-number
                :autofocus="false"
                :default-value="props.defaultValue"
                :disabled="props.on"
                :max="9999"
                :min="1"
                :show-button="false"
                :value="props.interval"
                style="max-width: 100px"
                @update:value="(val) => emit('update:interval', val)">
                <template #suffix>
                    {{ $t('common.unit_second') }}
                </template>
            </n-input-number>
        </n-form-item>
    </n-form>
</template>

<style lang="scss" scoped></style>
