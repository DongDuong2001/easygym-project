-- Catalog module: products, courses and their prices. Only the catalog package reads or writes this schema.
-- +goose Up
create schema catalog;

create table catalog.categories (
	slug text primary key,
	label text not null,
	position int not null default 0
);

create table catalog.products (
	slug text primary key,
	name text not null,
	category text not null references catalog.categories (slug),
	tag text not null default '',
	image text not null,
	images text[] not null default '{}',
	reviews int not null default 0 check (reviews >= 0),
	description text not null default '',
	detail text not null default '',
	position int not null default 0,
	active boolean not null default true,
	created_at timestamptz not null default now(),
	updated_at timestamptz not null default now()
);

create table catalog.product_variants (
	product_slug text not null references catalog.products (slug) on delete cascade on update cascade,
	label text not null,
	price bigint not null check (price >= 0),
	position int not null default 0,
	primary key (product_slug, label)
);

create table catalog.product_questions (
	product_slug text not null references catalog.products (slug) on delete cascade on update cascade,
	position int not null,
	question text not null,
	answer text not null,
	primary key (product_slug, position)
);

create table catalog.courses (
	slug text primary key,
	name text not null,
	mode text not null default '',
	card text not null,
	image text not null,
	summary text not null default '',
	body text not null default '',
	position int not null default 0,
	active boolean not null default true,
	created_at timestamptz not null default now(),
	updated_at timestamptz not null default now()
);

create table catalog.course_plans (
	course_slug text not null references catalog.courses (slug) on delete cascade on update cascade,
	label text not null,
	price bigint not null check (price >= 0),
	position int not null default 0,
	primary key (course_slug, label)
);

-- +goose Down
drop schema catalog cascade;
