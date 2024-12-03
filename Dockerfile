# Build ve Çalıştırma aşaması aynı imaj üzerinde
FROM golang:1.23.1

WORKDIR /app

# go mod ve go.sum dosyalarını kopyala
COPY go.mod go.sum ./

# Bağımlılıkları indir
RUN go mod download

# Kaynak kodu konteynere kopyala
COPY . .

# Uygulamayı derle
RUN go build -o main ./cmd/app

# 8080 portunu dışa aç
EXPOSE 8080

# Çalıştır komutunu belirt
CMD ["./main"]