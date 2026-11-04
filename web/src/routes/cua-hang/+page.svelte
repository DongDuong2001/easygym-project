<script>
	import { afterNavigate, goto } from '$app/navigation';
	import Icon from '#lib/components/Icon.svelte';
	import ProductCard from '#lib/components/ProductCard.svelte';
	import { catalog } from '#lib/catalog.js';
	import { SORTS, filterProducts, readFilters } from '#lib/shop.js';

	const CATEGORIES = $derived(catalog().categories);
	const products = $derived(catalog().products);

	// Filters come from the query string. The page is prerendered, so the URL is read in the browser only.
	let f = $state(readFilters(new URLSearchParams()));
	let ready = $state(false);

	afterNavigate(() => {
		f = readFilters(new URLSearchParams(location.search));
		ready = true;
	});

	const results = $derived(filterProducts(products, f));
	const title = $derived(f.q ? `Kết quả cho “${f.q}”` : (CATEGORIES.find((c) => c.value === f.dm)?.label ?? 'Tất cả sản phẩm'));

	$effect(() => {
		if (!ready) return;
		const qs = new URLSearchParams(Object.entries(f).filter(([, v]) => v));
		if (`?${qs}` === (location.search || '?')) return;
		goto(qs.size ? `?${qs}` : location.pathname, { shallow: true, replace: true, reset: false });
	});
</script>

<svelte:head>
	<title>{title} | easyGYM</title>
</svelte:head>

<div class="container page">
	<p class="crumbs"><a href="/">Trang chủ</a> / <span>Cửa hàng</span></p>
	<h1 class="heading">{title}</h1>

	<div class="tabs" role="group" aria-label="Lọc theo nhóm sản phẩm">
		<button onclick={() => (f.dm = '')} aria-pressed={!f.dm}>Tất cả</button>
		{#each CATEGORIES as c}
			<button onclick={() => (f.dm = c.value)} aria-pressed={f.dm === c.value}>{c.label}</button>
		{/each}
	</div>

	<div class="bar">
		<label class="q">
			<Icon name="search" size={18} />
			<span class="sr-only">Tìm trong cửa hàng</span>
			<input type="search" placeholder="Tìm theo tên sản phẩm" bind:value={f.q} />
		</label>
		<p class="n" aria-live="polite">{results.length} sản phẩm</p>
		<label class="sort">
			<span class="sr-only">Sắp xếp</span>
			<select bind:value={f['sap-xep']}>
				{#each SORTS as s}<option value={s.value}>{s.label}</option>{/each}
			</select>
		</label>
	</div>

	{#if results.length}
		<h2 class="sr-only">Danh sách sản phẩm</h2>
		<div class="grid">
			{#each results as p (p.slug)}<ProductCard product={p} />{/each}
		</div>
	{:else}
		<div class="empty">
			<p>Chưa có sản phẩm phù hợp với từ khóa này.</p>
			<button class="btn" onclick={() => (f = readFilters(new URLSearchParams()))}>Xóa bộ lọc</button>
		</div>
	{/if}
</div>

<style>
	.bar {
		display: flex;
		align-items: center;
		gap: 16px;
		margin-bottom: 24px;
		padding-block: 12px;
		border-block: 1px solid var(--line);
	}

	.q {
		display: flex;
		align-items: center;
		gap: 10px;
		flex: 0 1 340px;
		color: var(--muted);
	}

	.q input {
		flex: 1;
		min-width: 0;
		height: 44px;
		background: none;
		border: 0;
		outline: 0;
		color: var(--ink);
	}

	.n {
		margin-left: auto;
		color: var(--muted);
		font-size: var(--t-sm);
		white-space: nowrap;
	}

	.sort select {
		height: 44px;
		padding: 0 12px;
		background: var(--surface);
		border: 1px solid var(--line);
		font-size: var(--t-sm);
	}

	@media (max-width: 560px) {
		.bar {
			flex-wrap: wrap;
		}

		.q {
			flex-basis: 100%;
		}

		.n {
			margin-left: 0;
		}

		.sort {
			margin-left: auto;
		}
	}
</style>
