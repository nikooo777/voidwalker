-- +migrate Up

-- +migrate StatementBegin
CREATE SCHEMA IF NOT EXISTS `voidwalker` DEFAULT CHARACTER SET latin1 COLLATE latin1_general_ci;
-- +migrate StatementEnd
-- +migrate StatementBegin
CREATE TABLE `thumbnail`
(
    `id`                   bigint(20) unsigned                NOT NULL AUTO_INCREMENT,
    `name`                 varchar(150) UNIQUE                NOT NULL,
    `compressed_name`      varchar(150) UNIQUE                NULL,
    `compressed_mime_type` varchar(150)                       NULL,
    `compressed`           tinyint(1)                         NOT NULL,
    `created_at`           datetime default CURRENT_TIMESTAMP not null,
    `updated_at`           datetime default CURRENT_TIMESTAMP not null on update CURRENT_TIMESTAMP,
    index (`name`),
    index (`created_at`),
    index (`updated_at`),
    PRIMARY KEY (`id`)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_unicode_ci;
-- +migrate StatementEnd
