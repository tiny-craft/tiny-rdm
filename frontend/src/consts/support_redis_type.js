export const types = {
    STRING: 'STRING',
    HASH: 'HASH',
    LIST: 'LIST',
    SET: 'SET',
    ZSET: 'ZSET',
}

export const typesColor = {
    [types.STRING]: '#8256DC',
    [types.HASH]: '#2983ED',
    [types.LIST]: '#26A15E',
    [types.SET]: '#EE9F33',
    [types.ZSET]: '#CE3352',
}

export const typesBgColor = {
    [types.STRING]: '#F2EDFB',
    [types.HASH]: '#E4F0FC',
    [types.LIST]: '#E3F3EB',
    [types.SET]: '#FDF1DF',
    [types.ZSET]: '#FAEAED',
}

// export const typesName = Object.fromEntries(Object.entries(types).map(([key, value]) => [key, value.name]))

export const validType = (t) => {
    return types.hasOwnProperty(t)
}
