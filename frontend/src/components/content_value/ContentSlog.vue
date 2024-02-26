<script setup>
import { computed, h, nextTick, onMounted, onUnmounted, reactive, ref } from 'vue'
import Refresh from '@/components/icons/Refresh.vue'
import { debounce, isEmpty, map, size, split } from 'lodash'
import { useI18n } from 'vue-i18n'
import { NIcon, useThemeVars } from 'naive-ui'
import dayjs from 'dayjs'
import useBrowserStore from 'stores/browser.js'
import { timeout } from '@/utils/promise.js'
import AutoRefreshForm from '@/components/common/AutoRefreshForm.vue'

const themeVars = useThemeVars()

const browserStore = useBrowserStore()
const i18n = useI18n()
const props = defineProps({
    server: {
        type: String,
    },
})

const autoRefresh = reactive({
    on: false,
    interval: 5,
})

const data = reactive({
    list: [],
    sortOrder: 'descend',
    listLimit: 20,
    loading: false,
    client: '',
    keyword: '',
})

const tableRef = ref(null)

const columns = computed(() => [
    {
        title: () => i18n.t('slog.exec_time'),
        key: 'timestamp',
        sortOrder: data.sortOrder,
        sorter: 'default',
        width: 180,
        align: 'center',
        titleAlign: 'center',
        render: ({ timestamp }, index) => {
            return dayjs(timestamp).format('YYYY-MM-DD HH:mm:ss')
        },
    },
    {
        title: () => i18n.t('slog.client'),
        key: 'client',
        filterOptionValue: data.client,
        resizable: true,
        filter: (value, row) => {
            return value === '' || row.client === value.toString() || row.addr === value.toString()
        },
        width: 200,
        align: 'center',
        titleAlign: 'center',
        ellipsis: {
            tooltip: {
                style: {
                    maxWidth: '50vw',
                    maxHeight: '50vh',
                },
                scrollable: true,
            },
        },
        render: ({ client, addr }, index) => {
            let content = ''
            if (!isEmpty(client)) {
                content += client
            }
            if (!isEmpty(addr)) {
                if (!isEmpty(content)) {
                    content += ' - '
                }
                content += addr
            }
            return content
        },
    },
    {
        title: () => i18n.t('slog.cmd'),
        key: 'cmd',
        titleAlign: 'center',
        filterOptionValue: data.keyword,
        resizable: true,
        filter: (value, row) => {
            return value === '' || !!~row.cmd.indexOf(value.toString())
        },
        render: ({ cmd }, index) => {
            const cmdList = split(cmd, '\n')
            if (size(cmdList) > 1) {
                return h(
                    'div',
                    null,
                    map(cmdList, (c) => h('div', { class: 'cmd-line' }, c)),
                )
            }
            return h('div', { class: 'cmd-line' }, cmd)
        },
    },
    {
        title: () => i18n.t('slog.cost_time'),
        key: 'cost',
        width: 100,
        align: 'center',
        titleAlign: 'center',
        render: ({ cost }, index) => {
            const ms = dayjs.duration(cost).asMilliseconds()
            if (ms < 1000) {
                return `${ms} ms`
            } else {
                return `${Math.floor(ms / 1000)} s`
            }
        },
    },
])

const _loadSlowLog = () => {
    data.loading = true
    browserStore
        .getSlowLog(props.server, data.listLimit)
        .then((list) => {
            data.list = list || []
        })
        .finally(async () => {
            data.loading = false
            await nextTick()
            tableRef.value?.scrollTo({ position: data.sortOrder === 'ascend' ? 'bottom' : 'top' })
        })
}
const loadSlowLog = debounce(_loadSlowLog, 1000, { leading: true, trailing: true })

const startAutoRefresh = async () => {
    let lastExec = Date.now()
    do {
        if (!autoRefresh.on) {
            break
        }
        await timeout(100)
        if (data.loading || Date.now() - lastExec < autoRefresh.interval * 1000) {
            continue
        }
        lastExec = Date.now()
        loadSlowLog()
    } while (true)
    stopAutoRefresh()
}

const stopAutoRefresh = () => {
    autoRefresh.on = false
}

onMounted(() => loadSlowLog())

onUnmounted(() => stopAutoRefresh())

const onToggleRefresh = (on) => {
    if (on) {
        startAutoRefresh()
    } else {
        stopAutoRefresh()
    }
}

const onListLimitChanged = (limit) => {
    loadSlowLog()
}
</script>

<template>
    <div class="content-log content-container content-value fill-height flex-box-v">
        <n-form :disabled="data.loading" class="flex-item" inline>
            <n-form-item :label="$t('slog.limit')">
                <n-input-number
                    v-model:value="data.listLimit"
                    :max="9999"
                    style="width: 120px"
                    @update:value="onListLimitChanged" />
            </n-form-item>
            <n-form-item :label="$t('slog.filter')">
                <n-input v-model:value="data.keyword" clearable placeholder="" />
            </n-form-item>
            <n-form-item label="&nbsp;">
                <n-popover :delay="500" keep-alive-on-hover placement="bottom" trigger="hover">
                    <template #trigger>
                        <n-button :loading="data.loading" circle size="small" tertiary @click="_loadSlowLog">
                            <template #icon>
                                <n-icon :size="props.size">
                                    <refresh
                                        :class="{ 'auto-rotate': autoRefresh.on }"
                                        :color="autoRefresh.on ? themeVars.primaryColor : undefined"
                                        :stroke-width="autoRefresh.on ? 6 : 3" />
                                </n-icon>
                            </template>
                        </n-button>
                    </template>
                    <auto-refresh-form
                        v-model:interval="autoRefresh.interval"
                        v-model:on="autoRefresh.on"
                        :default-value="5"
                        :loading="data.loading"
                        @toggle="onToggleRefresh" />
                </n-popover>
            </n-form-item>
        </n-form>
        <n-data-table
            ref="tableRef"
            :columns="columns"
            :data="data.list"
            :loading="data.loading"
            class="flex-item-expand"
            flex-height
            striped
            @update:sorter="({ order }) => (data.sortOrder = order)" />
    </div>
</template>

<style lang="scss" scoped>
@import '@/styles/content';
</style>
