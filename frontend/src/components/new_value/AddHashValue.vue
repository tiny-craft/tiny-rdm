<script setup>
import { ref } from 'vue'
import { flatMap, reject } from 'lodash'
import Add from '@/components/icons/Add.vue'
import Delete from '@/components/icons/Delete.vue'
import IconButton from '@/components/common/IconButton.vue'
import { useI18n } from 'vue-i18n'

const props = defineProps({
    type: Number,
    value: Array,
})
const emit = defineEmits(['update:value', 'update:type'])

const i18n = useI18n()
const updateOption = [
    {
        value: 0,
        label: i18n.t('dialogue.field.overwrite_field'),
    },
    {
        value: 1,
        label: i18n.t('dialogue.field.ignore_field'),
    },
]

/**
 * @typedef Hash
 * @property {string} key
 * @property {string} [value]
 */
const kvList = ref([{ key: '', value: '' }])

/**
 *
 * @param {Hash[]} val
 */
const onUpdate = (val) => {
    val = reject(val, { key: '' })
    emit(
        'update:value',
        flatMap(val, (item) => [item.key, item.value]),
    )
}
</script>

<template>
    <n-form-item :label="$t('dialogue.field.conflict_handle')">
        <n-radio-group :value="props.type" @update:value="(val) => emit('update:type', val)">
            <n-radio-button v-for="(op, i) in updateOption" :key="i" :label="op.label" :value="op.value" />
        </n-radio-group>
    </n-form-item>
    <n-form-item
        :label="$t('dialogue.field.element') + ' (' + $t('common.field') + ':' + $t('common.value') + ')'"
        required>
        <n-dynamic-input
            v-model:value="kvList"
            :key-placeholder="$t('dialogue.field.enter_field')"
            :value-placeholder="$t('dialogue.field.enter_value')"
            preset="pair"
            @update:value="onUpdate">
            <template #action="{ index, create, remove, move }">
                <icon-button v-if="kvList.length > 1" :icon="Delete" size="18" @click="() => remove(index)" />
                <icon-button :icon="Add" size="18" @click="() => create(index)" />
            </template>
        </n-dynamic-input>
    </n-form-item>
</template>

<style lang="scss" scoped></style>
