<script setup>
import { computed } from 'vue'
import { typesBgColor, typesColor, typesShortName } from '@/consts/support_redis_type.js'
import Binary from '@/components/icons/Binary.vue'
import { get, toUpper } from 'lodash'
import { useThemeVars } from 'naive-ui'
import Loading from '@/components/icons/Loading.vue'

const props = defineProps({
    type: {
        type: String,
        default: 'STRING',
    },
    defaultLabel: String,
    binaryKey: Boolean,
    size: String,
    short: Boolean,
    point: Boolean,
    pointSize: {
        type: Number,
        default: 14,
    },
    round: Boolean,
    inverse: Boolean,
    loading: Boolean,
})

const themeVars = useThemeVars()

const fontColor = computed(() => {
    if (props.inverse) {
        return props.loading ? themeVars.value.tagColor : typesBgColor[props.type]
    } else {
        return props.loading ? themeVars.value.textColorBase : typesColor[props.type]
    }
})

const backgroundColor = computed(() => {
    if (props.inverse) {
        return props.loading ? themeVars.value.textColorBase : typesColor[props.type]
    } else {
        return props.loading ? themeVars.value.tagColor : typesBgColor[props.type]
    }
})

const label = computed(() => {
    if (props.short) {
        return get(typesShortName, toUpper(props.type), props.defaultLabel || 'N')
    }
    return toUpper(props.type)
})
</script>

<template>
    <div
        v-if="props.point"
        :class="{ 'redis-type-tag-loading': props.loading }"
        :style="{
            backgroundColor: fontColor,
            width: Math.max(props.pointSize, 5) + 'px',
            height: Math.max(props.pointSize, 5) + 'px',
        }"
        class="redis-type-tag-round redis-type-tag-point" />
    <n-tag
        v-else
        :class="{
            'redis-type-tag-normal': !props.short && props.size !== 'small',
            'redis-type-tag-small': !props.short && props.size === 'small',
            'redis-type-tag-round': props.round,
            'redis-type-tag-loading': props.loading,
            'redis-type-tag': props.short,
        }"
        :color="{ color: backgroundColor, textColor: fontColor }"
        :size="props.size"
        bordered
        strong>
        <b v-if="!props.loading">{{ label }}</b>
        <n-icon v-else-if="props.short" size="14">
            <loading stroke-width="4" />
        </n-icon>
        <b v-else>LOADING</b>
        <template #icon>
            <n-icon v-if="binaryKey" :component="Binary" size="18" />
        </template>
    </n-tag>
</template>

<style lang="scss">
.redis-type-tag-round {
    border-radius: 9999px;
}

.redis-type-tag-normal {
    padding: 0 12px;
}

.redis-type-tag-small {
    padding: 0 5px;
}

.redis-type-tag-loading {
    animation: fadeInOut 2s infinite;
}

@keyframes fadeInOut {
    0% {
        opacity: 0.4;
    }
    50% {
        opacity: 1;
    }
    100% {
        opacity: 0.4;
    }
}

.redis-type-tag {
    width: 22px;
    height: 22px;
    justify-content: center;
    align-items: center;
    text-align: center;
    vertical-align: middle;
}
</style>
