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
CREATE INDEX item_id_idx ON link_items_stores (item_id);
CREATE INDEX store_id_idx ON link_items_stores (item_id);

CREATE TABLE prices
(
    id            serial primary key,
    created       date,
    item_store_id int,
    FOREIGN KEY (item_store_id) REFERENCES link_items_stores (id),
    price         float4
);
CREATE INDEX item_store_idx ON prices (item_store_id);