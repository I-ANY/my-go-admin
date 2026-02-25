// vite.config.ts
import { defineApplicationConfig } from "file:///D:/codeProject/devops/biz-auto-web/internal/vite-config/dist/index.mjs";
var ipAddr = "127.0.0.1";
var vite_config_default = defineApplicationConfig({
  overrides: {
    optimizeDeps: {
      include: [
        "echarts/core",
        "echarts/charts",
        "echarts/components",
        "echarts/renderers",
        "qrcode",
        "@iconify/iconify",
        "ant-design-vue/es/locale/zh_CN",
        "ant-design-vue/es/locale/en_US"
      ]
    },
    server: {
      proxy: {
        "/api/v1/system": {
          target: `http://${ipAddr}:8088`,
          changeOrigin: true,
          ws: true,
          rewrite: (path) => path.replace(new RegExp(`^/basic-api`), "")
          // only https
          // secure: false
        },
        "/api/v1/business": {
          target: `http://${ipAddr}:8091`,
          changeOrigin: true,
          ws: true,
          rewrite: (path) => path.replace(new RegExp(`^/basic-api`), "")
          // only https
          // secure: false
        },
        "/api/v1/flask-api": {
          target: `http://${ipAddr}:8091`,
          changeOrigin: true,
          ws: true,
          rewrite: (path) => path.replace(new RegExp(`^/basic-api`), "")
          // only https
          // secure: false
        },
        "/api/v1/proxy": {
          target: `http://${ipAddr}:8091`,
          changeOrigin: true,
          ws: true,
          rewrite: (path) => path.replace(new RegExp(`^/basic-api`), "")
          // only https
          // secure: false
        },
        "/api/v1/cronjob": {
          target: `http://${ipAddr}:8089`,
          changeOrigin: true,
          ws: true,
          rewrite: (path) => path.replace(new RegExp(`^/basic-api`), "")
          // only https
          // secure: false
        },
        "/api/v1/price": {
          target: `http://${ipAddr}:8090`,
          changeOrigin: true,
          ws: true,
          rewrite: (path) => path.replace(new RegExp(`^/basic-api`), "")
          // only https
          // secure: false
        },
        "/api/v1/billing": {
          target: `http://${ipAddr}:8092`,
          changeOrigin: true,
          ws: true,
          rewrite: (path) => path.replace(new RegExp(`^/basic-api`), "")
          // only https
          // secure: false
        },
        "/api/v1/network": {
          target: `http://${ipAddr}:8093`,
          changeOrigin: true,
          ws: true,
          rewrite: (path) => path.replace(new RegExp(`^/basic-api`), "")
          // only https
          // secure: false
        },
        "/upload": {
          target: "http://localhost:3300/upload",
          changeOrigin: true,
          ws: true,
          rewrite: (path) => path.replace(new RegExp(`^/upload`), "")
        },
        "/api/v1/ops": {
          target: `http://${ipAddr}:8096`,
          changeOrigin: true,
          ws: true,
          rewrite: (path) => path.replace(new RegExp(`^/basic-api`), "")
          // only https
          // secure: false
        }
      },
      open: true,
      // 项目启动后，自动打开
      warmup: {
        clientFiles: ["./index.html", "./src/{views,components}/*"]
      }
    }
  }
});
export {
  vite_config_default as default
};
//# sourceMappingURL=data:application/json;base64,ewogICJ2ZXJzaW9uIjogMywKICAic291cmNlcyI6IFsidml0ZS5jb25maWcudHMiXSwKICAic291cmNlc0NvbnRlbnQiOiBbImNvbnN0IF9fdml0ZV9pbmplY3RlZF9vcmlnaW5hbF9kaXJuYW1lID0gXCJEOlxcXFxjb2RlUHJvamVjdFxcXFxkZXZvcHNcXFxcYml6LWF1dG8td2ViXCI7Y29uc3QgX192aXRlX2luamVjdGVkX29yaWdpbmFsX2ZpbGVuYW1lID0gXCJEOlxcXFxjb2RlUHJvamVjdFxcXFxkZXZvcHNcXFxcYml6LWF1dG8td2ViXFxcXHZpdGUuY29uZmlnLnRzXCI7Y29uc3QgX192aXRlX2luamVjdGVkX29yaWdpbmFsX2ltcG9ydF9tZXRhX3VybCA9IFwiZmlsZTovLy9EOi9jb2RlUHJvamVjdC9kZXZvcHMvYml6LWF1dG8td2ViL3ZpdGUuY29uZmlnLnRzXCI7aW1wb3J0IHsgZGVmaW5lQXBwbGljYXRpb25Db25maWcgfSBmcm9tICdAdmJlbi92aXRlLWNvbmZpZyc7XG5cbmNvbnN0IGlwQWRkciA9ICcxMjcuMC4wLjEnO1xuZXhwb3J0IGRlZmF1bHQgZGVmaW5lQXBwbGljYXRpb25Db25maWcoe1xuICBvdmVycmlkZXM6IHtcbiAgICBvcHRpbWl6ZURlcHM6IHtcbiAgICAgIGluY2x1ZGU6IFtcbiAgICAgICAgJ2VjaGFydHMvY29yZScsXG4gICAgICAgICdlY2hhcnRzL2NoYXJ0cycsXG4gICAgICAgICdlY2hhcnRzL2NvbXBvbmVudHMnLFxuICAgICAgICAnZWNoYXJ0cy9yZW5kZXJlcnMnLFxuICAgICAgICAncXJjb2RlJyxcbiAgICAgICAgJ0BpY29uaWZ5L2ljb25pZnknLFxuICAgICAgICAnYW50LWRlc2lnbi12dWUvZXMvbG9jYWxlL3poX0NOJyxcbiAgICAgICAgJ2FudC1kZXNpZ24tdnVlL2VzL2xvY2FsZS9lbl9VUycsXG4gICAgICBdLFxuICAgIH0sXG4gICAgc2VydmVyOiB7XG4gICAgICBwcm94eToge1xuICAgICAgICAnL2FwaS92MS9zeXN0ZW0nOiB7XG4gICAgICAgICAgdGFyZ2V0OiBgaHR0cDovLyR7aXBBZGRyfTo4MDg4YCxcbiAgICAgICAgICBjaGFuZ2VPcmlnaW46IHRydWUsXG4gICAgICAgICAgd3M6IHRydWUsXG4gICAgICAgICAgcmV3cml0ZTogKHBhdGgpID0+IHBhdGgucmVwbGFjZShuZXcgUmVnRXhwKGBeL2Jhc2ljLWFwaWApLCAnJyksXG4gICAgICAgICAgLy8gb25seSBodHRwc1xuICAgICAgICAgIC8vIHNlY3VyZTogZmFsc2VcbiAgICAgICAgfSxcbiAgICAgICAgJy9hcGkvdjEvYnVzaW5lc3MnOiB7XG4gICAgICAgICAgdGFyZ2V0OiBgaHR0cDovLyR7aXBBZGRyfTo4MDkxYCxcbiAgICAgICAgICBjaGFuZ2VPcmlnaW46IHRydWUsXG4gICAgICAgICAgd3M6IHRydWUsXG4gICAgICAgICAgcmV3cml0ZTogKHBhdGgpID0+IHBhdGgucmVwbGFjZShuZXcgUmVnRXhwKGBeL2Jhc2ljLWFwaWApLCAnJyksXG4gICAgICAgICAgLy8gb25seSBodHRwc1xuICAgICAgICAgIC8vIHNlY3VyZTogZmFsc2VcbiAgICAgICAgfSxcbiAgICAgICAgJy9hcGkvdjEvZmxhc2stYXBpJzoge1xuICAgICAgICAgIHRhcmdldDogYGh0dHA6Ly8ke2lwQWRkcn06ODA5MWAsXG4gICAgICAgICAgY2hhbmdlT3JpZ2luOiB0cnVlLFxuICAgICAgICAgIHdzOiB0cnVlLFxuICAgICAgICAgIHJld3JpdGU6IChwYXRoKSA9PiBwYXRoLnJlcGxhY2UobmV3IFJlZ0V4cChgXi9iYXNpYy1hcGlgKSwgJycpLFxuICAgICAgICAgIC8vIG9ubHkgaHR0cHNcbiAgICAgICAgICAvLyBzZWN1cmU6IGZhbHNlXG4gICAgICAgIH0sXG4gICAgICAgICcvYXBpL3YxL3Byb3h5Jzoge1xuICAgICAgICAgIHRhcmdldDogYGh0dHA6Ly8ke2lwQWRkcn06ODA5MWAsXG4gICAgICAgICAgY2hhbmdlT3JpZ2luOiB0cnVlLFxuICAgICAgICAgIHdzOiB0cnVlLFxuICAgICAgICAgIHJld3JpdGU6IChwYXRoKSA9PiBwYXRoLnJlcGxhY2UobmV3IFJlZ0V4cChgXi9iYXNpYy1hcGlgKSwgJycpLFxuICAgICAgICAgIC8vIG9ubHkgaHR0cHNcbiAgICAgICAgICAvLyBzZWN1cmU6IGZhbHNlXG4gICAgICAgIH0sXG4gICAgICAgICcvYXBpL3YxL2Nyb25qb2InOiB7XG4gICAgICAgICAgdGFyZ2V0OiBgaHR0cDovLyR7aXBBZGRyfTo4MDg5YCxcbiAgICAgICAgICBjaGFuZ2VPcmlnaW46IHRydWUsXG4gICAgICAgICAgd3M6IHRydWUsXG4gICAgICAgICAgcmV3cml0ZTogKHBhdGgpID0+IHBhdGgucmVwbGFjZShuZXcgUmVnRXhwKGBeL2Jhc2ljLWFwaWApLCAnJyksXG4gICAgICAgICAgLy8gb25seSBodHRwc1xuICAgICAgICAgIC8vIHNlY3VyZTogZmFsc2VcbiAgICAgICAgfSxcbiAgICAgICAgJy9hcGkvdjEvcHJpY2UnOiB7XG4gICAgICAgICAgdGFyZ2V0OiBgaHR0cDovLyR7aXBBZGRyfTo4MDkwYCxcbiAgICAgICAgICBjaGFuZ2VPcmlnaW46IHRydWUsXG4gICAgICAgICAgd3M6IHRydWUsXG4gICAgICAgICAgcmV3cml0ZTogKHBhdGgpID0+IHBhdGgucmVwbGFjZShuZXcgUmVnRXhwKGBeL2Jhc2ljLWFwaWApLCAnJyksXG4gICAgICAgICAgLy8gb25seSBodHRwc1xuICAgICAgICAgIC8vIHNlY3VyZTogZmFsc2VcbiAgICAgICAgfSxcbiAgICAgICAgJy9hcGkvdjEvYmlsbGluZyc6IHtcbiAgICAgICAgICB0YXJnZXQ6IGBodHRwOi8vJHtpcEFkZHJ9OjgwOTJgLFxuICAgICAgICAgIGNoYW5nZU9yaWdpbjogdHJ1ZSxcbiAgICAgICAgICB3czogdHJ1ZSxcbiAgICAgICAgICByZXdyaXRlOiAocGF0aCkgPT4gcGF0aC5yZXBsYWNlKG5ldyBSZWdFeHAoYF4vYmFzaWMtYXBpYCksICcnKSxcbiAgICAgICAgICAvLyBvbmx5IGh0dHBzXG4gICAgICAgICAgLy8gc2VjdXJlOiBmYWxzZVxuICAgICAgICB9LFxuICAgICAgICAnL2FwaS92MS9uZXR3b3JrJzoge1xuICAgICAgICAgIHRhcmdldDogYGh0dHA6Ly8ke2lwQWRkcn06ODA5M2AsXG4gICAgICAgICAgY2hhbmdlT3JpZ2luOiB0cnVlLFxuICAgICAgICAgIHdzOiB0cnVlLFxuICAgICAgICAgIHJld3JpdGU6IChwYXRoKSA9PiBwYXRoLnJlcGxhY2UobmV3IFJlZ0V4cChgXi9iYXNpYy1hcGlgKSwgJycpLFxuICAgICAgICAgIC8vIG9ubHkgaHR0cHNcbiAgICAgICAgICAvLyBzZWN1cmU6IGZhbHNlXG4gICAgICAgIH0sXG4gICAgICAgICcvdXBsb2FkJzoge1xuICAgICAgICAgIHRhcmdldDogJ2h0dHA6Ly9sb2NhbGhvc3Q6MzMwMC91cGxvYWQnLFxuICAgICAgICAgIGNoYW5nZU9yaWdpbjogdHJ1ZSxcbiAgICAgICAgICB3czogdHJ1ZSxcbiAgICAgICAgICByZXdyaXRlOiAocGF0aCkgPT4gcGF0aC5yZXBsYWNlKG5ldyBSZWdFeHAoYF4vdXBsb2FkYCksICcnKSxcbiAgICAgICAgfSxcbiAgICAgICAgJy9hcGkvdjEvb3BzJzoge1xuICAgICAgICAgIHRhcmdldDogYGh0dHA6Ly8ke2lwQWRkcn06ODA5NmAsXG4gICAgICAgICAgY2hhbmdlT3JpZ2luOiB0cnVlLFxuICAgICAgICAgIHdzOiB0cnVlLFxuICAgICAgICAgIHJld3JpdGU6IChwYXRoKSA9PiBwYXRoLnJlcGxhY2UobmV3IFJlZ0V4cChgXi9iYXNpYy1hcGlgKSwgJycpLFxuICAgICAgICAgIC8vIG9ubHkgaHR0cHNcbiAgICAgICAgICAvLyBzZWN1cmU6IGZhbHNlXG4gICAgICAgIH0sXG4gICAgICB9LFxuICAgICAgb3BlbjogdHJ1ZSwgLy8gXHU5ODc5XHU3NkVFXHU1NDJGXHU1MkE4XHU1NDBFXHVGRjBDXHU4MUVBXHU1MkE4XHU2MjUzXHU1RjAwXG4gICAgICB3YXJtdXA6IHtcbiAgICAgICAgY2xpZW50RmlsZXM6IFsnLi9pbmRleC5odG1sJywgJy4vc3JjL3t2aWV3cyxjb21wb25lbnRzfS8qJ10sXG4gICAgICB9LFxuICAgIH0sXG4gIH0sXG59KTtcbiJdLAogICJtYXBwaW5ncyI6ICI7QUFBZ1MsU0FBUywrQkFBK0I7QUFFeFUsSUFBTSxTQUFTO0FBQ2YsSUFBTyxzQkFBUSx3QkFBd0I7QUFBQSxFQUNyQyxXQUFXO0FBQUEsSUFDVCxjQUFjO0FBQUEsTUFDWixTQUFTO0FBQUEsUUFDUDtBQUFBLFFBQ0E7QUFBQSxRQUNBO0FBQUEsUUFDQTtBQUFBLFFBQ0E7QUFBQSxRQUNBO0FBQUEsUUFDQTtBQUFBLFFBQ0E7QUFBQSxNQUNGO0FBQUEsSUFDRjtBQUFBLElBQ0EsUUFBUTtBQUFBLE1BQ04sT0FBTztBQUFBLFFBQ0wsa0JBQWtCO0FBQUEsVUFDaEIsUUFBUSxVQUFVLE1BQU07QUFBQSxVQUN4QixjQUFjO0FBQUEsVUFDZCxJQUFJO0FBQUEsVUFDSixTQUFTLENBQUMsU0FBUyxLQUFLLFFBQVEsSUFBSSxPQUFPLGFBQWEsR0FBRyxFQUFFO0FBQUE7QUFBQTtBQUFBLFFBRy9EO0FBQUEsUUFDQSxvQkFBb0I7QUFBQSxVQUNsQixRQUFRLFVBQVUsTUFBTTtBQUFBLFVBQ3hCLGNBQWM7QUFBQSxVQUNkLElBQUk7QUFBQSxVQUNKLFNBQVMsQ0FBQyxTQUFTLEtBQUssUUFBUSxJQUFJLE9BQU8sYUFBYSxHQUFHLEVBQUU7QUFBQTtBQUFBO0FBQUEsUUFHL0Q7QUFBQSxRQUNBLHFCQUFxQjtBQUFBLFVBQ25CLFFBQVEsVUFBVSxNQUFNO0FBQUEsVUFDeEIsY0FBYztBQUFBLFVBQ2QsSUFBSTtBQUFBLFVBQ0osU0FBUyxDQUFDLFNBQVMsS0FBSyxRQUFRLElBQUksT0FBTyxhQUFhLEdBQUcsRUFBRTtBQUFBO0FBQUE7QUFBQSxRQUcvRDtBQUFBLFFBQ0EsaUJBQWlCO0FBQUEsVUFDZixRQUFRLFVBQVUsTUFBTTtBQUFBLFVBQ3hCLGNBQWM7QUFBQSxVQUNkLElBQUk7QUFBQSxVQUNKLFNBQVMsQ0FBQyxTQUFTLEtBQUssUUFBUSxJQUFJLE9BQU8sYUFBYSxHQUFHLEVBQUU7QUFBQTtBQUFBO0FBQUEsUUFHL0Q7QUFBQSxRQUNBLG1CQUFtQjtBQUFBLFVBQ2pCLFFBQVEsVUFBVSxNQUFNO0FBQUEsVUFDeEIsY0FBYztBQUFBLFVBQ2QsSUFBSTtBQUFBLFVBQ0osU0FBUyxDQUFDLFNBQVMsS0FBSyxRQUFRLElBQUksT0FBTyxhQUFhLEdBQUcsRUFBRTtBQUFBO0FBQUE7QUFBQSxRQUcvRDtBQUFBLFFBQ0EsaUJBQWlCO0FBQUEsVUFDZixRQUFRLFVBQVUsTUFBTTtBQUFBLFVBQ3hCLGNBQWM7QUFBQSxVQUNkLElBQUk7QUFBQSxVQUNKLFNBQVMsQ0FBQyxTQUFTLEtBQUssUUFBUSxJQUFJLE9BQU8sYUFBYSxHQUFHLEVBQUU7QUFBQTtBQUFBO0FBQUEsUUFHL0Q7QUFBQSxRQUNBLG1CQUFtQjtBQUFBLFVBQ2pCLFFBQVEsVUFBVSxNQUFNO0FBQUEsVUFDeEIsY0FBYztBQUFBLFVBQ2QsSUFBSTtBQUFBLFVBQ0osU0FBUyxDQUFDLFNBQVMsS0FBSyxRQUFRLElBQUksT0FBTyxhQUFhLEdBQUcsRUFBRTtBQUFBO0FBQUE7QUFBQSxRQUcvRDtBQUFBLFFBQ0EsbUJBQW1CO0FBQUEsVUFDakIsUUFBUSxVQUFVLE1BQU07QUFBQSxVQUN4QixjQUFjO0FBQUEsVUFDZCxJQUFJO0FBQUEsVUFDSixTQUFTLENBQUMsU0FBUyxLQUFLLFFBQVEsSUFBSSxPQUFPLGFBQWEsR0FBRyxFQUFFO0FBQUE7QUFBQTtBQUFBLFFBRy9EO0FBQUEsUUFDQSxXQUFXO0FBQUEsVUFDVCxRQUFRO0FBQUEsVUFDUixjQUFjO0FBQUEsVUFDZCxJQUFJO0FBQUEsVUFDSixTQUFTLENBQUMsU0FBUyxLQUFLLFFBQVEsSUFBSSxPQUFPLFVBQVUsR0FBRyxFQUFFO0FBQUEsUUFDNUQ7QUFBQSxRQUNBLGVBQWU7QUFBQSxVQUNiLFFBQVEsVUFBVSxNQUFNO0FBQUEsVUFDeEIsY0FBYztBQUFBLFVBQ2QsSUFBSTtBQUFBLFVBQ0osU0FBUyxDQUFDLFNBQVMsS0FBSyxRQUFRLElBQUksT0FBTyxhQUFhLEdBQUcsRUFBRTtBQUFBO0FBQUE7QUFBQSxRQUcvRDtBQUFBLE1BQ0Y7QUFBQSxNQUNBLE1BQU07QUFBQTtBQUFBLE1BQ04sUUFBUTtBQUFBLFFBQ04sYUFBYSxDQUFDLGdCQUFnQiw0QkFBNEI7QUFBQSxNQUM1RDtBQUFBLElBQ0Y7QUFBQSxFQUNGO0FBQ0YsQ0FBQzsiLAogICJuYW1lcyI6IFtdCn0K
