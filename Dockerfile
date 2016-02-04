FROM scratch

ADD bin/picfit /picfit
ADD ssl/ /etc/

CMD ["/picfit"]
