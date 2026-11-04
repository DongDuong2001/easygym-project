<script>
	import PostCard from '#lib/components/PostCard.svelte';
	import ProductCard from '#lib/components/ProductCard.svelte';
	import { catalog } from '#lib/catalog.js';

	let { data } = $props();

	const posts = $derived(data.posts);
	const foods = $derived(catalog().products.filter((p) => p.category !== 'ao'));
</script>

<svelte:head>
	<title>easyMEAL | easyGYM</title>
</svelte:head>

<div class="container page">
	<p class="crumbs"><a href="/">Trang chủ</a> / <span>easyMEAL</span></p>
	<h1 class="heading">easyMEAL</h1>
	<p class="lead">Ăn sạch, sống dễ cho người bận rộn.</p>

	<div class="posts">
		{#each posts as p (p.slug)}<PostCard post={p} />{/each}
	</div>

	<section>
		<h2 class="heading">Sản phẩm easyMEAL</h2>
		<div class="grid">
			{#each foods as p (p.slug)}<ProductCard product={p} />{/each}
		</div>
	</section>
</div>

<style>
	.posts {
		display: grid;
		grid-template-columns: repeat(4, minmax(0, 1fr));
		gap: 24px var(--gap);
	}

	section {
		margin-top: var(--section);
	}

	@media (max-width: 1024px) {
		.posts {
			grid-template-columns: repeat(2, minmax(0, 1fr));
		}
	}

	@media (max-width: 560px) {
		.posts {
			grid-template-columns: minmax(0, 1fr);
		}
	}
</style>
