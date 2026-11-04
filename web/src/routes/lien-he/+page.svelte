<script>
	import BrandIcon from '#lib/components/BrandIcon.svelte';
	import Icon from '#lib/components/Icon.svelte';
	import MapEmbed from '#lib/components/MapEmbed.svelte';
	import SignupForm from '#lib/components/SignupForm.svelte';
	import { brand, social } from '#lib/site.js';
</script>

<svelte:head>
	<title>Liên hệ | easyGYM</title>
</svelte:head>

<div class="container page">
	<p class="crumbs"><a href="/">Trang chủ</a> / <span>Liên hệ</span></p>
	<h1 class="heading">Liên hệ easyGYM</h1>

	<div class="layout">
		<div>
			<ul class="info">
				<li><Icon name="pin" /><a href={brand.map} target="_blank" rel="noreferrer">{brand.address}</a></li>
				<li><Icon name="phone" /><a href={`tel:${brand.phone}`}>{brand.phone}</a></li>
				<li><Icon name="mail" /><a href={`mailto:${brand.email}`}>{brand.email}</a></li>
				<li><Icon name="clock" /><span>{brand.hours}</span></li>
			</ul>
			<ul class="socials">
				{#each social as s}
					<li>
						<a href={s.href} target="_blank" rel="noreferrer">
							<span class="logo"><BrandIcon name={s.icon} size={s.icon === 'zalo' ? 34 : 24} /></span>
							{s.label}
						</a>
					</li>
				{/each}
			</ul>
		</div>

		<section id="dang-ky">
			<h2>Đăng ký trải nghiệm</h2>
			<SignupForm id="dk-lien-he" />
		</section>
	</div>

	<div class="where">
		<MapEmbed height={480} />
	</div>
</div>

<style>
	.layout {
		display: grid;
		grid-template-columns: minmax(0, 1fr) minmax(0, 1fr);
		gap: var(--gap);
		align-items: start;
	}

	.info {
		border-top: 1px solid var(--line);
	}

	.info li {
		display: flex;
		align-items: center;
		gap: 12px;
		min-height: 56px;
		padding-block: 8px;
		border-bottom: 1px solid var(--line);
	}

	.info :global(svg) {
		flex: none;
		color: var(--brand);
	}

	.info a {
		display: inline-flex;
		align-items: center;
		min-height: 44px;
	}

	.info a:hover {
		text-decoration: underline;
	}

	.socials {
		display: grid;
		grid-template-columns: repeat(2, minmax(0, 1fr));
		gap: 8px;
		margin-top: 16px;
	}

	.socials a {
		display: flex;
		align-items: center;
		gap: 12px;
		min-height: 56px;
		padding: 6px 12px 6px 6px;
		background: var(--soft);
		font-weight: 600;
	}

	.socials a:hover {
		color: var(--brand);
	}

	.logo {
		display: grid;
		place-items: center;
		flex: none;
		width: 44px;
		height: 44px;
		background: var(--surface);
		border-radius: 50%;
	}

	section {
		display: grid;
		gap: 16px;
		padding: clamp(20px, 3vw, 32px);
		background: var(--soft);
	}

	h2 {
		font-size: var(--t-lg);
		font-weight: 600;
		text-transform: uppercase;
	}

	.where {
		margin-top: var(--gap);
	}

	@media (max-width: 860px) {
		.layout {
			grid-template-columns: minmax(0, 1fr);
		}
	}
</style>
