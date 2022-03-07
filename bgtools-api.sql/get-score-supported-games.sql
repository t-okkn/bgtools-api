SELECT
  `bg`.`id` AS `game_id`,
  `bg`.`title`,
  `bg`.`min_players`,
  `bg`.`max_players`,
  `col`.`color`
FROM `M_BOARDGAME` AS `bg`
INNER JOIN `M_COLOR` AS `col`
  ON `bg`.`id` = `col`.`game_id`
WHERE CAST(`bg`.`score_tool` AS UNSIGNED) = 1;