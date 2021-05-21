FROM debian
COPY ./dist /app
CMD chmod +x /app/trader
ENTRYPOINT /app/trader