CREATE TABLE users(
id int unique not null CONSTRAINT users_pk PRIMARY KEY,
balance int not null);

CREATE TABLE products(
id int not null unique CONSTRAINT products_pk PRIMARY KEY,
cost int not null);

CREATE TABLE orders(
id int unique not null CONSTRAINT orders_pk PRIMARY KEY,
user_id int not null references users (id),
product_id int not null references products (id),
cost int not null,
started_at timestamp default now());

CREATE TABLE reports(
id int unique not null,
user_id int not null references users (id),
product_id int not null references products (id),
cost int not null,
started_at timestamp not null);
