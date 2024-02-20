/**
 * all redis type
 * @enum {string}
 */
export const types = {
    STRING: 'STRING',
    HASH: 'HASH',
    LIST: 'LIST',
    SET: 'SET',
    ZSET: 'ZSET',
    STREAM: 'STREAM',
    JSON: 'JSON',
}

export const typesShortName = {
    STRING: 'S',
    HASH: 'H',
    LIST: 'L',
    SET: 'E',
    ZSET: 'Z',
    STREAM: 'X',
    JSON: 'J',
}

/**
 * mark color for redis types
 * @enum {string}
 */
export const typesColor = {
    [types.STRING]: '#8B5CF6',
    [types.HASH]: '#3B82F6',
    [types.LIST]: '#10B981',
    [types.SET]: '#F59E0B',
    [types.ZSET]: '#EF4444',
    [types.STREAM]: '#EC4899',
    [types.JSON]: '#828766',
}

/**
 * background mark color for redis types
 * @enum {string}
 */
export const typesBgColor = {
    [types.STRING]: '#F2EDFB',
    [types.HASH]: '#E4F0FC',
    [types.LIST]: '#E3F3EB',
    [types.SET]: '#FDF1DF',
    [types.ZSET]: '#FAEAED',
    [types.STREAM]: '#FDE6F1',
    [types.JSON]: '#ECECD9',
}

// export const typesName = Object.fromEntries(Object.entries(types).map(([key, value]) => [key, value.name]))

export const validType = (t) => {
    return types.hasOwnProperty(t)
}

/**
 * icon type in browser tree
 * @enum {string}
 */
export const typesIconStyle = {
    SHORT: 0,
    FULL: 1,
    POINT: 2,
    ICON: 3,
}
