<script>
	import Icon from './Icon.svelte';
	import { img } from '#lib/site.js';

	// easyGYM's own ten benefits, clockwise from the top right of the ring.
	const benefits = [
		{ icon: 'heart', title: 'Tăng sức khỏe tim mạch', text: 'Giảm nguy cơ cao huyết áp, đột quỵ' },
		{ icon: 'flame', title: 'Giảm mỡ – tăng cơ', text: 'Gym giúp đốt cháy calo hiệu quả.' },
		{ icon: 'bolt', title: 'Tăng sức bền và thể lực', text: 'Cơ thể bạn dẻo dai, ít mệt hơn khi vận động' },
		{ icon: 'bone', title: 'Tăng mật độ xương', text: 'Giúp ngăn ngừa loãng xương, đặc biệt khi lớn tuổi' },
		{ icon: 'smile', title: 'Giảm stress và lo âu', text: 'Cơ thể tiết ra endorphin, giúp tinh thần thoải mái hơn.' },
		{ icon: 'star', title: 'Tăng sự tự tin', text: 'Vóc dáng cải thiện và cảm giác khỏe mạnh' },
		{ icon: 'target', title: 'Tăng cường tập trung', text: 'Giúp tinh thần minh mẫn, làm việc hiệu quả' },
		{ icon: 'moon', title: 'Quản lý cảm xúc tốt hơn', text: 'Giúp ngủ sâu, giảm mất ngủ' },
		{ icon: 'shield', title: 'Tăng cường miễn dịch', text: 'Giúp cơ thể khỏe mạnh, ít bệnh tật' },
		{ icon: 'spine', title: 'Cải thiện tư thế', text: 'Giữ cột sống thẳng, giảm đau lưng' }
	];

	// Evenly spaced from 72° above 3 o'clock, so no label ever sits straight above or below the figure.
	const place = (i) => {
		const a = ((-72 + i * 36) * Math.PI) / 180;
		return { x: 50 + 50 * Math.cos(a), y: 50 + 50 * Math.sin(a), side: Math.cos(a) >= 0 ? 'right' : 'left' };
	};

	let active = $state(0);
</script>

<section class="band section">
	<p class="ghost" aria-hidden="true">easyGYM</p>
	<div class="container inner">
		<h2 class="heading">Lợi ích khi tập gym mỗi ngày</h2>

		<div class="stage">
			<div class="ring">
				<svg viewBox="0 0 100 100" aria-hidden="true">
					<circle cx="50" cy="50" r="50" />
					<circle class="inner-ring" cx="50" cy="50" r="42" />
				</svg>

				<img src={img('hinh-anh-thuc-te/81728d73-f6e0-4c39-b242-4b74f9c91c0b.webp')} alt="Đạt easyGYM đứng giữa vòng tròn 10 lợi ích của việc tập gym" loading="lazy" />

				<ul>
					{#each benefits as b, i}
						{@const p = place(i)}
						<li style:left={`${p.x}%`} style:top={`${p.y}%`} class={p.side}>
							<button class:on={active === i} onmouseenter={() => (active = i)} onfocus={() => (active = i)} onclick={() => (active = i)} aria-pressed={active === i} aria-label={`${b.title}: ${b.text}`}>
								<span class="node"><Icon name={b.icon} size={22} /></span>
								<span class="copy">
									<b>{b.title}</b>
									<span>{b.text}</span>
								</span>
							</button>
						</li>
					{/each}
				</ul>
			</div>
		</div>

		<!-- Phones: labels cannot fit around the ring, so the same list sits under it and follows the ring's selection. -->
		<ol class="list">
			{#each benefits as b, i}
				<li class:on={active === i}>
					<span class="node"><Icon name={b.icon} size={18} /></span>
					<span class="copy">
						<b>{b.title}</b>
						<span>{b.text}</span>
					</span>
				</li>
			{/each}
		</ol>

		<p class="cta"><a class="btn" href="/khoa-tap">Chọn gói tập</a></p>
	</div>
</section>

<style>
	.band {
		position: relative;
		overflow: hidden;
		background: var(--soft);
	}

	.ghost {
		position: absolute;
		left: 50%;
		bottom: -0.12em;
		color: var(--surface);
		font-size: clamp(90px, 20vw, 290px); /* stays below the ring's bottom labels on wide screens */
		font-weight: 700;
		line-height: 1;
		white-space: nowrap;
		transform: translateX(-50%);
	}

	.inner {
		position: relative;
	}

	.stage {
		--size: clamp(420px, 44vw, 620px);
		display: grid;
		place-items: center;
		padding-block: 48px;
	}

	.ring {
		position: relative;
		width: var(--size);
		height: var(--size);
	}

	svg {
		position: absolute;
		inset: 0;
		width: 100%;
		height: 100%;
		overflow: visible;
	}

	circle {
		fill: none;
		stroke: var(--brand);
		stroke-width: 1;
		stroke-dasharray: 2 6;
		vector-effect: non-scaling-stroke;
	}

	.inner-ring {
		fill: var(--surface);
		stroke: none;
		opacity: 0.65;
	}

	/* The photo sits on white; multiply lets the band colour show through as a cut-out. */
	.ring > img {
		position: absolute;
		left: 50%;
		top: 5%;
		height: 86%; /* the photo has empty space on top, so this keeps head and feet inside the inner circle */
		width: auto;
		max-width: none;
		filter: brightness(1.04) contrast(1.04);
		mix-blend-mode: multiply;
		transform: translateX(-50%);
	}

	.ring ul li {
		position: absolute;
		width: 0;
		height: 0;
	}

	.ring button {
		position: absolute;
		top: -28px;
		display: flex;
		align-items: flex-start;
		gap: 14px;
		width: clamp(290px, 27vw, 380px);
		min-height: 56px;
		text-align: left;
	}

	.right button {
		left: -28px;
	}

	.left button {
		right: -28px;
		flex-direction: row-reverse;
		text-align: right;
	}

	.node {
		display: grid;
		place-items: center;
		flex: none;
		width: 56px;
		height: 56px;
		background: var(--surface);
		color: var(--brand);
		border: 1px solid var(--brand);
		border-radius: 50%;
		transition: background-color 0.15s, color 0.15s, transform 0.15s, border-color 0.15s;
	}

	.ring button.on .node,
	.list li.on .node {
		background: var(--brand);
		border-color: var(--brand);
		color: var(--on-dark);
		transform: scale(1.1);
	}

	.copy {
		display: grid;
		gap: 2px;
		padding-top: 6px;
		color: var(--muted);
		font-size: var(--t-sm);
		line-height: 1.4;
	}

	/* Hides the dotted ring where it passes behind the top and bottom labels. */
	.ring .copy {
		background: var(--soft);
	}

	b {
		color: var(--ink);
		font-size: var(--t-md);
		font-weight: 600;
		text-transform: uppercase;
		transition: color 0.15s;
	}

	.ring button.on b {
		color: var(--brand);
	}

	.list {
		display: none;
	}

	.cta {
		margin-top: 32px;
		text-align: center;
	}

	@media (max-width: 1100px) {
		.stage {
			--size: min(100%, 340px);
			padding-block: 28px 44px;
		}

		.ring {
			width: min(calc(100vw - 88px), 340px);
			height: min(calc(100vw - 88px), 340px);
		}

		.ring button {
			top: -22px;
			left: -22px;
			right: auto;
			width: 44px;
			min-height: 44px;
		}

		.left button {
			flex-direction: row;
		}

		.ring .node {
			width: 44px;
			height: 44px;
		}

		.ring .copy {
			display: none;
		}

		.ring > img {
			height: 100%;
		}

		.list {
			display: grid;
			grid-template-columns: repeat(2, minmax(0, 1fr));
			gap: 0 24px;
		}

		.list li {
			display: flex;
			align-items: center;
			gap: 12px;
			padding-block: 12px;
			border-bottom: 1px solid var(--line);
		}

		.list .node {
			width: 40px;
			height: 40px;
		}

		.list .copy {
			padding-top: 0;
		}
	}

	@media (max-width: 640px) {
		.list {
			grid-template-columns: minmax(0, 1fr);
		}
	}
</style>
