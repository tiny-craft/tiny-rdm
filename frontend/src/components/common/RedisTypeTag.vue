<script setup>
import { computed } from 'vue'
import { typesColor, validType } from '../../consts/support_redis_type.js'

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
        default: 'white',
    },
    size: String,
})

const backgroundColor = computed(() => {
    return typesColor[props.type]
})
</script>

<template>
    <n-tag
        :bordered="false"
        :color="{ color: backgroundColor, textColor: props.color }"
        :size="props.size"
        :class="[props.size === 'small' ? 'redis-type-tag-small' : 'redis-type-tag']"
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

.redis-type-tag-small {
    padding: 0 5px;
}
</style>
