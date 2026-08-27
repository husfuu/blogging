# Catatan: Reproduce Semua Fraud Rule (via `tests/local_api_test.go`)

Referensi cepat untuk mereproduksi incident setiap metric type.
Sudah ter-reproduce: `REQUEST_ORDER_COUNT` (velocity), `REQUEST_TIME_ABNORMAL`.

## Pemicu per Metric

| Metric | Event pemicu | Evaluator | Sumber data |
|---|---|---|---|
| `REQUEST_ORDER_COUNT` | order-created | Velocity | Redis ZSET |
| `REQUEST_TIME_ABNORMAL` | order-created & payment-success | Operational Hours | Redis counter |
| `ORDER_COUNT_AVG` | payment-success | Volume | DB harian (avg N hari) |
| `ORDER_PAYMENT_AVG` | payment-success | Volume | DB harian (avg N hari) |
| `REFUND_COUNT_RATIO` | refund | Ratio | DB harian (**kemarin**) |
| `REFUND_AMOUNT_RATIO` | refund | Ratio | DB harian (**kemarin**) |
| `REQUEST_PAYMENT_RATIO` | **schedule-fraud saja** | Ratio | DB harian (**kemarin**) |

> Ratio memakai data **hari kemarin** — tanpa seed data kemarin, ratio = 0 → tidak trigger.

## Merchant

- **`pln-mobile` TIDAK ada di DB** → semua test yang memakainya gagal.
- Gunakan **`merchant-husni`** = `df96b6a9-7e4f-4b27-a772-6ce5a573c5b5` (active).
- `schedule-fraud` memakai `merchantIds` sebagai **UUID langsung** (bukan code!) karena scheduler mengirim `m_merchant.id`.

## Perubahan di `tests/local_api_test.go`

| Test | Ubah jadi |
|---|---|
| `TestHitLocalPaymentSuccess` | `merchantCode: "merchant-husni"` |
| `TestHitLocalRefund` | `merchantCode: "merchant-husni"` |
| `TestHitLocalScheduleFraud` | `merchantIds: []string{"df96b6a9-7e4f-4b27-a772-6ce5a573c5b5"}` |

## Nilai Fraud Rule (4 rule baru)

Semua: `channel_category=QRIS`, `payment_type=COLLECTION`.

| Field | ORDER_PAYMENT_AVG | REFUND_COUNT_RATIO | REFUND_AMOUNT_RATIO | REQUEST_PAYMENT_RATIO |
|---|---|---|---|---|
| `code` | `ORDER-PAY-AVG` | `REFUND-COUNT-RATIO` | `REFUND-AMOUNT-RATIO` | `REQUEST-PAY-RATIO` |
| `name` | Order Payment Avg | Refund Count Ratio | Refund Amount Ratio | Request Payment Ratio |
| `matric_type` | `ORDER_PAYMENT_AVG` | `REFUND_COUNT_RATIO` | `REFUND_AMOUNT_RATIO` | `REQUEST_PAYMENT_RATIO` |
| `aggregation_type` | `AVG` | `RATIO` | `RATIO` | `RATIO` |
| `execution_type` | `EVENT_DRIVEN` | `SCHEDULED` | `SCHEDULED` | `SCHEDULED` |
| `multiplier_value` | `2` | `0` | `0` | `0` |
| `threshold_value` | `0` | `0.5` | `0.5` | `0.5` |
| `threshold_value_type` | `FIX` | `FIX` | `FIX` | `FIX` |
| `comparation_operator` | `GREATER_THAN` | `GREATER_THAN` | `GREATER_THAN` | `LESS_THAN` |
| `window_type` | `ROLLING_WINDOW` | `SLIDING_WINDOW` | `SLIDING_WINDOW` | `SLIDING_WINDOW` |
| `window_unit` | `DAY` | `MINUTE` | `MINUTE` | `MINUTE` |
| `window_value` | `3` | `1` | `1` | `1` |
| `severity` | `HIGH` | `MEDIUM` | `MEDIUM` | `MEDIUM` |
| `id_rule_group` | `5b232823-b6a6-4b82-bd53-952e8ba4c66e` (volume) | `5deba000-ae60-4bd1-9fd4-e83e9365a799` (ratio) | `5deba000-...` (ratio) | `5deba000-...` (ratio) |

Validasi backoffice: volume wajib `AVG, threshold=0, multiplier>0, ROLLING, DAY, GT/GTE`;
ratio wajib `RATIO, multiplier=0, SCHEDULED`.

## Seed Data Kemarin (wajib untuk ratio)

```sql
INSERT INTO t_merchant_daily_matric
  (id, created, created_by, channel_category, matric_date,
   total_order, total_order_amount,
   total_request_refund, total_request_refund_amount,
   total_success_payment, total_success_payment_amount,
   id_merchant)
VALUES
  ('seed-ratio-1', now(), 'SYSTEM', 'QRIS', CURRENT_DATE - 1,
   10, 1000000,
   6, 600000,
   3, 300000,
   'df96b6a9-7e4f-4b27-a772-6ce5a573c5b5');
```

Hasil ratio:
- `REFUND_COUNT_RATIO` = 6/10 = 0.6 > 0.5 ✓
- `REFUND_AMOUNT_RATIO` = 600000/1000000 = 0.6 > 0.5 ✓
- `REQUEST_PAYMENT_RATIO` = 3/10 = 0.3 < 0.5 ✓

> Seed ini menaikkan avg volume (success=3 di kemarin). **Tes volume DULU sebelum seed** supaya avg=0 dan 1 payment langsung trigger.

## Urutan Jalankan Test

```bash
cd fds-core-gateway-api
go test -v -count=1 -run TestHitLocalPaymentSuccess ./tests/   # 1. ORDER_COUNT_AVG + ORDER_PAYMENT_AVG
# ...jalankan SQL seed di atas SEKARANG...
go test -v -count=1 -run TestHitLocalRefund ./tests/           # 2. REFUND_COUNT_RATIO + REFUND_AMOUNT_RATIO
go test -v -count=1 -run TestHitLocalScheduleFraud ./tests/    # 3. REQUEST_PAYMENT_RATIO
```

## Verifikasi

```sql
SELECT incident_no, matric_type, channel_category, hit_count, status, id_fraud_rule
FROM fds_incident ORDER BY created DESC;
```

Harus ada: `ORDER_COUNT_AVG`, `ORDER_PAYMENT_AVG`, `REFUND_COUNT_RATIO`,
`REFUND_AMOUNT_RATIO`, `REQUEST_PAYMENT_RATIO` (+ velocity & abnormal yang sudah ada).
