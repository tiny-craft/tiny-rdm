<script setup>
import { get, isEmpty, map, mapValues, pickBy, split, sum, toArray, toNumber } from 'lodash'
import { computed, onMounted, onUnmounted, ref } from 'vue'
import IconButton from '@/components/common/IconButton.vue'
import Filter from '@/components/icons/Filter.vue'
import Refresh from '@/components/icons/Refresh.vue'
import useBrowserStore from 'stores/browser.js'

const props = defineProps({
    server: String,
})

const browserStore = useBrowserStore()
const serverInfo = ref({})
const autoRefresh = ref(false)
const loading = ref(false) // loading status for refresh
const autoLoading = ref(false) // loading status for auto refresh

/**
 * refresh server status info
 * @param {boolean} [force] force refresh will show loading indicator
 * @returns {Promise<void>}
 */
const refreshInfo = async (force) => {
    if (force) {
        loading.value = true
    } else {
        autoLoading.value = true
    }
    if (!isEmpty(props.server) && browserStore.isConnected(props.server)) {
        try {
            serverInfo.value = await browserStore.getServerInfo(props.server)
        } finally {
            loading.value = false
            autoLoading.value = false
        }
    }
}

let intervalID
onMounted(() => {
    refreshInfo()
    intervalID = setInterval(() => {
        if (autoRefresh.value === true) {
            refreshInfo()
        }
    }, 5000)
})

onUnmounted(() => {
    clearInterval(intervalID)
})

const scrollRef = ref(null)
const redisVersion = computed(() => {
    return get(serverInfo.value, 'Server.redis_version', '')
})

const redisMode = computed(() => {
    return get(serverInfo.value, 'Server.redis_mode', '')
})

const role = computed(() => {
    return get(serverInfo.value, 'Replication.role', '')
})

const timeUnit = ['common.unit_minute', 'common.unit_hour', 'common.unit_day']
const uptime = computed(() => {
    let seconds = get(serverInfo.value, 'Server.uptime_in_seconds', 0)
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
    let size = get(serverInfo.value, 'Memory.used_memory', 0)
    let unitIndex = 0

    while (size >= 1024 && unitIndex < units.length - 1) {
        size /= 1024
        unitIndex++
    }

    return [size.toFixed(2), units[unitIndex]]
})

const totalKeys = computed(() => {
    const regex = /^db\d+$/
    const result = pickBy(serverInfo.value['Keyspace'], (value, key) => {
        return regex.test(key)
    })
    const nums = mapValues(result, (v) => {
        const keys = split(v, ',', 1)[0]
        const num = split(keys, '=', 2)[1]
        return toNumber(num)
    })
    return sum(toArray(nums))
})
const infoFilter = ref('')
</script>

<template>
    <n-scrollbar ref="scrollRef">
        <n-back-top :listen-to="scrollRef" />
        <n-space :size="5" :wrap-item="false" style="padding: 5px" vertical>
            <n-card embedded>
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
                        {{ $t('status.auto_refresh') }}
                        <n-switch v-model:value="autoRefresh" :loading="autoLoading" />
                        <n-tooltip>
                            {{ $t('status.refresh') }}
                            <template #trigger>
                                <n-button
                                    :loading="autoLoading"
                                    circle
                                    size="small"
                                    tertiary
                                    @click="refreshInfo(true)">
                                    <template #icon>
                                        <n-icon :component="Refresh" />
                                    </template>
                                </n-button>
                            </template>
                        </n-tooltip>
                    </n-space>
                </template>
                <n-spin :show="loading">
                    <n-grid style="min-width: 500px" x-gap="5">
                        <n-gi :span="6">
                            <n-statistic :label="$t('status.uptime')" :value="uptime[0]">
                                <template #suffix>{{ $t(uptime[1]) }}</template>
                            </n-statistic>
                        </n-gi>
                        <n-gi :span="6">
                            <n-statistic
                                :label="$t('status.connected_clients')"
                                :value="get(serverInfo, 'Clients.connected_clients', 0)" />
                        </n-gi>
                        <n-gi :span="6">
                            <n-statistic :value="totalKeys">
                                <template #label>
                                    {{ $t('status.total_keys') }}
                                </template>
                            </n-statistic>
                        </n-gi>
                        <n-gi :span="6">
                            <n-statistic :label="$t('status.memory_used')" :value="usedMemory[0]">
                                <template #suffix>{{ usedMemory[1] }}</template>
                            </n-statistic>
                        </n-gi>
                    </n-grid>
                </n-spin>
            </n-card>
            <n-card :title="$t('status.all_info')" embedded>
                <template #header-extra>
                    <n-input v-model:value="infoFilter" :input-props="{ spellcheck: 'false' }" clearable placeholder="">
                        <template #prefix>
                            <icon-button :icon="Filter" size="18" />
                        </template>
                    </n-input>
                </template>
                <n-spin :show="loading">
                    <n-tabs default-value="CPU" placement="left" type="line">
                        <n-tab-pane v-for="(v, k) in serverInfo" :key="k" :disabled="isEmpty(v)" :name="k">
                            <n-data-table
                                :columns="[
                                    {
                                        title: $t('common.key'),
                                        key: 'key',
                                        defaultSortOrder: 'ascend',
                                        sorter: 'default',
                                        minWidth: 100,
                                        filterOptionValue: infoFilter,
                                        filter(value, row) {
                                            return !!~row.key.indexOf(value.toString())
                                        },
                                    },
                                    { title: $t('common.value'), key: 'value' },
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
