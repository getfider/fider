FROM scratch
ADD wchy-api /
ENV PORT 8080
CMD ["/wchy-api"]
EXPOSE 8080