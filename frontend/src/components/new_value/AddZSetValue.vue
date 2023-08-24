<script setup>
import { defineOptions, ref } from 'vue'
import { isEmpty, reject } from 'lodash'
import Add from '@/components/icons/Add.vue'
import Delete from '@/components/icons/Delete.vue'
import IconButton from '@/components/common/IconButton.vue'
import { useI18n } from 'vue-i18n'

const props = defineProps({
    type: Number,
    value: Object,
})
defineOptions({
    inheritAttrs: false,
})
const emit = defineEmits(['update:value', 'update:type'])

const i18n = useI18n()
const updateOption = [
    {
        value: 0,
        label: i18n.t('overwrite_field'),
    },
    {
        value: 1,
        label: i18n.t('ignore_field'),
    },
]

/**
 * @typedef ZSetItem
 * @property {string} value
 * @property {string} score
 */
const zset = ref([{ value: '', score: 0 }])
const onCreate = () => {
    return {
        value: '',
        score: 0,
    }
}
/**
 * update input items
 */
const onUpdate = () => {
    const val = reject(zset.value, (v) => v == null || isEmpty(v.value))
    const result = {}
    for (const elem of val) {
        result[elem.value] = elem.score
    }
    emit('update:value', result)
}
</script>

<template>
    <n-form-item :label="$t('type')">
        <n-radio-group :value="props.type" @update:value="(val) => emit('update:type', val)">
            <n-radio-button v-for="(op, i) in updateOption" :key="i" :label="op.label" :value="op.value" />
        </n-radio-group>
    </n-form-item>
    <n-form-item :label="$t('element')" required>
        <n-dynamic-input v-model:value="zset" @create="onCreate" @update:value="onUpdate">
            <template #default="{ value }">
                <n-input
                    v-model:value="value.value"
                    :placeholder="$t('enter_elem')"
                    type="text"
                    @update:value="onUpdate" />
                <n-input-number v-model:value="value.score" :placeholder="$t('enter_score')" @update:value="onUpdate" />
            </template>
            <template #action="{ index, create, remove, move }">
                <icon-button v-if="zset.length > 1" :icon="Delete" size="18" @click="() => remove(index)" />
                <icon-button :icon="Add" size="18" @click="() => create(index)" />
            </template>
        </n-dynamic-input>
    </n-form-item>
</template>

<style lang="scss"></style>
