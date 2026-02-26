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
        <path d="M23.9917 6H6V42H24" stroke="currentColor" stroke-linecap="round" stroke-linejoin="round"/>
        <path d="M33 33L42 24L33 15" stroke="currentColor" stroke-linecap="round" stroke-linejoin="round"/>
        <path d="M16 23.9917H42" stroke="currentColor" stroke-linecap="round" stroke-linejoin="round"/>
    </svg>
</template>

<style lang="scss" scoped></style>
