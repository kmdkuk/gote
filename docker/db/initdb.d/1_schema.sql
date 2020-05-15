USE `network-monitoring`;
---- drop ----
DROP TABLE IF EXISTS disconnect_dates;
---- create ----
CREATE TABLE IF NOT EXISTS disconnect_dates (
  id INT NOT NULL AUTO_INCREMENT,
  start_at datetime DEFAULT 0,
  end_at datetime DEFAULT 0,
  PRIMARY KEY (id)
) DEFAULT CHARSET = utf8 COLLATE = utf8_bin;
