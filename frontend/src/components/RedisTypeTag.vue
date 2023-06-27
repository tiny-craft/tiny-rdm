<script setup>
import { computed } from 'vue'
import { types, validType } from '../consts/support_redis_type.js'

const props = defineProps({
    type: {
        type: String,
        validator(value) {
            return validType(value)
        },
        default: 'STRING',
    },
    color: {
        type: String,
        default: '',
    },
    size: String,
})

const color = {
    [types.STRING]: '#626aef',
    [types.HASH]: '#576bfa',
    [types.LIST]: '#34b285',
    [types.SET]: '#bb7d52',
    [types.ZSET]: '#d053a5',
}

const backgroundColor = computed(() => {
    return color[props.type]
})
</script>

<template>
    <n-tag
        :bordered="false"
        :color="{ color: backgroundColor, textColor: 'white' }"
        :size="props.size"
        class="redis-type-tag"
        strong
    >
        {{ props.type }}
    </n-tag>
    <!--  <div class="redis-type-tag flex-box-h" :style="{backgroundColor: backgroundColor}">{{ props.type }}</div>-->
</template>

<style lang="scss">
.redis-type-tag {
    padding: 0 12px;
}
</style>
