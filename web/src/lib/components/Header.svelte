<script>
	import { afterNavigate, goto } from '$app/navigation';
	import Icon from './Icon.svelte';
	import { catalog } from '#lib/catalog.js';
	import { formatVND, fromPrice, suggest } from '#lib/shop.js';
	import { brand, logo, nav } from '#lib/site.js';
	import { store } from '#lib/store.svelte.js';

	const products = $derived(catalog().products);

	let menu = $state(false);
	let searching = $state(false);
	let focused = $state(false);
	let q = $state('');

	const hits = $derived(focused ? suggest(products, q) : []);
	const count = $derived(store.cart.reduce((s, l) => s + l.qty, 0));

	afterNavigate(() => {
		menu = false;
		searching = false;
		focused = false;
	});

	function submit(e) {
		e.preventDefault();
		const term = q.trim();
		goto(term ? `/cua-hang?q=${encodeURIComponent(term)}` : '/cua-hang');
	}

	function closeAll(e) {
		if (e.key !== 'Escape') return;
		menu = false;
		searching = false;
		focused = false;
	}
</script>

<svelte:window onkeydown={closeAll} />

<a class="skip" href="#noi-dung">Tới nội dung chính</a>

<aside class="top" aria-label="Thông tin liên hệ">
	<div class="container row">
		<p><a href={`tel:${brand.phone}`}>{brand.phone}</a><a href={`mailto:${brand.email}`}>{brand.email}</a></p>
		<p class="mid">Đăng ký trải nghiệm: <a href="/lien-he#dang-ky"><b>Để lại số điện thoại</b></a></p>
		<p class="end"><a href={brand.map} target="_blank" rel="noreferrer">Gò Vấp, TP.HCM</a><span>8:00 - 20:00</span></p>
	</div>
</aside>

<header class="bar">
	<div class="container row">
		<button class="icon burger" onclick={() => (menu = !menu)} aria-expanded={menu} aria-controls="menu-chinh" aria-label={menu ? 'Đóng menu' : 'Mở menu'}>
			<Icon name={menu ? 'x' : 'menu'} size={24} />
		</button>

		<a class="logo" href="/" aria-label="easyGYM, về trang chủ">
			<img src={logo} alt="" width="48" height="48" />
			<span>easyGYM</span>
		</a>

		<nav id="menu-chinh" class:open={menu} aria-label="Menu chính">
			{#each nav as n}
				<div class="item">
					<a href={n.href}>
						{n.label}
						{#if n.tag}<small>{n.tag}</small>{/if}
					</a>
					{#if n.children}
						<ul class="sub">
							{#each n.children as c}<li><a href={c.href}>{c.label}</a></li>{/each}
						</ul>
					{/if}
				</div>
			{/each}
		</nav>

		<div class="tools">
			<form
				class="search"
				class:open={searching}
				role="search"
				onsubmit={submit}
				onfocusin={() => (focused = true)}
				onfocusout={(e) => (focused = e.currentTarget.contains(e.relatedTarget))}
			>
				<button class="go" aria-label="Tìm kiếm"><Icon name="search" size={18} /></button>
				<label class="sr-only" for="tim-kiem">Tìm sản phẩm</label>
				<input id="tim-kiem" type="search" placeholder="Bạn đang tìm gì?" autocomplete="off" bind:value={q} />
				{#if hits.length}
					<ul class="hits">
						{#each hits as p}
							<li>
								<a href={`/san-pham/${p.slug}`}>
									<img src={p.image} alt="" width="44" height="44" />
									<span>{p.name}<small>Từ {formatVND(fromPrice(p.variants))}</small></span>
								</a>
							</li>
						{/each}
					</ul>
				{/if}
			</form>

			<button class="icon mobile-only" onclick={() => (searching = !searching)} aria-expanded={searching} aria-label="Mở ô tìm kiếm">
				<Icon name="search" size={22} />
			</button>
			<a class="icon desktop-only" href="/yeu-thich" aria-label={`Yêu thích, ${store.wishes.length} sản phẩm`}>
				<Icon name="heart" size={22} />
				{#if store.wishes.length}<span class="count">{store.wishes.length}</span>{/if}
			</a>
			<a class="icon" href="/gio-hang" aria-label={`Giỏ hàng, ${count} sản phẩm`}>
				<Icon name="cart" size={22} />
				{#if count}<span class="count">{count}</span>{/if}
			</a>
		</div>
	</div>
</header>

<style>
	.skip {
		position: absolute;
		left: 8px;
		top: -60px;
		z-index: 60;
		padding: 12px 16px;
		background: var(--ink);
		color: var(--on-dark);
	}

	.skip:focus {
		top: 8px;
	}

	.top {
		background: var(--ink);
		color: var(--on-dark);
		font-size: var(--t-xs);
	}

	.top .row {
		display: grid;
		grid-template-columns: 1fr auto 1fr;
		align-items: center;
		min-height: 36px;
	}

	.top p {
		display: flex;
		gap: 20px;
	}

	.top .mid {
		gap: 4px;
	}

	.top .end {
		justify-content: flex-end;
	}

	.top a:hover {
		text-decoration: underline;
	}

	.bar {
		position: sticky;
		top: 0;
		z-index: 30;
		background: var(--surface);
		border-bottom: 1px solid var(--line);
	}

	.bar .row {
		display: flex;
		align-items: center;
		gap: 24px;
		height: var(--bar);
	}

	.logo {
		display: flex;
		align-items: center;
		gap: 10px;
		font-size: var(--t-lg);
		font-weight: 700;
	}

	.logo img {
		border-radius: 50%;
	}

	nav {
		display: flex;
		flex: 1;
		justify-content: center;
		gap: 32px;
	}

	nav a {
		position: relative;
		display: flex;
		align-items: center;
		min-height: 44px;
		font-size: var(--t-sm);
		font-weight: 600;
		text-transform: uppercase;
	}

	nav a:hover {
		color: var(--brand);
	}

	.item {
		position: relative;
	}

	.sub {
		position: absolute;
		top: 100%;
		left: -16px;
		z-index: 2;
		display: none;
		width: 300px;
		padding-block: 8px;
		background: var(--surface);
		border: 1px solid var(--line);
		box-shadow: 0 12px 28px rgb(13 13 63 / 0.12);
	}

	.item:hover .sub,
	.item:focus-within .sub {
		display: block;
	}

	nav .sub a {
		padding: 6px 16px;
		font-weight: 500;
		line-height: 1.35;
		text-transform: none;
	}

	nav .sub a:hover {
		background: var(--soft);
	}

	nav small {
		position: absolute;
		left: 100%;
		top: -8px;
		margin-left: -22px;
		padding: 1px 5px;
		background: var(--accent);
		color: var(--on-dark);
		font-size: 10px;
		font-weight: 600;
		line-height: 1.4;
		white-space: nowrap;
		text-transform: none;
	}

	.tools {
		display: flex;
		align-items: center;
		gap: 4px;
	}

	.search {
		position: relative;
		display: flex;
		align-items: center;
		width: 240px;
		margin-right: 8px;
		border: 1px solid var(--line);
		border-radius: 999px;
	}

	.search:focus-within {
		border-color: var(--brand);
	}

	.go {
		display: grid;
		place-items: center;
		flex: none;
		width: 44px;
		height: 44px;
		color: var(--muted);
	}

	.search input {
		flex: 1;
		min-width: 0;
		height: 44px;
		padding: 0 16px 0 0;
		background: none;
		border: 0;
		outline: 0;
		font-size: var(--t-sm);
	}

	.hits {
		position: absolute;
		top: calc(100% + 6px);
		right: 0;
		width: min(340px, calc(100vw - 32px));
		background: var(--surface);
		border: 1px solid var(--line);
		box-shadow: 0 12px 28px rgb(13 13 63 / 0.12);
	}

	.hits a {
		display: flex;
		align-items: center;
		gap: 12px;
		min-height: 60px;
		padding: 8px 12px;
		font-size: var(--t-sm);
		font-weight: 500;
	}

	.hits a:hover {
		background: var(--soft);
	}

	.hits img {
		object-fit: cover;
	}

	.hits small {
		display: block;
		color: var(--muted);
		font-size: var(--t-xs);
		font-weight: 400;
	}

	.icon {
		position: relative;
		display: grid;
		place-items: center;
		width: 44px;
		height: 44px;
	}

	.icon:hover {
		color: var(--brand);
	}

	.count {
		position: absolute;
		top: 4px;
		right: 2px;
		display: grid;
		place-items: center;
		min-width: 18px;
		height: 18px;
		padding: 0 4px;
		background: var(--accent);
		color: var(--on-dark);
		border-radius: 999px;
		font-size: 11px;
		font-weight: 700;
	}

	.burger,
	.mobile-only {
		display: none;
	}

	@media (max-width: 1100px) {
		.top .mid {
			display: none;
		}

		.top .row {
			grid-template-columns: 1fr auto;
		}

		nav {
			gap: 20px;
		}

		.search {
			width: 160px;
		}
	}

	@media (max-width: 900px) {
		.top {
			display: none;
		}

		.bar .row {
			position: relative;
			justify-content: space-between;
			gap: 8px;
		}

		.burger,
		.mobile-only {
			display: grid;
		}

		.logo {
			position: absolute;
			left: 50%;
			transform: translateX(-50%);
		}

		.logo span {
			display: none;
		}

		nav {
			position: absolute;
			top: 100%;
			left: 0;
			right: 0;
			display: none;
			flex-direction: column;
			/* centred overflow puts the first items above the scroll start, out of reach */
			justify-content: flex-start;
			gap: 0;
			padding: 8px var(--gutter) 16px;
			background: var(--surface);
			border-bottom: 1px solid var(--line);
		}

		nav.open {
			display: flex;
			max-height: calc(100dvh - var(--bar) - 56px);
			overflow-y: auto;
		}

		nav a {
			min-height: 48px;
			font-size: var(--t-md);
		}

		.sub {
			position: static;
			display: block;
			width: auto;
			padding-block: 0 8px;
			border: 0;
			box-shadow: none;
		}

		nav .sub a {
			min-height: 44px;
			padding-inline: 12px 0;
			font-size: var(--t-sm);
		}

		nav small {
			position: static;
			margin-left: 8px;
			top: auto;
		}

		.search {
			position: absolute;
			top: 100%;
			left: 0;
			right: 0;
			display: none;
			width: auto;
			margin: 0;
			background: var(--surface);
			border: 0;
			border-bottom: 1px solid var(--line);
			border-radius: 0;
			padding-inline: var(--gutter);
		}

		.search.open {
			display: flex;
		}

		.search input {
			font-size: var(--t-md);
		}

		.hits {
			left: var(--gutter);
			right: var(--gutter);
			width: auto;
		}

		.desktop-only {
			display: none;
		}
	}
</style>
