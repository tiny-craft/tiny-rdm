/**
 * string format types
 * @enum {string}
 */
export const formatTypes = {
    RAW: 'Raw',
    JSON: 'JSON',
    YAML: 'YAML',
    XML: 'XML',
    HEX: 'Hex',
    BINARY: 'Binary',
}

/**
 * string decode types
 * @enum {string}
 */
export const decodeTypes = {
    NONE: 'None',
    BASE64: 'Base64',
    GZIP: 'GZip',
    DEFLATE: 'Deflate',
    ZSTD: 'ZStd',
    BROTLI: 'Brotli',
    MSGPACK: 'Msgpack',
    PHP: 'PHP',
    PICKLE: 'Pickle',
    // Java: 'Java',
}
