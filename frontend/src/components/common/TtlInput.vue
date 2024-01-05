<script setup>
import { useI18n } from 'vue-i18n'
import { computed } from 'vue'

const props = defineProps({
    value: {
        type: Number,
        default: -1,
    },
    unit: {
        type: Number,
        default: 1,
    },
})

const emit = defineEmits(['update:value', 'update:unit'])

const i18n = useI18n()
const unit = computed(() => [
    { value: 1, label: i18n.t('common.second') },
    {
        value: 60,
        label: i18n.t('common.minute'),
    },
    {
        value: 3600,
        label: i18n.t('common.hour'),
    },
    {
        value: 86400,
        label: i18n.t('common.day'),
    },
])

const unitValue = computed(() => {
    switch (props.unit) {
        case 60:
            return 60
        case 3600:
            return 3600
        case 86400:
            return 86400
        default:
            return 1
    }
})
</script>

<template>
    <n-input-group>
        <n-input-number
            :max="Number.MAX_SAFE_INTEGER"
            :min="-1"
            :show-button="false"
            :value="props.value"
            class="flex-item-expand"
            @update:value="(val) => emit('update:value', val)" />
        <n-select
            :options="unit"
            :value="unitValue"
            style="max-width: 150px"
            @update:value="(val) => emit('update:unit', val)" />
    </n-input-group>
</template>

<style lang="scss" scoped></style>
