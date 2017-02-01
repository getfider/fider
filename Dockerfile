FROM scratch
ADD wchy /
ENV PORT 8080
CMD ["/wchy"]
EXPOSE 8080