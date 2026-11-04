<script>
	import PostCard from '#lib/components/PostCard.svelte';

	let { data } = $props();

	const p = $derived(data.post);
	const related = $derived(data.posts.filter((x) => x.slug !== p.slug && x.cat === p.cat).slice(0, 3));
</script>

<svelte:head>
	<title>{p.title} | easyGYM</title>
	<meta name="description" content={p.excerpt} />
</svelte:head>

<div class="container page">
	<p class="crumbs"><a href="/">Trang chủ</a> / <a href="/tin-tuc">Tin tức</a></p>

	<article>
		<p class="date">{p.date}</p>
		<h1>{p.title}</h1>
		{#if !data.body.includes(p.image)}<img class="cover" src={p.image} alt={p.title} />{/if}
		<div class="prose">{@html data.body}</div>
	</article>

	<section class="more">
		<h2 class="heading">Bài viết liên quan</h2>
		<div class="posts">
			{#each related as r (r.slug)}<PostCard post={r} />{/each}
		</div>
	</section>
</div>

<style>
	article {
		max-width: 760px;
		margin-inline: auto;
	}

	.date {
		color: var(--muted);
		font-size: var(--t-sm);
		font-weight: 600;
	}

	h1 {
		margin-block: 8px 20px;
		font-size: clamp(28px, 4vw, 44px);
		font-weight: 600;
	}

	.cover {
		width: 100%;
		margin-bottom: 24px;
		background: var(--soft);
	}

	.more {
		margin-top: var(--section);
	}

	.posts {
		display: grid;
		grid-template-columns: repeat(3, minmax(0, 1fr));
		gap: 24px var(--gap);
	}

	@media (max-width: 800px) {
		.posts {
			grid-template-columns: minmax(0, 1fr);
		}
	}
</style>
