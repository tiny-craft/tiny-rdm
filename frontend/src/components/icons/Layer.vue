<script setup>
const props = defineProps({
    strokeWidth: {
        type: [Number, String],
        default: 3,
    },
    fillColor: {
        type: String,
        default: '#f2c55c',
    },
})
</script>

<template>
    <svg
        stroke="currentColor"
        stroke-linecap="round"
        stroke-linejoin="round"
        stroke-width="1.5"
        viewBox="0 0 24 24"
        xmlns="http://www.w3.org/2000/svg">
        <path
            :fill="props.fillColor"
            d="M10 20H4a2 2 0 0 1-2-2V5c0-1.1.9-2 2-2h3.93a2 2 0 0 1 1.66.9l.82 1.2a2 2 0 0 0 1.66.9H20a2 2 0 0 1 2 2v2" />
        <circle :fill="props.fillColor" cx="16" cy="20" r="2" />
        <path d="m22 14-4.5 4.5" />
        <path d="m21 15 1 1" />
    </svg>
</template>

<style lang="scss" scoped></style>
