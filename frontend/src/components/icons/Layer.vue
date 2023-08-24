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
        xmlns="http://www.w3.org/2000/svg"
        viewBox="0 0 24 24"
        stroke="currentColor"
        stroke-width="1.5"
        stroke-linecap="round"
        stroke-linejoin="round">
        <path
            d="M10 20H4a2 2 0 0 1-2-2V5c0-1.1.9-2 2-2h3.93a2 2 0 0 1 1.66.9l.82 1.2a2 2 0 0 0 1.66.9H20a2 2 0 0 1 2 2v2"
            :fill="props.fillColor" />
        <circle cx="16" cy="20" r="2" :fill="props.fillColor" />
        <path d="m22 14-4.5 4.5" />
        <path d="m21 15 1 1" />
    </svg>
</template>

<style lang="scss" scoped></style>
