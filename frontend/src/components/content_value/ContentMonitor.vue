<script setup>
import { computed, nextTick, onMounted, onUnmounted, reactive, ref } from 'vue'
import { debounce, filter, get, includes, isEmpty, join } from 'lodash'
import { useI18n } from 'vue-i18n'
import { useThemeVars } from 'naive-ui'
import useBrowserStore from 'stores/browser.js'
import Play from '@/components/icons/Play.vue'
import Pause from '@/components/icons/Pause.vue'
import { ExportLog, StartMonitor, StopMonitor } from 'wailsjs/go/services/monitorService.js'
import { ClipboardSetText, EventsOff, EventsOn } from 'wailsjs/runtime/runtime.js'
import Copy from '@/components/icons/Copy.vue'
import Export from '@/components/icons/Export.vue'
import Delete from '@/components/icons/Delete.vue'
import IconButton from '@/components/common/IconButton.vue'
import Bottom from '@/components/icons/Bottom.vue'

const themeVars = useThemeVars()

const browserStore = useBrowserStore()
const i18n = useI18n()
const props = defineProps({
    server: {
        type: String,
    },
})

const data = reactive({
    monitorEvent: '',
    list: [],
    listLimit: 20,
    keyword: '',
    autoShowLast: true,
})

const listRef = ref(null)

onMounted(() => {
    // try to stop prev monitor first
    onStopMonitor()
})

onUnmounted(() => {
    onStopMonitor()
})

const isMonitoring = computed(() => {
    return !isEmpty(data.monitorEvent)
})

const displayList = computed(() => {
    if (!isEmpty(data.keyword)) {
        return filter(data.list, (line) => includes(line, data.keyword))
    }
    return data.list
})

const _scrollToBottom = () => {
    nextTick(() => {
        listRef.value?.scrollTo({ position: 'bottom' })
    })
}
const scrollToBottom = debounce(_scrollToBottom, 1000, { leading: true, trailing: true })

const onStartMonitor = async () => {
    if (isMonitoring.value) {
        return
    }

    const { data: ret, success, msg } = await StartMonitor(props.server)
    if (!success) {
        $message.error(msg)
        return
    }
    data.monitorEvent = get(ret, 'eventName')
    EventsOn(data.monitorEvent, (content) => {
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
const onStopMonitor = async () => {
    const { success, msg } = await StopMonitor(props.server)
    if (!success) {
        $message.error(msg)
        return
    }

    EventsOff(data.monitorEvent)
    data.monitorEvent = ''
}

const onCopyLog = async () => {
    try {
        const content = join(data.list, '\n')
        const succ = await ClipboardSetText(content)
        if (succ) {
            $message.success(i18n.t('interface.copy_succ'))
        }
    } catch (e) {
        $message.error(e.message)
    }
}

const onExportLog = () => {
    ExportLog(data.list)
}

const onCleanLog = () => {
    data.list = []
}
</script>

<template>
    <div class="content-log content-container fill-height flex-box-v">
        <n-form class="flex-item" label-align="left" label-placement="left" label-width="auto" size="small">
            <n-form-item :feedback="$t('monitor.warning')" :label="$t('monitor.actions')">
                <n-space :wrap="false" :wrap-item="false" style="width: 100%">
                    <n-button
                        v-if="!isMonitoring"
                        :focusable="false"
                        secondary
                        strong
                        type="success"
                        @click="onStartMonitor">
                        <template #icon>
                            <n-icon :component="Play" size="18" />
                        </template>
                        {{ $t('monitor.start') }}
                    </n-button>
                    <n-button v-else :focusable="false" secondary strong type="warning" @click="onStopMonitor">
                        <template #icon>
                            <n-icon :component="Pause" size="18" />
                        </template>
                        {{ $t('monitor.stop') }}
                    </n-button>
                    <n-button-group>
                        <icon-button
                            :icon="Copy"
                            border
                            size="18"
                            stroke-width="3.5"
                            t-tooltip="monitor.copy_log"
                            @click="onCopyLog" />
                        <icon-button
                            :icon="Export"
                            border
                            size="18"
                            stroke-width="3.5"
                            t-tooltip="monitor.save_log"
                            @click="onExportLog" />
                    </n-button-group>
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
                        t-tooltip="monitor.clean_log"
                        @click="onCleanLog" />
                </n-space>
            </n-form-item>
            <n-form-item :label="$t('monitor.search')">
                <n-input v-model:value="data.keyword" clearable placeholder="" />
            </n-form-item>
        </n-form>
        <n-virtual-list ref="listRef" :item-size="25" :items="displayList" class="list-wrapper">
            <template #default="{ item }">
                <div class="line-item content-value">
                    <b>&gt;</b>
                    {{ item }}
                </div>
            </template>
        </n-virtual-list>
    </div>
</template>

<style lang="scss" scoped>
@import '@/styles/content';

.line-item {
    margin-bottom: 5px;
}

.list-wrapper {
    background-color: v-bind('themeVars.codeColor');
    border: solid 1px v-bind('themeVars.borderColor');
    border-radius: 3px;
    padding: 5px 10px;
    box-sizing: border-box;
}
</style>
