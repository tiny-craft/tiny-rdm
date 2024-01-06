<script setup>
import IconButton from './IconButton.vue'
import Delete from '@/components/icons/Delete.vue'
import Edit from '@/components/icons/Edit.vue'
import Close from '@/components/icons/Close.vue'
import Save from '@/components/icons/Save.vue'
import Copy from '@/components/icons/Copy.vue'

const props = defineProps({
    bindKey: String,
    editing: Boolean,
    readonly: Boolean,
})

const emit = defineEmits(['edit', 'delete', 'copy', 'save', 'cancel'])
</script>

<template>
    <div v-if="props.editing" class="flex-box-h edit-column-func">
        <icon-button :icon="Save" @click="emit('save')" />
        <icon-button :icon="Close" @click="emit('cancel')" />
    </div>
    <div v-else class="flex-box-h edit-column-func">
        <icon-button :icon="Copy" :title="$t('interface.copy_value')" @click="emit('copy')" />
        <icon-button v-if="!props.readonly" :icon="Edit" :title="$t('interface.edit_row')" @click="emit('edit')" />
        <n-popconfirm
            :negative-text="$t('common.cancel')"
            :positive-text="$t('common.confirm')"
            @positive-click="emit('delete')">
            <template #trigger>
                <icon-button :icon="Delete" :title="$t('interface.delete_row')" />
            </template>
            {{ $t('dialogue.remove_tip', { name: props.bindKey }) }}
        </n-popconfirm>
    </div>
</template>

<style lang="scss">
.edit-column-func {
    align-items: center;
    justify-content: center;
    gap: 10px;
}
</style>
