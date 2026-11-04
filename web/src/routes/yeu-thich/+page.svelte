<script>
	import ProductCard from '#lib/components/ProductCard.svelte';
	import { catalog } from '#lib/catalog.js';
	import { store } from '#lib/store.svelte.js';

	const products = $derived(catalog().products);

	const list = $derived(products.filter((p) => store.wishes.includes(p.slug)));
</script>

<svelte:head>
	<title>Yêu thích | easyGYM</title>
</svelte:head>

<div class="container page">
	<p class="crumbs"><a href="/">Trang chủ</a> / <span>Yêu thích</span></p>
	<h1 class="heading">Sản phẩm yêu thích</h1>

	{#if list.length}
		<h2 class="sr-only">Danh sách sản phẩm</h2>
		<div class="grid">
			{#each list as p (p.slug)}<ProductCard product={p} />{/each}
		</div>
	{:else}
		<div class="empty">
			<p>Bấm biểu tượng trái tim trên sản phẩm để lưu lại đây.</p>
			<a class="btn" href="/cua-hang">Xem sản phẩm</a>
		</div>
	{/if}
</div>
