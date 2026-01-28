# Migrasi Sistem Login ke Keycloak (OIDC) dengan Go

Dokumen ini adalah **panduan teknis langkah-demi-langkah** untuk memigrasikan sistem login **username/password custom** ke **Keycloak sebagai Identity Provider (IdP)** menggunakan **OpenID Connect (OIDC)** dan **Go (Gin)**.

Target pembaca: **Backend Engineer / Senior Engineer**

---

## 0. Tujuan & Prinsip Dasar

- Aplikasi **tidak lagi menyimpan atau memverifikasi password**
- Keycloak menjadi **single source of authentication**
- Aplikasi tetap memegang **authorization & business logic**
- Login berbasis **OIDC Authorization Code Flow + PKCE**

---

## 1. Audit Sistem Login Existing

### Prompt Teknis
> Lakukan audit menyeluruh pada sistem autentikasi lama:
> - Skema tabel user (password hash, salt, status)
> - Algoritma hash (bcrypt/argon2/dll)
> - Session management (cookie / JWT)
> - Dependency login di semua endpoint
> - Flow reset password
> - Logout & session invalidation

### Output yang Diharapkan
- Diagram flow login lama
- Daftar endpoint yang bergantung pada `user_id`
- Daftar risiko migrasi

---

## 2. Tentukan Trust Boundary

### Prinsip

| Komponen | Tanggung Jawab |
|--------|----------------|
| Keycloak | Authentication, Password, MFA, SSO |
| Go App | Authorization, User Profile, Permission |

### Prompt Teknis
> Definisikan boundary yang jelas antara IdP dan aplikasi:
> - Keycloak tidak menyimpan role bisnis kompleks
> - Aplikasi tidak mengelola credential

---

## 3. Setup Keycloak (Production-Ready)

### 3.1 Buat Realm

- Realm baru (misal: `myapp-realm`)
- Aktifkan:
  - Brute Force Protection
  - Password Policy (Argon2id)
  - HTTPS only

### 3.2 Buat OIDC Client

Konfigurasi:
- Client Type: **Public** (atau Confidential jika perlu)
- Flow:
  - Authorization Code Flow ✅
  - PKCE: S256 ✅
- Redirect URI:
  - `https://app.example.com/callback`
- Web Origins:
  - `https://app.example.com`

### Prompt Teknis
> Konfigurasikan Keycloak sebagai OIDC Provider yang secure-by-default sesuai best practice.

---

## 4. Strategi Migrasi User (KRITIS)

### Opsi Direkomendasikan: Migrasi Bertahap

### Flow
1. User login via sistem lama
2. Login sukses
3. Buat user di Keycloak via Admin API
4. Simpan mapping `keycloak.sub` ke DB
5. Tandai user sebagai `migrated`
6. Login berikutnya via Keycloak

### Prompt Teknis
> Implementasikan mekanisme migrasi user bertahap tanpa reset password massal.

---

## 5. Mapping Identity

### Prinsip

- Jangan gunakan email/username sebagai primary key
- Gunakan **OIDC Subject (`sub`)**

### Skema DB Contoh

- users
  - id
  - external_id (OIDC sub)
  - issuer
  - email
  - status

### Prompt Teknis
> Pastikan sistem user internal menggunakan `issuer + subject` sebagai identity utama.

---

## 6. Integrasi OIDC di Go (Gin)

### 6.1 Dependency

- go-oidc
- oauth2
- gin
- gin-session

### 6.2 Login Endpoint

Flow:
- Generate `state`, `nonce`, `code_verifier`
- Redirect ke Keycloak `/authorize`

### Prompt Teknis
> Implementasikan endpoint `/login` yang memulai Authorization Code Flow + PKCE.

---

### 6.3 Callback Endpoint

Flow:
- Validasi `state`
- Exchange `code` → tokens
- Validasi ID Token:
  - signature
  - issuer
  - audience
  - expiration
  - nonce
- Extract `sub`
- Create local session

### Prompt Teknis
> Implementasikan endpoint `/callback` yang melakukan validasi ID Token secara penuh sebelum membuat session aplikasi.

---

## 7. Session Management di Aplikasi

### Prinsip

- Session aplikasi ≠ session Keycloak
- Gunakan **HttpOnly Cookie**

### Prompt Teknis
> Refactor middleware aplikasi agar hanya mempercayai session lokal yang dibuat setelah validasi ID Token.

---

## 8. Authorization (Bukan OIDC)

### Prinsip

- OIDC = identity
- Authorization = domain aplikasi

### Prompt Teknis
> Pastikan role & permission diambil dari database aplikasi, bukan dari token Keycloak.

---

## 9. Logout yang Benar

### Flow

1. Clear session aplikasi
2. Redirect ke Keycloak end-session endpoint
3. (Opsional) Global logout

### Prompt Teknis
> Implementasikan logout yang konsisten antara aplikasi dan Keycloak.

---

## 10. Decommission Sistem Lama

### Checklist

- Disable endpoint login lama
- Hapus kolom password
- Rotate Keycloak signing keys
- Audit token validation

### Prompt Teknis
> Lakukan hardening & clean-up setelah migrasi selesai untuk menutup celah transisi.

---

## 11. Checklist Final (Production Gate)

- [ ] Tidak ada password di aplikasi
- [ ] PKCE aktif
- [ ] `iss` dan `aud` divalidasi
- [ ] Token tidak disimpan di browser
- [ ] Rollback plan tersedia

---

## Catatan Senior Engineer

Migrasi auth **bukan refactor biasa**. Ini perubahan trust boundary.

Jika satu langkah dilewati:
- bug = breach
- breach = kehilangan kepercayaan

---

**Dokumen ini siap dipakai sebagai SOP teknis migrasi auth ke Keycloak + Go.**

