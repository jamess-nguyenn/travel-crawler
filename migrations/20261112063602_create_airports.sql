-- +goose Up
-- +goose StatementBegin
CREATE TABLE airports (
    id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    iata_code VARCHAR(100) NOT NULL,
    city VARCHAR(255) NOT NULL,
    country VARCHAR(255) NOT NULL,
    metadata JSON NULL,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    created_by BIGINT UNSIGNED NOT NULL,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    updated_by BIGINT UNSIGNED NOT NULL,
    deleted_at DATETIME NULL,
    deleted_by BIGINT UNSIGNED NULL,
    PRIMARY KEY (id),
    UNIQUE KEY uq_airports_iata (iata_code),
    KEY idx_airports_created_deleted_id (created_at, deleted_at, id)
) AUTO_INCREMENT = 10001;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE airports;
-- +goose StatementEnd
