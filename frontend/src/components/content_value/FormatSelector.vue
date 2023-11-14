<script setup>
import { decodeTypes, formatTypes } from '@/consts/value_view_type.js'
import Code from '@/components/icons/Code.vue'
import Conversion from '@/components/icons/Conversion.vue'
import DropdownSelector from '@/components/common/DropdownSelector.vue'
import { some } from 'lodash'

const props = defineProps({
    decode: {
        type: String,
        default: decodeTypes.NONE,
    },
    format: {
        type: String,
        default: formatTypes.RAW,
    },
    disabled: Boolean,
})

const emit = defineEmits(['formatChanged', 'update:decode', 'update:format'])
const onFormatChanged = (selDecode, selFormat) => {
    if (!some(decodeTypes, (val) => val === selDecode)) {
        selDecode = decodeTypes.NONE
    }
    if (!some(formatTypes, (val) => val === selFormat)) {
        selFormat = formatTypes.RAW
    }
    emit('formatChanged', selDecode, selFormat)
    if (selDecode !== props.decode) {
        emit('update:decode', selDecode)
    }
    if (selFormat !== props.format) {
        emit('update:format', selFormat)
    }
}
</script>

<template>
    <n-space :size="0" :wrap="false" :wrap-item="false" align="center" justify="start">
        <dropdown-selector
            :default="formatTypes.RAW"
            :disabled="props.disabled"
            :icon="Code"
            :options="formatTypes"
            :tooltip="$t('interface.view_as')"
            :value="props.format"
            @update:value="(f) => onFormatChanged(props.decode, f)" />
        <n-divider vertical />
        <dropdown-selector
            :default="decodeTypes.NONE"
            :disabled="props.disabled"
            :icon="Conversion"
            :options="decodeTypes"
            :tooltip="$t('interface.decode_with')"
            :value="props.decode"
            @update:value="(d) => onFormatChanged(d, props.format)" />
    </n-space>
</template>

<style lang="scss" scoped></style>
