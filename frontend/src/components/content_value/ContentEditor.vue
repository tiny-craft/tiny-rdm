<script setup>
import { computed, nextTick, onMounted, onUnmounted, ref, watch } from 'vue'
import * as monaco from 'monaco-editor'
import usePreferencesStore from 'stores/preferences.js'
import { useThemeVars } from 'naive-ui'
import { isEmpty } from 'lodash'

const props = defineProps({
    content: {
        type: String,
    },
    language: {
        type: String,
        default: 'json',
    },
    readonly: {
        type: String,
    },
    loading: {
        type: Boolean,
    },
    border: {
        type: Boolean,
        default: false,
    },
    offsetKey: {
        type: String,
        default: '',
    },
    keepOffset: {
        type: Boolean,
        default: false,
    },
})

const emit = defineEmits(['reset', 'input', 'save'])

const themeVars = useThemeVars()
/** @type {HTMLElement|null} */
const editorRef = ref(null)
/** @type monaco.editor.IStandaloneCodeEditor */
let editorNode = null
const scrollOffset = { top: 0, left: 0 }

const updateScroll = () => {
    if (editorNode != null) {
        if (props.keepOffset && !isEmpty(props.offsetKey)) {
            editorNode.setScrollPosition({ scrollTop: scrollOffset.top, scrollLeft: scrollOffset.left })
        } else {
            // reset offset if not needed
            editorNode.setScrollPosition({ scrollTop: 0, scrollLeft: 0 })
        }
    }
}

const destroyEditor = () => {
    if (editorNode != null && editorNode.dispose != null) {
        const model = editorNode.getModel()
        if (model != null) {
            model.dispose()
        }
        editorNode.dispose()
        editorNode = null
    }
}

const readonlyValue = computed(() => {
    return props.readonly || props.loading
})

const pref = usePreferencesStore()
onMounted(async () => {
    if (editorRef.value != null) {
        const { fontSize, fontFamily = ['monaco'] } = pref.editorFont
        editorNode = monaco.editor.create(editorRef.value, {
            // value: props.content,
            theme: pref.isDark ? 'rdm-dark' : 'rdm-light',
            language: props.language,
            lineNumbers: pref.showLineNum ? 'on' : 'off',
            links: pref.editorLinks,
            readOnly: readonlyValue.value,
            colorDecorators: true,
            accessibilitySupport: 'off',
            wordWrap: 'on',
            tabSize: 2,
            folding: pref.showFolding,
            dragAndDrop: pref.dropText,
            fontFamily,
            fontSize,
            scrollBeyondLastLine: false,
            automaticLayout: true,
            scrollbar: {
                useShadows: false,
                verticalScrollbarSize: '10px',
            },
            // formatOnType: true,
            contextmenu: false,
            lineNumbersMinChars: 2,
            lineDecorationsWidth: 0,
            minimap: {
                enabled: false,
            },
            selectionHighlight: false,
            renderLineHighlight: 'gutter',
        })

        // add shortcut for save
        editorNode.addCommand(monaco.KeyMod.CtrlCmd | monaco.KeyCode.KeyS, (event) => {
            emit('save')
        })

        editorNode.onDidScrollChange((event) => {
            // save scroll offset when changes, ie. content changes
            if (props.keepOffset && !event.scrollHeightChanged) {
                scrollOffset.top = event.scrollTop
                scrollOffset.left = event.scrollLeft
            }
        })

        editorNode.onDidLayoutChange((event) => {
            updateScroll()
        })

        // editorNode.onDidChangeModelLanguageConfiguration(() => {
        //     editorNode?.getAction('editor.action.formatDocument')?.run()
        // })

        if (editorNode.onDidChangeModelContent) {
            editorNode.onDidChangeModelContent(() => {
                emit('input', editorNode.getValue())
            })
        }
    }
})

watch(
    () => props.content,
    async (content) => {
        if (editorNode != null) {
            editorNode.setValue(content)
            await nextTick(() => emit('reset', content))
            updateScroll()
        }
    },
)

watch(
    () => props.offsetKey,
    () => {
        // reset scroll offset when key changed
        if (editorNode != null) {
            scrollOffset.top = 0
            scrollOffset.left = 0
            editorNode.setScrollPosition({ scrollTop: 0, scrollLeft: 0 })
        }
    },
)

watch(
    () => readonlyValue.value,
    (readOnly) => {
        if (editorNode != null) {
            editorNode.updateOptions({
                readOnly,
            })
        }
    },
)

watch(
    () => props.language,
    (language) => {
        if (editorNode != null) {
            const model = editorNode.getModel()
            if (model != null) {
                monaco.editor.setModelLanguage(model, language)
            }
        }
    },
)

watch(
    () => pref.isDark,
    (dark) => {
        if (editorNode != null) {
            editorNode.updateOptions({
                theme: dark ? 'rdm-dark' : 'rdm-light',
            })
        }
    },
)

watch(
    () => pref.editor,
    ({ showLineNum = true, showFolding = true, dropText = true, links = true }) => {
        if (editorNode != null) {
            const { fontSize, fontFamily } = pref.editorFont
            editorNode.updateOptions({
                fontSize,
                fontFamily,
                lineNumbers: showLineNum ? 'on' : 'off',
                folding: showFolding,
                dragAndDrop: dropText,
                links,
            })
        }
    },
    { deep: true },
)

onUnmounted(() => {
    destroyEditor()
})
</script>

<template>
    <div :class="{ 'editor-border': props.border === true }" style="position: relative">
        <div ref="editorRef" class="editor-inst" />
    </div>
</template>

<style lang="scss" scoped>
.editor-border {
    border: 1px solid v-bind('themeVars.borderColor');
    border-radius: v-bind('themeVars.borderRadius');
    padding: 3px;
    box-sizing: border-box;
}

.editor-inst {
    position: absolute;
    top: 2px;
    bottom: 2px;
    left: 2px;
    right: 2px;
}
</style>
