<script>
	import Faq from '#lib/components/Faq.svelte';
	import SignupForm from '#lib/components/SignupForm.svelte';

	let { data } = $props();
</script>

<svelte:head>
	<title>FAQ | easyGYM</title>
</svelte:head>

<div class="container page">
	<p class="crumbs"><a href="/">Trang chủ</a> / <span>FAQ</span></p>
	<h1 class="heading">Câu hỏi thường gặp</h1>

	<div class="layout">
		<div class="groups">
			{#each data.groups as g}
				<section>
					<h2><a href={`/san-pham/${g.slug}`}>{g.name}</a></h2>
					<Faq items={g.items} />
				</section>
			{/each}
		</div>
		<aside id="dang-ky">
			<h2>Gửi câu hỏi cho chúng tôi</h2>
			<SignupForm id="dk-faq" />
		</aside>
	</div>
</div>

<style>
	.layout {
		display: grid;
		grid-template-columns: minmax(0, 1fr) 360px;
		gap: clamp(24px, 4vw, 56px);
		align-items: start;
	}

	.groups {
		display: grid;
		gap: 32px;
	}

	h2 {
		margin-bottom: 12px;
		font-size: var(--t-lg);
		font-weight: 600;
		text-transform: uppercase;
	}

	h2 a {
		display: inline-block;
		padding-block: 11px;
	}

	h2 a:hover {
		color: var(--brand);
	}

	aside {
		position: sticky;
		top: calc(var(--bar) + 16px);
		padding: 20px;
		background: var(--soft);
	}

	@media (max-width: 960px) {
		.layout {
			grid-template-columns: minmax(0, 1fr);
		}

		aside {
			position: static;
		}
	}
</style>
