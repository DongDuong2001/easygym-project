<script>
	import BrandIcon from './BrandIcon.svelte';
	import Icon from './Icon.svelte';
	import { social } from '#lib/site.js';

	let { videos } = $props();

	const channel = social.find((s) => s.icon === 'youtube').href;
	// Only one player at a time: the YouTube iframe loads on click, the rest stay as thumbnails.
	let playing = $state('');
</script>

<div class="videos">
	{#each videos as v, i}
		<article class:feature={i === 0}>
			{#if playing === v.id}
				<iframe src={`https://www.youtube.com/embed/${v.id}?autoplay=1&rel=0`} title={v.title} allow="autoplay; accelerometer; clipboard-write; encrypted-media; gyroscope; picture-in-picture; web-share" referrerpolicy="strict-origin-when-cross-origin" allowfullscreen></iframe>
			{:else}
				<button onclick={() => (playing = v.id)} aria-label={`Phát video: ${v.title}`}>
					<img src={`https://i.ytimg.com/vi/${v.id}/maxresdefault.jpg`} alt="" loading="lazy" width="1280" height="720" />
					{#if i === 0}<span class="play"><Icon name="play" size={30} /></span>{/if}
					<span class="cap">
						{#if i > 0}<span class="play"><Icon name="play" size={18} /></span>{/if}
						<span>
							<span class="no">{String(i + 1).padStart(2, '0')}</span>
							<span class="title">{v.title}</span>
						</span>
					</span>
				</button>
			{/if}
		</article>
	{/each}
	<a class="more" href={channel} target="_blank" rel="noreferrer">
		<BrandIcon name="youtube" size={44} />
		<span>Xem thêm trên YouTube</span>
	</a>
</div>

<style>
	.videos {
		display: grid;
		grid-template-columns: repeat(3, minmax(0, 1fr));
		gap: var(--gap);
	}

	article,
	.more {
		min-width: 0;
		aspect-ratio: 16 / 9;
		background: var(--navy);
	}

	.feature {
		grid-column: span 2;
		grid-row: span 2;
		aspect-ratio: auto;
	}

	iframe,
	button {
		display: block;
		width: 100%;
		height: 100%;
		border: 0;
	}

	button {
		position: relative;
		overflow: hidden;
		padding: 0;
		background: var(--navy);
		color: #fff;
		text-align: left;
		cursor: pointer;
	}

	img {
		width: 100%;
		height: 100%;
		object-fit: cover;
		transition: transform 0.5s ease;
	}

	button::after {
		content: '';
		position: absolute;
		inset: 0;
		background: linear-gradient(to top, rgb(13 13 63 / 0.92), rgb(13 13 63 / 0.15) 55%, rgb(13 13 63 / 0.05));
	}

	.play {
		position: relative;
		display: grid;
		flex: none;
		place-items: center;
		width: 48px;
		height: 48px;
		background: var(--accent);
		transition: transform 0.25s ease;
	}

	.feature .play {
		position: absolute;
		z-index: 1;
		top: 50%;
		left: 50%;
		width: 88px;
		height: 88px;
		transform: translate(-50%, -50%);
	}

	.play::before {
		content: '';
		position: absolute;
		inset: -6px;
		border: 1px solid rgb(255 255 255 / 0.45);
		transition: inset 0.25s ease;
	}

	.feature .play::before {
		inset: -10px;
	}

	button:hover img {
		transform: scale(1.05);
	}

	button:hover .play {
		transform: scale(1.06);
	}

	.feature button:hover .play {
		transform: translate(-50%, -50%) scale(1.06);
	}

	button:hover .play::before {
		inset: -10px;
	}

	.feature button:hover .play::before {
		inset: -16px;
	}

	.cap {
		position: absolute;
		z-index: 1;
		right: 0;
		bottom: 0;
		left: 0;
		display: flex;
		align-items: center;
		gap: 18px;
		padding: 16px 16px 16px 22px;
	}

	.feature .cap {
		padding: 28px;
	}

	.no {
		display: block;
		margin-bottom: 2px;
		color: #ffc2cf;
		font-size: var(--t-sm);
		font-weight: 700;
		letter-spacing: 0.08em;
	}

	.title {
		display: -webkit-box;
		overflow: hidden;
		-webkit-box-orient: vertical;
		-webkit-line-clamp: 2;
		font-size: var(--t-md);
		font-weight: 600;
		line-height: 1.3;
	}

	.feature .title {
		font-size: var(--t-xl);
		line-height: 1.15;
		text-transform: uppercase;
	}

	.more {
		display: flex;
		flex-direction: column;
		align-items: center;
		justify-content: center;
		gap: 10px;
		background: var(--surface);
		color: var(--ink);
		font-size: var(--t-sm);
		font-weight: 600;
		text-transform: uppercase;
	}

	.more:hover {
		background: var(--soft);
	}

	@media (prefers-reduced-motion: reduce) {
		img,
		.play,
		.play::before {
			transition: none;
		}
	}

	@media (max-width: 900px) {
		.videos {
			grid-template-columns: repeat(2, minmax(0, 1fr));
		}

		.feature {
			grid-row: auto;
			aspect-ratio: 16 / 9;
		}
	}

	@media (max-width: 560px) {
		.videos {
			grid-template-columns: minmax(0, 1fr);
		}

		.feature {
			grid-column: auto;
		}

		.feature .cap {
			padding: 16px;
		}

		.feature .title {
			font-size: var(--t-lg);
		}

		.feature .play {
			top: 36%;
			width: 56px;
			height: 56px;
		}
	}
</style>
