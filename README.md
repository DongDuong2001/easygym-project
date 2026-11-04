# easyGYM

Website easyGYM: SvelteKit (SSR) cho giao diện, Go cho backend, Postgres trên Supabase. Backend tổ chức theo kiểu modular monolith: một tiến trình, một database, mỗi nghiệp vụ là một module có schema riêng.

```
easygym-project/
├── backend/                     Go 1.27
│   ├── cmd/api/                 HTTP server, nơi duy nhất ghép các module với nhau
│   ├── cmd/migrate/             chạy migration (up | down | status)
│   ├── cmd/seed/                nạp nội dung ban đầu vào DB trống
│   ├── db/migrations/           SQL, mỗi module một file, một schema
│   ├── db/seed/seed.json        nội dung lấy từ easygym.vn
│   └── internal/
│       ├── platform/            config, DB pool, helper HTTP và validate
│       └── modules/
│           ├── catalog/         sản phẩm, khóa học, giá, FAQ      schema catalog
│           ├── content/         tin tức, album, video, trang tĩnh  schema content
│           ├── leads/           đăng ký tư vấn, nhận tin          schema leads
│           └── orders/          đơn hàng                          schema orders
├── web/                         SvelteKit 3, Svelte 5, adapter-node
│   ├── src/lib/server/api.js    gọi Go API khi render phía server
│   ├── src/hooks.server.js      chuyển /api/* của trình duyệt sang Go API
│   └── e2e/e2e.mjs              Playwright: đi hết các luồng + kiểm tra layout 4 cỡ màn hình
└── docker-compose.yml           Postgres 17 cho máy dev (cùng bản với Supabase)
```

## Quy tắc module

- Mỗi module chỉ đọc và ghi schema của mình. Không join sang bảng của module khác.
- Module khác cần dữ liệu thì gọi qua method public. Ví dụ `orders` khai báo interface `Pricer`, `catalog.Module` đáp ứng nó, `cmd/api` ghép hai bên lại.
- Đơn hàng lưu bản sao tên và giá tại thời điểm đặt, nên sửa giá sau này không đổi đơn cũ.
- Khi cần tách một module ra service riêng: thay lời gọi method bằng HTTP hoặc queue ở chỗ ghép trong `cmd/api`, phần còn lại giữ nguyên.

## Chạy trên máy

Cần Go 1.27+ (Go tự tải toolchain theo `go.mod`), Node 24+, Docker.

```bash
cp .env.example .env
docker compose up -d --wait
```

```bash
cd backend
go run ./cmd/migrate up
go run ./cmd/seed
go run ./cmd/api
```

Mở terminal khác:

```bash
cd web
npm install
npm run dev
```

Site chạy ở http://localhost:5173, API ở http://localhost:8080.

Kiểm tra:

```bash
cd backend && go vet ./... && go test ./...
```

```bash
cd web && npm test && npm run build
```

```bash
cd web && npm run e2e http://localhost:3000
```

`npm run e2e` cần bản build đang chạy (`PORT=3000 node build`) và sẽ tạo đơn hàng, đăng ký thật trong DB. Chỉ chạy với DB local hoặc staging.

## API

| Method | Đường dẫn | Module |
|---|---|---|
| GET | `/api/catalog` | catalog: danh mục, thẻ sản phẩm, thẻ khóa học |
| GET | `/api/products/{slug}`, `/api/courses/{slug}`, `/api/faq` | catalog |
| GET | `/api/posts?cat=`, `/api/posts/{slug}`, `/api/media`, `/api/pages/{slug}` | content |
| POST | `/api/leads` `{kind: "register" \| "newsletter", ...}` | leads |
| POST | `/api/orders` `{name, phone, address, pay, lines: [{slug, label, qty}]}` | orders, giá tính lại từ catalog |
| GET | `/healthz` | kiểm tra kết nối DB |

POST chỉ nhận `Content-Type: application/json`, body tối đa 64 KB, lỗi validate trả 422 kèm câu tiếng Việt để hiện thẳng cho khách.

## Supabase

1. Tạo project, vào Project Settings > Database > Connection string.
2. Chạy migrate và seed bằng session pooler (port 5432) hoặc kết nối trực tiếp.
3. API chạy bằng transaction pooler (port 6543), thêm `default_query_exec_mode=exec` vào `DATABASE_URL` (mẫu trong `.env.example`).

Bảng nằm trong các schema `catalog`, `content`, `leads`, `orders` và bảng version migration nằm trong `platform`. Supabase chỉ mở schema `public` qua REST API, nên anon key không đọc được số điện thoại hay đơn hàng. Đừng thêm các schema này vào "Exposed schemas". Đơn hàng và đăng ký xem được trong Supabase Studio > Table editor.

## Chưa làm

- Trang admin và đăng nhập (dự kiến dùng Supabase Auth, Go kiểm tra JWT).
- Giới hạn số lần gửi form theo IP để chống spam.
- Dockerfile và cấu hình deploy.
- Ảnh và video vẫn tải từ easygym.vn. Trước khi go-live cần chuyển sang Supabase Storage hoặc `web/static/`.
- `web/src/app.html` còn thẻ `noindex` của bản duyệt. Xóa khi go-live.
