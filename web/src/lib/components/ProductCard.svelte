<script>
	import Icon from './Icon.svelte';
	import { formatVND, fromPrice } from '#lib/shop.js';
	import { addToCart, store, toggleWish } from '#lib/store.svelte.js';

	let { product: p } = $props();

	const wished = $derived(store.wishes.includes(p.slug));
	const many = $derived(new Set(p.variants.map((v) => v.price)).size > 1);
</script>

<article class="card">
	<div class="media">
		<a href={`/san-pham/${p.slug}`} tabindex="-1" aria-hidden="true">
			<img src={p.image} alt="" loading="lazy" />
		</a>
		<p class="strip">{p.tag}</p>
		<button class="wish" class:on={wished} onclick={() => toggleWish(p.slug)} aria-pressed={wished} aria-label={`Yêu thích ${p.name}`}>
			<Icon name="heart" fill={wished ? 'currentColor' : 'none'} />
		</button>
	</div>

	<h3><a href={`/san-pham/${p.slug}`}>{p.name}</a></h3>
	<p class="stars" role="img" aria-label={`${p.reviews} đánh giá`}>
		{#each [1, 2, 3, 4, 5] as _}<Icon name="star" size={13} fill="currentColor" />{/each}
		<span aria-hidden="true">({p.reviews})</span>
	</p>
	<div class="foot">
		<p class="price">{many ? 'Từ ' : ''}{formatVND(fromPrice(p.variants))}</p>
		<button class="add" onclick={() => addToCart(p.slug)} aria-label={`Thêm ${p.name} ${p.variants[0].label} vào giỏ`.replace('  ', ' ')}>
			<Icon name="plus" />
		</button>
	</div>
</article>

<style>
	.card {
		display: flex;
		flex-direction: column;
		min-width: 0;
	}

	.media {
		position: relative;
		aspect-ratio: 4 / 5;
		margin-bottom: 12px;
		overflow: hidden;
		background: var(--soft);
	}

	.media a,
	.media img {
		display: block;
		width: 100%;
		height: 100%;
	}

	.media img {
		object-fit: cover;
		transition: transform 0.35s;
	}

	.card:hover img {
		transform: scale(1.03);
	}

	.strip {
		position: absolute;
		inset: auto 0 0;
		padding: 5px 12px;
		background: var(--ink);
		color: var(--on-dark);
		font-size: var(--t-xs);
		font-weight: 600;
		text-transform: uppercase;
		white-space: nowrap;
		overflow: hidden;
		text-overflow: ellipsis;
	}

	.wish {
		position: absolute;
		top: 8px;
		right: 8px;
		display: grid;
		place-items: center;
		width: 44px;
		height: 44px;
		background: var(--surface);
	}

	.wish:hover,
	.wish.on {
		color: var(--accent);
	}

	h3 {
		font-size: var(--t-sm);
		font-weight: 600;
		line-height: 1.35;
		text-transform: uppercase;
	}

	/* Two-line height for every title: a full tap target, and prices line up across the row. */
	h3 a {
		display: block;
		min-height: 44px;
	}

	h3 a:hover {
		color: var(--brand);
	}

	.stars {
		display: flex;
		align-items: center;
		gap: 1px;
		margin-top: 6px;
		color: var(--ink);
	}

	.stars span {
		margin-left: 6px;
		color: var(--muted);
		font-size: var(--t-xs);
	}

	.foot {
		display: flex;
		align-items: center;
		justify-content: space-between;
		gap: 8px;
		margin-top: auto;
		padding-top: 2px;
	}

	.price {
		font-weight: 600;
	}

	.add {
		display: grid;
		place-items: center;
		flex: none;
		width: 44px;
		height: 44px;
		border: 1px solid var(--line);
	}

	.add:hover {
		background: var(--ink);
		border-color: var(--ink);
		color: var(--on-dark);
	}
</style>
