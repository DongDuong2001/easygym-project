<script>
	import '../app.css';
	import { onMount } from 'svelte';
	import BrandIcon from '#lib/components/BrandIcon.svelte';
	import Footer from '#lib/components/Footer.svelte';
	import Header from '#lib/components/Header.svelte';
	import Icon from '#lib/components/Icon.svelte';
	import SignupForm from '#lib/components/SignupForm.svelte';
	import { brand, logo } from '#lib/site.js';
	import { loadStore, toast, ui } from '#lib/store.svelte.js';

	let { children } = $props();
	let showToast = $state(false);
	let box = $state();

	onMount(loadStore);

	$effect(() => {
		if (ui.opened) box.showModal();
	});

	$effect(() => {
		if (!toast.at) return;
		showToast = true;
		const id = setTimeout(() => (showToast = false), 2800);
		return () => clearTimeout(id);
	});
</script>

<svelte:head>
	<link rel="icon" href={logo} />
	<title>easyGYM | Khóa tập gym và thực phẩm easyMEAL</title>
	<meta name="description" content="Khóa tập 1-1, tập nhóm, coach online cùng easyGYM và thực phẩm easyMEAL: kefir, bánh Bông Lean, Chocolean." />
</svelte:head>

<Header />
<main id="noi-dung">
	{@render children()}
</main>
<Footer />

<div class="toast" class:show={showToast} role="status" aria-live="polite">
	{#if showToast}
		<Icon name="check" size={18} />
		{toast.text}
		{#if toast.cart}<a href="/gio-hang">Xem giỏ hàng</a>{/if}
	{/if}
</div>

<!-- A slim bar instead of the old floating dock, so nothing sits on top of prices or text. -->
<nav class="reach" aria-label="Liên hệ nhanh">
	<a href={`tel:${brand.phone}`}><Icon name="phone" size={18} /> Gọi</a>
	<a href={brand.zalo} target="_blank" rel="noreferrer" aria-label="Nhắn Zalo"><BrandIcon name="zalo" size={30} /></a>
	<button class="main" onclick={() => box.showModal()}>Đăng ký tập</button>
</nav>

<dialog bind:this={box} aria-labelledby="dk-tieu-de">
	<button class="close" onclick={() => box.close()} aria-label="Đóng"><Icon name="x" size={22} /></button>
	<h2 id="dk-tieu-de">Đăng ký tham gia khóa học</h2>
	{#key ui.opened}
		<SignupForm id="dk-popup" preset={ui.course} />
	{/key}
</dialog>

<style>
	main {
		min-height: 60vh;
	}

	.toast {
		position: fixed;
		left: 50%;
		bottom: 24px;
		z-index: 50;
		display: flex;
		align-items: center;
		gap: 10px;
		max-width: calc(100vw - 32px);
		padding: 12px 18px;
		background: var(--ink);
		color: var(--on-dark);
		font-size: var(--t-sm);
		font-weight: 500;
		white-space: nowrap;
		opacity: 0;
		pointer-events: none;
		transform: translate(-50%, 12px);
		transition: opacity 0.2s, transform 0.2s;
	}

	.toast.show {
		opacity: 1;
		pointer-events: auto;
		transform: translate(-50%, 0);
	}

	.toast a {
		margin-left: 6px;
		text-decoration: underline;
	}

	.reach {
		display: none;
	}

	dialog {
		width: min(calc(100vw - 32px), 460px);
		max-height: calc(100dvh - 32px);
		padding: 24px;
		background: var(--surface);
		color: var(--ink);
		border: 0;
	}

	dialog::backdrop {
		background: rgb(13 13 63 / 0.7);
	}

	dialog h2 {
		margin-bottom: 16px;
		padding-right: 44px;
		font-size: var(--t-lg);
		font-weight: 600;
		text-transform: uppercase;
	}

	.close {
		position: absolute;
		top: 12px;
		right: 12px;
		display: grid;
		place-items: center;
		width: 44px;
		height: 44px;
	}

	@media (max-width: 900px) {
		:global(body) {
			padding-bottom: 56px;
		}

		.toast {
			bottom: 68px;
		}

		.reach {
			position: fixed;
			inset: auto 0 0;
			z-index: 40;
			display: grid;
			grid-template-columns: 1fr 1fr 1.6fr;
			height: 56px;
			background: var(--surface);
			border-top: 1px solid var(--line);
		}

		.reach a,
		.reach button {
			display: flex;
			align-items: center;
			justify-content: center;
			gap: 6px;
			font-size: var(--t-sm);
			font-weight: 600;
		}

		.reach .main {
			background: var(--ink);
			color: var(--on-dark);
			text-transform: uppercase;
		}
	}
</style>
