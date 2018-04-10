FROM alpine:3.7
COPY ./build/daguerreFlag.amd64 /app
CMD /app/daguerreFlag