SET FOREIGN_KEY_CHECKS = 0;
SET NAMES utf8mb4;

CREATE TABLE IF NOT EXISTS `user`
(
    `id`       int unsigned NOT NULL AUTO_INCREMENT PRIMARY KEY,
    `username` varchar(255) NOT NULL,
    `password` binary(40)   NOT NULL,
    UNIQUE KEY `u_username` (`username`)
) ENGINE = InnoDB;

CREATE TABLE  IF NOT EXISTS `user_data`
(
    `id`      int unsigned NOT NULL AUTO_INCREMENT PRIMARY KEY,
    `user_id` int unsigned NOT NULL,
    `type`    int unsigned NOT NULL,
    `data`    blob         NOT NULL,
    `meta`    text         NOT NULL,
    FOREIGN KEY (`user_id`) REFERENCES `user` (`id`)
) ENGINE = InnoDB;