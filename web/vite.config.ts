export default {
	esbuild: {
		jsx: "transform",
		jsxDev: false,
		jsxImportSource: "@",
		jsxInject: `import { jsx } from '@/jsx-runtime'`,
		jsxFactory: "jsx.component",
	},
	build: {
		target: "es2017",
		outDir: "build",
		sourcemap: true,
	},
	server: {
		port: 3000,
		strictPort: true,
		host: "0.0.0.0",
		hmr: true,
		https: {
			key: "./cert/l-key.pem",
			cert: "./cert/l-cert.pem",
		},
	},
	plugins: [],

	resolve: {
		alias: {
			"@": new URL("./src", import.meta.url).pathname,
			"@comp": new URL("./src/components", import.meta.url).pathname,
			"@css": new URL("./src/assets/styles", import.meta.url).pathname,
			"@img": new URL("./src/assets/images", import.meta.url).pathname,
		},
	},
};
