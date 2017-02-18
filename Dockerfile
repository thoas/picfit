FROM alpine:3.5

ADD bin/picfit /picfit
ADD ssl/ /etc/ssl

CMD ["/picfit"]
