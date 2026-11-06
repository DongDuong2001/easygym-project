-- Orders placed while signed in remember who placed them, so customers can see their order history.
-- +goose Up
alter table orders.orders add column user_id text;
create index orders_by_user on orders.orders (user_id, created_at desc) where user_id is not null;

-- +goose Down
alter table orders.orders drop column user_id;
