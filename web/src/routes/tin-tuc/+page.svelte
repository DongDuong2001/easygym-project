<script>
	import PostCard from '#lib/components/PostCard.svelte';

	let { data } = $props();

	const TABS = [
		{ value: '', label: 'Tất cả' },
		{ value: 'tin-tuc', label: 'Tin tức' },
		{ value: 'easymeal', label: 'easyMEAL' }
	];

	let tab = $state('');

	const list = $derived(data.posts.filter((p) => !tab || p.cat === tab));
</script>

<svelte:head>
	<title>Tin tức | easyGYM</title>
</svelte:head>

<div class="container page">
	<p class="crumbs"><a href="/">Trang chủ</a> / <span>Tin tức</span></p>
	<h1 class="heading">Tin tức</h1>

	<div class="tabs" role="group" aria-label="Lọc bài viết">
		{#each TABS as t}
			<button onclick={() => (tab = t.value)} aria-pressed={tab === t.value}>{t.label}</button>
		{/each}
	</div>

	<h2 class="sr-only">Danh sách bài viết</h2>
	<div class="posts">
		{#each list as p (p.slug)}<PostCard post={p} />{/each}
	</div>
</div>

<style>
	.posts {
		display: grid;
		grid-template-columns: repeat(3, minmax(0, 1fr));
		gap: 32px var(--gap);
	}

	@media (max-width: 900px) {
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
