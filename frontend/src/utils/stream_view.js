const isPlainObject = (value) => {
    return Object.prototype.toString.call(value) === '[object Object]'
}

export const formatStreamCellValue = (value) => {
    if (value == null) {
        return ''
    }
    if (Array.isArray(value) || isPlainObject(value)) {
        return JSON.stringify(value, null, 2)
    }
    if (typeof value !== 'string') {
        return String(value)
    }

    const trimmed = value.trim()
    const looksLikeStructuredJson =
        (trimmed.startsWith('{') && trimmed.endsWith('}')) ||
        (trimmed.startsWith('[') && trimmed.endsWith(']'))

    if (!looksLikeStructuredJson) {
        return value
    }

    try {
        const parsed = JSON.parse(trimmed)
        return Array.isArray(parsed) || isPlainObject(parsed) ? JSON.stringify(parsed, null, 2) : value
    } catch {
        return value
    }
}

export const collectStreamColumns = (rows = []) => {
    const keys = ['id']
    const seen = new Set(keys)

    for (const row of rows) {
        for (const fieldName of Object.keys(row?.v || {})) {
            if (!seen.has(fieldName)) {
                seen.add(fieldName)
                keys.push(fieldName)
            }
        }
    }

    return keys.map((key) => ({
        key,
        title: key === 'id' ? 'ID' : key,
    }))
}

export const buildStreamFilterText = (entry = {}) => {
    return Object.entries(entry)
        .map(([key, value]) => `${key}: ${formatStreamCellValue(value)}`)
        .join('\n')
}

export const filterStreamRows = (rows = [], filter = '') => {
    if (!filter) {
        return rows
    }

    return rows.filter((row) => buildStreamFilterText(row?.v || {}).includes(filter))
}

export const getStreamFieldValue = (entry = {}, fieldName = '') => {
    return fieldName === 'id' ? entry?.id || '' : entry?.v?.[fieldName] ?? ''
}
