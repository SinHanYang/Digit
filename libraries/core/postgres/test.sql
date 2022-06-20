/***** create and use database *****/

/***** info *****/

/***** table creation and insertion *****/

DROP TABLE IF EXISTS Member;
DROP TABLE IF EXISTS Member__backend;

CREATE TABLE Member (
    GroupName varchar(20) NOT NULL,
    StageName varchar(20) NOT NULL,
    RealName varchar(10) NOT NULL,
    Posistion varchar(10),
    Nationality varchar(50),
    Age int
);

INSERT INTO 
    Member
VALUES 
    ('ONF', 'Hyojin', '金曉珍', 'Vocal', 'Korean', '27'),
    ('ONF', 'J-US', '李昇俊', 'Dancer', 'Korean', '27'),
    ('ONF', 'E-TION', '李昌潤', 'Vocal', 'Korean', '27'),
    ('ONF', 'Wyatt', '沈宰營', 'Rapper', 'Korean', '27'),
    ('ONF', 'MK', '朴珉均', 'Vocal', 'Korean', '26'),
    ('ONF', 'U', '水口裕斗', 'Dancer', 'Japan', '22'),
    ('BTS', 'Jin', '金碩珍', 'Vocal', 'Korean', '29'),
    ('BTS', 'SUGA', '閔玧其', 'Rapper', 'Korean', '29'),
    ('BTS', 'j-hope', '鄭號錫', 'Rapper', 'Korean', '28'),
    ('BTS', 'RM', '金南俊', 'Rapper', 'Korean', '27'),
    ('BTS', 'Jimin', '朴智旻', 'Vocal', 'Korean', '26'),
    ('BTS', 'V', '金泰亨', 'Vocal', 'Korean', '26'),
    ('BTS', 'Jung Kook', '田柾國', 'Vocal', 'Korean', '24'),
    ('BTOB', '恩光', '徐恩光', 'Vocal', 'Korean', '31'),
    ('BTOB', '旼赫', '李旼赫', 'Rapper', 'Korean', '31'),
    ('BTOB', '昌燮', '李昌燮', 'Vocal', 'Korean', '31'),
    ('BTOB', '炫植', '任炫植', 'Vocal', 'Korean', '30'),
    ('BTOB', 'Peniel', '辛東根', 'Rapper', 'America', '29'),
    ('BTOB', '星材', '陸星材', 'Vocal', 'Korean', '27'),
    ('BTOB', '鎰勳', '鄭鎰勳', 'Rapper', 'Korean', '27'),
    ('Twice', 'Tzuyu', '周子瑜', 'Dancer', 'Taiwan', '22'),
    ('Twice', 'Momo', '平井桃', 'Dancer', 'Japan', '25');

INSERT INTO Member VALUES ('Twice', 'Jeongyeon', '俞定延', 'Vocal', 'Taiwan', '25');
INSERT INTO Member VALUES ('Twice', 'Sana', '湊崎紗夏', 'Vocal', 'Japan', '25');
INSERT INTO Member VALUES ('Twice', 'Nayeon', '林娜璉', 'Dancer', 'Korean', '26');