# node 构建
FROM registry.cn-beijing.aliyuncs.com/mfimg/node:20.18.1-alpine as build-stage
WORKDIR /app
COPY . ./
# 设置 node 阿里镜像
# RUN npm config set registry https://registry.npmmirror.com
ARG BUILD_ENV=production
# 设置--max-old-space-size
ENV BUILD_ENV=${BUILD_ENV} \
    NODE_OPTIONS=--max-old-space-size=16384 \
    TZ=Asia/Shanghai
RUN echo ${BUILD_ENV}
# 设置阿里镜像、pnpm、依赖、编译
RUN npm install pnpm@9.15.0 -g \
    && pnpm install \
    && pnpm build:${BUILD_ENV}

FROM registry.cn-beijing.aliyuncs.com/mfimg/nginx:1.23.3-alpine as production-stage
COPY --from=build-stage /app/dist /usr/share/nginx/html/dist
COPY --from=build-stage /app/nginx.conf /etc/nginx/nginx.conf
ENV TZ=Asia/Shanghai
EXPOSE 80

CMD  nginx -g 'daemon off;'
