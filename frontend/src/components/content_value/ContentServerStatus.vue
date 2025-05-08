<script setup>
import {
    cloneDeep,
    flatMap,
    get,
    isEmpty,
    map,
    mapValues,
    pickBy,
    random,
    slice,
    split,
    sum,
    toArray,
    toNumber,
} from 'lodash'
import { computed, h, onMounted, onUnmounted, reactive, ref, shallowRef, toRaw, watch } from 'vue'
import IconButton from '@/components/common/IconButton.vue'
import Filter from '@/components/icons/Filter.vue'
import Refresh from '@/components/icons/Refresh.vue'
import useBrowserStore from 'stores/browser.js'
import { timeout } from '@/utils/promise.js'
import AutoRefreshForm from '@/components/common/AutoRefreshForm.vue'
import { NButton, NIcon, NSpace, useThemeVars } from 'naive-ui'
import { Line } from 'vue-chartjs'
import dayjs from 'dayjs'
import { convertBytes, formatBytes } from '@/utils/byte_convert.js'
import usePreferencesStore from 'stores/preferences.js'
import { useI18n } from 'vue-i18n'
import useConnectionStore from 'stores/connections.js'
import { toHumanReadable } from '@/utils/date.js'

const props = defineProps({
    server: String,
    pause: Boolean,
})

const browserStore = useBrowserStore()
const prefStore = usePreferencesStore()
const connectionStore = useConnectionStore()
const i18n = useI18n()
const themeVars = useThemeVars()
const serverInfo = ref({})
const pageState = reactive({
    autoRefresh: false,
    refreshInterval: 5,
    loading: false, // loading status for refresh
    autoLoading: false, // loading status for auto refresh
})
const statusHistory = 5

/**
 *
 * @param origin
 * @param {string[]} [labels]
 * @param {number[][]} [datalist]
 * @return {unknown}
 */
const generateData = (origin, labels, datalist) => {
    let ret = toRaw(origin)
    ret.labels = labels || ret.labels
    if (datalist && datalist.length > 0) {
        for (let i = 0; i < datalist.length; i++) {
            ret.datasets[i].data = datalist[i]
        }
    }
    return cloneDeep(ret)
}

/**
 * refresh server status info
 * @param {boolean} [force] force refresh will show loading indicator
 * @returns {Promise<void>}
 */
const refreshInfo = async (force) => {
    if (force) {
        pageState.loading = true
    } else {
        pageState.autoLoading = true
    }
    if (!isEmpty(props.server) && browserStore.isConnected(props.server)) {
        try {
            const info = await browserStore.getServerInfo(props.server, true)
            if (!isEmpty(info)) {
                serverInfo.value = info
                _updateChart(info)
            }
        } finally {
            pageState.loading = false
            pageState.autoLoading = false
        }
    }
}

const _updateChart = (info) => {
    let timeLabels = toRaw(cmdRate.value.labels)
    timeLabels = timeLabels.concat(dayjs().format('HH:mm:ss'))
    timeLabels = slice(timeLabels, Math.max(0, timeLabels.length - statusHistory))

    // commands per seconds
    {
        let dataset = toRaw(cmdRate.value.datasets[0].data)
        const cmd = parseInt(get(info, 'Stats.instantaneous_ops_per_sec', '0'))
        dataset = dataset.concat(cmd)
        dataset = slice(dataset, Math.max(0, dataset.length - statusHistory))
        cmdRate.value = generateData(cmdRate.value, timeLabels, [dataset])
    }

    // connected clients
    {
        let dataset = toRaw(connectedClients.value.datasets[0].data)
        const count = parseInt(get(info, 'Clients.connected_clients', '0'))
        dataset = dataset.concat(count)
        dataset = slice(dataset, Math.max(0, dataset.length - statusHistory))
        connectedClients.value = generateData(connectedClients.value, timeLabels, [dataset])
    }

    // memory usage
    {
        let dataset = toRaw(memoryUsage.value.datasets[0].data)
        let size = parseInt(get(info, 'Memory.used_memory', '0'))
        dataset = dataset.concat(size)
        dataset = slice(dataset, Math.max(0, dataset.length - statusHistory))
        memoryUsage.value = generateData(memoryUsage.value, timeLabels, [dataset])
    }

    // network input/output rate
    {
        let dataset1 = toRaw(networkRate.value.datasets[0].data)
        const input = parseInt(get(info, 'Stats.instantaneous_input_kbps', '0'))
        dataset1 = dataset1.concat(input)
        dataset1 = slice(dataset1, Math.max(0, dataset1.length - statusHistory))

        let dataset2 = toRaw(networkRate.value.datasets[1].data)
        const output = parseInt(get(info, 'Stats.instantaneous_output_kbps', '0'))
        dataset2 = dataset2.concat(output)
        dataset2 = slice(dataset2, Math.max(0, dataset2.length - statusHistory))
        networkRate.value = generateData(networkRate.value, timeLabels, [dataset1, dataset2])
    }
}

/**
 * for mock activity data only
 * @private
 */
const _mockChart = () => {
    const timeLabels = []
    for (let i = 0; i < 5; i++) {
        timeLabels.push(dayjs().add(5, 'seconds').format('HH:mm:ss'))
    }

    // commands per seconds
    {
        const dataset = []
        for (let i = 0; i < 5; i++) {
            dataset.push(random(10, 200))
        }
        cmdRate.value = generateData(cmdRate.value, timeLabels, [dataset])
    }

    // connected clients
    {
        const dataset = []
        for (let i = 0; i < 5; i++) {
            dataset.push(random(10, 20))
        }
        connectedClients.value = generateData(connectedClients.value, timeLabels, [dataset])
    }

    // memory usage
    {
        const dataset = []
        for (let i = 0; i < 5; i++) {
            dataset.push(random(120 * 1024 * 1024, 200 * 1024 * 1024))
        }
        memoryUsage.value = generateData(memoryUsage.value, timeLabels, [dataset])
    }

    // network input/output rate
    {
        const dataset1 = []
        for (let i = 0; i < 5; i++) {
            dataset1.push(random(100, 1500))
        }

        const dataset2 = []
        for (let i = 0; i < 5; i++) {
            dataset2.push(random(200, 3000))
        }

        networkRate.value = generateData(networkRate.value, timeLabels, [dataset1, dataset2])
    }
}

const isLoading = computed(() => {
    return pageState.loading || pageState.autoLoading
})

const startAutoRefresh = async () => {
    // connectionStore.getRefreshInterval()
    let lastExec = Date.now()
    do {
        if (!pageState.autoRefresh) {
            break
        }
        await timeout(100)
        if (
            props.pause ||
            pageState.loading ||
            pageState.autoLoading ||
            Date.now() - lastExec < pageState.refreshInterval * 1000
        ) {
            continue
        }
        lastExec = Date.now()
        await refreshInfo()
    } while (true)
    stopAutoRefresh()
}

const stopAutoRefresh = () => {
    pageState.autoRefresh = false
}

const onToggleRefresh = (on) => {
    if (on) {
        tabVal.value = 'activity'
        connectionStore.saveRefreshInterval(props.server, pageState.refreshInterval || 5)
        startAutoRefresh()
    } else {
        connectionStore.saveRefreshInterval(props.server, -1)
        stopAutoRefresh()
    }
}

onMounted(() => {
    const interval = connectionStore.getRefreshInterval(props.server)
    if (interval >= 0) {
        pageState.autoRefresh = true
        pageState.refreshInterval = interval === 0 ? 5 : interval
        onToggleRefresh(true)
    } else {
        setTimeout(refreshInfo, 3000)
        // setTimeout(_mockChart, 1000)
    }
    refreshInfo()
})

onUnmounted(() => {
    stopAutoRefresh()
})

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
    let seconds = parseInt(get(serverInfo.value, 'Server.uptime_in_seconds', '0'))
    seconds /= 60
    if (seconds < 60) {
        // minutes
        return { value: Math.floor(seconds), unit: timeUnit[0] }
    }
    seconds /= 60
    if (seconds < 60) {
        // hours
        return { value: Math.floor(seconds), unit: timeUnit[1] }
    }
    return { value: Math.floor(seconds / 24), unit: timeUnit[2] }
})

const usedMemory = computed(() => {
    let size = parseInt(get(serverInfo.value, 'Memory.used_memory', '0'))
    const { value, unit } = convertBytes(size)
    return [value, unit]
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

const tabVal = ref('activity')
const infoFilter = reactive({
    keyword: '',
    group: 'CPU',
})

const info = computed(() => {
    if (!isEmpty(infoFilter.group)) {
        const val = serverInfo.value[infoFilter.group]
        if (!isEmpty(val)) {
            return map(val, (v, k) => ({
                key: k,
                value: v,
            }))
        }
    }

    return flatMap(serverInfo.value, (value, key) => {
        return map(value, (v, k) => ({
            group: key,
            key: k,
            value: v,
        }))
    })
})

const onFilterGroup = (group) => {
    if (group === infoFilter.group) {
        infoFilter.group = ''
    } else {
        infoFilter.group = group
    }
}

watch(
    () => prefStore.currentLanguage,
    () => {
        // force update labels of charts
        cmdRate.value.datasets[0].label = i18n.t('status.act_cmd')
        cmdRate.value = generateData(cmdRate.value)
        connectedClients.value.datasets[0].label = i18n.t('status.connected_clients')
        connectedClients.value = generateData(connectedClients.value)
        memoryUsage.value.datasets[0].label = i18n.t('status.memory_used')
        memoryUsage.value = generateData(memoryUsage.value)
        networkRate.value.datasets[0].label = i18n.t('status.act_network_input')
        networkRate.value.datasets[1].label = i18n.t('status.act_network_output')
        networkRate.value = generateData(networkRate.value)
    },
)

const chartBGColor = [
    'rgba(255, 99, 132, 0.2)',
    'rgba(255, 159, 64, 0.2)',
    'rgba(153, 102, 255, 0.2)',
    'rgba(75, 192, 192, 0.2)',
    'rgba(54, 162, 235, 0.2)',
]

const chartBorderColor = [
    'rgb(255, 99, 132)',
    'rgb(255, 159, 64)',
    'rgb(153, 102, 255)',
    'rgb(75, 192, 192)',
    'rgb(54, 162, 235)',
]

const cmdRate = shallowRef({
    labels: [],
    datasets: [
        {
            label: i18n.t('status.act_cmd'),
            data: [],
            fill: true,
            backgroundColor: chartBGColor[0],
            borderColor: chartBorderColor[0],
            tension: 0.4,
        },
    ],
})

const connectedClients = shallowRef({
    labels: [],
    datasets: [
        {
            label: i18n.t('status.connected_clients'),
            data: [],
            fill: true,
            backgroundColor: chartBGColor[1],
            borderColor: chartBorderColor[1],
            tension: 0.4,
        },
    ],
})

const memoryUsage = shallowRef({
    labels: [],
    datasets: [
        {
            label: i18n.t('status.memory_used'),
            data: [],
            fill: true,
            backgroundColor: chartBGColor[2],
            borderColor: chartBorderColor[2],
            tension: 0.4,
        },
    ],
})

const networkRate = shallowRef({
    labels: [],
    datasets: [
        {
            label: i18n.t('status.act_network_input'),
            data: [],
            fill: true,
            backgroundColor: chartBGColor[3],
            borderColor: chartBorderColor[3],
            tension: 0.4,
        },
        {
            label: i18n.t('status.act_network_output'),
            data: [],
            fill: true,
            backgroundColor: chartBGColor[4],
            borderColor: chartBorderColor[4],
            tension: 0.4,
        },
    ],
})

const chartOption = computed(() => {
    return {
        animation: false,
        responsive: true,
        maintainAspectRatio: false,
        events: [],
        scales: {
            x: {
                grid: {
                    color: themeVars.value.borderColor,
                },
                ticks: {
                    color: themeVars.value.textColor3,
                },
            },
            y: {
                beginAtZero: true,
                stepSize: 1024,
                suggestedMin: 0,
                grid: {
                    color: themeVars.value.borderColor,
                },
                ticks: {
                    color: themeVars.value.textColor3,
                    precision: 0,
                },
            },
        },
        plugins: {
            legend: {
                labels: {
                    color: themeVars.value.textColor2,
                },
            },
        },
    }
})

const byteChartOption = computed(() => {
    return {
        animation: false,
        responsive: true,
        maintainAspectRatio: false,
        events: [],
        scales: {
            x: {
                grid: {
                    color: themeVars.value.borderColor,
                },
                ticks: {
                    color: themeVars.value.textColor3,
                },
            },
            y: {
                beginAtZero: true,
                stepSize: 1024,
                suggestedMin: 0,
                grid: {
                    color: themeVars.value.borderColor,
                },
                ticks: {
                    color: themeVars.value.textColor3,
                    precision: 0,
                    // format display y axios tag
                    callback: function (value, index, values) {
                        return formatBytes(value, 1)
                    },
                },
            },
        },
        plugins: {
            legend: {
                labels: {
                    color: themeVars.value.textColor2,
                },
            },
        },
    }
})

const clientInfo = reactive({
    loading: false,
    content: [],
})
const onShowClients = async (show) => {
    if (show) {
        try {
            clientInfo.loading = true
            clientInfo.content = await browserStore.getClientList(props.server)
        } finally {
            clientInfo.loading = false
        }
    }
}

const clientTableColumns = computed(() => {
    return [
        {
            key: 'title',
            title: () => {
                return h(NSpace, { wrap: false, wrapItem: false, justify: 'center' }, () => [
                    h('span', { style: { fontWeight: '550', fontSize: '15px' } }, i18n.t('status.client.title')),
                    h(IconButton, {
                        icon: Refresh,
                        size: 16,
                        onClick: () => onShowClients(true),
                    }),
                ])
            },
            align: 'center',
            titleAlign: 'center',
            children: [
                {
                    key: 'no',
                    title: '#',
                    width: 60,
                    align: 'center',
                    titleAlign: 'center',
                    render: (row, index) => {
                        return index + 1
                    },
                },
                {
                    key: 'addr',
                    title: () => i18n.t('status.client.addr'),
                    sorter: 'default',
                    align: 'center',
                    titleAlign: 'center',
                },
                {
                    key: 'db',
                    title: () => i18n.t('status.client.db'),
                    align: 'center',
                    titleAlign: 'center',
                },
                {
                    key: 'age',
                    title: () => i18n.t('status.client.age'),
                    sorter: (row1, row2) => row1.age - row2.age,
                    defaultSortOrder: 'descend',
                    align: 'center',
                    titleAlign: 'center',
                    render: ({ age }, index) => {
                        return toHumanReadable(age)
                    },
                },
                {
                    key: 'idle',
                    title: () => i18n.t('status.client.idle'),
                    sorter: (row1, row2) => row1.idle - row2.idle,
                    align: 'center',
                    titleAlign: 'center',
                    render: ({ idle }, index) => {
                        return toHumanReadable(idle)
                    },
                },
            ],
        },
    ]
})
</script>

<template>
    <n-space :size="5" :wrap-item="false" style="padding: 5px; box-sizing: border-box; height: 100%" vertical>
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
                    <n-tooltip v-if="role">
                        Role
                        <template #trigger>
                            <n-tag size="small" type="primary">{{ role }}</n-tag>
                        </template>
                    </n-tooltip>
                </n-space>
            </template>
            <template #header-extra>
                <n-popover keep-alive-on-hover placement="bottom-end" trigger="hover">
                    <template #trigger>
                        <n-button
                            :loading="pageState.loading"
                            :type="isLoading ? 'primary' : 'default'"
                            circle
                            size="small"
                            tertiary
                            @click="refreshInfo(true)">
                            <template #icon>
                                <n-icon :size="props.size">
                                    <refresh
                                        :class="{
                                            'auto-rotate': pageState.autoRefresh || isLoading,
                                        }"
                                        :color="pageState.autoRefresh ? themeVars.primaryColor : undefined"
                                        :stroke-width="pageState.autoRefresh ? 6 : 3" />
                                </n-icon>
                            </template>
                        </n-button>
                    </template>
                    <auto-refresh-form
                        v-model:interval="pageState.refreshInterval"
                        v-model:on="pageState.autoRefresh"
                        :default-value="5"
                        :loading="pageState.autoLoading"
                        @toggle="onToggleRefresh" />
                </n-popover>
            </template>
            <n-grid style="min-width: 500px" x-gap="5">
                <n-gi :span="6">
                    <n-statistic :label="$t('status.uptime')" :value="uptime.value">
                        <template #suffix>{{ $t(uptime.unit) }}</template>
                    </n-statistic>
                </n-gi>
                <n-gi :span="6">
                    <n-statistic
                        :label="$t('status.connected_clients')"
                        :value="get(serverInfo, 'Clients.connected_clients', '0')">
                        <template #suffix>
                            <n-tooltip
                                :content-style="{ backgroundColor: themeVars.tableColor }"
                                trigger="click"
                                width="70vw"
                                @update-show="onShowClients">
                                <template #trigger>
                                    <n-button :bordered="false" size="small">&LowerRightArrow;</n-button>
                                </template>
                                <n-data-table
                                    :columns="clientTableColumns"
                                    :data="clientInfo.content"
                                    :loading="clientInfo.loading"
                                    :single-column="false"
                                    :single-line="false"
                                    max-height="50vh"
                                    size="small"
                                    striped />
                            </n-tooltip>
                        </template>
                    </n-statistic>
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
        </n-card>
        <n-card class="flex-item-expand" content-style="padding: 0; height: 100%;" embedded style="overflow: hidden">
            <n-tabs
                v-model:value="tabVal"
                :tabs-padding="20"
                pane-style="padding: 10px; box-sizing: border-box; display: flex; flex-direction: column; flex-grow: 1;"
                size="large"
                style="height: 100%; overflow: hidden"
                type="line">
                <template #suffix>
                    <div v-if="tabVal === 'info'" style="padding-right: 10px">
                        <n-input v-model:value="infoFilter.keyword" clearable placeholder="">
                            <template #prefix>
                                <icon-button :icon="Filter" size="18" />
                            </template>
                        </n-input>
                    </div>
                </template>

                <!-- activity tab pane -->
                <n-tab-pane
                    :tab="$t('status.activity_status')"
                    class="line-chart"
                    display-directive="show:lazy"
                    name="activity">
                    <div class="line-chart">
                        <div class="line-chart-item">
                            <Line :data="cmdRate" :options="chartOption" />
                        </div>
                        <div class="line-chart-item">
                            <Line :data="connectedClients" :options="chartOption" />
                        </div>
                        <div class="line-chart-item">
                            <Line :data="memoryUsage" :options="byteChartOption" />
                        </div>
                        <div class="line-chart-item">
                            <Line :data="networkRate" :options="byteChartOption" />
                        </div>
                    </div>
                </n-tab-pane>

                <!-- info tab pane -->
                <n-tab-pane :tab="$t('status.server_info')" name="info">
                    <n-space :wrap="false" :wrap-item="false" class="flex-item-expand">
                        <n-space align="end" item-style="padding: 0 5px;" vertical>
                            <n-button
                                v-for="(v, k) in serverInfo"
                                :key="k"
                                :disabled="isEmpty(v)"
                                :focusable="false"
                                :type="infoFilter.group === k ? 'primary' : 'default'"
                                secondary
                                size="small"
                                @click="onFilterGroup(k)">
                                <span style="min-width: 80px">{{ k }}</span>
                            </n-button>
                        </n-space>
                        <n-data-table
                            :columns="[
                                {
                                    title: $t('common.key'),
                                    key: 'key',
                                    defaultSortOrder: 'ascend',
                                    minWidth: 80,
                                    titleAlign: 'center',
                                    filterOptionValue: infoFilter.keyword,
                                    filter(value, row) {
                                        return !!~row.key.indexOf(value.toString())
                                    },
                                },
                                { title: $t('common.value'), titleAlign: 'center', key: 'value' },
                            ]"
                            :data="info"
                            :loading="pageState.loading"
                            :single-line="false"
                            class="flex-item-expand"
                            flex-height
                            striped />
                    </n-space>
                </n-tab-pane>
            </n-tabs>
        </n-card>
    </n-space>
</template>

<style lang="scss" scoped>
@use '@/styles/content';

.line-chart {
    display: flex;
    flex-wrap: wrap;
    width: 100%;
    height: 100%;

    &-item {
        display: flex;
        flex-direction: column;
        justify-content: center;
        align-items: center;
        width: 50%;
        height: 50%;
        padding: 10px;
        box-sizing: border-box;
    }
}
</style>
