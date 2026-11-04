<script>
	import BrandIcon from './BrandIcon.svelte';
	import { catalog } from '#lib/catalog.js';
	import { brand, library, logo, nav, social } from '#lib/site.js';
	import { addLead, focusError, say } from '#lib/store.svelte.js';

	const courses = $derived(catalog().courses);

	let email = $state('');
	let error = $state('');
	let sending = $state(false);

	async function subscribe(e) {
		e.preventDefault();
		if (!/^\S+@\S+\.\S+$/.test(email)) {
			error = 'Email chưa đúng định dạng';
			return focusError(e.currentTarget);
		}
		if (sending) return;
		sending = true;
		try {
			await addLead({ kind: 'newsletter', email, source: 'tin-moi' });
			email = '';
			error = '';
			say('Đã đăng ký nhận tin');
		} catch (err) {
			error = err.message;
		} finally {
			sending = false;
		}
	}
</script>

<footer class="dark">
	<div class="container cols">
		<div class="about">
			<a class="logo" href="/" aria-label="easyGYM, về trang chủ">
				<img src={logo} alt="" width="56" height="56" loading="lazy" />
				<span>easyGYM</span>
			</a>
			<p>{brand.line}</p>
		</div>

		<nav aria-label="Trang">
			<h2>Khám phá</h2>
			<ul>
				{#each nav as n}<li><a href={n.href}>{n.label}</a></li>{/each}
				{#each library as n}<li><a href={n.href}>{n.label}</a></li>{/each}
			</ul>
		</nav>

		<nav aria-label="Khóa học">
			<h2>Khóa học</h2>
			<ul>
				{#each courses as c}<li><a href={`/khoa-tap/${c.slug}`}>{c.name}</a></li>{/each}
			</ul>
		</nav>

		<div>
			<h2>Liên hệ</h2>
			<ul>
				<li><a href={brand.map} target="_blank" rel="noreferrer">{brand.address}</a></li>
				<li><a href={`tel:${brand.phone}`}>{brand.phone}</a></li>
				<li><a href={`mailto:${brand.email}`}>{brand.email}</a></li>
			</ul>
			<div class="socials">
				{#each social as s}
					<a href={s.href} target="_blank" rel="noreferrer" aria-label={s.label} title={s.label}><BrandIcon name={s.icon} size={s.icon === 'zalo' ? 32 : 22} /></a>
				{/each}
			</div>

			<h2 class="news">Đăng ký nhận tin</h2>
			<p class="hint">Để lại email để nhận thông tin về các khóa học và bài tập cùng nhiều ưu đãi hấp dẫn khác</p>
			<form onsubmit={subscribe} novalidate>
				<label class="sr-only" for="email-tin">Email</label>
				<input id="email-tin" type="email" autocomplete="email" placeholder="Email của bạn" bind:value={email} aria-invalid={Boolean(error)} aria-describedby="email-loi" />
				<button class="btn light" disabled={sending}>Gửi</button>
			</form>
			<p class="err" id="email-loi" aria-live="polite">{error}</p>
		</div>
	</div>

	<div class="container note">© {new Date().getFullYear()} easyGYM. Bảo lưu mọi quyền.</div>
</footer>

<style>
	footer {
		background: var(--ink);
		color: var(--on-dark-2);
		font-size: var(--t-sm);
	}

	.cols {
		display: grid;
		grid-template-columns: 1.3fr 0.9fr 1.3fr 1.6fr;
		gap: 32px var(--gutter);
		padding-block: 56px 32px;
	}

	.logo {
		display: flex;
		align-items: center;
		gap: 12px;
		margin-bottom: 16px;
		color: var(--on-dark);
		font-size: var(--t-lg);
		font-weight: 700;
	}

	.logo img {
		border-radius: 50%;
	}

	.about p {
		max-width: 34ch;
	}

	h2 {
		margin-bottom: 12px;
		color: var(--on-dark);
		font-size: var(--t-sm);
		font-weight: 600;
		text-transform: uppercase;
	}

	.news {
		margin-top: 20px;
		margin-bottom: 6px;
	}

	li a {
		display: inline-flex;
		align-items: center;
		min-height: 36px;
	}

	li a:hover {
		color: var(--on-dark);
		text-decoration: underline;
	}

	.socials {
		display: flex;
		flex-wrap: wrap;
		gap: 8px;
		margin-top: 12px;
	}

	.socials a {
		display: grid;
		place-items: center;
		width: 44px;
		height: 44px;
		background: var(--surface);
		border-radius: 50%;
		transition: transform 0.15s;
	}

	.socials a:hover {
		transform: translateY(-2px);
	}

	form {
		display: flex;
		gap: 8px;
		margin-top: 10px;
	}

	input {
		flex: 1;
		min-width: 0;
		height: 48px;
		padding: 0 12px;
		background: var(--surface);
		color: var(--ink);
		border: 0;
	}

	.err {
		min-height: 20px;
		margin-top: 4px;
		color: #ffc2cf;
	}

	.note {
		padding-block: 16px;
		border-top: 1px solid var(--navy);
		font-size: var(--t-xs);
	}

	@media (max-width: 1000px) {
		.cols {
			grid-template-columns: 1fr 1fr;
		}
	}

	@media (max-width: 560px) {
		.cols {
			grid-template-columns: 1fr;
		}
	}
</style>
