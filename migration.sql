CREATE TABLE stores
(
    id      serial PRIMARY KEY,
    title   varchar(80) not null,
    website varchar(80) not null
);

CREATE TABLE items
(
    id    serial PRIMARY KEY,
    title text not null
);

CREATE TABLE link_items_stores
(
    id        serial primary key,
    item_id   int not null,
    FOREIGN KEY (item_id) REFERENCES items (id),
    store_id  int not null,
    FOREIGN KEY (store_id) REFERENCES stores (id),
    url       text
);

CREATE TABLE prices
(
    id            serial primary key,
    created       date,
    item_store_id int,
    FOREIGN KEY (item_store_id) REFERENCES link_items_stores (id),
    price         float4
);
CREATE INDEX item_id_idx ON link_items_stores (item_id);
CREATE INDEX store_id_idx ON link_items_stores (item_id);
CREATE INDEX item_store_idx ON prices (item_store_id);

INSERT INTO items (title)
VALUES ('Apple MacBook Pro 14" MKGP3 (M1 Pro 8C CPU, 14C GPU, 2021) 16 ГБ, 512 ГБ SSD, серый космос');
INSERT INTO stores(title, website)
VALUES ('pitergsm', 'https://pitergsm.ru/');
INSERT INTO link_items_stores (item_id, store_id, url)
VALUES (1, 1, 'https://pitergsm.ru/catalog/tablets-and-laptops/mac/macbook-pro/macbook-pro-14-2021/13263/');
INSERT INTO prices(created, item_store_id, price) VALUES ('2022-11-21', 1, 130490);

INSERT INTO items (title)
VALUES ('Apple MacBook Pro 14" MKGR3 (M1 Pro 8C CPU, 14C GPU, 2021) 16 ГБ, 512 ГБ SSD, серебристый');
INSERT INTO link_items_stores (item_id, store_id, url)
VALUES (2, 1, 'https://pitergsm.ru/catalog/tablets-and-laptops/mac/macbook-pro/macbook-pro-14-2021/13262/');
INSERT INTO prices(created, item_store_id, price) VALUES ('2022-11-21', 2, 130490);


select *
from prices;


drop table prices;
drop table link_items_stores;
drop table items;
drop table stores;
