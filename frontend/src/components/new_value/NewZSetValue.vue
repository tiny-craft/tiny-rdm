<script setup>
import { ref } from 'vue'
import { flatMap, isEmpty, reject } from 'lodash'
import Add from '../icons/Add.vue'
import Delete from '../icons/Delete.vue'
import IconButton from '../common/IconButton.vue'

const props = defineProps({
    value: Array,
})
const emit = defineEmits(['update:value'])

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
    emit(
        'update:value',
        flatMap(val, (item) => [item.value, item.score.toString()]),
    )
}

defineExpose({
    validate: () => {
        return !isEmpty(props.value)
    },
})
</script>

<template>
    <n-form-item :label="$t('element')" required>
        <n-dynamic-input v-model:value="zset" @create="onCreate" @update:value="onUpdate">
            <template #default="{ value }">
                <n-input
                    v-model:value="value.value"
                    :placeholder="$t('enter_member')"
                    type="text"
                    @update:value="onUpdate"
                />
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
