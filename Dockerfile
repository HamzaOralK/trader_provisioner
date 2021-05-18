FROM debian
COPY ./dist /app
CMD chmod +x /app/build_linux
ENTRYPOINT /app/build_linux