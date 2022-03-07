DROP TABLE IF EXISTS `users`;
CREATE TABLE `users`
(
    `id`           int          NOT NULL AUTO_INCREMENT,
    `username`     varchar(100) NOT NULL,
    `approved`     tinyint(1) NOT NULL DEFAULT 0,
    `admin`        tinyint(1) NOT NULL DEFAULT 0,
    `legacy_count` int          NOT NULL DEFAULT 0,
    PRIMARY KEY (`id`),
    UNIQUE KEY `username` (`username`)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_bin;

DROP TABLE IF EXISTS `user_classification`;
CREATE TABLE `user_classification`
(
    `id`             int NOT NULL AUTO_INCREMENT,
    `user_id`        int NOT NULL,
    `comment`        varchar(1024) NULL,
    `classification` int NOT NULL,
    `edit_id`        int NOT NULL,
    PRIMARY KEY (`id`),
    INDEX            `user_id` (`user_id`),
    INDEX            `edit_id` (`edit_id`),
    INDEX            `classification` (`classification`),
    UNIQUE KEY `user_edit` (`user_id`, `edit_id`)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_bin;

DROP TABLE IF EXISTS `edit_group`;
CREATE TABLE `edit_group`
(
    `id`     int          NOT NULL AUTO_INCREMENT,
    `name`   varchar(255) NOT NULL,
    `weight` int          NOT NULL,
    PRIMARY KEY (`id`)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_bin;

DROP TABLE IF EXISTS `edit`;
CREATE TABLE `edit`
(
    `id`             int NOT NULL,
    `required`       int NOT NULL,
    `classification` int NOT NULL,
    PRIMARY KEY (`id`)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_bin;

DROP TABLE IF EXISTS `edit_edit_group`;
CREATE TABLE `edit_edit_group`
(
    `edit_id`             int NOT NULL,
    `edit_group_id`  int NOT NULL,
    PRIMARY KEY (`edit_id`, `edit_group_id`),
    INDEX            `edit_id` (`edit_id`),
    INDEX            `edit_group_id` (`edit_group_id`)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_bin;

DROP TABLE IF EXISTS `edit_training_data`;
CREATE TABLE `edit_training_data`
(
    `edit_id`             int NOT NULL,
    `training_data`       longblob NOT NULL,
    PRIMARY KEY (`edit_id`)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_bin;
