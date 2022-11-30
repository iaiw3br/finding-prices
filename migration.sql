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

CREATE TABLE item_in_store
(
    id        serial primary key,
    item_id   int not null,
    FOREIGN KEY (item_id) REFERENCES items (id),
    store_id  int not null,
    FOREIGN KEY (store_id) REFERENCES stores (id),
    url       text
);
CREATE INDEX item_id_idx ON item_in_store (item_id);
CREATE INDEX store_id_idx ON item_in_store (item_id);

CREATE TABLE prices
(
    id            serial primary key,
    created       timestamp without time zone,
    item_store_id int,
    FOREIGN KEY (item_store_id) REFERENCES item_in_store (id),
    price         float4
);
CREATE INDEX item_store_idx ON prices (item_store_id);


-- test data
INSERT INTO public.stores (title, website) VALUES ('pitergsm', 'https://pitergsm.ru/');
INSERT INTO public.items (title) VALUES ('Apple MacBook Pro 14" MKGP3 (M1 Pro 8C CPU, 14C GPU, 2021) 16 ГБ, 512 ГБ SSD, серый космос');
INSERT INTO public.item_in_store (item_id, store_id, url) VALUES (1, 1, 'https://pitergsm.ru/catalog/tablets-and-laptops/mac/macbook-pro/macbook-pro-14-2021/13263/');