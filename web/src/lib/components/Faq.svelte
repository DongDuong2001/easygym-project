<script>
	import Icon from './Icon.svelte';

	let { items } = $props();
</script>

<div class="faq">
	{#each items as item}
		<details>
			<summary>{item.q}<Icon name="down" /></summary>
			<div class="prose">{@html item.a}</div>
		</details>
	{/each}
</div>

<style>
	.faq {
		border-top: 1px solid var(--line);
	}

	details {
		border-bottom: 1px solid var(--line);
	}

	summary {
		display: flex;
		align-items: center;
		justify-content: space-between;
		gap: 16px;
		min-height: 52px;
		padding-block: 12px;
		font-weight: 600;
		cursor: pointer;
		list-style: none;
	}

	summary::-webkit-details-marker {
		display: none;
	}

	summary :global(svg) {
		flex: none;
		transition: transform 0.15s;
	}

	details[open] summary :global(svg) {
		transform: rotate(180deg);
	}

	.prose {
		padding-bottom: 16px;
	}
</style>
