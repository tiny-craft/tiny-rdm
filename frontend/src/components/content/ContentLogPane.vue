<script setup>
import { computed, nextTick, reactive, ref } from 'vue'
import IconButton from '@/components/common/IconButton.vue'
import Refresh from '@/components/icons/Refresh.vue'
import useConnectionStore from 'stores/connections.js'
import { map, uniqBy } from 'lodash'
import { useI18n } from 'vue-i18n'
import dayjs from 'dayjs'

const connectionStore = useConnectionStore()
const i18n = useI18n()
const data = reactive({
    loading: false,
    server: '',
    keyword: '',
    history: [],
})
const filterServerOption = computed(() => {
    const serverSet = uniqBy(data.history, 'server')
    const options = map(serverSet, ({ server }) => ({
        label: server,
        value: server,
    }))
    options.splice(0, 0, {
        label: i18n.t('all'),
        value: '',
    })
    return options
})

const tableRef = ref(null)

const loadHistory = () => {
    data.loading = true
    connectionStore
        .getCmdHistory()
        .then((list) => {
            data.history = list || []
        })
        .finally(() => {
            data.loading = false
            tableRef.value?.scrollTo({ top: 999999 })
        })
}

defineExpose({
    refresh: () => nextTick().then(loadHistory),
})
</script>

<template>
    <n-card
        :title="$t('launch_log')"
        class="content-container flex-box-v"
        content-style="display: flex;flex-direction: column; overflow: hidden;"
    >
        <n-form :disabled="data.loading" class="flex-item" inline>
            <n-form-item :label="$t('filter_server')">
                <n-select
                    v-model:value="data.server"
                    :consistent-menu-width="false"
                    :options="filterServerOption"
                    style="min-width: 100px"
                />
            </n-form-item>
            <n-form-item :label="$t('filter_keyword')">
                <n-input v-model:value="data.keyword" clearable placeholder="" />
            </n-form-item>
            <n-form-item>
                <icon-button :icon="Refresh" border t-tooltip="refresh" @click="loadHistory" />
            </n-form-item>
        </n-form>
        <div class="fill-height flex-box-h" style="user-select: text">
            <n-data-table
                ref="tableRef"
                :columns="[
                    {
                        title: $t('exec_time'),
                        key: 'timestamp',
                        defaultSortOrder: 'ascend',
                        sorter: 'default',
                        width: 180,
                        align: 'center',
                        titleAlign: 'center',
                        render({ timestamp }, index) {
                            return dayjs(timestamp).locale('zh-cn').format('YYYY-MM-DD hh:mm:ss')
                        },
                    },
                    {
                        title: $t('server'),
                        key: 'server',
                        filterOptionValue: data.server,
                        filter(value, row) {
                            return value === '' || row.server === value.toString()
                        },
                        width: 150,
                        align: 'center',
                        titleAlign: 'center',
                        ellipsis: true,
                    },
                    {
                        title: $t('cmd'),
                        key: 'cmd',
                        titleAlign: 'center',
                        filterOptionValue: data.keyword,
                        resizable: true,
                        filter(value, row) {
                            return value === '' || !!~row.cmd.indexOf(value.toString())
                        },
                    },
                    {
                        title: $t('cost_time'),
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
                :data="data.history"
                class="flex-item-expand"
                flex-height
            />
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
