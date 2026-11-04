-- Leads module: course sign-ups, contact requests and newsletter emails.
-- +goose Up
create schema leads;

create table leads.leads (
	id bigint generated always as identity primary key,
	code text not null unique,
	kind text not null check (kind in ('register', 'newsletter')),
	name text not null default '',
	phone text not null default '',
	email text not null default '',
	course text not null default '',
	plan text not null default '',
	note text not null default '',
	source text not null default '',
	status text not null default 'new' check (status in ('new', 'contacted', 'done', 'spam')),
	created_at timestamptz not null default now()
);

create index leads_inbox on leads.leads (status, created_at desc);

-- +goose Down
drop schema leads cascade;
