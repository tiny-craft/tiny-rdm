<script setup>
const emit = defineEmits(['update:modelValue'])

const props = defineProps({
    inverse: {
        type: Boolean,
        default: false,
    },
    strokeWidth: {
        type: [Number, String],
        default: 3,
    },
})
</script>

<template>
    <svg fill="none" viewBox="0 0 48 48" xmlns="http://www.w3.org/2000/svg">
        <path
            :fill="props.inverse ? 'currentColor' : 'none'"
            :stroke-width="props.strokeWidth"
            d="M41 4H7C5.34315 4 4 5.34315 4 7V41C4 42.6569 5.34315 44 7 44H41C42.6569 44 44 42.6569 44 41V7C44 5.34315 42.6569 4 41 4Z"
            stroke="currentColor"
            stroke-linecap="round"
            stroke-linejoin="round" />
        <path
            :stroke="props.inverse ? '#FFF' : 'currentColor'"
            :stroke-width="props.strokeWidth"
            d="M4 32H44"
            stroke-linecap="round" />
        <path
            :stroke="props.inverse ? '#FFF' : 'currentColor'"
            :stroke-width="props.strokeWidth"
            d="M10 38H11"
            stroke-linecap="round"
            stroke-linejoin="round" />
        <path
            :stroke="props.inverse ? '#FFF' : 'currentColor'"
            :stroke-width="props.strokeWidth"
            d="M26 38H38"
            stroke-linecap="round"
            stroke-linejoin="round" />
        <path
            :stroke-width="props.strokeWidth"
            d="M44 37V27"
            stroke="currentColor"
            stroke-linecap="round"
            stroke-linejoin="round" />
        <path
            :stroke-width="props.strokeWidth"
            d="M4 37V27"
            stroke="currentColor"
            stroke-linecap="round"
            stroke-linejoin="round" />
    </svg>
</template>

<style lang="scss" scoped></style>
