<script>
	import Icon from './Icon.svelte';
	import { brand } from '#lib/site.js';

	let { height = 420 } = $props();

	// shortcut: keyless embed URL, switch to the Maps Embed API with a key if Google ever retires it.
	const src = `https://maps.google.com/maps?q=${brand.geo}&z=17&hl=vi&output=embed`;
	const directions = `https://www.google.com/maps/dir/?api=1&destination=${brand.geo}`;
</script>

<figure class="map" style:--h={`${height}px`}>
	<iframe {src} title="Bản đồ đường đến easyGYM" loading="lazy" referrerpolicy="no-referrer-when-downgrade"></iframe>
	<figcaption>
		<span class="pin"><Icon name="pin" size={22} /></span>
		<span class="text">
			<b>{brand.name}</b>
			<span>{brand.address}</span>
			<span>{brand.hours}</span>
		</span>
		<a class="btn light" href={directions} target="_blank" rel="noreferrer">Chỉ đường <Icon name="arrow" size={18} /></a>
	</figcaption>
</figure>

<style>
	.map {
		position: relative;
		height: var(--h);
		background: var(--soft);
	}

	iframe {
		display: block;
		width: 100%;
		height: 100%;
		border: 0;
		filter: grayscale(0.75) contrast(1.05);
		transition: filter 0.3s ease;
	}

	.map:hover iframe {
		filter: none;
	}

	figcaption {
		position: absolute;
		bottom: 16px;
		left: 16px;
		display: flex;
		align-items: center;
		gap: 16px;
		max-width: min(520px, calc(100% - 32px));
		padding: 16px;
		background: var(--ink);
		color: var(--on-dark);
		box-shadow: 0 12px 32px rgb(13 13 63 / 0.25);
	}

	.pin {
		display: grid;
		flex: none;
		place-items: center;
		width: 48px;
		height: 48px;
		background: var(--accent);
	}

	.text {
		display: grid;
		gap: 2px;
		min-width: 0;
		font-size: var(--t-sm);
		color: var(--on-dark-2);
	}

	.text b {
		color: var(--on-dark);
		font-size: var(--t-md);
		font-weight: 700;
		text-transform: uppercase;
	}

	.btn {
		flex: none;
	}

	@media (prefers-reduced-motion: reduce) {
		iframe {
			transition: none;
		}
	}

	/* On phones the card sits under the map so it never hides the pin. */
	@media (max-width: 640px) {
		.map {
			height: auto;
		}

		iframe {
			height: 300px;
		}

		figcaption {
			position: static;
			flex-wrap: wrap;
			max-width: none;
			box-shadow: none;
		}

		.text {
			flex: 1 1 200px;
		}

		.btn {
			width: 100%;
			justify-content: center;
		}
	}
</style>
