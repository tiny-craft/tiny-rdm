<script setup>
const props = defineProps({
    inverse: {
        type: Boolean,
        default: false,
    },
    strokeWidth: {
        type: [Number, String],
        default: 3,
    },
    strokeColor: {
        type: String,
        default: '#FFF',
    },
})
</script>

<template>
    <svg v-if="props.inverse" fill="none" viewBox="0 0 48 48" xmlns="http://www.w3.org/2000/svg">
        <path
            d="M42 8H6C4.89543 8 4 8.89543 4 10V38C4 39.1046 4.89543 40 6 40H42C43.1046 40 44 39.1046 44 38V10C44 8.89543 43.1046 8 42 8Z"
            fill="currentColor"
            stroke="currentColor"
            stroke-width="3" />
        <path :stroke="props.strokeColor" :stroke-width="props.strokeWidth" d="M24 17V31" stroke-linecap="round" />
        <path :stroke="props.strokeColor" :stroke-width="props.strokeWidth" d="M32 24V31" stroke-linecap="round" />
        <path :stroke="props.strokeColor" :stroke-width="props.strokeWidth" d="M16 22V31" stroke-linecap="round" />
    </svg>
    <svg v-else fill="none" viewBox="0 0 48 48" xmlns="http://www.w3.org/2000/svg">
        <rect
            :stroke-width="props.strokeWidth"
            fill="none"
            height="36"
            rx="2"
            stroke="currentColor"
            stroke-linecap="round"
            stroke-linejoin="round"
            width="40"
            x="4"
            y="6" />
        <path
            :stroke-width="props.strokeWidth"
            d="M32 25V32"
            stroke="currentColor"
            stroke-linecap="round"
            stroke-linejoin="round" />
        <path
            :stroke-width="props.strokeWidth"
            d="M24 16V32"
            stroke="currentColor"
            stroke-linecap="round"
            stroke-linejoin="round" />
        <path
            :stroke-width="props.strokeWidth"
            d="M16 20V32"
            stroke="currentColor"
            stroke-linecap="round"
            stroke-linejoin="round" />
    </svg>
</template>

<style lang="scss" scoped></style>
