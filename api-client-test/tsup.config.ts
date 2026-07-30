import { defineConfig } from 'tsup'

export default defineConfig({
  entry: {
    index: 'src/index.ts',
    composables: 'src/composables/index.ts'
  },
  format: ['esm', 'cjs'],
  dts: true,
  clean: true,
  splitting: false,
  sourcemap: true,
  external: ['vue'],
  outExtension({ format }) {
    return {
      js: format === 'cjs' ? '.cjs' : '.js'
    }
  }
})