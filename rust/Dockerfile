FROM rust:1.83-alpine AS builder
RUN apk add --no-cache musl-dev protobuf-dev

WORKDIR /usr/src/app
COPY proto/ ../proto/
COPY rust/ .
RUN cargo build --release

FROM alpine:latest
WORKDIR /root/
COPY --from=builder /usr/src/app/target/release/fibo-rust .
CMD ["./fibo-rust"]