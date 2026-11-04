<script>
	import { goto } from '$app/navigation';
	import Faq from '#lib/components/Faq.svelte';
	import Icon from '#lib/components/Icon.svelte';
	import ProductCard from '#lib/components/ProductCard.svelte';
	import { catalog } from '#lib/catalog.js';
	import { formatVND } from '#lib/shop.js';
	import { addToCart, store, toggleWish } from '#lib/store.svelte.js';

	const CATEGORIES = $derived(catalog().categories);
	const products = $derived(catalog().products);

	let { data } = $props();

	const p = $derived(data.product);
	const info = $derived(data.info);
	const images = $derived(info.images.length ? info.images : [p.image]);
	const wished = $derived(store.wishes.includes(p.slug));
	const related = $derived([...products.filter((x) => x.category === p.category), ...products.filter((x) => x.category !== p.category)].filter((x) => x.slug !== p.slug).slice(0, 4));

	let variant = $state(0);
	let qty = $state(1);
	let shot = $state(0);

	$effect(() => {
		p.slug;
		variant = 0;
		qty = 1;
		shot = 0;
	});

	function buyNow() {
		addToCart(p.slug, variant, qty);
		goto('/gio-hang');
	}
</script>

<svelte:head>
	<title>{p.name} | easyGYM</title>
	<meta name="description" content={`${p.name}: ${p.tag}.`} />
</svelte:head>

<div class="container page">
	<p class="crumbs"><a href="/">Trang chủ</a> / <a href="/cua-hang">Cửa hàng</a> / <a href={`/cua-hang?dm=${p.category}`}>{CATEGORIES.find((c) => c.value === p.category).label}</a></p>

	<div class="detail">
		<div>
			<div class="media"><img src={images[shot]} alt={p.name} /></div>
			{#if images.length > 1}
				<div class="thumbs">
					{#each images as src, i}
						<button class:on={shot === i} onclick={() => (shot = i)} aria-label={`Xem ảnh ${i + 1}`}><img {src} alt="" loading="lazy" /></button>
					{/each}
				</div>
			{/if}
		</div>

		<div class="info">
			<p class="tag">{p.tag}</p>
			<h1>{p.name}</h1>
			<p class="stars" role="img" aria-label={`${p.reviews} đánh giá`}>
				{#each [1, 2, 3, 4, 5] as _}<Icon name="star" size={15} fill="currentColor" />{/each}
				<span aria-hidden="true">{p.reviews} đánh giá</span>
			</p>
			<p class="price">{formatVND(p.variants[variant].price)}</p>

			{#if info.desc}<div class="prose">{@html info.desc}</div>{/if}

			{#if p.variants.length > 1}
				<fieldset>
					<legend>Chọn loại</legend>
					<div class="opts">
						{#each p.variants as v, i}
							<label><input type="radio" name="loai" value={i} bind:group={variant} />{v.label}</label>
						{/each}
					</div>
				</fieldset>
			{/if}

			<div class="buy">
				<div class="qty" role="group" aria-label="Số lượng">
					<button onclick={() => (qty = Math.max(1, qty - 1))} aria-label="Giảm số lượng"><Icon name="minus" size={16} /></button>
					<output aria-live="polite">{qty}</output>
					<button onclick={() => (qty = Math.min(99, qty + 1))} aria-label="Tăng số lượng"><Icon name="plus" size={16} /></button>
				</div>
				<button class="btn grow" onclick={() => addToCart(p.slug, variant, qty)}>Thêm vào giỏ</button>
				<button class="wish" class:on={wished} onclick={() => toggleWish(p.slug)} aria-pressed={wished} aria-label={`Yêu thích ${p.name}`}>
					<Icon name="heart" fill={wished ? 'currentColor' : 'none'} />
				</button>
			</div>
			<button class="btn line block" onclick={buyNow}>Mua ngay</button>
		</div>
	</div>

	<div class="below">
		{#if info.detail}
			<section>
				<h2>Nội dung chi tiết</h2>
				<div class="prose">{@html info.detail}</div>
			</section>
		{/if}
		{#if info.questions.length}
			<section>
				<h2>Câu hỏi thường gặp</h2>
				<Faq items={info.questions} />
			</section>
		{/if}
	</div>

	<section class="more">
		<h2 class="heading">Có thể bạn thích</h2>
		<div class="grid">
			{#each related as r (r.slug)}<ProductCard product={r} />{/each}
		</div>
	</section>
</div>

<style>
	.detail {
		display: grid;
		grid-template-columns: minmax(0, 520px) minmax(0, 1fr); /* source photos are under 500px wide */
		gap: clamp(24px, 4vw, 64px);
		align-items: start;
	}

	.media {
		aspect-ratio: 4 / 5;
		width: 100%;
		background: var(--soft);
	}

	.media img {
		width: 100%;
		height: 100%;
		object-fit: cover;
	}

	.thumbs {
		display: grid;
		grid-template-columns: repeat(6, minmax(0, 1fr));
		gap: 6px;
		margin-top: 6px;
	}

	.thumbs button {
		aspect-ratio: 1;
		min-height: 44px;
		overflow: hidden;
		background: var(--soft);
		border: 2px solid transparent;
	}

	.thumbs button.on {
		border-color: var(--ink);
	}

	.thumbs img {
		width: 100%;
		height: 100%;
		object-fit: cover;
	}

	.info {
		display: grid;
		gap: 16px;
		max-width: 560px;
	}

	.tag {
		font-size: var(--t-xs);
		font-weight: 600;
		text-transform: uppercase;
		color: var(--muted);
	}

	h1 {
		margin-top: -10px;
		font-size: var(--t-xl);
		font-weight: 600;
		text-transform: uppercase;
	}

	.stars {
		display: flex;
		align-items: center;
		gap: 2px;
	}

	.stars span {
		margin-left: 8px;
		color: var(--muted);
		font-size: var(--t-sm);
	}

	.price {
		font-size: 28px;
		font-weight: 700;
		line-height: 1.1;
	}

	.buy {
		display: flex;
		gap: 8px;
	}

	.grow {
		flex: 1;
		padding-inline: 12px;
	}

	.wish {
		display: grid;
		place-items: center;
		flex: none;
		width: 48px;
		border: 1px solid var(--line);
	}

	.wish:hover,
	.wish.on {
		color: var(--accent);
	}

	.below {
		display: grid;
		grid-template-columns: repeat(2, minmax(0, 1fr));
		gap: 40px clamp(24px, 4vw, 64px);
		align-items: start;
		margin-top: var(--section);
	}

	.below h2 {
		margin-bottom: 16px;
		font-size: var(--t-lg);
		font-weight: 600;
		text-transform: uppercase;
	}

	.below section:only-child {
		grid-column: 1 / -1;
	}

	.more {
		margin-top: var(--section);
	}

	@media (max-width: 860px) {
		.detail,
		.below {
			grid-template-columns: minmax(0, 1fr);
		}

		.media {
			max-width: 520px;
			margin-inline: auto;
		}

		.info {
			max-width: none;
		}
	}
</style>
