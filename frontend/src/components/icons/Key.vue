<script setup>
const props = defineProps({
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
        fill="none"
        stroke="currentColor"
        stroke-width="1.4"
        stroke-linecap="round"
        stroke-linejoin="round">
        <path :fill="props.fillColor" d="M2 18v3c0 .6.4 1 1 1h4v-3h3v-3h2l1.4-1.4a6.5 6.5 0 1 0-4-4Z" />
        <circle cx="16.5" cy="7.5" r="1.5" />
    </svg>
</template>

<style lang="scss" scoped></style>
