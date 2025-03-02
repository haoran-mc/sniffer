CREATE TABLE default.traffic
(
    `ts` DateTime('Asia/Shanghai') COMMENT '时间(秒)',
    `host` LowCardinality(String),
    `method` LowCardinality(String),
    `url` String,
    `ip` String,
    `status` String COMMENT '响应状态码',
    `req_body` String COMMENT '请求体',
    `res_body` String COMMENT '响应体'
)
ENGINE = MergeTree
PARTITION BY toYYYYMMDD(ts)
ORDER BY (ts)
TTL ts + toIntervalYear(1)
SETTINGS index_granularity = 8192, ttl_only_drop_parts = 0
