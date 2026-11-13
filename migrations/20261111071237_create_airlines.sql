-- +goose Up
-- +goose StatementBegin
CREATE TABLE airlines (
    id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    provider_code VARCHAR(100) NOT NULL,
    iata_code VARCHAR(100) NOT NULL,
    name VARCHAR(255) NOT NULL,
    alliance VARCHAR(100) NULL COMMENT 'star_alliance, oneworld, skyteam, vanilla_alliance',
    active BOOLEAN NOT NULL DEFAULT TRUE,
    logo VARCHAR(500) NULL,
    metadata JSON NULL,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    created_by BIGINT UNSIGNED NOT NULL,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    updated_by BIGINT UNSIGNED NOT NULL,
    deleted_at DATETIME NULL,
    deleted_by BIGINT UNSIGNED NULL,
    PRIMARY KEY (id),
    UNIQUE KEY uq_airlines_iata_provider (iata_code, provider_code),
    KEY idx_airlines_created_deleted_id (created_at, deleted_at, id)
) AUTO_INCREMENT = 10001;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE airlines;
-- +goose StatementEnd
