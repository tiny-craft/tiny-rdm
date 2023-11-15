<script setup>
import { computed, h, nextTick, reactive, ref } from 'vue'
import IconButton from '@/components/common/IconButton.vue'
import Refresh from '@/components/icons/Refresh.vue'
import { map, size, split, uniqBy } from 'lodash'
import { useI18n } from 'vue-i18n'
import Delete from '@/components/icons/Delete.vue'
import dayjs from 'dayjs'
import { useThemeVars } from 'naive-ui'
import useBrowserStore from 'stores/browser.js'

const themeVars = useThemeVars()

const browserStore = useBrowserStore()
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
        label: i18n.t('common.all'),
        value: '',
    })
    return options
})

const tableRef = ref(null)

const loadHistory = () => {
    data.loading = true
    browserStore
        .getCmdHistory()
        .then((list) => {
            data.history = list || []
        })
        .finally(() => {
            data.loading = false
            tableRef.value?.scrollTo({ top: 999999 })
        })
}

const cleanHistory = async () => {
    $dialog.warning(i18n.t('log.confirm_clean_log'), () => {
        data.loading = true
        browserStore
            .cleanCmdHistory()
            .then((success) => {
                if (success) {
                    data.history = []
                    tableRef.value?.scrollTo({ top: 0 })
                    $message.success(i18n.t('common.success'))
                }
            })
            .finally(() => {
                data.loading = false
            })
    })
}

defineExpose({
    refresh: () => nextTick().then(loadHistory),
})
</script>

<template>
    <n-card
        :bordered="false"
        :theme-overrides="{ borderRadius: '0px' }"
        :title="$t('log.title')"
        class="content-container flex-box-v"
        content-style="display: flex;flex-direction: column; overflow: hidden; backgroundColor: gray">
        <n-form :disabled="data.loading" class="flex-item" inline>
            <n-form-item :label="$t('log.filter_server')">
                <n-select
                    v-model:value="data.server"
                    :consistent-menu-width="false"
                    :options="filterServerOption"
                    style="min-width: 100px" />
            </n-form-item>
            <n-form-item :label="$t('log.filter_keyword')">
                <n-input v-model:value="data.keyword" clearable placeholder="" />
            </n-form-item>
            <n-form-item label="&nbsp;">
                <icon-button :icon="Refresh" border t-tooltip="log.refresh" @click="loadHistory" />
            </n-form-item>
            <n-form-item label="&nbsp;">
                <icon-button :icon="Delete" border t-tooltip="log.clean_log" @click="cleanHistory" />
            </n-form-item>
        </n-form>
        <div class="content-value fill-height flex-box-h">
            <n-data-table
                ref="tableRef"
                :columns="[
                    {
                        title: $t('log.exec_time'),
                        key: 'timestamp',
                        defaultSortOrder: 'ascend',
                        sorter: 'default',
                        width: 180,
                        align: 'center',
                        titleAlign: 'center',
                        render: ({ timestamp }, index) => {
                            return dayjs(timestamp).format('YYYY-MM-DD HH:mm:ss')
                        },
                    },
                    {
                        title: $t('log.server'),
                        key: 'server',
                        filterOptionValue: data.server,
                        filter: (value, row) => {
                            return value === '' || row.server === value.toString()
                        },
                        width: 150,
                        align: 'center',
                        titleAlign: 'center',
                        ellipsis: true,
                    },
                    {
                        title: $t('log.cmd'),
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
                                    map(cmdList, (c) => h('div', null, c)),
                                )
                            }
                            return cmd
                        },
                    },
                    {
                        title: $t('log.cost_time'),
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
                ]"
                :data="data.history"
                class="flex-item-expand"
                flex-height />
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
