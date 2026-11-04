-- Content module: news, albums, videos and static pages.
-- +goose Up
create schema content;

create table content.posts (
	slug text primary key,
	title text not null,
	excerpt text not null default '',
	image text not null default '',
	category text not null default 'tin-tuc',
	published_on date not null,
	body text not null default '',
	published boolean not null default true,
	created_at timestamptz not null default now(),
	updated_at timestamptz not null default now()
);

create index posts_listing on content.posts (category, published_on desc) where published;

create table content.albums (
	slug text primary key,
	title text not null,
	position int not null default 0
);

create table content.album_images (
	album_slug text not null references content.albums (slug) on delete cascade on update cascade,
	position int not null,
	url text not null,
	primary key (album_slug, position)
);

create table content.videos (
	youtube_id text primary key,
	title text not null,
	position int not null default 0
);

create table content.pages (
	slug text primary key,
	body text not null,
	updated_at timestamptz not null default now()
);

-- +goose Down
drop schema content cascade;
