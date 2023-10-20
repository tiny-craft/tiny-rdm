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
    <svg v-if="props.modelValue" fill="none" viewBox="0 0 48 48" xmlns="http://www.w3.org/2000/svg">
        <rect
            :fill="props.fillColor"
            :stroke="props.fillColor"
            :stroke-width="props.strokeWidth"
            height="8"
            stroke-linecap="round"
            stroke-linejoin="round"
            width="8"
            x="4"
            y="34" />
        <rect
            :fill="props.fillColor"
            :stroke="props.fillColor"
            :stroke-width="props.strokeWidth"
            height="12"
            stroke-linecap="round"
            stroke-linejoin="round"
            width="32"
            x="8"
            y="6" />
        <path
            :stroke="props.fillColor"
            :stroke-width="props.strokeWidth"
            d="M24 34V18"
            stroke-linecap="round"
            stroke-linejoin="round" />
        <path
            :stroke="props.fillColor"
            :stroke-width="props.strokeWidth"
            d="M8 34V26H40V34"
            stroke-linecap="round"
            stroke-linejoin="round" />
        <rect
            :fill="props.fillColor"
            :stroke="props.fillColor"
            :stroke-width="props.strokeWidth"
            height="8"
            stroke-linecap="round"
            stroke-linejoin="round"
            width="8"
            x="36"
            y="34" />
        <rect
            :fill="props.fillColor"
            :stroke="props.fillColor"
            :stroke-width="props.strokeWidth"
            height="8"
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
        <path
            :stroke-width="props.strokeWidth"
            d="M5 24L43 24"
            stroke="currentColor"
            stroke-linecap="round"
            stroke-linejoin="round" />
        <path
            :stroke-width="props.strokeWidth"
            d="M28 4H20C18.8954 4 18 4.89543 18 6V14C18 15.1046 18.8954 16 20 16H28C29.1046 16 30 15.1046 30 14V6C30 4.89543 29.1046 4 28 4Z"
            fill="none"
            stroke="currentColor"
            stroke-linejoin="round" />
        <path
            :stroke-width="props.strokeWidth"
            d="M16 32H8C6.89543 32 6 32.8954 6 34V42C6 43.1046 6.89543 44 8 44H16C17.1046 44 18 43.1046 18 42V34C18 32.8954 17.1046 32 16 32Z"
            fill="none"
            stroke="currentColor"
            stroke-linejoin="round" />
        <path
            :stroke-width="props.strokeWidth"
            d="M40 32H32C30.8954 32 30 32.8954 30 34V42C30 43.1046 30.8954 44 32 44H40C41.1046 44 42 43.1046 42 42V34C42 32.8954 41.1046 32 40 32Z"
            fill="none"
            stroke="currentColor"
            stroke-linejoin="round" />
        <path
            :stroke-width="props.strokeWidth"
            d="M24 24V16"
            stroke="currentColor"
            stroke-linecap="round"
            stroke-linejoin="round" />
        <path
            :stroke-width="props.strokeWidth"
            d="M36 32V24"
            stroke="currentColor"
            stroke-linecap="round"
            stroke-linejoin="round" />
        <path
            :stroke-width="props.strokeWidth"
            d="M12 32V24"
            stroke="currentColor"
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
