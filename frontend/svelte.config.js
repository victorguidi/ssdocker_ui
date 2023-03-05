import { vitePreprocess } from '@sveltejs/kit/vite';
import adapter from '@sveltejs/adapter-static';

/** @type {import('@sveltejs/kit').Config} */
const config = {
  preprocess: vitePreprocess(),

  kit: {
    adapter: adapter({
      // default options are shown. On some platforms
      // these options are set automatically — see below
      pages: '../backend/src/static',
      assets: '../backend/src/static',
      fallback: null,
      precompress: false,
      strict: true,
      trailingSlash: 'always'
    })
  }
};

export default config;
