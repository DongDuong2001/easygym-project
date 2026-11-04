<script>
	let { post: p } = $props();
</script>

<article class="card">
	<a class="media" href={`/tin-tuc/${p.slug}`} tabindex="-1" aria-hidden="true">
		<img src={p.image} alt="" loading="lazy" />
	</a>
	<p class="date">{p.date}</p>
	<h3><a href={`/tin-tuc/${p.slug}`}>{p.title}</a></h3>
	<p class="desc">{p.excerpt}</p>
</article>

<style>
	.card {
		display: flex;
		flex-direction: column;
		min-width: 0;
	}

	.media {
		display: block;
		aspect-ratio: 3 / 2;
		margin-bottom: 12px;
		overflow: hidden;
		background: var(--soft);
	}

	.media img {
		width: 100%;
		height: 100%;
		object-fit: cover;
		transition: transform 0.35s;
	}

	.card:hover img {
		transform: scale(1.03);
	}

	.date {
		color: var(--muted);
		font-size: var(--t-xs);
		font-weight: 600;
	}

	h3 {
		margin-top: 4px;
		font-size: var(--t-md);
		font-weight: 600;
		line-height: 1.3;
	}

	h3 a {
		display: block;
		min-height: 44px;
	}

	h3 a:hover {
		color: var(--brand);
	}

	.desc {
		display: -webkit-box;
		margin-top: 6px;
		overflow: hidden;
		color: var(--muted);
		font-size: var(--t-sm);
		-webkit-box-orient: vertical;
		-webkit-line-clamp: 3;
		line-clamp: 3;
	}
</style>
