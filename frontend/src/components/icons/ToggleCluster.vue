<script setup>
const emit = defineEmits(['update:modelValue'])

const props = defineProps({
    modelValue: {
        type: Boolean,
        default: false,
    },
    strokeWidth: {
        type: [Number, String],
        default: 3,
    },
    fillColor: {
        type: String,
        default: '#dc423c',
    },
})
</script>

<template>
    <svg v-if="props.modelValue" viewBox="0 0 48 48" fill="none" xmlns="http://www.w3.org/2000/svg">
        <rect
            x="4"
            y="34"
            width="8"
            height="8"
            :fill="props.fillColor"
            :stroke="props.fillColor"
            :stroke-width="props.strokeWidth"
            stroke-linecap="round"
            stroke-linejoin="round" />
        <rect
            x="8"
            y="6"
            width="32"
            height="12"
            :fill="props.fillColor"
            :stroke="props.fillColor"
            :stroke-width="props.strokeWidth"
            stroke-linecap="round"
            stroke-linejoin="round" />
        <path
            d="M24 34V18"
            :stroke="props.fillColor"
            :stroke-width="props.strokeWidth"
            stroke-linecap="round"
            stroke-linejoin="round" />
        <path
            d="M8 34V26H40V34"
            :stroke="props.fillColor"
            :stroke-width="props.strokeWidth"
            stroke-linecap="round"
            stroke-linejoin="round" />
        <rect
            x="36"
            y="34"
            width="8"
            height="8"
            :fill="props.fillColor"
            :stroke="props.fillColor"
            :stroke-width="props.strokeWidth"
            stroke-linecap="round"
            stroke-linejoin="round" />
        <rect
            x="20"
            y="34"
            width="8"
            height="8"
            :fill="props.fillColor"
            :stroke="props.fillColor"
            :stroke-width="props.strokeWidth"
            stroke-linecap="round"
            stroke-linejoin="round" />
        <path
            d="M14 12H16"
            stroke="#FFF"
            :stroke-width="props.strokeWidth"
            stroke-linecap="round"
            stroke-linejoin="round" />
    </svg>
    <svg v-else viewBox="0 0 48 48" fill="none" xmlns="http://www.w3.org/2000/svg">
        <path
            d="M5 24L43 24"
            stroke="currentColor"
            :stroke-width="props.strokeWidth"
            stroke-linecap="round"
            stroke-linejoin="round" />
        <path
            d="M28 4H20C18.8954 4 18 4.89543 18 6V14C18 15.1046 18.8954 16 20 16H28C29.1046 16 30 15.1046 30 14V6C30 4.89543 29.1046 4 28 4Z"
            fill="none"
            stroke="currentColor"
            :stroke-width="props.strokeWidth"
            stroke-linejoin="round" />
        <path
            d="M16 32H8C6.89543 32 6 32.8954 6 34V42C6 43.1046 6.89543 44 8 44H16C17.1046 44 18 43.1046 18 42V34C18 32.8954 17.1046 32 16 32Z"
            fill="none"
            stroke="currentColor"
            :stroke-width="props.strokeWidth"
            stroke-linejoin="round" />
        <path
            d="M40 32H32C30.8954 32 30 32.8954 30 34V42C30 43.1046 30.8954 44 32 44H40C41.1046 44 42 43.1046 42 42V34C42 32.8954 41.1046 32 40 32Z"
            fill="none"
            stroke="currentColor"
            :stroke-width="props.strokeWidth"
            stroke-linejoin="round" />
        <path
            d="M24 24V16"
            stroke="currentColor"
            :stroke-width="props.strokeWidth"
            stroke-linecap="round"
            stroke-linejoin="round" />
        <path
            d="M36 32V24"
            stroke="currentColor"
            :stroke-width="props.strokeWidth"
            stroke-linecap="round"
            stroke-linejoin="round" />
        <path
            d="M12 32V24"
            stroke="currentColor"
            :stroke-width="props.strokeWidth"
            stroke-linecap="round"
            stroke-linejoin="round" />
    </svg>
    <svg v-else viewBox="0 0 48 48" fill="none" xmlns="http://www.w3.org/2000/svg">
        <rect
            x="4"
            y="34"
            width="8"
            height="8"
            fill="none"
            stroke="currentColor"
            :stroke-width="props.strokeWidth"
            stroke-linecap="round"
            stroke-linejoin="round" />
        <rect
            x="8"
            y="6"
            width="32"
            height="12"
            fill="none"
            stroke="currentColor"
            :stroke-width="props.strokeWidth"
            stroke-linecap="round"
            stroke-linejoin="round" />
        <path
            d="M24 34V18"
            stroke="currentColor"
            :stroke-width="props.strokeWidth"
            stroke-linecap="round"
            stroke-linejoin="round" />
        <path
            d="M8 34V26H40V34"
            stroke="currentColor"
            :stroke-width="props.strokeWidth"
            stroke-linecap="round"
            stroke-linejoin="round" />
        <rect
            x="36"
            y="34"
            width="8"
            height="8"
            fill="none"
            stroke="currentColor"
            :stroke-width="props.strokeWidth"
            stroke-linecap="round"
            stroke-linejoin="round" />
        <rect
            x="20"
            y="34"
            width="8"
            height="8"
            fill="none"
            stroke="currentColor"
            :stroke-width="props.strokeWidth"
            stroke-linecap="round"
            stroke-linejoin="round" />
        <path
            d="M14 12H16"
            stroke="currentColor"
            :stroke-width="props.strokeWidth"
            stroke-linecap="round"
            stroke-linejoin="round" />
    </svg>
</template>

<style lang="scss" scoped></style>
