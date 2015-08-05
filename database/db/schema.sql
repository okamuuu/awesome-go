DROP TABLE IF EXISTS user;
CREATE TABLE user (
    id            bigint(20) PRIMARY KEY AUTO_INCREMENT,
    name          varchar(255) COMMENT '名前'
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='記事' AUTO_INCREMENT=1;
