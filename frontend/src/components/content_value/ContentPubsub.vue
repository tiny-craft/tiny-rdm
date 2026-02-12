<script setup>
import { computed, h, nextTick, onMounted, onUnmounted, reactive, ref } from 'vue'
import { debounce, get, isEmpty, size, uniq } from 'lodash'
import { useI18n } from 'vue-i18n'
import { useThemeVars } from 'naive-ui'
import useBrowserStore from 'stores/browser.js'
import { ClipboardSetText, EventsOff, EventsOn } from 'wailsjs/runtime/runtime.js'
import dayjs from 'dayjs'
import Publish from '@/components/icons/Publish.vue'
import Subscribe from '@/components/icons/Subscribe.vue'
import Pause from '@/components/icons/Pause.vue'
import Delete from '@/components/icons/Delete.vue'
import { Publish as PublishSend, StartSubscribe, StopSubscribe } from 'wailsjs/go/services/pubsubService.js'
import Checked from '@/components/icons/Checked.vue'
import Bottom from '@/components/icons/Bottom.vue'
import IconButton from '@/components/common/IconButton.vue'
import EditableTableColumn from '@/components/common/EditableTableColumn.vue'

const themeVars = useThemeVars()

const browserStore = useBrowserStore()
const i18n = useI18n()
const props = defineProps({
    server: {
        type: String,
    },
})

const data = reactive({
    subscribeEvent: '',
    list: [],
    keyword: '',
    autoShowLast: true,
    ellipsisMessage: false,
    channelHistory: [],
})

const publishData = reactive({
    channel: '',
    message: '',
    received: 0,
    lastShowReceived: -1,
})

const tableRef = ref(null)

const columns = computed(() => [
    {
        title: () => i18n.t('pubsub.time'),
        key: 'timestamp',
        width: 200,
        align: 'center',
        titleAlign: 'center',
        render: ({ timestamp }, index) => {
            return dayjs(timestamp).format('YYYY-MM-DD HH:mm:ss.SSS')
        },
    },
    {
        title: () => i18n.t('pubsub.channel'),
        key: 'channel',
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
    },
    {
        title: () => i18n.t('pubsub.message'),
        key: 'message',
        titleAlign: 'center',
        filterOptionValue: data.keyword,
        resizable: true,
        className: 'content-value',
        ellipsis: data.ellipsisMessage
            ? {
                  tooltip: {
                      style: {
                          maxWidth: '50vw',
                          maxHeight: '50vh',
                      },
                      scrollable: true,
                  },
              }
            : undefined,
        filter: (value, row) => {
            return value === '' || !!~row.cmd.indexOf(value.toString())
        },
    },
    {
        title: () => i18n.t('interface.action'),
        key: 'action',
        width: 60,
        titleAlign: 'center',
        align: 'center',
        fixed: 'right',
        render: (row, index) => {
            return h(EditableTableColumn, {
                editing: false,
                readonly: true,
                canRefresh: false,
                canDelete: false,
                onCopy: async () => {
                    await ClipboardSetText(row.message)
                    $message.success(i18n.t('interface.copy_succ'))
                },
            })
        },
    },
])

onMounted(() => {
    // try to stop prev subscribe first
    onStopSubscribe()
})

onUnmounted(() => {
    onStopSubscribe()
})

const isSubscribing = computed(() => {
    return !isEmpty(data.subscribeEvent)
})

const publishEnable = computed(() => {
    return !isEmpty(publishData.channel)
})

const _scrollToBottom = () => {
    nextTick(() => {
        tableRef.value?.scrollTo({ position: 'bottom' })
    })
}
const scrollToBottom = debounce(_scrollToBottom, 300, { leading: true, trailing: true })

const onStartSubscribe = async () => {
    if (isSubscribing.value) {
        return
    }

    const { data: ret, success, msg } = await StartSubscribe(props.server)
    if (!success) {
        $message.error(msg)
        return
    }
    data.subscribeEvent = get(ret, 'eventName')
    EventsOn(data.subscribeEvent, (content) => {
        if (content instanceof Array) {
            data.list.push(...content)
        } else {
            data.list.push(content)
        }
        if (data.autoShowLast) {
            scrollToBottom()
        }
    })
}
const onStopSubscribe = async () => {
    const { success, msg } = await StopSubscribe(props.server)
    if (!success) {
        $message.error(msg)
        return
    }

    EventsOff(data.subscribeEvent)
    data.subscribeEvent = ''
}

const onCleanLog = () => {
    data.list = []
}

const onPublish = async () => {
    if (isEmpty(publishData.channel)) {
        return
    }

    const {
        success,
        msg,
        data: { received = 0 },
    } = await PublishSend(props.server, publishData.channel, publishData.message || '')
    if (!success) {
        publishData.received = 0
        if (!isEmpty(msg)) {
            $message.error(msg)
        }
        return
    }
    publishData.message = ''
    publishData.received = received
    publishData.lastShowReceived = Date.now()
    // save channel history
    data.channelHistory = uniq(data.channelHistory.concat(publishData.channel))

    // hide send status after 2 seconds
    setTimeout(() => {
        if (publishData.lastShowReceived > 0 && Date.now() - publishData.lastShowReceived > 2000) {
            publishData.lastShowReceived = -1
        }
    }, 2100)
}
</script>

<template>
    <div class="content-log content-container fill-height flex-box-v">
        <n-form class="flex-item" label-align="left" label-placement="left" label-width="auto" size="small">
            <n-form-item :show-label="false">
                <n-space :wrap="false" :wrap-item="false" style="width: 100%">
                    <n-button
                        v-if="!isSubscribing"
                        :focusable="false"
                        secondary
                        strong
                        type="success"
                        @click="onStartSubscribe">
                        <template #icon>
                            <n-icon :component="Subscribe" size="18" />
                        </template>
                        {{ $t('pubsub.subscribe') }}
                    </n-button>
                    <n-button v-else :focusable="false" secondary strong type="warning" @click="onStopSubscribe">
                        <template #icon>
                            <n-icon :component="Pause" size="18" />
                        </template>
                        {{ $t('pubsub.unsubscribe') }}
                    </n-button>
                    <icon-button
                        :icon="Bottom"
                        :secondary="data.autoShowLast"
                        :type="data.autoShowLast ? 'primary' : 'default'"
                        border
                        size="18"
                        stroke-width="3.5"
                        t-tooltip="monitor.always_show_last"
                        @click="data.autoShowLast = !data.autoShowLast" />
                    <div class="flex-item-expand" />
                    <icon-button
                        :icon="Delete"
                        border
                        size="18"
                        stroke-width="3.5"
                        t-tooltip="pubsub.clear"
                        @click="onCleanLog" />
                </n-space>
            </n-form-item>
        </n-form>
        <n-data-table
            ref="tableRef"
            :columns="columns"
            :data="data.list"
            :loading="data.loading"
            class="flex-item-expand"
            flex-height
            size="small"
            virtual-scroll />
        <div class="total-message">{{ $t('pubsub.receive_message', { total: size(data.list) }) }}</div>
        <div class="flex-box-h publish-input">
            <n-input-group>
                <n-auto-complete
                    v-model:value="publishData.channel"
                    :get-show="() => true"
                    :options="data.channelHistory"
                    :placeholder="$t('pubsub.channel')"
                    style="width: 35%; max-width: 200px"
                    @keydown.enter="onPublish" />
                <n-input
                    v-model:value="publishData.message"
                    :placeholder="$t('pubsub.message')"
                    @keydown.enter="onPublish">
                    <template #suffix>
                        <transition mode="out-in" name="fade">
                            <n-tag v-show="publishData.lastShowReceived > 0" bordered size="small" type="success">
                                <template #icon>
                                    <n-icon :component="Checked" size="16" />
                                </template>
                                {{ publishData.received }}
                            </n-tag>
                        </transition>
                    </template>
                </n-input>
            </n-input-group>
            <n-button :disabled="!publishEnable" type="info" @click="onPublish">
                <template #icon>
                    <n-icon :component="Publish" size="18" />
                </template>
                {{ $t('pubsub.publish') }}
            </n-button>
        </div>
    </div>
</template>

<style lang="scss" scoped>
@use '@/styles/content';

.total-message {
    margin: 10px 0 0;
}

.publish-input {
    margin: 10px 0 0;
    gap: 10px;
}
</style>
