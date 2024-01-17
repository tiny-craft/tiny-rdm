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
        <g clip-path="url(#icon-25e88d94353e4f38)">
            <path
                :stroke-width="props.strokeWidth"
                d="M42 20V39C42 40.6569 40.6569 42 39 42H9C7.34315 42 6 40.6569 6 39V9C6 7.34315 7.34315 6 9 6H30"
                stroke="currentColor"
                stroke-linecap="round"
                stroke-linejoin="round" />
            <path
                :stroke-width="props.strokeWidth"
                d="M16 20L26 28L41 7"
                stroke="currentColor"
                stroke-linecap="round"
                stroke-linejoin="round" />
        </g>
        <defs>
            <clipPath id="icon-25e88d94353e4f38">
                <rect fill="currentColor" height="48" width="48" />
            </clipPath>
        </defs>
    </svg>
</template>

<style lang="scss" scoped></style>
