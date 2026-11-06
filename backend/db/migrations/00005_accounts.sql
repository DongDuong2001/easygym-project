-- Accounts module: people who signed in with Google, their role and their wishlist.
-- id is the sign-in provider's user id (the token's "sub"); the provider keeps the login itself.
-- +goose Up
create schema accounts;

create table accounts.users (
	id text primary key,
	email text not null default '',
	name text not null default '',
	role text not null default 'customer' check (role in ('customer', 'staff', 'admin')),
	created_at timestamptz not null default now(),
	last_seen timestamptz not null default now()
);

-- Product slugs are copied, not referenced: modules do not join across schemas.
create table accounts.wishlist (
	user_id text not null references accounts.users (id) on delete cascade,
	product_slug text not null,
	position int not null,
	primary key (user_id, product_slug)
);

-- +goose Down
drop schema accounts cascade;
