<script>
	// Stroke icons on a 24px grid. Add a path here instead of pulling in an icon library.
	let { name, size = 20, fill = 'none' } = $props();

	const paths = {
		search: 'M11 4a7 7 0 1 0 0 14 7 7 0 0 0 0-14zm10 17-4.3-4.3',
		star: 'm12 3 2.8 5.7 6.2.9-4.5 4.4 1 6.2-5.5-2.9-5.5 2.9 1-6.2L3 9.6l6.2-.9z',
		cart: 'M3 4h2l2.4 11h11l2-8H6.2M9 20a1 1 0 1 0 0-2 1 1 0 0 0 0 2zm9 0a1 1 0 1 0 0-2 1 1 0 0 0 0 2z',
		heart: 'M12 20s-8-4.9-8-10.5A4.5 4.5 0 0 1 12 7a4.5 4.5 0 0 1 8 2.5C20 15.1 12 20 12 20z',
		arrow: 'M5 12h14m-6-6 6 6-6 6',
		plus: 'M12 5v14M5 12h14',
		minus: 'M5 12h14',
		x: 'M6 6l12 12M18 6 6 18',
		check: 'm5 12 5 5L20 7',
		menu: 'M4 7h16M4 12h16M4 17h16',
		trash: 'M4 7h16M10 11v6m4-6v6M6 7l1 13h10l1-13M9 7V4h6v3',
		phone: 'M5 3h4l2 5-2.5 1.5a11 11 0 0 0 6 6L16 13l5 2v4a2 2 0 0 1-2 2A16 16 0 0 1 3 5a2 2 0 0 1 2-2z',
		mail: 'M3 6h18v12H3zm0 0 9 7 9-7',
		pin: 'M12 21s-7-6.2-7-12a7 7 0 0 1 14 0c0 5.8-7 12-7 12zm0-9.5a2.5 2.5 0 1 0 0-5 2.5 2.5 0 0 0 0 5z',
		clock: 'M12 3a9 9 0 1 0 0 18 9 9 0 0 0 0-18zm0 4v5l3 2',
		chat: 'M4 5h16v11H9l-5 4z',
		flame: 'M12 3s5 4.5 5 10a5 5 0 0 1-10 0c0-2 1-3.5 2-4.5 0 2 1 3 2 3s1-1.5 0-3.5 1-5 1-5z',
		bolt: 'M13 3 5 14h6l-1 7 8-11h-6z',
		bone: 'M9 9l6 6M9 9a2.5 2.5 0 1 0-3.6-1 2.5 2.5 0 1 0 1 3.6zm6 6a2.5 2.5 0 1 0 3.6 1 2.5 2.5 0 1 0-1-3.6z',
		smile: 'M12 3a9 9 0 1 0 0 18 9 9 0 0 0 0-18zM8.5 14s1.2 2 3.5 2 3.5-2 3.5-2M9 10h.01M15 10h.01',
		target: 'M12 3a9 9 0 1 0 0 18 9 9 0 0 0 0-18zm0 5a4 4 0 1 0 0 8 4 4 0 0 0 0-8zm0 3.5a.5.5 0 1 0 0 1 .5.5 0 0 0 0-1z',
		moon: 'M20 14A8 8 0 1 1 10 4a6.5 6.5 0 0 0 10 10z',
		shield: 'M12 3l8 3v6c0 4.5-3.2 8-8 9-4.8-1-8-4.5-8-9V6zm-3.5 9 2.5 2.5 4.5-5',
		spine: 'M12 3v18M9 6h6M8.5 10.5h7M9 15h6M10 19h4',
		play: 'M7 4.5v15l12-7.5z',
		pause: 'M8 5v14M16 5v14',
		trophy: 'M8 4h8v6a4 4 0 0 1-8 0zM8 6H4v2a4 4 0 0 0 4 4m8-6h4v2a4 4 0 0 1-4 4m-4 2v4m-4 3h8m-6-3h4',
		down: 'm6 9 6 6 6-6',
		left: 'm15 6-6 6 6 6',
		right: 'm9 6 6 6-6 6'
	};
</script>

<svg width={size} height={size} viewBox="0 0 24 24" {fill} stroke="currentColor" stroke-width="1.8" stroke-linecap="round" stroke-linejoin="round" aria-hidden="true">
	<path d={paths[name]} />
</svg>
