FROM scratch
ADD wechy /
ENV PORT 8080
CMD ["/wechy"]
EXPOSE 8080