<script setup>
import { get, isEmpty, map, mapValues, pickBy, split, sum, toArray, toNumber } from 'lodash'
import { computed, ref } from 'vue'
import Help from '@/components/icons/Help.vue'
import IconButton from '@/components/common/IconButton.vue'
import Filter from '@/components/icons/Filter.vue'
import Refresh from '@/components/icons/Refresh.vue'

const props = defineProps({
    server: String,
    info: Object,
    autoRefresh: false,
    loading: false,
})

const emit = defineEmits(['update:autoRefresh', 'refresh'])

const scrollRef = ref(null)
const redisVersion = computed(() => {
    return get(props.info, 'Server.redis_version', '')
})

const redisMode = computed(() => {
    return get(props.info, 'Server.redis_mode', '')
})

const role = computed(() => {
    return get(props.info, 'Replication.role', '')
})

const timeUnit = ['unit_minute', 'unit_hour', 'unit_day']
const uptime = computed(() => {
    let seconds = get(props.info, 'Server.uptime_in_seconds', 0)
    seconds /= 60
    if (seconds < 60) {
        // minutes
        return [Math.floor(seconds), timeUnit[0]]
    }
    seconds /= 60
    if (seconds < 60) {
        // hours
        return [Math.floor(seconds), timeUnit[1]]
    }
    return [Math.floor(seconds / 24), timeUnit[2]]
})

const units = ['B', 'KB', 'MB', 'GB', 'TB']
const usedMemory = computed(() => {
    let size = get(props.info, 'Memory.used_memory', 0)
    let unitIndex = 0

    while (size >= 1024 && unitIndex < units.length - 1) {
        size /= 1024
        unitIndex++
    }

    return [size.toFixed(2), units[unitIndex]]
})

const totalKeys = computed(() => {
    const regex = /^db\d+$/
    const result = pickBy(props.info['Keyspace'], (value, key) => {
        return regex.test(key)
    })
    const nums = mapValues(result, (v) => {
        const keys = split(v, ',', 1)[0]
        const num = split(keys, '=', 2)[1]
        return toNumber(num)
    })
    return sum(toArray(nums))
})
const infoList = computed(() => map(props.info, (value, key) => ({ value, key })))
const infoTab = ref('')
const infoFilter = ref('')
</script>

<template>
    <n-scrollbar ref="scrollRef">
        <n-back-top :listen-to="scrollRef" />
        <n-space vertical>
            <n-card>
                <template #header>
                    <n-space :wrap-item="false" align="center" inline size="small">
                        {{ props.server }}
                        <n-tooltip v-if="redisVersion">
                            Redis Version
                            <template #trigger>
                                <n-tag size="small" type="primary">v{{ redisVersion }}</n-tag>
                            </template>
                        </n-tooltip>
                        <n-tooltip v-if="redisMode">
                            Mode
                            <template #trigger>
                                <n-tag size="small" type="primary">{{ redisMode }}</n-tag>
                            </template>
                        </n-tooltip>
                        <n-tooltip v-if="redisMode">
                            Role
                            <template #trigger>
                                <n-tag size="small" type="primary">{{ role }}</n-tag>
                            </template>
                        </n-tooltip>
                    </n-space>
                </template>
                <template #header-extra>
                    <n-space align="center" inline>
                        {{ $t('auto_refresh') }}
                        <n-switch :value="props.autoRefresh" @update:value="(v) => emit('update:autoRefresh', v)" />
                        <n-tooltip>
                            {{ $t('refresh') }}
                            <template #trigger>
                                <n-button circle size="small" tertiary @click="emit('refresh')">
                                    <template #icon>
                                        <n-icon :component="Refresh" />
                                    </template>
                                </n-button>
                            </template>
                        </n-tooltip>
                    </n-space>
                </template>
                <n-spin :show="props.loading">
                    <n-grid style="min-width: 500px" x-gap="5">
                        <n-gi :span="6">
                            <n-statistic :label="$t('uptime')" :value="uptime[0]">
                                <template #suffix> {{ $t(uptime[1]) }}</template>
                            </n-statistic>
                        </n-gi>
                        <n-gi :span="6">
                            <n-statistic
                                :label="$t('connected_clients')"
                                :value="get(props.info, 'Clients.connected_clients', 0)" />
                        </n-gi>
                        <n-gi :span="6">
                            <n-statistic :value="totalKeys">
                                <template #label>
                                    {{ $t('total_keys') }}
                                </template>
                            </n-statistic>
                        </n-gi>
                        <n-gi :span="6">
                            <n-statistic :label="$t('memory_used')" :value="usedMemory[0]">
                                <template #suffix> {{ usedMemory[1] }}</template>
                            </n-statistic>
                        </n-gi>
                    </n-grid>
                </n-spin>
            </n-card>
            <n-card :title="$t('all_info')">
                <template #header-extra>
                    <n-input v-model:value="infoFilter" clearable placeholder="">
                        <template #prefix>
                            <icon-button :icon="Filter" size="18" />
                        </template>
                    </n-input>
                </template>
                <n-spin :show="props.loading">
                    <n-tabs default-value="CPU" placement="left" type="line">
                        <n-tab-pane v-for="(v, k) in props.info" :key="k" :disabled="isEmpty(v)" :name="k">
                            <n-data-table
                                :columns="[
                                    {
                                        title: $t('key'),
                                        key: 'key',
                                        defaultSortOrder: 'ascend',
                                        sorter: 'default',
                                        minWidth: 100,
                                        filterOptionValue: infoFilter,
                                        filter(value, row) {
                                            return !!~row.key.indexOf(value.toString())
                                        },
                                    },
                                    { title: $t('value'), key: 'value' },
                                ]"
                                :data="map(v, (value, key) => ({ value, key }))" />
                        </n-tab-pane>
                    </n-tabs>
                </n-spin>
            </n-card>
        </n-space>
    </n-scrollbar>
</template>

<style lang="scss" scoped></style>
