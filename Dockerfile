FROM alpine

ADD bin/picfit /picfit
ADD ssl/ /etc/ssl

CMD ["/picfit"]
