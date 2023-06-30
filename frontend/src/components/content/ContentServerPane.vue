<script setup>
import useDialog from '../../stores/dialog.js'
import AddLink from '../icons/AddLink.vue'

const dialogStore = useDialog()
</script>

<template>
    <div class="content-container flex-box-v">
        <!-- TODO: replace icon to app icon -->
        <n-empty :description="$t('empty_server_content')">
            <template #extra>
                <n-button @click="dialogStore.openNewDialog()">
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
@import 'content';

.content-container {
    justify-content: center;
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
