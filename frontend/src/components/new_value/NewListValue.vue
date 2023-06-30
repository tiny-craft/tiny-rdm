<script setup>
import { ref } from 'vue'
import { compact } from 'lodash'
import Add from '../icons/Add.vue'
import Delete from '../icons/Delete.vue'
import IconButton from '../common/IconButton.vue'

const props = defineProps({
    value: Array,
})
const emit = defineEmits(['update:value'])

const list = ref([''])
const onUpdate = (val) => {
    val = compact(val)
    emit('update:value', val)
}
</script>

<template>
    <n-form-item :label="$t('element')" required>
        <n-dynamic-input v-model:value="list" :placeholder="$t('enter_elem')" @update:value="onUpdate">
            <template #action="{ index, create, remove, move }">
                <icon-button v-if="list.length > 1" :icon="Delete" size="18" @click="() => remove(index)" />
                <icon-button :icon="Add" size="18" @click="() => create(index)" />
            </template>
        </n-dynamic-input>
    </n-form-item>
</template>

<style lang="scss" scoped></style>
