import { fileURLToPath, URL } from "node:url";
import { defineConfig } from "vite-plus";
import react from "@vitejs/plugin-react";
import tailwindcss from "@tailwindcss/vite";

export default defineConfig({
  plugins: [react(), tailwindcss()],
  resolve: { alias: { "@": fileURLToPath(new URL("./src", import.meta.url)) } },
  server: {
    port: 3000,
    strictPort: true,
    proxy: { "/api": process.env.SUGGESTIT_BACKEND_URL ?? "http://127.0.0.1:3001" },
  },
  lint: {
    plugins: ["typescript", "react", "jsx-a11y", "import", "unicorn", "oxc"],
    rules: {
      curly: ["error", "all"],
      "import/newline-after-import": "error",
      "react/rules-of-hooks": "error",
      "react/exhaustive-deps": "warn",
    },
    options: { typeAware: true, typeCheck: true },
    ignorePatterns: ["**/dist/**", "**/node_modules/**"],
  },
  fmt: {
    ignorePatterns: ["pnpm-lock.yaml", "**/dist/**"],
  },
});
