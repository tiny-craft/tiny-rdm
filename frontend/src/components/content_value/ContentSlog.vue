<script setup>
import { h, onMounted, onUnmounted, reactive, ref } from 'vue'
import Refresh from '@/components/icons/Refresh.vue'
import { debounce, isEmpty, map, size, split } from 'lodash'
import { useI18n } from 'vue-i18n'
import { useThemeVars } from 'naive-ui'
import dayjs from 'dayjs'
import useBrowserStore from 'stores/browser.js'

const themeVars = useThemeVars()

const browserStore = useBrowserStore()
const i18n = useI18n()
const props = defineProps({
    server: {
        type: String,
    },
    db: {
        type: Number,
        default: 0,
    },
})

const data = reactive({
    list: [],
    sortOrder: 'descend',
    listLimit: 20,
    loading: false,
    autoLoading: false,
    client: '',
    keyword: '',
})

const tableRef = ref(null)

const _loadSlowLog = () => {
    data.loading = true
    browserStore
        .getSlowLog(props.server, props.db, data.listLimit)
        .then((list) => {
            data.list = list || []
        })
        .finally(() => {
            data.loading = false
            tableRef.value?.scrollTo({ top: data.sortOrder === 'ascend' ? 999999 : 0 })
        })
}
const loadSlowLog = debounce(_loadSlowLog, 1000, { leading: true, trailing: true })

let intervalID
onMounted(() => {
    loadSlowLog()
    intervalID = setInterval(() => {
        if (data.autoLoading === true) {
            loadSlowLog()
        }
    }, 5000)
})

onUnmounted(() => {
    clearInterval(intervalID)
})

const onListLimitChanged = (limit) => {
    loadSlowLog()
}
</script>

<template>
    <n-card
        :bordered="false"
        :theme-overrides="{ borderRadius: '0px' }"
        :title="$t('slog.title')"
        class="content-container flex-box-v"
        content-style="display: flex;flex-direction: column; overflow: hidden; backgroundColor: gray">
        <n-form :disabled="data.loading" class="flex-item" inline>
            <n-form-item :label="$t('slog.limit')">
                <n-input-number
                    v-model:value="data.listLimit"
                    :max="9999"
                    style="width: 120px"
                    @update:value="onListLimitChanged" />
            </n-form-item>
            <n-form-item :label="$t('slog.auto_refresh')">
                <n-switch v-model:value="data.autoLoading" :loading="data.loading" />
            </n-form-item>
            <n-form-item label="&nbsp;">
                <n-tooltip>
                    {{ $t('slog.refresh') }}
                    <template #trigger>
                        <n-button :loading="data.loading" circle size="small" tertiary @click="_loadSlowLog">
                            <template #icon>
                                <n-icon :component="Refresh" />
                            </template>
                        </n-button>
                    </template>
                </n-tooltip>
            </n-form-item>
            <n-form-item :label="$t('slog.filter')">
                <n-input v-model:value="data.keyword" clearable placeholder="" />
            </n-form-item>
        </n-form>
        <div class="content-value fill-height flex-box-h">
            <n-data-table
                ref="tableRef"
                :columns="[
                    {
                        title: $t('slog.exec_time'),
                        key: 'timestamp',
                        sortOrder: data.sortOrder,
                        sorter: 'default',
                        width: 180,
                        align: 'center',
                        titleAlign: 'center',
                        render({ timestamp }, index) {
                            return dayjs(timestamp).format('YYYY-MM-DD HH:mm:ss')
                        },
                    },
                    {
                        title: $t('slog.client'),
                        key: 'client',
                        filterOptionValue: data.client,
                        resizable: true,
                        filter(value, row) {
                            return value === '' || row.client === value.toString() || row.addr === value.toString()
                        },
                        width: 200,
                        align: 'center',
                        titleAlign: 'center',
                        ellipsis: true,
                        render({ client, addr }, index) {
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
                        title: $t('slog.cmd'),
                        key: 'cmd',
                        titleAlign: 'center',
                        filterOptionValue: data.keyword,
                        resizable: true,
                        width: 100,
                        filter(value, row) {
                            return value === '' || !!~row.cmd.indexOf(value.toString())
                        },
                        render({ cmd }, index) {
                            const cmdList = split(cmd, '\n')
                            if (size(cmdList) > 1) {
                                return h(
                                    'div',
                                    null,
                                    map(cmdList, (c) => h('div', null, c)),
                                )
                            }
                            return cmd
                        },
                    },
                    {
                        title: $t('slog.cost_time'),
                        key: 'cost',
                        width: 100,
                        align: 'center',
                        titleAlign: 'center',
                        render({ cost }, index) {
                            const ms = dayjs.duration(cost).asMilliseconds()
                            if (ms < 1000) {
                                return `${ms} ms`
                            } else {
                                return `${Math.floor(ms / 1000)} s`
                            }
                        },
                    },
                ]"
                :data="data.list"
                class="flex-item-expand"
                flex-height
                @update:sorter="({ order }) => (data.sortOrder = order)" />
        </div>
    </n-card>
</template>

<style lang="scss" scoped>
@import '@/styles/content';

.content-container {
    padding: 5px;
    box-sizing: border-box;
}
</style>
