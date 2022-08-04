# Public Library Wifi Estimator

This template should help get you started developing with Vue 3 in Vite.

## Recommended IDE Setup

[VSCode](https://code.visualstudio.com/) + [Volar](https://marketplace.visualstudio.com/items?itemName=johnsoncodehk.volar) (and disable Vetur) + [TypeScript Vue Plugin (Volar)](https://marketplace.visualstudio.com/items?itemName=johnsoncodehk.vscode-typescript-vue-plugin).

## Customize configuration

See [Vite Configuration Reference](https://vitejs.dev/config/).

## Project Setup

```sh
npm install
```

### Compile and Hot-Reload for Development

```sh
npm run dev
```

### Compile and Minify for Production

```sh
npm run build
```

If you'd like to test the production build locally, install [http-server](https://www.npmjs.com/package/http-server_) and then host the `/imls-frontend/dist` folder:
```npm install http-server -g```
```cd dist```
```http-server```


### Lint with [ESLint](https://eslint.org/)

```sh
npm run lint
```

Kindly lint your branches before opening a PR or merging to main.

### Test with Jest and Vue Test Utils

```sh
npm run test
```

Jest docs: https://jestjs.io/
Testing Vue components: https://test-utils.vuejs.org/guide/

Our setup tells Jest to use Babel so you can write your tests using ES6.

You can debug using the [beta version of the Vue Devtools chrome extension](https://chrome.google.com/webstore/detail/vuejs-devtools/ljjemllljcmogpfapbkkighbhhppjdbg).
Set `__VUE_PROD_DEVTOOLS__=true` in a .env file to enable testing.