FROM scratch

ENV PORT 8000
EXPOSE $PORT

COPY url-shortener /
CMD ["/url-shortener"]
