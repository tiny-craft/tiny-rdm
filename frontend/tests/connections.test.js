import assert from 'node:assert/strict'
import { readFileSync } from 'node:fs'
import test from 'node:test'
import vm from 'node:vm'
import { ConnectionType } from '../src/consts/connection_type.js'

// Exercise the actual component/store scripts without requiring a browser or Wails.
const loadScript = (path, bindings, exports) => {
    const source = readFileSync(new URL(path, import.meta.url), 'utf8')
    let script = source.includes('<script setup>') ? source.split('<script setup>')[1].split('</script>')[0] : source
    const imports = {}
    script = script.replace(/^import\s+([\s\S]*?)\s+from\s+['"][^'"]+['"]\s*$/gm, (_, names) => {
        for (const name of names.replace(/[{}]/g, '').split(',')) {
            imports[name.trim()] = () => ({})
        }
        return ''
    })
    script = script.replace(/^export default .+$/gm, '')
    const exposed = exports.map((name) => `${name}: typeof ${name} === 'undefined' ? undefined : ${name}`)
    return vm.runInNewContext(`${script}\n;({ ${exposed.join(', ')} })`, { ...imports, ...bindings })
}

const server = (name, group = '') => ({
    key: `${group}/${name}`,
    label: name,
    name,
    type: ConnectionType.Server,
})
const group = (name, children = []) => ({ key: `${name}/`, label: name, type: ConnectionType.Group, children })
const plain = (value) => JSON.parse(JSON.stringify(value))

const loadTree = (connections) => {
    let saves = 0
    const store = {
        connections,
        connectionSortOrder: 'original',
        saveConnectionSorted: async () => ({ success: true }),
    }
    const tree = loadScript(
        '../src/components/sidebar/ConnectionTree.vue',
        {
            ConnectionType,
            defineProps: () => ({}),
            ref: (value) => ({ value }),
            reactive: (value) => value,
            useConnectionStore: () => store,
            debounce: () => () => saves++,
        },
        ['handleDrop', 'selectedKeys', 'allowDrop', 'draggingNode'],
    )
    return { ...tree, store, saves: () => saves }
}

for (const position of ['before', 'after']) {
    test(`groups can move ${position} another group with their children intact`, () => {
        const a = group('A', [server('a', 'A')])
        const b = group('B', [server('b', 'B')])
        const tree = loadTree(position === 'before' ? [a, b] : [b, a])
        tree.handleDrop({ node: a, dragNode: b, dropPosition: position })
        assert.deepEqual(
            plain(tree.store.connections.map((node) => node.label)),
            position === 'before' ? ['B', 'A'] : ['A', 'B'],
        )
        assert.equal(a.children.length, 1)
        assert.equal(b.children.length, 1)
        assert.equal(tree.saves(), 1)
    })
}

for (const draggedType of ['group', 'server']) {
    test(`${draggedType} cannot be dropped inside a server`, () => {
        const target = server('root')
        const dragged = draggedType === 'group' ? group('A', [server('a', 'A')]) : server('a')
        const tree = loadTree([target, dragged])
        const before = plain(tree.store.connections)
        tree.handleDrop({ node: target, dragNode: dragged, dropPosition: 'inside' })
        assert.deepEqual(plain(tree.store.connections), before)
        assert.equal(tree.saves(), 0)
    })
}

test('a group cannot be nested or placed beside a connection inside another group', () => {
    const a = group('A', [server('a', 'A')])
    const b = group('B', [server('b', 'B')])
    const tree = loadTree([a, b])
    const before = plain(tree.store.connections)
    tree.handleDrop({ node: b, dragNode: a, dropPosition: 'inside' })
    for (const position of ['before', 'inside', 'after']) {
        tree.handleDrop({ node: b.children[0], dragNode: a, dropPosition: position })
    }
    assert.deepEqual(plain(tree.store.connections), before)
    assert.equal(tree.saves(), 0)
})

test('moving a connection into a group updates its key and selection', () => {
    const item = server('cache')
    const target = group('production')
    const tree = loadTree([item, target])
    tree.selectedKeys.value = [item.key]
    tree.handleDrop({ node: target, dragNode: item, dropPosition: 'inside' })
    assert.equal(tree.store.connections.length, 1)
    assert.equal(target.children[0], item)
    assert.equal(item.key, 'production/cache')
    assert.deepEqual(plain(tree.selectedKeys.value), ['production/cache'])
})

test('moving a connection back to the root updates its key', () => {
    const item = server('cache', 'production')
    const target = group('production', [item])
    const tree = loadTree([target])
    tree.handleDrop({ node: target, dragNode: item, dropPosition: 'after' })
    assert.equal(item.key, '/cache')
    assert.equal(target.children.length, 0)
    assert.equal(tree.store.connections[1], item)
})

test('moving between groups beside a sibling updates its key', () => {
    const item = server('cache', 'A')
    const sibling = server('other', 'B')
    const a = group('A', [item])
    const b = group('B', [sibling])
    const tree = loadTree([a, b])
    tree.handleDrop({ node: sibling, dragNode: item, dropPosition: 'before' })
    assert.equal(item.key, 'B/cache')
    assert.equal(a.children.length, 0)
    assert.equal(b.children[0], item)
})

test('a missing target or a self drop never removes the dragged connection', () => {
    const item = server('cache')
    const tree = loadTree([item])
    tree.handleDrop({ node: server('missing'), dragNode: item, dropPosition: 'after' })
    tree.handleDrop({ node: item, dragNode: item, dropPosition: 'inside' })
    assert.equal(tree.store.connections.length, 1)
    assert.equal(tree.store.connections[0], item)
    assert.equal(tree.saves(), 0)
})

test('drag indicators reject nesting groups and dropping inside connections', () => {
    const a = group('A')
    const b = group('B', [server('cache', 'B')])
    const root = server('root')
    const tree = loadTree([a, b, root])
    tree.draggingNode.value = a
    assert.equal(tree.allowDrop({ node: b, dropPosition: 'before' }), true)
    assert.equal(tree.allowDrop({ node: b, dropPosition: 'inside' }), false)
    assert.equal(tree.allowDrop({ node: root, dropPosition: 'inside' }), false)
    assert.equal(tree.allowDrop({ node: b.children[0], dropPosition: 'before' }), false)
})

test('all self-drop positions leave a connection in place', () => {
    const item = server('cache')
    const tree = loadTree([item])
    for (const dropPosition of ['before', 'inside', 'after']) {
        tree.handleDrop({ node: item, dragNode: item, dropPosition })
    }
    assert.equal(tree.store.connections[0], item)
    assert.equal(tree.saves(), 0)
})

test('same-level moves adjust the target index after removing the dragged node', () => {
    const a = server('a')
    const b = server('b')
    const c = server('c')
    const tree = loadTree([a, b, c])
    tree.handleDrop({ node: c, dragNode: a, dropPosition: 'after' })
    assert.deepEqual(plain(tree.store.connections.map((node) => node.name)), ['b', 'c', 'a'])
    tree.handleDrop({ node: b, dragNode: a, dropPosition: 'before' })
    assert.deepEqual(plain(tree.store.connections.map((node) => node.name)), ['a', 'b', 'c'])
})

const loadStore = (connections, save, preferences = { behavior: {} }) => {
    const { useConnectionStore: definition } = loadScript(
        '../src/stores/connections.js',
        {
            ConnectionType,
            defineStore: (_, options) => options,
            SaveSortedConnection: save,
            usePreferencesStore: () => preferences,
        },
        ['useConnectionStore'],
    )
    const store = Object.assign(definition.state(), definition.actions, { connections })
    for (const [name, getter] of Object.entries(definition.getters)) {
        Object.defineProperty(store, name, { get: () => getter.call(store, store) })
    }
    return store
}

test('saving a reordered tree waits for persistence and returns its result', async () => {
    let finish
    const store = loadStore([server('cache')], () => new Promise((resolve) => (finish = resolve)))
    const pending = store.saveConnectionSorted()
    finish({ success: false, msg: 'disk full' })
    assert.deepEqual(await pending, { success: false, msg: 'disk full' })
})

for (const order of ['asc', 'desc']) {
    test(`name ${order} display preserves the configured order and uses natural numeric ordering`, () => {
        const connections = [
            server('redis-10'),
            group('B', [server('cache-10', 'B'), server('cache-2', 'B')]),
            server('redis-2'),
            group('A'),
        ]
        const before = plain(connections)
        const preferences = { behavior: { connectionSortOrder: order } }
        const store = loadStore(
            connections,
            () => assert.fail('display sorting must not write connections'),
            preferences,
        )
        assert.ok(store.displayConnections)
        const expected = ['A', 'B', 'redis-2', 'redis-10']
        assert.deepEqual(
            plain(store.displayConnections.map((node) => node.label)),
            order === 'asc' ? expected : expected.reverse(),
        )
        const children = store.displayConnections.find((node) => node.label === 'B').children
        assert.deepEqual(
            plain(children.map((node) => node.name)),
            order === 'asc' ? ['cache-2', 'cache-10'] : ['cache-10', 'cache-2'],
        )
        assert.deepEqual(plain(store.connections), before)
        preferences.behavior.connectionSortOrder = 'original'
        assert.equal(store.displayConnections, connections)
    })
}

test('old preferences default to the configured order', () => {
    const connections = [server('b'), server('a')]
    const store = loadStore(connections)
    assert.equal(store.connectionSortOrder, 'original')
    assert.equal(store.displayConnections, connections)
})

test('selecting a sort mode updates only its preference', async () => {
    const events = []
    const pane = loadScript(
        '../src/components/sidebar/ConnectionPane.vue',
        {
            ref: (value) => ({ value }),
            computed: (getter) => ({
                get value() {
                    return getter()
                },
            }),
            usePreferencesStore: () => ({
                setConnectionSortOrder: async (order) => {
                    events.push(order)
                    return { success: true }
                },
            }),
        },
        ['onSelectOptions'],
    )
    await pane.onSelectOptions('asc')
    assert.deepEqual(events, ['asc'])
})

test('a delayed drag save uses the configured order after switching to a sorted display', async () => {
    const connections = [server('b'), server('a')]
    const preferences = { behavior: { connectionSortOrder: 'original' } }
    let saved
    const store = loadStore(
        connections,
        async (value) => {
            saved = value
            return { success: true }
        },
        preferences,
    )
    preferences.behavior.connectionSortOrder = 'asc'
    assert.deepEqual(plain(store.displayConnections.map((node) => node.name)), ['a', 'b'])
    await store.saveConnectionSorted()
    assert.deepEqual(plain(saved.map((node) => node.name)), ['b', 'a'])
})

test('dragging is disabled in a name-sorted view without changing the configured order', () => {
    const a = server('a')
    const b = server('b')
    const tree = loadTree([a, b])
    tree.store.connectionSortOrder = 'asc'
    tree.handleDrop({ node: a, dragNode: b, dropPosition: 'before' })
    assert.deepEqual(plain(tree.store.connections.map((node) => node.name)), ['a', 'b'])
    assert.equal(tree.saves(), 0)
})

const loadPreferences = (save) => {
    const { usePreferencesStore: definition } = loadScript(
        '../src/stores/preferences.js',
        {
            defineStore: (_, options) => options,
            SetPreferences: () => assert.fail('changing a sort mode must not replace unrelated preferences'),
            UpdatePreferences: save,
            pick: (value, keys) => Object.fromEntries(keys.map((key) => [key, value[key]])),
            TextAlignType: { Center: 0 },
        },
        ['usePreferencesStore'],
    )
    return Object.assign(definition.state(), definition.actions)
}

for (const order of ['original', 'asc', 'desc']) {
    test(`the ${order} mode is persisted in preferences`, async () => {
        let saved
        const preferences = loadPreferences(async (value) => {
            saved = value
            return { success: true }
        })
        assert.equal(typeof preferences.setConnectionSortOrder, 'function')
        preferences.behavior.connectionSortOrder = 'desc'
        const result = await preferences.setConnectionSortOrder(order)
        assert.equal(result.success, true)
        assert.deepEqual(plain(saved), { 'behavior.connectionSortOrder': order })
        assert.equal(preferences.behavior.connectionSortOrder, order)
    })
}

test('a failed preference save keeps the previous mode', async () => {
    const preferences = loadPreferences(async () => ({ success: false, msg: 'disk full' }))
    assert.equal(typeof preferences.setConnectionSortOrder, 'function')
    preferences.behavior.connectionSortOrder = 'desc'
    const result = await preferences.setConnectionSortOrder('original')
    assert.equal(result.success, false)
    assert.equal(preferences.behavior.connectionSortOrder, 'desc')
})
