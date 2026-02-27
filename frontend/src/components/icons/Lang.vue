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
            d="M28.2857 37H39.7143M42 42L39.7143 37L42 42ZM26 42L28.2857 37L26 42ZM28.2857 37L34 24L39.7143 37H28.2857Z"
            stroke="currentColor"
            stroke-linecap="round"
            stroke-linejoin="round" />
        <path
            :stroke-width="props.strokeWidth"
            d="M16 6L17 9"
            stroke="currentColor"
            stroke-linecap="round"
            stroke-linejoin="round" />
        <path
            :stroke-width="props.strokeWidth"
            d="M6 11H28"
            stroke="currentColor"
            stroke-linecap="round"
            stroke-linejoin="round" />
        <path
            :stroke-width="props.strokeWidth"
            d="M10 16C10 16 11.7895 22.2609 16.2632 25.7391C20.7368 29.2174 28 32 28 32"
            stroke="currentColor"
            stroke-linecap="round"
            stroke-linejoin="round" />
        <path
            :stroke-width="props.strokeWidth"
            d="M24 11C24 11 22.2105 19.2174 17.7368 23.7826C13.2632 28.3478 6 32 6 32"
            stroke="currentColor"
            stroke-linecap="round"
            stroke-linejoin="round" />
    </svg>
</template>

<style lang="scss" scoped></style>
