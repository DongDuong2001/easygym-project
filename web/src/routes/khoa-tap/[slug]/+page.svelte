<script>
	import SignupForm from '#lib/components/SignupForm.svelte';
	import { formatVND } from '#lib/shop.js';

	let { data } = $props();

	const c = $derived(data.course);

	let plan = $state(0);

	$effect(() => {
		c.slug;
		plan = 0;
	});
</script>

<svelte:head>
	<title>{c.name} | easyGYM</title>
	<meta name="description" content={c.desc} />
</svelte:head>


<div class="container page">
	<p class="crumbs"><a href="/">Trang chủ</a> / <a href="/khoa-tap">Khóa tập</a></p>

	<div class="detail">
		<div class="media"><img src={c.image} alt={c.name} /></div>

		<div class="info">
			<p class="tag">{c.mode}</p>
			<h1>{c.name}</h1>
			<p>{c.desc}</p>

			{#if c.plans.length > 1}
				<fieldset>
					<legend>Chọn thời gian</legend>
					<div class="opts">
						{#each c.plans as p, i}
							<label><input type="radio" name="goi" value={i} bind:group={plan} />{p.label}</label>
						{/each}
					</div>
				</fieldset>
			{/if}

			<p class="price">{formatVND(c.plans[plan].price)} <small>/ {c.plans[plan].label}</small></p>

			<section id="dang-ky">
				<h2>Đăng ký khóa này</h2>
				<SignupForm id="dk-khoa" course={c.slug} plan={c.plans[plan].label} />
			</section>
		</div>
	</div>

	{#if data.body}
		<article class="body">
			<h2>Chi tiết khóa học</h2>
			<div class="prose">{@html data.body}</div>
		</article>
	{/if}
</div>

<style>
	.detail {
		display: grid;
		grid-template-columns: minmax(0, 1.05fr) minmax(0, 1fr);
		gap: clamp(24px, 4vw, 64px);
		align-items: start;
	}

	.media {
		aspect-ratio: 4 / 5;
		background: var(--soft);
	}

	.media img {
		width: 100%;
		height: 100%;
		object-fit: cover;
		object-position: top;
	}

	.info {
		display: grid;
		gap: 16px;
		max-width: 560px;
	}

	.tag {
		color: var(--muted);
		font-size: var(--t-xs);
		font-weight: 600;
		text-transform: uppercase;
	}

	h1 {
		margin-top: -10px;
		font-size: var(--t-xl);
		font-weight: 600;
		text-transform: uppercase;
	}

	.price {
		font-size: 28px;
		font-weight: 700;
		line-height: 1.1;
	}

	small {
		color: var(--muted);
		font-size: var(--t-md);
		font-weight: 400;
	}

	.body {
		margin-top: var(--section);
		padding-top: 24px;
		border-top: 1px solid var(--line);
	}

	.body h2 {
		margin-bottom: 16px;
	}

	section {
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

	@media (max-width: 860px) {
		.detail {
			grid-template-columns: minmax(0, 1fr);
		}

		.info {
			max-width: none;
		}
	}
</style>
