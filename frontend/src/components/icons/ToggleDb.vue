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
        <path
            :stroke="props.fillColor"
            :stroke-width="props.strokeWidth"
            d="M44.0001 11C44.0001 11 44 36.0623 44 38C44 41.3137 35.0457 44 24 44C12.9543 44 4.00003 41.3137 4.00003 38C4.00003 36.1423 4 11 4 11"
            stroke-linecap="round"
            stroke-linejoin="round"
        />
        <path
            :stroke="props.fillColor"
            :stroke-width="props.strokeWidth"
            d="M44 29C44 32.3137 35.0457 35 24 35C12.9543 35 4 32.3137 4 29"
            stroke-linecap="round"
            stroke-linejoin="round"
        />
        <path
            :stroke="props.fillColor"
            :stroke-width="props.strokeWidth"
            d="M44 20C44 23.3137 35.0457 26 24 26C12.9543 26 4 23.3137 4 20"
            stroke-linecap="round"
            stroke-linejoin="round"
        />
        <ellipse
            :stroke="props.fillColor"
            :stroke-width="props.strokeWidth"
            cx="24"
            cy="10"
            fill="#dc423c"
            rx="20"
            ry="6"
            stroke-linecap="round"
            stroke-linejoin="round"
        />
    </svg>
    <svg v-else fill="none" viewBox="0 0 48 48" xmlns="http://www.w3.org/2000/svg">
        <path
            :stroke-width="props.strokeWidth"
            d="M44.0001 11C44.0001 11 44 36.0623 44 38C44 41.3137 35.0457 44 24 44C12.9543 44 4.00003 41.3137 4.00003 38C4.00003 36.1423 4 11 4 11"
            stroke="currentColor"
            stroke-linecap="round"
            stroke-linejoin="round"
        />
        <path
            :stroke-width="props.strokeWidth"
            d="M44 29C44 32.3137 35.0457 35 24 35C12.9543 35 4 32.3137 4 29"
            stroke="currentColor"
            stroke-linecap="round"
            stroke-linejoin="round"
        />
        <path
            :stroke-width="props.strokeWidth"
            d="M44 20C44 23.3137 35.0457 26 24 26C12.9543 26 4 23.3137 4 20"
            stroke="currentColor"
            stroke-linecap="round"
            stroke-linejoin="round"
        />
        <ellipse
            :stroke-width="props.strokeWidth"
            cx="24"
            cy="10"
            fill="none"
            rx="20"
            ry="6"
            stroke="currentColor"
            stroke-linecap="round"
            stroke-linejoin="round"
        />
    </svg>
</template>

<style lang="scss" scoped></style>
