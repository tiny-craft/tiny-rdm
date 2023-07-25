export const types = {
    STRING: 'STRING',
    HASH: 'HASH',
    LIST: 'LIST',
    SET: 'SET',
    ZSET: 'ZSET',
    STREAM: 'STREAM',
}

export const typesColor = {
    [types.STRING]: '#8256DC',
    [types.HASH]: '#0171F5',
    [types.LIST]: '#23C338',
    [types.SET]: '#F29E33',
    [types.ZSET]: '#F53227',
    [types.STREAM]: '#F5C201',
}

export const typesBgColor = {
    [types.STRING]: '#F2EDFB',
    [types.HASH]: '#E4F0FC',
    [types.LIST]: '#E3F3EB',
    [types.SET]: '#FDF1DF',
    [types.ZSET]: '#FAEAED',
    [types.STREAM]: '#FFF8DF',
}

// export const typesName = Object.fromEntries(Object.entries(types).map(([key, value]) => [key, value.name]))

export const validType = (t) => {
    return types.hasOwnProperty(t)
}
