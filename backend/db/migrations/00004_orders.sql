-- Orders module. Lines keep a copy of name and price at order time, so they do not reference catalog tables.
-- +goose Up
create schema orders;

create table orders.orders (
	id bigint generated always as identity primary key,
	code text not null unique,
	name text not null,
	phone text not null,
	address text not null,
	payment text not null check (payment in ('cod', 'bank')),
	note text not null default '',
	total bigint not null check (total >= 0),
	status text not null default 'new' check (status in ('new', 'confirmed', 'shipping', 'done', 'cancelled')),
	created_at timestamptz not null default now()
);

create index orders_inbox on orders.orders (status, created_at desc);

create table orders.order_lines (
	order_id bigint not null references orders.orders (id) on delete cascade,
	position int not null,
	product_slug text not null,
	product_name text not null,
	variant_label text not null,
	unit_price bigint not null check (unit_price >= 0),
	qty int not null check (qty between 1 and 99),
	primary key (order_id, position)
);

-- +goose Down
drop schema orders cascade;
