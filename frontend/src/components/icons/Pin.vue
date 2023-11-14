<script setup>
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
    <svg fill="none" viewBox="0 0 48 48" xmlns="http://www.w3.org/2000/svg">
        <g clip-path="url(#icon-36c2b93b4e0022c)">
            <path
                :fill="inverse ? 'currentColor' : 'none'"
                :stroke-width="props.strokeWidth"
                d="M10.6963 17.5042C13.3347 14.8657 16.4701 14.9387 19.8781 16.8076L32.62 9.74509L31.8989 4.78683L43.2126 16.1005L38.2656 15.3907L31.1918 28.1214C32.9752 31.7589 33.1337 34.6647 30.4953 37.3032C30.4953 37.3032 26.235 33.0429 22.7171 29.525L6.44305 41.5564L18.4382 25.2461C14.9202 21.7281 10.6963 17.5042 10.6963 17.5042Z"
                stroke="currentColor"
                stroke-linejoin="round" />
        </g>
        <defs>
            <clipPath id="icon-36c2b93b4e0022c">
                <rect fill="#FFF" height="48" width="48" />
            </clipPath>
        </defs>
    </svg>
</template>

<style lang="scss" scoped></style>
