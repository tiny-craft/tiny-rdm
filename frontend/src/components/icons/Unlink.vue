<script setup>
const props = defineProps({
    strokeWidth: {
        type: [Number, String],
        default: 3,
    },
})
</script>

<template>
    <svg fill="none" viewBox="0 0 48 48" xmlns="http://www.w3.org/2000/svg">
        <path
            :stroke-width="strokeWidth"
            d="M25.8927 16.0307L18.1145 8.2525C15.2506 5.38866 10.7031 5.29302 7.9572 8.0389C5.21132 10.7848 5.30696 15.3323 8.1708 18.1962L15.949 25.9744"
            stroke="currentColor"
            stroke-linecap="round"
            stroke-linejoin="round" />
        <path
            :stroke-width="strokeWidth"
            d="M31.9161 22.0707L39.6943 29.8489C42.5581 32.7127 42.9291 37.1233 39.9079 40.0062C36.8867 42.8891 32.6144 42.6564 29.7506 39.7926L21.9724 32.0144"
            stroke="currentColor"
            stroke-linecap="round"
            stroke-linejoin="round" />
        <path
            :stroke-width="strokeWidth"
            d="M21.2384 21.0759L17.3493 17.1868"
            stroke="currentColor"
            stroke-linecap="round"
            stroke-linejoin="round" />
        <path
            :stroke-width="strokeWidth"
            d="M30.3131 30.1504L26.424 26.2613"
            stroke="currentColor"
            stroke-linecap="round"
            stroke-linejoin="round" />
    </svg>
</template>

<style lang="scss" scoped></style>
