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
            height="36"
            rx="3"
            stroke="currentColor"
            stroke-linejoin="round"
            width="36"
            x="6"
            y="6" />
        <rect
            :stroke-width="props.strokeWidth"
            fill="#FFF"
            height="8"
            stroke="#FFF"
            stroke-linejoin="round"
            width="8"
            x="13"
            y="13" />
        <path
            :stroke-width="props.strokeWidth"
            d="M27 13L35 13"
            stroke="#FFF"
            stroke-linecap="round"
            stroke-linejoin="round" />
        <path
            :stroke-width="props.strokeWidth"
            d="M27 20L35 20"
            stroke="#FFF"
            stroke-linecap="round"
            stroke-linejoin="round" />
        <path
            :stroke-width="props.strokeWidth"
            d="M13 28L35 28"
            stroke="#FFF"
            stroke-linecap="round"
            stroke-linejoin="round" />
        <path
            :stroke-width="props.strokeWidth"
            d="M13 35H35"
            stroke="#FFF"
            stroke-linecap="round"
            stroke-linejoin="round" />
    </svg>
    <svg v-else fill="none" viewBox="0 0 48 48" xmlns="http://www.w3.org/2000/svg">
        <rect
            :stroke-width="props.strokeWidth"
            fill="none"
            height="36"
            rx="3"
            stroke="currentColor"
            stroke-linejoin="round"
            width="36"
            x="6"
            y="6" />
        <rect
            :stroke-width="props.strokeWidth"
            fill="none"
            height="8"
            stroke="currentColor"
            stroke-linejoin="round"
            width="8"
            x="13"
            y="13" />
        <path
            :stroke-width="props.strokeWidth"
            d="M27 13L35 13"
            stroke="currentColor"
            stroke-linecap="round"
            stroke-linejoin="round" />
        <path
            :stroke-width="props.strokeWidth"
            d="M27 20L35 20"
            stroke="currentColor"
            stroke-linecap="round"
            stroke-linejoin="round" />
        <path
            :stroke-width="props.strokeWidth"
            d="M13 28L35 28"
            stroke="currentColor"
            stroke-linecap="round"
            stroke-linejoin="round" />
        <path
            :stroke-width="props.strokeWidth"
            d="M13 35H35"
            stroke="currentColor"
            stroke-linecap="round"
            stroke-linejoin="round" />
    </svg>
</template>

<style lang="scss" scoped></style>
