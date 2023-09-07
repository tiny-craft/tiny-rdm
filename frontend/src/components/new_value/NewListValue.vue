<script setup>
import { ref } from 'vue'
import { compact, isEmpty } from 'lodash'
import Add from '@/components/icons/Add.vue'
import Delete from '@/components/icons/Delete.vue'
import IconButton from '@/components/common/IconButton.vue'

const props = defineProps({
    value: Array,
})
const emit = defineEmits(['update:value'])

const list = ref([''])
const onUpdate = (val) => {
    val = compact(val)
    emit('update:value', val)
}

defineExpose({
    validate: () => {
        return !isEmpty(props.value)
    },
})
</script>

<template>
    <n-form-item :label="$t('dialogue.field.element')" required>
        <n-dynamic-input v-model:value="list" :placeholder="$t('dialogue.field.enter_elem')" @update:value="onUpdate">
            <template #action="{ index, create, remove, move }">
                <icon-button v-if="list.length > 1" :icon="Delete" size="18" @click="() => remove(index)" />
                <icon-button :icon="Add" size="18" @click="() => create(index)" />
            </template>
        </n-dynamic-input>
    </n-form-item>
</template>

<style lang="scss" scoped></style>
