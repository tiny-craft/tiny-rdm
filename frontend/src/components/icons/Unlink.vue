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
            d="M37 22.0001L34 25.0001L23 14.0001L26 11.0001C27.5 9.50002 33 7.00005 37 11.0001C41 15.0001 38.5 20.5 37 22.0001Z"
            fill="none"
            stroke="currentColor"
            stroke-linecap="round"
            stroke-linejoin="round"
        />
        <path
            :stroke-width="strokeWidth"
            d="M42 6L37 11"
            stroke="currentColor"
            stroke-linecap="round"
            stroke-linejoin="round"
        />
        <path
            :stroke-width="strokeWidth"
            d="M11 25.9999L14 22.9999L25 33.9999L22 36.9999C20.5 38.5 15 41 11 36.9999C7 32.9999 9.5 27.5 11 25.9999Z"
            fill="none"
            stroke="currentColor"
            stroke-linecap="round"
            stroke-linejoin="round"
        />
        <path
            :stroke-width="strokeWidth"
            d="M23 32L27 28"
            stroke="currentColor"
            stroke-linecap="round"
            stroke-linejoin="round"
        />
        <path
            :stroke-width="strokeWidth"
            d="M6 42L11 37"
            stroke="currentColor"
            stroke-linecap="round"
            stroke-linejoin="round"
        />
        <path
            :stroke-width="strokeWidth"
            d="M16 25L20 21"
            stroke="currentColor"
            stroke-linecap="round"
            stroke-linejoin="round"
        />
    </svg>
</template>

<style lang="scss" scoped></style>
