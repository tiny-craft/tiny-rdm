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
    <svg v-if="props.inverse" viewBox="0 0 48 48" fill="none" xmlns="http://www.w3.org/2000/svg">
        <rect
            x="6"
            y="6"
            width="36"
            height="36"
            rx="3"
            fill="currentColor"
            stroke="currentColor"
            :stroke-width="props.strokeWidth"
            stroke-linejoin="round" />
        <rect
            x="13"
            y="13"
            width="8"
            height="8"
            fill="#FFF"
            stroke="#FFF"
            :stroke-width="props.strokeWidth"
            stroke-linejoin="round" />
        <path
            d="M27 13L35 13"
            stroke="#FFF"
            :stroke-width="props.strokeWidth"
            stroke-linecap="round"
            stroke-linejoin="round" />
        <path
            d="M27 20L35 20"
            stroke="#FFF"
            :stroke-width="props.strokeWidth"
            stroke-linecap="round"
            stroke-linejoin="round" />
        <path
            d="M13 28L35 28"
            stroke="#FFF"
            :stroke-width="props.strokeWidth"
            stroke-linecap="round"
            stroke-linejoin="round" />
        <path
            d="M13 35H35"
            stroke="#FFF"
            :stroke-width="props.strokeWidth"
            stroke-linecap="round"
            stroke-linejoin="round" />
    </svg>
    <svg v-else viewBox="0 0 48 48" fill="none" xmlns="http://www.w3.org/2000/svg">
        <rect
            x="6"
            y="6"
            width="36"
            height="36"
            rx="3"
            fill="none"
            stroke="currentColor"
            :stroke-width="props.strokeWidth"
            stroke-linejoin="round" />
        <rect
            x="13"
            y="13"
            width="8"
            height="8"
            fill="none"
            stroke="currentColor"
            :stroke-width="props.strokeWidth"
            stroke-linejoin="round" />
        <path
            d="M27 13L35 13"
            stroke="currentColor"
            :stroke-width="props.strokeWidth"
            stroke-linecap="round"
            stroke-linejoin="round" />
        <path
            d="M27 20L35 20"
            stroke="currentColor"
            :stroke-width="props.strokeWidth"
            stroke-linecap="round"
            stroke-linejoin="round" />
        <path
            d="M13 28L35 28"
            stroke="currentColor"
            :stroke-width="props.strokeWidth"
            stroke-linecap="round"
            stroke-linejoin="round" />
        <path
            d="M13 35H35"
            stroke="currentColor"
            :stroke-width="props.strokeWidth"
            stroke-linecap="round"
            stroke-linejoin="round" />
    </svg>
</template>

<style lang="scss" scoped></style>
