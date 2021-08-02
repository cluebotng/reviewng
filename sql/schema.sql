# User tables
DROP TABLE IF EXISTS `users`;
CREATE TABLE `users`
(
    `id`           int          NOT NULL AUTO_INCREMENT,
    `username`     varchar(512) NOT NULL,
    `approved`     tinyint(1)   NOT NULL DEFAULT 0,
    `admin`        tinyint(1)   NOT NULL DEFAULT 0,
    `legacy_count` int          NOT NULL DEFAULT 0,
    PRIMARY KEY (`id`),
    UNIQUE KEY `username` (`username`)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_0900_ai_ci;
