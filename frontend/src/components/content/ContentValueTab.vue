<script setup>
import { ref } from 'vue'
import { ConnectionType } from '../../consts/connection_type.js'
import Close from '../icons/Close.vue'
import useConnectionStore from '../../stores/connections.js'

const emit = defineEmits(['switchTab', 'closeTab', 'update:modelValue'])

const props = defineProps({
    selectedIndex: {
        type: Number,
        default: 0,
    },
    modelValue: {
        type: Object,
        default: [
            {
                // label: 'tab1',
                // key: 'key',
                // bgColor: 'white',
            },
        ],
    },
})

const connectionStore = useConnectionStore()
const onCurrentSelectChange = ({ type, group = '', server = '', db = 0, key = '' }) => {
    console.log(`group: ${group}\n server: ${server}\n db: ${db}\n key: ${key}`)
    if (type === ConnectionType.RedisValue) {
        // load and update content value
    }
}
// watch(() => databaseStore.currentSelect, throttle(onCurrentSelectChange, 1000))

const items = ref(props.modelValue)
const selIndex = ref(props.selectedIndex)

const onClickTab = (idx, key) => {
    if (idx !== selIndex.value) {
        selIndex.value = idx
        emit('update:modelValue', idx, key)
    }
}

const onCloseTab = (idx, key) => {
    const removed = items.value.splice(idx, 1)
    if (removed.length <= 0) {
        return
    }

    // Update select index if removed index equal current selected
    if (selIndex.value === idx) {
        selIndex.value -= 1
        if (selIndex.value < 0 && items.value.length > 0) {
            selIndex.value = 0
        }
    }
    emit('update:modelValue', items)
    emit('closeTab', idx, key)
}
</script>

<template>
    <!-- TODO: 检查标签是否太多, 左右两边显示左右切换翻页按钮 -->
    <div class="content-tab flex-box-h">
        <div
            v-for="(item, i) in props.modelValue"
            :key="item.key"
            :class="{ 'content-tab_selected': selIndex === i }"
            :style="{ backgroundColor: item.bgColor || '' }"
            :title="item.label"
            class="content-tab_item flex-item-expand icon-btn flex-box-h"
            @click="onClickTab(i, item.key)"
        >
            <n-icon :component="Close" class="content-tab_item-close" size="20" @click.stop="onCloseTab(i, item.key)" />
            <div class="content-tab_item-label ellipsis flex-item-expand">
                {{ item.label }}
            </div>
        </div>
    </div>
</template>

<style lang="scss" scoped>
.content-tab {
    align-items: center;
    //justify-content: center;
    width: 100%;
    height: 40px;
    overflow: hidden;
    font-size: 14px;

    &_item {
        flex: 1 0;
        overflow: hidden;
        align-items: center;
        justify-content: center;
        gap: 3px;
        height: 100%;
        box-sizing: border-box;
        background-color: var(--bg-color-page);
        color: var(--text-color-secondary);
        padding: 0 5px;

        //border-top: var(--el-border-color) 1px solid;
        border-right: var(--border-color) 1px solid;
        transition: all var(--transition-duration-fast) var(--transition-function-ease-in-out-bezier);

        &-label {
            text-align: center;
        }

        &-close {
            //display: none;
            display: inline-flex;
            width: 0;
            transition: width 0.3s;

            &:hover {
                background-color: rgb(176, 177, 182, 0.4);
            }
        }

        &:hover {
            .content-tab_item-close {
                //display: block;
                width: 20px;
                transition: width 0.3s;
            }
        }
    }

    &_selected {
        border-top: #409eff 4px solid !important;
        background-color: #ffffff;
        color: #303133;
    }
}
</style>
