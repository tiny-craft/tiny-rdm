export const types = {
    STRING: 'STRING',
    HASH: 'HASH',
    LIST: 'LIST',
    SET: 'SET',
    ZSET: 'ZSET',
}

// export const typesName = Object.fromEntries(Object.entries(types).map(([key, value]) => [key, value.name]))

export const validType = (t) => {
    return types.hasOwnProperty(t)
}
