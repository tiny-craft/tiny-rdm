<script setup>
const props = defineProps({
    size: {
        type: [Number, String],
        default: 12,
    },
})
</script>

<template>
    <svg aria-hidden="false" :width="props.size" :height="props.size" viewBox="0 0 12 12">
        <rect fill="currentColor" width="10" height="1" x="1" y="6"></rect>
    </svg>
</template>

<style lang="scss" scoped></style>
