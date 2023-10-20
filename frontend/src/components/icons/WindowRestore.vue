<script setup>
const props = defineProps({
    size: {
        type: [Number, String],
        default: 14,
    },
})
</script>

<template>
    <svg :height="props.size" :width="props.size" fill="none" viewBox="0 0 48 48">
        <path
            d="M13 12.4316V7.8125C13 6.2592 14.2592 5 15.8125 5H40.1875C41.7408 5 43 6.2592 43 7.8125V32.1875C43 33.7408 41.7408 35 40.1875 35H35.5163"
            stroke="currentColor"
            stroke-linecap="round"
            stroke-linejoin="round"
            stroke-width="4" />
        <path
            d="M32.1875 13H7.8125C6.2592 13 5 14.2592 5 15.8125V40.1875C5 41.7408 6.2592 43 7.8125 43H32.1875C33.7408 43 35 41.7408 35 40.1875V15.8125C35 14.2592 33.7408 13 32.1875 13Z"
            fill="none"
            stroke="currentColor"
            stroke-linejoin="round"
            stroke-width="4" />
    </svg>
</template>

<style lang="scss" scoped></style>
