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
        <mask
            id="icon-122098f7f10b972"
            height="48"
            maskUnits="userSpaceOnUse"
            style="mask-type: alpha"
            width="48"
            x="0"
            y="0">
            <path d="M48 0H0V48H48V0Z" fill="currentColor" />
        </mask>
        <g mask="url(#icon-122098f7f10b972)">
            <path
                :stroke-width="props.strokeWidth"
                d="M6 24.0083V42H42V24"
                stroke="currentColor"
                stroke-linecap="round"
                stroke-linejoin="round" />
            <path
                :stroke-width="props.strokeWidth"
                d="M33 15L24 6L15 15"
                stroke="currentColor"
                stroke-linecap="round"
                stroke-linejoin="round" />
            <path
                :stroke-width="props.strokeWidth"
                d="M23.9917 32V6"
                stroke="currentColor"
                stroke-linecap="round"
                stroke-linejoin="round" />
        </g>
    </svg>
</template>

<style lang="scss" scoped></style>
