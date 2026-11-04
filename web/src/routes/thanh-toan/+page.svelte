<script>
	import Icon from '#lib/components/Icon.svelte';
	import { catalog } from '#lib/catalog.js';
	import { formatVND, isPhone, priceCart } from '#lib/shop.js';
	import { focusError, placeOrder, store } from '#lib/store.svelte.js';

	const products = $derived(catalog().products);

	const PAYMENTS = [
		{ value: 'cod', label: 'Thanh toán khi nhận hàng' },
		{ value: 'bank', label: 'Chuyển khoản, nhân viên gửi thông tin sau' }
	];

	const cart = $derived(priceCart(store.cart, products));

	let form = $state({ name: '', phone: '', address: '', pay: 'cod' });
	let tried = $state(false);
	let placed = $state(null);
	let sending = $state(false);
	let failed = $state('');

	const nameErr = $derived(tried && !form.name.trim() ? 'Vui lòng nhập họ tên' : '');
	const phoneErr = $derived(tried && !isPhone(form.phone) ? 'Số điện thoại gồm 10 số, bắt đầu bằng 0' : '');
	const addressErr = $derived(tried && !form.address.trim() ? 'Vui lòng nhập địa chỉ giao hàng' : '');

	async function submit(e) {
		e.preventDefault();
		tried = true;
		if (!form.name.trim() || !isPhone(form.phone) || !form.address.trim()) return focusError(e.currentTarget);
		if (sending) return;
		sending = true;
		failed = '';
		try {
			placed = await placeOrder({ ...form, lines: cart.items.map((i) => ({ slug: i.slug, label: i.label, qty: i.qty })) });
		} catch (err) {
			failed = err.message;
		} finally {
			sending = false;
		}
	}
</script>

<svelte:head>
	<title>Thanh toán | easyGYM</title>
</svelte:head>

<div class="container page">
	<p class="crumbs"><a href="/">Trang chủ</a> / <a href="/gio-hang">Giỏ hàng</a> / <span>Thanh toán</span></p>

	{#if placed}
		<div class="empty" role="status">
			<Icon name="check" size={40} />
			<h1>Đã nhận đơn {placed.code}</h1>
			<p>Tổng {formatVND(placed.total)}, chưa gồm phí giao hàng. easyGYM sẽ gọi xác nhận đơn trong giờ làm việc.</p>
			<a class="btn" href="/cua-hang">Tiếp tục mua sắm</a>
		</div>
	{:else if !cart.items.length}
		<div class="empty">
			<h1>Chưa có sản phẩm để thanh toán</h1>
			<a class="btn" href="/cua-hang">Xem sản phẩm</a>
		</div>
	{:else}
		<h1 class="heading">Thanh toán</h1>
		<form class="layout" onsubmit={submit} novalidate>
			<div class="fields">
				<label class="field">
					Họ và tên
					<input bind:value={form.name} autocomplete="name" aria-invalid={Boolean(nameErr)} aria-describedby="tt-ten" />
					<span class="err" id="tt-ten">{nameErr}</span>
				</label>
				<label class="field">
					Số điện thoại
					<input type="tel" bind:value={form.phone} autocomplete="tel" inputmode="tel" aria-invalid={Boolean(phoneErr)} aria-describedby="tt-sdt" />
					<span class="err" id="tt-sdt">{phoneErr}</span>
				</label>
				<label class="field">
					Địa chỉ giao hàng
					<input bind:value={form.address} autocomplete="street-address" aria-invalid={Boolean(addressErr)} aria-describedby="tt-dc" />
					<span class="err" id="tt-dc">{addressErr}</span>
				</label>
				<fieldset>
					<legend>Hình thức thanh toán</legend>
					<div class="opts">
						{#each PAYMENTS as p}
							<label><input type="radio" name="pay" value={p.value} bind:group={form.pay} />{p.label}</label>
						{/each}
					</div>
				</fieldset>
			</div>

			<aside>
				<h2>Đơn hàng</h2>
				<ul>
					{#each cart.items as i}
						<li><span>{i.product.name}{i.label ? `, ${i.label}` : ''} × {i.qty}</span><b>{formatVND(i.total)}</b></li>
					{/each}
				</ul>
				<p class="row"><span>Tổng cộng</span><b>{formatVND(cart.total)}</b></p>
				<p class="hint">Chưa gồm phí giao hàng.</p>
				{#if failed}<p class="err" role="alert">{failed}</p>{/if}
				<button class="btn block" disabled={sending}>{sending ? 'Đang gửi...' : 'Đặt hàng'}</button>
			</aside>
		</form>
	{/if}
</div>

<style>
	.layout {
		display: grid;
		grid-template-columns: minmax(0, 1fr) 380px;
		gap: 32px;
		align-items: start;
	}

	.fields {
		display: grid;
		gap: 16px;
		max-width: 640px;
	}

	.err:empty {
		display: none;
	}

	aside {
		display: grid;
		gap: 12px;
		padding: 20px;
		background: var(--soft);
	}

	h2 {
		font-size: var(--t-lg);
		font-weight: 600;
		text-transform: uppercase;
	}

	li,
	.row {
		display: flex;
		justify-content: space-between;
		gap: 16px;
	}

	li {
		padding-block: 8px;
		border-bottom: 1px solid var(--line);
		font-size: var(--t-sm);
	}

	li b {
		white-space: nowrap;
	}

	.hint {
		color: var(--muted);
		font-size: var(--t-sm);
	}

	@media (max-width: 960px) {
		.layout {
			grid-template-columns: minmax(0, 1fr);
		}
	}
</style>
