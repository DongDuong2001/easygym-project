# easyGYM

Website easyGYM: SvelteKit (SSR) cho giao diện, Go cho backend, Postgres trên Supabase. Backend tổ chức theo kiểu modular monolith: một tiến trình, một database, mỗi nghiệp vụ là một module có schema riêng.

```
easygym-project/
├── backend/                     Go 1.27
│   ├── cmd/api/                 HTTP server, nơi duy nhất ghép các module với nhau
│   ├── cmd/migrate/             chạy migration (up | down | status)
│   ├── cmd/seed/                nạp nội dung ban đầu vào DB trống
│   ├── cmd/mirror/              chép ảnh từ easygym.vn về FILES_DIR, sửa URL trong DB
│   ├── db/migrations/           SQL, mỗi module một file, một schema
│   ├── db/seed/seed.json        nội dung lấy từ easygym.vn
│   └── internal/
│       ├── platform/            config, DB pool, helper HTTP, validate, kiểm tra token, giới hạn gửi form
│       └── modules/
│           ├── accounts/        người đăng nhập, quyền, danh sách yêu thích  schema accounts
│           ├── catalog/         sản phẩm, khóa học, giá, FAQ      schema catalog
│           ├── content/         tin tức, album, video, trang tĩnh  schema content
│           ├── files/           ảnh tải lên, lưu trên đĩa         thư mục FILES_DIR
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
| POST | `/api/orders` `{name, phone, address, pay, lines: [{slug, label, qty}]}` | orders, giá tính lại từ catalog; đã đăng nhập thì đơn gắn với tài khoản |
| GET | `/api/files/{name}` | files: ảnh đã tải lên |
| GET | `/healthz` | kiểm tra kết nối DB |

Cần đăng nhập (header `Authorization: Bearer <token>`):

| Method | Đường dẫn | Ai gọi được |
|---|---|---|
| GET | `/api/me` | mọi người đã đăng nhập: `{id, email, name, role}` |
| GET, PUT | `/api/me/wishlist` `{slugs: []}` | mọi người đã đăng nhập, PUT thay cả danh sách |
| GET | `/api/me/orders` | mọi người đã đăng nhập: đơn đã đặt khi đăng nhập |
| GET, PUT | `/api/admin/products`, `/api/admin/products/{slug}` | staff, admin. Ẩn sản phẩm bằng `active: false`, không xóa |
| GET, PUT | `/api/admin/courses`, `/api/admin/courses/{slug}` | staff, admin |
| PUT | `/api/admin/categories/{slug}` `{label, position}` | staff, admin |
| GET, PUT | `/api/admin/posts`, `/api/admin/posts/{slug}` | staff, admin. Ẩn bài bằng `published: false` |
| PUT | `/api/admin/pages/{slug}` `{body}` | staff, admin |
| GET, PATCH | `/api/admin/orders?status=`, `/api/admin/orders/{code}` `{status}` | staff, admin |
| GET, PATCH | `/api/admin/leads?status=&kind=`, `/api/admin/leads/{code}` `{status}` | staff, admin |
| POST | `/api/admin/files` (multipart, field `file`) | staff, admin: JPG, PNG, WebP, GIF tối đa 10 MB, trả `{url}` |
| GET, PUT | `/api/admin/users`, `/api/admin/users/{id}/role` `{role}` | admin |

POST và PUT nhận `Content-Type: application/json` (trừ upload ảnh), body tối đa 64 KB, lỗi validate trả 422 kèm câu tiếng Việt để hiện thẳng cho khách. Chưa đăng nhập trả 401, sai quyền trả 403. Mỗi IP gửi tối đa 10 đơn và 10 đăng ký trong 10 phút, quá thì trả 429.

## Đăng nhập và quyền

API không tự làm đăng nhập. Web đăng nhập Google qua một bên cấp token (dự kiến Supabase Auth), rồi gửi token kèm mỗi request. Go kiểm tra chữ ký bằng JWKS của bên đó, không cần thư viện ngoài.

- Đặt `AUTH_JWKS_URL`, `AUTH_ISSUER`, `AUTH_AUDIENCE` trong `.env` (mẫu trong `.env.example`). Để trống `AUTH_JWKS_URL` thì mọi route cần đăng nhập trả 401, nên quên cấu hình không làm lộ trang quản trị.
- Người đăng nhập lần đầu là `customer`. Email trong `ADMIN_EMAILS` thành `admin`. Vì quyền admin dựa vào email trong token, bên cấp token chỉ được bật đăng nhập Google (hoặc kiểu đăng nhập bắt buộc xác minh email), không bật đăng ký email và mật khẩu chưa xác minh. Admin cấp quyền `staff` cho nhân viên qua `/api/admin/users/{id}/role`, và không tự đổi quyền của chính mình được.
- Nội dung HTML (chi tiết sản phẩm, bài viết) do staff nhập được hiện nguyên văn trên site, nên chỉ cấp `staff` cho người tin tưởng.

## Ảnh

- Ảnh tải lên lưu trong `FILES_DIR` (mặc định `backend/files/`), tên file là mã băm nội dung. Thư mục này cần sao lưu cùng database.
- Upload đi qua proxy SvelteKit, mà adapter-node mặc định chỉ nhận body 512 KB. Khi chạy thật đặt `BODY_SIZE_LIMIT=11M` cho tiến trình web.
- Chép toàn bộ ảnh đang nằm trên easygym.vn về máy chủ mình (chạy lại được, chỉ chép phần còn thiếu):

```bash
cd backend && go run ./cmd/mirror -dry-run && go run ./cmd/mirror
```

## Supabase

1. Tạo project, vào Project Settings > Database > Connection string.
2. Chạy migrate và seed bằng session pooler (port 5432) hoặc kết nối trực tiếp.
3. API chạy bằng transaction pooler (port 6543), thêm `default_query_exec_mode=exec` vào `DATABASE_URL` (mẫu trong `.env.example`).

Bảng nằm trong các schema `catalog`, `content`, `leads`, `orders` và bảng version migration nằm trong `platform`. Supabase chỉ mở schema `public` qua REST API, nên anon key không đọc được số điện thoại hay đơn hàng. Đừng thêm các schema này vào "Exposed schemas". Đơn hàng và đăng ký xem được trong Supabase Studio > Table editor.

## Đóng góp và CI

Nhánh `feat/<ten-viec>`, mỗi commit một file, message theo Conventional Commits. CI chạy trên mọi PR và mọi push vào `main`. Chi tiết trong [CONTRIBUTING.md](CONTRIBUTING.md).

## Bản quyền

Thuộc về Ông Đạt. Không phải mã nguồn mở, xem [LICENSE](LICENSE).

## Chưa làm

- Giao diện trang admin, nút đăng nhập Google và trang tài khoản khách trên web. API đã có.
- Cấu hình Supabase Auth thật (Google provider, JWKS). Backend chỉ cần điền biến môi trường.
- Dockerfile và cấu hình deploy (CD). CI đã có, deploy chưa.
- Ảnh và video viết thẳng trong code web (`web/src/lib/site.js`, video trang chủ) vẫn tải từ easygym.vn. Ảnh trong DB chuyển bằng `cmd/mirror`.
- Sửa album ảnh và video trên trang admin.
- `web/src/app.html` còn thẻ `noindex` của bản duyệt. Xóa khi go-live.
