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
            d="M12 9.92704V7C12 5.34315 13.3431 4 15 4H41C42.6569 4 44 5.34315 44 7V33C44 34.6569 42.6569 36 41 36H38.0174"
            stroke="currentColor"
        />
        <rect
            :stroke-width="props.strokeWidth"
            fill="none"
            height="34"
            rx="3"
            stroke="currentColor"
            stroke-linejoin="round"
            width="34"
            x="4"
            y="10"
        />
        <path
            :stroke-width="props.strokeWidth"
            d="M18.4394 23.1101L23.7319 17.6006C25.1835 16.1489 27.5691 16.1809 29.0602 17.672C30.5513 19.1631 30.5833 21.5487 29.1316 23.0003L27.2215 25.0231"
            stroke="currentColor"
            stroke-linecap="round"
            stroke-linejoin="round"
        />
        <path
            :stroke-width="props.strokeWidth"
            d="M13.4661 28.7472C12.9558 29.2575 11.9006 30.2765 11.9006 30.2765C10.4489 31.7281 10.4095 34.3155 11.9006 35.8066C13.3917 37.2977 15.7772 37.3296 17.2289 35.878L22.3931 31.1896"
            stroke="currentColor"
            stroke-linecap="round"
            stroke-linejoin="round"
        />
        <path
            :stroke-width="props.strokeWidth"
            d="M18.6626 28.3284C17.97 27.6358 17.5922 26.7502 17.5317 25.8548C17.4619 24.8226 17.8138 23.7774 18.5912 23.0001"
            stroke="currentColor"
            stroke-linecap="round"
            stroke-linejoin="round"
        />
        <path
            :stroke-width="props.strokeWidth"
            d="M22.3213 25.8613C23.8124 27.3524 23.8444 29.738 22.3927 31.1896"
            stroke="currentColor"
            stroke-linecap="round"
            stroke-linejoin="round"
        />
    </svg>
</template>

<style lang="scss" scoped></style>
