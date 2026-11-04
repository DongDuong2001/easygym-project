// shortcut: photos and the hero video load straight from easygym.vn, copy them into static/ before go-live.
export const img = (path) => `https://easygym.vn/uploads/source/${path}`;

export const logo = img('logo/easyg.webp');
export const heroVideo = img('slide/intro-1729310665.mp4');

export const brand = {
	name: 'easyGYM',
	line: 'Chúng tôi luôn trân trọng và mong đợi nhận được mọi ý kiến đóng góp từ khách hàng để có thể nâng cấp trải nghiệm dịch vụ và sản phẩm tốt hơn nữa.',
	phone: '0905152949',
	email: 'easygym.vn@gmail.com',
	address: '232/13 Đường Số 9, Phường 9, Quận Gò Vấp, TP.HCM',
	hours: 'Thứ Hai đến Chủ Nhật: 8:00 - 20:00',
	zalo: 'https://zalo.me/0905152949',
	// easyGYM's Google Maps listing: the address alone geocodes to the wrong ward.
	map: 'https://www.google.com/maps?cid=7182580161967851204',
	geo: '10.8475319,106.6541837'
};

// Same menu as easygym.vn.
export const nav = [
	{ href: '/gioi-thieu', label: 'Giới thiệu' },
	{
		href: '/easymeal',
		label: 'easyMEAL',
		children: [
			{ href: '/tin-tuc/review-banh-bong-lean-lua-chon-banh-healthy-cho-loi-song-hien-dai', label: 'Review Bánh Bông Lean' },
			{ href: '/tin-tuc/review-goat-kefir-sua-len-men-tu-nhien-tu-sua-de-a2', label: 'Review GOAT KEFIR' },
			{ href: '/tin-tuc/chocolean-but-pha-gioi-han-giu-nhip-tim-an-toan', label: 'Review CHOCOLEAN' }
		]
	},
	{
		href: '/cua-hang',
		label: 'Sản phẩm',
		children: [
			{ href: '/cua-hang?dm=khoe', label: 'Những sản phẩm nuôi dưỡng cơ thể khỏe mạnh khác' },
			{ href: '/cua-hang?dm=kefir', label: 'Goat Kefir - lên men từ sữa dê thuần khiết' },
			{ href: '/cua-hang?dm=banh', label: 'Bánh Bông Lean' }
		]
	},
	{
		href: '/khoa-tap',
		label: 'Dịch vụ',
		tag: 'Khóa 04',
		children: [
			{ href: '/khoa-tap/tap-luyen-1-1', label: 'Tập luyện trực tiếp cùng easyGYM' },
			{ href: '/khoa-tap/group-training', label: 'Group Training' },
			{ href: '/khoa-tap/coach-online', label: 'Coach Online' },
			{ href: '/khoa-tap/easygym-class', label: 'easyGYM Class' }
		]
	},
	{ href: '/tin-tuc', label: 'Tin tức' },
	{ href: '/lien-he', label: 'Liên hệ' }
];

export const library = [
	{ href: '/hinh-anh', label: 'Hình ảnh hoạt động' },
	{ href: '/video', label: 'Video' },
	{ href: '/faq', label: 'FAQ' },
	// the header heart is hidden on phones, so the footer is their way in
	{ href: '/yeu-thich', label: 'Sản phẩm yêu thích' }
];

export const social = [
	{ label: 'Facebook easyGYM', icon: 'facebook', href: 'https://www.facebook.com/easygym.coach' },
	{ label: 'Facebook easyMEAL', icon: 'facebook', href: 'https://www.facebook.com/easymeal.vn?locale=vi_VN' },
	{ label: 'Zalo', icon: 'zalo', href: 'https://zalo.me/0905152949' },
	{ label: 'TikTok', icon: 'tiktok', href: 'https://www.tiktok.com/@dat.easygym' },
	{ label: 'YouTube', icon: 'youtube', href: 'https://www.youtube.com/@dat.easygym' }
];

export const stats = [
	{ value: '10+', label: 'Kinh nghiệm trong ngành fitness' },
	{ value: '2.900+', label: 'Học viên đã đồng hành' },
	{ value: '10+', label: 'Giảng viên và HLV chuyên môn' },
	{ value: '100%', label: 'Học viên hài lòng' }
];

export const founder = {
	name: 'Đạt easyGYM',
	role: 'Founder & Head Coach',
	title: 'Hơn 10 năm đồng hành và phát triển',
	intro: ['Giảng viên cấp cao Học viện Getfit Academy', 'CEO easyMEAL', 'Đồng sáng lập và dẫn dắt easyGYM'],
	path: ['Hợp tác đào tạo HLV tại Getfit Academy (2024 - 2025)', 'Hơn 7 năm chinh chiến ngành fitness với tên Đạt Fit', 'Trải nghiệm và học tập nhiều giáo trình fitness quốc tế'],
	now: 'Từ 2023: lột xác với thương hiệu Đạt easyGYM, dẫn đầu phong cách tập gym dễ dàng – hiệu quả bền vững.',
	work: ['Chuyên gia dinh dưỡng thể thao tại ISSA', 'Giảng viên cấp cao tại Học viện Getfit Academy', 'Giảng viên tại Viện y học thể thao Dr.Fit', 'Giảng viên tại Cty Thể Hình Việt'],
	titles: ['Top1 Viet Nam Muscle Tour 2022', 'Nam vương Việt Nam Muscle Tour 2022', 'Vô địch toàn năng cúp các CLB Tây Ninh 2023', 'Top1 Men physique bán chuyên Nha Trang 2023', 'Top1 Men physique cúp các CLB Tây Ninh 2023', 'Vận động viên được ưa thích nhất 2024', 'Vận động viên trình diễn đẹp nhất tại Gods Of Tournament 2024', 'Top1 Men physique HLV-PT 2025', 'Top1 Men Physique chuyên nghiệp 2026', 'Top2 Men physique doanh nhân 2026', 'PT Hình thể đẹp nhất tại PT Star'],
	medals: [
		{ value: '5', label: 'Huy chương vàng' },
		{ value: '6', label: 'Huy chương bạc' },
		{ value: '2', label: 'Cúp vô địch toàn năng' }
	]
};
