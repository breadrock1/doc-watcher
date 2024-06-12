CREATE TABLE IF NOT EXISTS metabase.documents
(
    folder_id            TEXT COLLATE pg_catalog."default" NOT NULL,
    folder_path          TEXT NOT NULL,
    content              TEXT COLLATE pg_catalog."default",
    document_id          TEXT COLLATE pg_catalog."default" NOT NULL UNIQUE,
    document_ssdeep      TEXT COLLATE pg_catalog."default",
    document_name        TEXT COLLATE pg_catalog."default" NOT NULL,
    document_path        TEXT,
    document_size        BIGINT NOT NULL,
    document_type        TEXT COLLATE pg_catalog."default",
    document_ext         TEXT COLLATE pg_catalog."default",
    document_perm        BIGINT,
    document_created     DATE,
    document_modified    DATE,
    class                TEXT COLLATE pg_catalog."default"
)

    TABLESPACE pg_default;

ALTER TABLE IF EXISTS metabase.documents
    OWNER to metabase;
