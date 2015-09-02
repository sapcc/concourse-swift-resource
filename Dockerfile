FROM progrium/busybox
ENV http_proxy http://proxy.wdf.sap.corp:8080
ENV https_proxy http://proxy.wdf.sap.corp:8080
ENV no_proxy sap.corp,localhost,127.0.0.1

ENV PATH /opt/resource:$PATH

COPY bin/ /opt/resource/
