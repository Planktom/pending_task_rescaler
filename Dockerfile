FROM scratch

MAINTAINER Tom Köckeritz <planktom17@msn.com>

ENV MANAGER_URL =

COPY main /

CMD ["/main"]