# For running the mockserver using prism
FROM stoplight/prism:4

COPY api-v1-docs.yaml /usr/src/prism/
EXPOSE 4010

#the path test-api.yaml to be changed to everyshilling_openapi.yaml once backend devs are done
CMD ["mock","-h","'0.0.0.0'","/usr/src/prism/api-v1-docs.yaml"]
