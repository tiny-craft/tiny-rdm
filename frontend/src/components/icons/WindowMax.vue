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
        <rect width="9" height="9" x="1.5" y="1.5" fill="none" stroke="currentColor"></rect>
    </svg>
</template>

<style lang="scss" scoped></style>
