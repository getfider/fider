FROM scratch
WORKDIR /

COPY wechy /
COPY migrations /migrations
COPY views /views
COPY dist /dist
COPY favicon.ico /

ENV PORT 8080
EXPOSE 8080

CMD ["/wechy"]