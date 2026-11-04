<script>
	import { openRegister } from '#lib/store.svelte.js';

	// Same card as easygym.vn: photo, title, summary and two actions. Prices stay on the detail page.
	let { course: c } = $props();
</script>

<article class="card">
	<a class="media" href={`/khoa-tap/${c.slug}`} tabindex="-1" aria-hidden="true">
		<img src={c.card} alt="" loading="lazy" />
	</a>
	<h3><a href={`/khoa-tap/${c.slug}`}>{c.name}</a></h3>
	<p class="desc">{c.desc}</p>
	<div class="actions">
		<a class="btn" href={`/khoa-tap/${c.slug}`}>Xem chi tiết</a>
		<button class="btn line" onclick={() => openRegister(c.slug)}>Đăng ký ngay</button>
	</div>
</article>

<style>
	.card {
		display: flex;
		flex-direction: column;
		min-width: 0;
	}

	.media {
		display: block;
		aspect-ratio: 4 / 5;
		margin-bottom: 12px;
		overflow: hidden;
		background: var(--soft);
	}

	.media img {
		width: 100%;
		height: 100%;
		object-fit: cover;
		object-position: top;
		transition: transform 0.35s;
	}

	.card:hover img {
		transform: scale(1.03);
	}

	h3 {
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
		margin-block: 6px 12px;
		overflow: hidden;
		color: var(--muted);
		font-size: var(--t-sm);
		-webkit-box-orient: vertical;
		-webkit-line-clamp: 3;
		line-clamp: 3;
	}

	.actions {
		display: grid;
		grid-template-columns: repeat(2, minmax(0, 1fr));
		gap: 8px;
		margin-top: auto;
	}

	.btn {
		padding-inline: 8px;
		letter-spacing: 0;
		white-space: nowrap;
	}

	@media (max-width: 720px) {
		.actions {
			grid-template-columns: minmax(0, 1fr);
		}
	}
</style>
