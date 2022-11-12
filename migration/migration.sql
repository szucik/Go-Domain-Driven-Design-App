CREATE TABLE `tradehelper`.`users`
(
    `id`        INT                                NOT NULL AUTO_INCREMENT,
    `email`     VARCHAR(45)                        NOT NULL,
    `username`  VARCHAR(45) NULL,
    `password`  VARCHAR(255)                       NOT NULL,
    `token_hash` VARCHAR(15)                        NOT NULL,
    `created`   DATETIME DEFAULT CURRENT_TIMESTAMP NOT NULL,
    `updated`   DATETIME ON UPDATE CURRENT_TIMESTAMP,

    PRIMARY KEY (`id`),
    UNIQUE INDEX `id_UNIQUE` (`id` ASC) VISIBLE,
    UNIQUE INDEX `email_UNIQUE` (`email` ASC) VISIBLE,
    UNIQUE INDEX `username_UNIQUE` (`username` ASC) VISIBLE
);



CREATE TABLE `tradehelper`.`portfolios`
(
    `id`                INT         NOT NULL AUTO_INCREMENT,
    `user_id`           INT         NOT NULL DEFAULT 0,
    `name`              VARCHAR(45) NOT NULL,
    `total_balance`     DECIMAL(27, 18),
    `total_cost`        DECIMAL(18, 2),
    `total_profit_loss` DECIMAL(18, 2),
    `profit_loss_day`   DECIMAL(18, 2),
    `created`           DATETIME             DEFAULT CURRENT_TIMESTAMP NOT NULL,
    PRIMARY KEY (`id`),
    UNIQUE INDEX `id_UNIQUE` (`id` ASC) VISIBLE,
    UNIQUE INDEX `name_UNIQUE` (`name` ASC) VISIBLE,
    CONSTRAINT `fk_portfolios_1`
        FOREIGN KEY (`user_id`)
            REFERENCES `tradehelper`.`users` (`id`)
            ON DELETE CASCADE
            ON UPDATE NO ACTION
);

CREATE TABLE `tradehelper`.`transactions`
(
    `id`           INT             NOT NULL AUTO_INCREMENT,
    `user_id`      INT             NOT NULL DEFAULT 0,
    `symbol`       VARCHAR(45)     NOT NULL DEFAULT '',
    `portfolio_id` INT             NOT NULL DEFAULT 0,
    `quantity`     DECIMAL(27, 18) NOT NULL DEFAULT 0,
    `price`        DECIMAL(18, 2)  NOT NULL DEFAULT 0,
    `created`      DATETIME                 DEFAULT CURRENT_TIMESTAMP NOT NULL,
    PRIMARY KEY (`id`),
    UNIQUE INDEX `id_UNIQUE` (`id` ASC) VISIBLE,
    UNIQUE INDEX `portfolio_id_UNIQUE` (`portfolio_id` ASC) VISIBLE,
    INDEX          `fk_transactions_1_idx` (`user_id` ASC) VISIBLE,
    CONSTRAINT `fk_transactions_id`
        FOREIGN KEY (`user_id`)
            REFERENCES `tradehelper`.`users` (`id`)
            ON DELETE CASCADE
            ON UPDATE NO ACTION,
    CONSTRAINT `fk_portfolios_2`
        FOREIGN KEY (`portfolio_id`)
            REFERENCES `tradehelper`.`portfolios` (`id`)
            ON DELETE CASCADE
            ON UPDATE NO ACTION
);

-- CREATE TABLE `tradehelper`.`cryptocurrencys`
-- (
--     `id`           INT                                NOT NULL AUTO_INCREMENT,
--     `project_name` VARCHAR(45)                        NOT NULL,
--     `slug`         VARCHAR(45)                        NOT NULL,
--     `symbol`       VARCHAR(45)                        NOT NULL,
--     `link`         VARCHAR(255) NULL,
--     `description`  VARCHAR(255),
--     `created`      DATETIME DEFAULT CURRENT_TIMESTAMP NOT NULL,
--     PRIMARY KEY (`id`),
--     UNIQUE INDEX `id_UNIQUE` (`id` ASC) VISIBLE,
--     UNIQUE INDEX `slug_UNIQUE` (`slug` ASC) VISIBLE
-- );