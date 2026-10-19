-- +goose Up
-- +goose StatementBegin
CREATE TABLE accommodations (
    id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    provider_code VARCHAR(100) NOT NULL,
    external_id VARCHAR(255) NOT NULL,
    external_parent_id VARCHAR(255) NULL,
    name VARCHAR(255) NOT NULL,
    description TEXT NOT NULL,
    property_type VARCHAR(100) NOT NULL,
    chain VARCHAR(255) NOT NULL,
    currency CHAR(3) NOT NULL,
    country CHAR(2) NOT NULL,
    city VARCHAR(255) NOT NULL,
    latitude DECIMAL(10,7) NULL,
    longitude DECIMAL(10,7) NULL,
    address VARCHAR(500) NOT NULL,
    zip VARCHAR(50) NOT NULL,
    main_photo VARCHAR(500) NULL,
    thumbnail VARCHAR(500) NULL,
    stars DECIMAL(2,1) NOT NULL,
    rating DECIMAL(3,1) NULL,
    review_count INT UNSIGNED NOT NULL DEFAULT 0,
    metadata JSON NULL,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    created_by BIGINT UNSIGNED NOT NULL,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    updated_by BIGINT UNSIGNED NOT NULL,
    deleted_at DATETIME NULL,
    deleted_by BIGINT UNSIGNED NULL,
    PRIMARY KEY (id),
    UNIQUE KEY uq_accommodations_provider_external (external_id, provider_code),
    KEY idx_accommodations_type_location_deleted (property_type, country, city, deleted_at),
    KEY idx_accommodations_location_deleted (country, city, deleted_at)
) AUTO_INCREMENT = 10001;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE accommodations;
-- +goose StatementEnd
