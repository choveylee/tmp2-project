CREATE TABLE IF NOT EXISTS knowledge_bases
(
    id          VARCHAR(24)  NOT NULL,
    code        VARCHAR(64)  NOT NULL,
    name        VARCHAR(128) NOT NULL,
    owner_id    VARCHAR(24)  NOT NULL,
    description TEXT         NULL,
    visible     INT          NOT NULL COMMENT '0: 私有 1: 内部 2: 公共',
    status      INT          NOT NULL COMMENT '0: 禁用 1: 启用',
    created_at  DATETIME     NOT NULL,
    updated_at  DATETIME     NOT NULL,
    deleted_at  DATETIME     NULL,
    active_code VARCHAR(64) GENERATED ALWAYS AS ( CASE WHEN deleted_at IS NULL THEN code END ) STORED COMMENT '未删除知识库的唯一编码',
    PRIMARY KEY (id),
    UNIQUE KEY uk_active_code (active_code),
    KEY idx_code (code),
    KEY idx_status (status),
    KEY idx_owner (owner_id)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_0900_ai_ci COMMENT ='知识库';

CREATE TABLE IF NOT EXISTS file_objects
(
    id               VARCHAR(24)     NOT NULL,
    bucket_name      VARCHAR(128)    NOT NULL,
    object_key       VARCHAR(512)    NOT NULL,
    file_name        VARCHAR(255)    NOT NULL,
    mime_type        VARCHAR(128)    NULL,
    file_ext         VARCHAR(32)     NULL,
    size_bytes       BIGINT UNSIGNED NOT NULL,
    sha256           CHAR(64)        NULL,
    storage_provider INT             NOT NULL COMMENT '0: seaweedfs',
    created_at       DATETIME        NOT NULL,
    updated_at       DATETIME        NOT NULL,
    PRIMARY KEY (id),
    UNIQUE KEY uk_bucket_key (bucket_name, object_key),
    KEY idx_file_sha256 (sha256)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_0900_ai_ci COMMENT ='文件对象';

CREATE TABLE IF NOT EXISTS documents
(
    id                VARCHAR(24)  NOT NULL,
    knowledge_base_id VARCHAR(24)  NULL,
    chat_session_id   VARCHAR(24)  NULL,
    scope_type        INT          NOT NULL COMMENT '0: 知识库 1: 附件',
    source_type       INT          NOT NULL COMMENT '0: 用户上传 1: api接口',
    title             VARCHAR(255) NOT NULL,
    summary           TEXT         NULL,
    tags              JSON         NULL,
    owner_id          VARCHAR(24)  NOT NULL,
    lang_code         VARCHAR(16)  NOT NULL DEFAULT 'zh-CN',
    cur_version_no    INT UNSIGNED NOT NULL DEFAULT 1,
    cur_version_id    VARCHAR(24)  NULL,
    process_status    INT          NOT NULL COMMENT '处理状态 0: 待处理 1: 已上传 2: 处理中 3: 已处理 4: 处理失败',
    status            INT          NOT NULL COMMENT '状态 0: 禁用 1: 启用',

    created_at        DATETIME     NOT NULL,
    updated_at        DATETIME     NOT NULL,
    deleted_at        DATETIME     NULL,
    PRIMARY KEY (id),
    KEY idx_doc_scope_status (scope_type, status),
    KEY idx_doc_kb (knowledge_base_id),
    KEY idx_doc_chat_session (chat_session_id),
    KEY idx_doc_owner (owner_id),
    KEY idx_doc_cur_version (cur_version_id)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_0900_ai_ci COMMENT ='文档';

CREATE TABLE IF NOT EXISTS document_versions
(
    id             VARCHAR(24)  NOT NULL,
    document_id    VARCHAR(24)  NOT NULL,
    version_no     INT UNSIGNED NOT NULL,
    file_object_id VARCHAR(24)  NOT NULL,
    parse_strategy INT          NOT NULL COMMENT '解析策略 0: 自动 1: Tika 2: OCR 3: Tika+Ocr',
    parser_type    INT          NOT NULL COMMENT '解析类型 0: 未知 1: Tika 2: OCR 3: 混合 4: 人工',
    plain_text     LONGTEXT     NULL,
    content_sha256 CHAR(64)     NULL,
    page_count     INT UNSIGNED NOT NULL,
    token_count    INT UNSIGNED NOT NULL,
    chunk_count    INT UNSIGNED NOT NULL,
    parse_status   INT          NOT NULL COMMENT '解析状态 0: 待处理 1: 处理中 2: 成功  3: 失败',
    parse_error    TEXT         NULL,
    ocr_status     INT          NOT NULL COMMENT 'OCR状态 0: 不需要 1: 待处理 2: 处理中 3: 成功 4: 失败',
    ocr_error      TEXT         NULL,
    created_at     DATETIME     NOT NULL,
    updated_at     DATETIME     NOT NULL,
    PRIMARY KEY (id),
    UNIQUE KEY uk_doc_version (document_id, version_no),
    KEY idx_doc_ver_doc (document_id),
    KEY idx_doc_ver_file (file_object_id),
    KEY idx_doc_ver_parse_status (parse_status),
    KEY idx_doc_ver_ocr_status (ocr_status)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_0900_ai_ci COMMENT ='文档版本';

CREATE TABLE IF NOT EXISTS document_chunks
(
    id                VARCHAR(24)  NOT NULL,
    document_id       VARCHAR(24)  NOT NULL,
    version_id        VARCHAR(24)  NOT NULL,
    chunk_no          INT UNSIGNED NOT NULL,
    chunk_type        INT          NOT NULL COMMENT '切片类型 0: 段落 1: 标题 2: 表格 3: 列表 4: 代码块 5: OCR识别',
    heading_path      VARCHAR(512) NULL,
    chunk_text        MEDIUMTEXT   NOT NULL,
    text_hash         CHAR(64)     NULL,
    page_from         INT UNSIGNED NULL,
    page_to           INT UNSIGNED NULL,
    char_start        INT UNSIGNED NULL,
    char_end          INT UNSIGNED NULL,
    token_count       INT UNSIGNED NOT NULL,
    index_status      INT          NOT NULL COMMENT '索引状态 0: 待索引 1: 索引中 2: 已索引 3: 失败',
    opensearch_doc_id VARCHAR(128) NULL,
    created_at        DATETIME     NOT NULL,
    updated_at        DATETIME     NOT NULL,
    deleted_at        DATETIME     NULL,
    PRIMARY KEY (id),
    UNIQUE KEY uk_chunk_version_no (version_id, chunk_no),
    KEY idx_chunk_doc (document_id),
    KEY idx_chunk_index_status (index_status),
    KEY idx_chunk_os_doc_id (opensearch_doc_id)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_0900_ai_ci COMMENT ='文档块';

CREATE TABLE IF NOT EXISTS ingest_jobs
(
    id            VARCHAR(24)  NOT NULL,
    document_id   VARCHAR(24)  NOT NULL,
    version_id    VARCHAR(24)  NOT NULL,
    job_type      INT          NOT NULL COMMENT '任务类型 0: 解析 1: OCR 2: 切分 3: 向量化 4: 索引 5: 重建索引',
    job_status    INT          NOT NULL COMMENT '任务状态 0: 待执行 1: 执行中 2: 已执行 3: 失败 4: 已取消',
    retry_count   INT UNSIGNED NOT NULL DEFAULT 0,
    worker_name   VARCHAR(128) NULL,
    started_at    DATETIME     NULL,
    finished_at   DATETIME     NULL,
    error_message TEXT         NULL,
    payload       JSON         NULL,
    created_at    DATETIME     NOT NULL,
    updated_at    DATETIME     NOT NULL,
    deleted_at    DATETIME     NULL,
    PRIMARY KEY (id),
    KEY idx_jobs_status_type (job_status, job_type),
    KEY idx_jobs_doc_version (document_id, version_id),
    KEY idx_jobs_created (created_at)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_0900_ai_ci COMMENT ='处理任务';

CREATE TABLE IF NOT EXISTS chat_sessions
(
    id            VARCHAR(24)  NOT NULL,
    kb_id         VARCHAR(24)  NULL,
    user_id       VARCHAR(24)  NOT NULL,
    session_title VARCHAR(255) NULL,
    status        INT          NOT NULL COMMENT '状态 0: 正常 1: 已归档',
    created_at    DATETIME     NOT NULL,
    updated_at    DATETIME     NOT NULL,
    deleted_at    DATETIME,
    PRIMARY KEY (id),
    KEY idx_sess_user (user_id),
    KEY idx_sess_kb (kb_id)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_0900_ai_ci COMMENT ='会话';

CREATE TABLE IF NOT EXISTS chat_messages
(
    id         VARCHAR(24)  NOT NULL,
    session_id VARCHAR(24)  NOT NULL,
    role       INT          NOT NULL COMMENT '角色 0: system, 1: user, 2: assistant, 3: tool',
    content    MEDIUMTEXT   NOT NULL,
    model_name VARCHAR(128) NULL,
    citations  JSON         NULL COMMENT '引用消息',
    latency_ms INT UNSIGNED NULL COMMENT '处理耗时',
    created_at DATETIME     NOT NULL,
    PRIMARY KEY (id),
    KEY idx_msg_created (session_id, created_at)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_0900_ai_ci COMMENT ='会话消息';

CREATE TABLE IF NOT EXISTS search_logs
(
    id           VARCHAR(24)  NOT NULL,
    kb_id        VARCHAR(24)  NULL,
    user_id      VARCHAR(24)  NULL,
    session_id   VARCHAR(24)  NULL,
    query_text   TEXT         NOT NULL,
    query_mode   INT          NOT NULL COMMENT '查询模式 0: 搜索 1: 问答',
    scope        INT          NOT NULL COMMENT '检索范围 0: 知识库 1: 附件 2: 全部',
    top_k        INT UNSIGNED NOT NULL DEFAULT 10,
    recall_count INT UNSIGNED NOT NULL,
    rerank_count INT UNSIGNED NOT NULL,
    result_count INT UNSIGNED NOT NULL,
    latency_ms   INT UNSIGNED NOT NULL,
    created_at   DATETIME     NOT NULL,
    PRIMARY KEY (id),
    KEY idx_log_created (created_at),
    KEY idx_log_user_created (user_id, created_at),
    KEY idx_log_session_created (session_id, created_at)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_0900_ai_ci COMMENT ='搜索日志';
