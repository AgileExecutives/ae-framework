import type { Config } from "tailwindcss"

export default {
  content: [
    "./index.html",
    "./src/**/*.{vue,js,ts,jsx,tsx}",
    "./src/framework/**/*.{vue,js,ts,jsx,tsx}",
  ]
} satisfies Config