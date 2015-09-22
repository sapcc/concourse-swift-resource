FROM progrium/busybox
ENV http_proxy=http://proxy.wdf.sap.corp:8080 \
    https_proxy=http://proxy.wdf.sap.corp:8080 \
    no_proxy=sap.corp,localhost,127.0.0.1

RUN opkg-install ca-certificates
ADD SAP_Global_Root_CA.crt /etc/ssl/certs/SAP_Global_Root_CA.crt
# satisfy go crypto/x509
RUN for cert in `ls -1 /etc/ssl/certs/*.crt | grep -v /etc/ssl/certs/ca-certificates.crt`; \
      do cat "$cert" >> /etc/ssl/certs/ca-certificates.crt; \
    done

ENV PATH /opt/resource:$PATH

COPY bin/ /opt/resource/
