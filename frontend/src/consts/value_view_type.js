/**
 * string format types
 * @enum {string}
 */
export const formatTypes = {
    RAW: 'Raw',
    JSON: 'JSON',
    UNICODE_JSON: 'Unicode JSON',
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
    LZ4: 'LZ4',
    BROTLI: 'Brotli',
    MSGPACK: 'Msgpack',
    PHP: 'PHP',
    PICKLE: 'Pickle',
    // Java: 'Java',
}
