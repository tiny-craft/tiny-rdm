<script setup>
const props = defineProps({
    title: String,
})
</script>

<template>
    <div class="db-tree-item flex-box-v">
        <div class="tree-item-title">{{ title }}</div>
        <div class="tree-item-addr">localhost:3306</div>
    </div>
</template>

<style lang="scss">
.db-tree-item {
    padding: 0 5px;

    .tree-item-title {
        font-size: 16px;
    }

    .tree-item-addr {
        font-size: 12px;
    }
}
</style>
