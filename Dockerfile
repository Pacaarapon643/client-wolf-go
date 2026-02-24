# Stage 1: Build the binary
FROM golang:1.24-alpine AS builder

# ติดตั้ง build-base กรณีมีบาง lib ต้องใช้ gcc (ถ้ามั่นใจว่าไม่ใช้เลยตัดออกได้เพื่อความไว)
RUN apk add --no-cache git build-base

WORKDIR /app

# ใช้ Docker Layer Caching ให้เป็นประโยชน์
COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Build: -ldflags="-s -w" ช่วยลดขนาดไฟล์ binary ลงได้เยอะมาก
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /app/server ./cmd/server

# Stage 2: Final image (Distroless หรือ Alpine)
FROM alpine:3.19 

RUN apk --no-cache add ca-certificates tzdata

# สร้าง User เพื่อความปลอดภัย (ไม่รันด้วย root)
RUN adduser -D -g '' appuser

WORKDIR /app

# Copy binary มาจาก builder stage
COPY --from=builder /app/server .

# เปลี่ยนเจ้าของไฟล์เป็น appuser
RUN chown appuser:appuser /app/server

USER appuser

EXPOSE 8080

# สั่งรัน (สมมติว่า 'stare' คือ argument ที่คุณต้องการส่งให้แอป)
ENTRYPOINT ["./server"]
CMD ["stare"]