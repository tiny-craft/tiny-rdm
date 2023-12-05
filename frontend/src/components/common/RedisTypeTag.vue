<script setup>
import { computed } from 'vue'
import { typesBgColor, typesColor, typesShortName } from '@/consts/support_redis_type.js'
import Binary from '@/components/icons/Binary.vue'
import { toUpper } from 'lodash'

const props = defineProps({
    type: {
        type: String,
        default: 'STRING',
    },
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
})

const fontColor = computed(() => {
    if (props.inverse) {
        return typesBgColor[props.type]
    } else {
        return typesColor[props.type]
    }
})

const backgroundColor = computed(() => {
    if (props.inverse) {
        return typesColor[props.type]
    } else {
        return typesBgColor[props.type]
    }
})

const label = computed(() => {
    if (props.short) {
        return typesShortName[toUpper(props.type)] || 'N'
    }
    return toUpper(props.type)
})
</script>

<template>
    <div
        v-if="props.point"
        :style="{
            backgroundColor: fontColor,
            width: Math.max(props.pointSize, 5) + 'px',
            height: Math.max(props.pointSize, 5) + 'px',
        }"
        class="redis-type-tag-round redis-type-tag-point"></div>
    <n-tag
        v-else
        :class="{
            'redis-type-tag-normal': !props.short && props.size !== 'small',
            'redis-type-tag-small': !props.short && props.size === 'small',
            'redis-type-tag-round': props.round,
        }"
        :color="{ color: backgroundColor, textColor: fontColor }"
        :size="props.size"
        bordered
        strong>
        <b>{{ label }}</b>
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
</style>
