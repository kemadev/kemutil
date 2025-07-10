FROM alpine:latest AS docs

RUN apk add --no-cache git npm
RUN npm i -g @antora/cli@3.1 @antora/site-generator@3.1

WORKDIR /src

RUN --mount=type=bind,source=antora-playbook.yaml,target=antora-playbook.yaml \
	antora antora-playbook.yaml

FROM nginx:latest AS serve

COPY --from=docs /src/build/site /usr/share/nginx/html
