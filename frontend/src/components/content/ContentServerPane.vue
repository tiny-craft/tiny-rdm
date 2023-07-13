<script setup>
import useDialog from '../../stores/dialog.js'
import AddLink from '../icons/AddLink.vue'
import { computed, onMounted, onUnmounted, ref, watch } from 'vue'
import useConnectionStore from '../../stores/connections.js'
import { isEmpty } from 'lodash'
import ContentServerStatus from '../content_value/ContentServerStatus.vue'

const dialogStore = useDialog()
const connectionStore = useConnectionStore()
const serverInfo = ref({})
const autoRefresh = ref(true)

const refreshInfo = async () => {
    const server = connectionStore.selectedServer
    if (!isEmpty(server) && connectionStore.isConnected(server)) {
        serverInfo.value = await connectionStore.getServerInfo(server)
    }
}

let intervalId
onMounted(() => {
    intervalId = setInterval(() => {
        if (autoRefresh.value) {
            refreshInfo()
        }
    }, 5000)
})

onUnmounted(() => {
    clearInterval(intervalId)
})

watch(() => connectionStore.selectedServer, refreshInfo)

const hasContent = computed(() => !isEmpty(serverInfo.value))
</script>

<template>
    <div class="content-container flex-box-v" :style="{ 'justify-content': hasContent ? 'flex-start' : 'center' }">
        <!-- TODO: replace icon to app icon -->
        <n-empty v-if="!hasContent" :description="$t('empty_server_content')">
            <template #extra>
                <n-button @click="dialogStore.openNewDialog()">
                    <template #icon>
                        <n-icon :component="AddLink" size="18" />
                    </template>
                    {{ $t('new_conn') }}
                </n-button>
            </template>
        </n-empty>
        <content-server-status
            v-else
            v-model:auto-refresh="autoRefresh"
            :server="connectionStore.selectedServer"
            :info="serverInfo"
        />
    </div>
</template>

<style lang="scss" scoped>
@import 'content';

.content-container {
    //justify-content: center;
    padding: 5px;
    box-sizing: border-box;
}

.color-preset-item {
    width: 24px;
    height: 24px;
    margin-right: 2px;
    border: white 3px solid;
    cursor: pointer;

    &_selected,
    &:hover {
        border-color: #cdd0d6;
    }
}
</style>
