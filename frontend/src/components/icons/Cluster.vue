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
    <svg v-if="props.inverse" fill="none" viewBox="0 0 48 48" xmlns="http://www.w3.org/2000/svg">
        <rect
            :stroke-width="props.strokeWidth"
            fill="currentColor"
            height="8"
            stroke="currentColor"
            stroke-linecap="round"
            stroke-linejoin="round"
            width="8"
            x="4"
            y="34" />
        <rect
            :stroke-width="props.strokeWidth"
            fill="currentColor"
            height="12"
            stroke="currentColor"
            stroke-linecap="round"
            stroke-linejoin="round"
            width="32"
            x="8"
            y="6" />
        <path
            :stroke-width="props.strokeWidth"
            d="M24 34V18"
            stroke="currentColor"
            stroke-linecap="round"
            stroke-linejoin="round" />
        <path
            :stroke-width="props.strokeWidth"
            d="M8 34V26H40V34"
            stroke="currentColor"
            stroke-linecap="round"
            stroke-linejoin="round" />
        <rect
            :stroke-width="props.strokeWidth"
            fill="currentColor"
            height="8"
            stroke="currentColor"
            stroke-linecap="round"
            stroke-linejoin="round"
            width="8"
            x="36"
            y="34" />
        <rect
            :stroke-width="props.strokeWidth"
            fill="currentColor"
            height="8"
            stroke="currentColor"
            stroke-linecap="round"
            stroke-linejoin="round"
            width="8"
            x="20"
            y="34" />
        <path
            :stroke-width="props.strokeWidth"
            d="M14 12H16"
            stroke="#FFF"
            stroke-linecap="round"
            stroke-linejoin="round" />
    </svg>
    <svg v-else fill="none" viewBox="0 0 48 48" xmlns="http://www.w3.org/2000/svg">
        <rect
            :stroke-width="props.strokeWidth"
            fill="none"
            height="8"
            stroke="currentColor"
            stroke-linecap="round"
            stroke-linejoin="round"
            width="8"
            x="4"
            y="34" />
        <rect
            :stroke-width="props.strokeWidth"
            fill="none"
            height="12"
            stroke="currentColor"
            stroke-linecap="round"
            stroke-linejoin="round"
            width="32"
            x="8"
            y="6" />
        <path
            :stroke-width="props.strokeWidth"
            d="M24 34V18"
            stroke="currentColor"
            stroke-linecap="round"
            stroke-linejoin="round" />
        <path
            :stroke-width="props.strokeWidth"
            d="M8 34V26H40V34"
            stroke="currentColor"
            stroke-linecap="round"
            stroke-linejoin="round" />
        <rect
            :stroke-width="props.strokeWidth"
            fill="none"
            height="8"
            stroke="currentColor"
            stroke-linecap="round"
            stroke-linejoin="round"
            width="8"
            x="36"
            y="34" />
        <rect
            :stroke-width="props.strokeWidth"
            fill="none"
            height="8"
            stroke="currentColor"
            stroke-linecap="round"
            stroke-linejoin="round"
            width="8"
            x="20"
            y="34" />
        <path
            :stroke-width="props.strokeWidth"
            d="M14 12H16"
            stroke="currentColor"
            stroke-linecap="round"
            stroke-linejoin="round" />
    </svg>
</template>

<style lang="scss" scoped></style>
