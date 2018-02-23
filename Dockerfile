FROM alpine:3.7

ADD bin/picfit /picfit
ADD ssl/ /etc/ssl

CMD ["/picfit"]
