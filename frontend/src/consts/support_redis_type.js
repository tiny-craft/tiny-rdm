export const types = {
    STRING: 'STRING',
    HASH: 'HASH',
    LIST: 'LIST',
    SET: 'SET',
    ZSET: 'ZSET',
}

export const typesColor = {
    [types.STRING]: '#5A96E3',
    [types.HASH]: '#9575DE',
    [types.LIST]: '#7A9D54',
    [types.SET]: '#F3AA60',
    [types.ZSET]: '#FF6666',
}

// export const typesName = Object.fromEntries(Object.entries(types).map(([key, value]) => [key, value.name]))

export const validType = (t) => {
    return types.hasOwnProperty(t)
}
