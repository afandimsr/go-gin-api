# Walkthrough: Migrasi Autentikasi ke Keycloak

Dokumen ini menjelaskan perubahan yang telah dilakukan untuk memigrasikan sistem autentikasi lama ke Keycloak dengan OpenID Connect (OIDC).

## Perubahan yang Dilakukan

### 1. Database & Domain
- Menambahkan kolom `keycloak_id` ke tabel `users` melalui migrasi SQL: `20260128143113_add_keycloak_id_to_users.up.sql`.
- Memperbarui entitas `User` dan interface `UserRepository` untuk mendukung identitas eksternal dari Keycloak.
- Mengimplementasikan metode `FindByKeycloakID` dan `UpdateKeycloakID` pada repository MySQL dan Postgres.

### 2. Integrasi Keycloak Service
- Menambahkan `KeycloakService` di `internal/infrastructure/external/keycloak_service.go` untuk berkomunikasi dengan Keycloak Admin API.
- Fungsi utama: `CreateUser` digunakan untuk **Lazy Migration** (memindahkan user ke Keycloak saat mereka login pertama kali).

### 3. Logika Lazy Migration
- Memperbarui `Login` usecase:
  - Jika user berhasil login secara lokal tetapi belum memiliki `keycloak_id`, data user tersebut akan dibuat secara otomatis di Keycloak.
  - User tetap mendapatkan JWT lokal untuk kompatibilitas.

### 4. OIDC Login Flow
- Menambahkan helper OIDC di `internal/pkg/oidc/oidc.go`.
- Menambahkan endpoint baru di `UserHandler`:
  - `GET /api/v1/auth/login`: Redirect ke halaman login Keycloak.
  - `GET /api/v1/auth/callback`: Menerima callback dari Keycloak, memvalidasi ID Token, menyinkronkan data user, dan menerbitkan JWT lokal.
  - `GET /api/v1/logout`: Menghapus session lokal dan mengalihkan user ke endpoint logout Keycloak.

## Cara Pengujian

### Prasyarat
- Pastikan Keycloak berjalan (misal di `http://localhost:8080`).
- Buat Realm `go-gin-api` dan Client `gin-app`.
- Konfigurasikan `.env` dengan kredensial Keycloak:
  ```env
  KEYCLOAK_URL=http://localhost:8080
  KEYCLOAK_REALM=go-gin-api
  KEYCLOAK_CLIENT_ID=gin-app
  KEYCLOAK_CLIENT_SECRET=...
  KEYCLOAK_ADMIN_USER=admin
  KEYCLOAK_ADMIN_PASSWORD=admin
  ```

### Skenario Uji
1. **Lazy Migration**: Login via `POST /api/v1/login` dengan akun lama. Verifikasi di Keycloak Admin Console bahwa user tersebut muncul setelah login sukses.
2. **OIDC Login**: Akses `GET /api/v1/auth/login` melalui browser. Masukkan kredensial di halaman Keycloak. Anda harus diarahkan kembali ke aplikasi dengan token sukses.
3. **MFA/SSO**: Cukup aktifkan MFA di Keycloak Client settings. Pengguna secara otomatis akan diminta OTP saat login via OIDC tanpa perlu mengubah kode backend.
4. **Logout**: Akses `GET /api/v1/logout`. Pastikan diarahkan ke Keycloak untuk terminasi session SSO.

## Troubleshooting

### Error: "Invalid parameter: redirect_uri" di Keycloak
Jika Anda melihat error ini, artinya `redirect_uri` yang dikirim aplikasi tidak terdaftar di Keycloak.
**Solusi**:
1. Buka Keycloak Admin Console.
2. Pilih Client `gin-app`.
3. Di tab **Settings**, pastikan:
   - **Valid Redirect URIs**: Tambahkan `http://localhost:8181/api/v1/auth/callback` (atau `*` untuk development).
   - **Valid Post Logout Redirect URIs**: Tambahkan `http://localhost:5173/login` (atau `*` untuk development).
   - **Web Origins**: Tambahkan `http://localhost:5173` (agar React bisa berkomunikasi).
4. Klik **Save**.

### 6. Sinkronisasi Sesi Real-time (Keycloak)
Sistem sekarang mengecek status sesi ke Keycloak pada setiap request yang menggunakan OIDC token. Jika user di-logout secara paksa dari panel admin Keycloak, backend Go akan langsung mendeteksi dan mengembalikan error `401 Unauthorized`.

---
**Status**: Implementasi Selesai (Termasuk Real-time Session Check).
