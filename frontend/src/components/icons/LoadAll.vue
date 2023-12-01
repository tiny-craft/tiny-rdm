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
            clip-rule="evenodd"
            d="M23.9999 31L12 19L19.9999 19L19.9999 8L27.9999 8L27.9999 19L35.9999 19L23.9999 31Z"
            fill="none"
            fill-rule="evenodd"
            stroke="currentColor"
            stroke-linecap="round"
            stroke-linejoin="round" />
        <path :stroke-width="props.strokeWidth" d="M42 38L6 38" stroke="currentColor" stroke-linecap="round" />
    </svg>
</template>

<style lang="scss" scoped></style>
