<script>
	import Icon from '#lib/components/Icon.svelte';
	import { catalog } from '#lib/catalog.js';
	import { formatVND, priceCart } from '#lib/shop.js';
	import { setQty, store } from '#lib/store.svelte.js';

	const products = $derived(catalog().products);

	const cart = $derived(priceCart(store.cart, products));
</script>

<svelte:head>
	<title>Giỏ hàng | easyGYM</title>
</svelte:head>

<div class="container page">
	<p class="crumbs"><a href="/">Trang chủ</a> / <span>Giỏ hàng</span></p>
	<h1 class="heading">Giỏ hàng</h1>

	{#if !cart.items.length}
		<div class="empty">
			<p>Giỏ hàng của bạn đang trống.</p>
			<a class="btn" href="/cua-hang">Xem sản phẩm</a>
		</div>
	{:else}
		<div class="layout">
			<ul class="lines">
				{#each cart.items as i (`${i.slug}-${i.variant}`)}
					<li>
						<a class="thumb" href={`/san-pham/${i.slug}`} tabindex="-1" aria-hidden="true"><img src={i.product.image} alt="" /></a>
						<div class="meta">
							<a class="name" href={`/san-pham/${i.slug}`}>{i.product.name}</a>
							{#if i.label}<p>{i.label}</p>{/if}
							<p>{formatVND(i.price)}</p>
						</div>
						<div class="qty" role="group" aria-label={`Số lượng ${i.product.name}`}>
							<button onclick={() => setQty(i.slug, i.variant, i.qty - 1)} aria-label="Giảm số lượng"><Icon name="minus" size={16} /></button>
							<output>{i.qty}</output>
							<button onclick={() => setQty(i.slug, i.variant, i.qty + 1)} aria-label="Tăng số lượng"><Icon name="plus" size={16} /></button>
						</div>
						<p class="total">{formatVND(i.total)}</p>
						<button class="remove" onclick={() => setQty(i.slug, i.variant, 0)} aria-label={`Xóa ${i.product.name}`}><Icon name="trash" /></button>
					</li>
				{/each}
			</ul>

			<aside>
				<p class="row"><span>Tạm tính ({cart.count} sản phẩm)</span><b>{formatVND(cart.total)}</b></p>
				<p class="hint">Phí giao hàng được báo khi xác nhận đơn.</p>
				<a class="btn block" href="/thanh-toan">Thanh toán</a>
				<a class="btn line block" href="/cua-hang">Mua thêm</a>
			</aside>
		</div>
	{/if}
</div>

<style>
	.layout {
		display: grid;
		grid-template-columns: minmax(0, 1fr) 340px;
		gap: 32px;
		align-items: start;
	}

	.lines {
		border-top: 1px solid var(--line);
	}

	li {
		display: grid;
		grid-template-columns: 88px minmax(0, 1fr) auto 120px 44px;
		align-items: center;
		gap: 16px;
		padding-block: 12px;
		border-bottom: 1px solid var(--line);
	}

	.thumb img {
		width: 88px;
		height: 110px;
		object-fit: cover;
		background: var(--soft);
	}

	.name {
		font-size: var(--t-sm);
		font-weight: 600;
		text-transform: uppercase;
	}

	.meta p {
		color: var(--muted);
		font-size: var(--t-sm);
	}

	.total {
		text-align: right;
		font-weight: 600;
	}

	.remove {
		display: grid;
		place-items: center;
		width: 44px;
		height: 44px;
		color: var(--muted);
	}

	.remove:hover {
		color: var(--accent);
	}

	aside {
		position: sticky;
		top: calc(var(--bar) + 16px);
		display: grid;
		gap: 12px;
		padding: 20px;
		background: var(--soft);
	}

	.row {
		display: flex;
		justify-content: space-between;
		gap: 12px;
	}

	.hint {
		color: var(--muted);
		font-size: var(--t-sm);
	}

	@media (max-width: 960px) {
		.layout {
			grid-template-columns: minmax(0, 1fr);
		}

		aside {
			position: static;
		}
	}

	@media (max-width: 640px) {
		li {
			grid-template-columns: 72px minmax(0, 1fr) 44px;
			gap: 8px 12px;
		}

		.thumb {
			grid-row: span 2;
		}

		.thumb img {
			width: 72px;
			height: 90px;
		}

		.meta {
			grid-column: 2;
		}

		.remove {
			grid-column: 3;
			grid-row: 1;
		}

		.qty {
			grid-column: 2;
			justify-self: start;
		}

		.total {
			grid-column: 2 / 4;
			grid-row: 2;
			align-self: center;
		}
	}
</style>
