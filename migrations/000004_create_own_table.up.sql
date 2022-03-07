CREATE TABLE IF NOT EXISTS `T_OWN` (
  `user_id` VARCHAR(8) NOT NULL DEFAULT '',
  `game_id` VARCHAR(8) NOT NULL DEFAULT '',
  UNIQUE id_pair (`user_id`, `game_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;