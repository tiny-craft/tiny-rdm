<script setup>
import AddLink from '@/components/icons/AddLink.vue'
import useDialogStore from 'stores/dialog.js'

const dialogStore = useDialogStore()
</script>

<template>
    <div class="content-container flex-box-v">
        <!-- TODO: replace icon to app icon -->
        <n-empty :description="$t('empty_server_content')">
            <template #extra>
                <n-button :focusable="false" @click="dialogStore.openNewDialog()">
                    <template #icon>
                        <n-icon :component="AddLink" size="18" />
                    </template>
                    {{ $t('new_conn') }}
                </n-button>
            </template>
        </n-empty>
    </div>
</template>

<style lang="scss" scoped>
@import '@/styles/content';

.content-container {
    justify-content: center;
    padding: 5px;
    box-sizing: border-box;
}

.color-preset-item {
    width: 24px;
    height: 24px;
    margin-right: 2px;
    border: white 3px solid;
    cursor: pointer;

    &_selected,
    &:hover {
        border-color: #cdd0d6;
    }
}
</style>
