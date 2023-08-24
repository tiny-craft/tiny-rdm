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
            :stroke-width="props.strokeWidth"
            d="M6 9L20.4 25.8178V38.4444L27.6 42V25.8178L42 9H6Z"
            fill="none"
            stroke="currentColor"
            stroke-linejoin="round" />
    </svg>
</template>
