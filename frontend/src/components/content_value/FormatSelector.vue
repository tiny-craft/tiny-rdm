<script setup>
import { decodeTypes, formatTypes } from '@/consts/value_view_type.js'
import Code from '@/components/icons/Code.vue'
import Conversion from '@/components/icons/Conversion.vue'
import DropdownSelector from '@/components/common/DropdownSelector.vue'
import { isEmpty, map, some } from 'lodash'
import { computed } from 'vue'
import usePreferencesStore from 'stores/preferences.js'

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

const prefStore = usePreferencesStore()

const formatTypeOption = computed(() => {
    return map(formatTypes, (t) => t)
})

const decodeTypeOption = computed(() => {
    const customTypes = []
    // has custom decoder
    if (!isEmpty(prefStore.decoder)) {
        for (const decoder of prefStore.decoder) {
            // types[decoder.name] = types[decoder.name] || decoder.name
            if (!decodeTypes.hasOwnProperty(decoder.name)) {
                customTypes.push(decoder.name)
            }
        }
    }
    return [map(decodeTypes, (t) => t), customTypes]
})

const emit = defineEmits(['formatChanged', 'update:decode', 'update:format'])
const onFormatChanged = (selDecode, selFormat) => {
    const [buildin, external] = decodeTypeOption.value
    if (!some([...buildin, ...external], (val) => val === selDecode)) {
        selDecode = decodeTypes.NONE
    }
    if (!some(formatTypes, (val) => val === selFormat)) {
        // set to auto chose format
        selFormat = ''
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
            :options="formatTypeOption"
            :tooltip="$t('interface.view_as')"
            :value="props.format"
            @update:value="(f) => onFormatChanged(props.decode, f)" />
        <n-divider vertical />
        <dropdown-selector
            :default="decodeTypes.NONE"
            :disabled="props.disabled"
            :icon="Conversion"
            :options="decodeTypeOption"
            :tooltip="$t('interface.decode_with')"
            :value="props.decode"
            @update:value="(d) => onFormatChanged(d, '')" />
    </n-space>
</template>

<style lang="scss" scoped></style>
