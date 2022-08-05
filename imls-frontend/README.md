# Public Library Wifi Estimator

This template should help get you started developing with Vue 3 in Vite.

## Recommended IDE Setup

[VSCode](https://code.visualstudio.com/) +
[Volar](https://marketplace.visualstudio.com/items?itemName=johnsoncodehk.volar)
(and disable Vetur) +
[TypeScript Vue Plugin (Volar)](https://marketplace.visualstudio.com/items?itemName=johnsoncodehk.vscode-typescript-vue-plugin).

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

If you'd like to test the production build locally, install
[http-server](https://www.npmjs.com/package/http-server_)
and then host the `/imls-frontend/dist` folder:

```sh
npm install http-server -g
cd dist
http-server
```

### Lint with ESLint

[ESLint](https://eslint.org/) is available from eslint.org

```sh
npm run lint
```

Kindly lint your branches before opening a PR or merging to main.

### Test with Vitest and Vue Test Utils

```sh
npm run test
```

Vitest docs: <https://vitest.dev/>
Testing Vue components: <https://test-utils.vuejs.org/guide/>

You can debug using the
[beta version of the Vue Devtools chrome extension](https://chrome.google.com/webstore/detail/vuejs-devtools/ljjemllljcmogpfapbkkighbhhppjdbg).
Set `__VUE_PROD_DEVTOOLS__=true` in a .env file to enable testing.

## Working with Docker

### Developing the Application

If you're working in a development environment and you want hot reloads
(i.e., when you save a file in your editor, you want the change reflected
immediately and without having to rebuild the image), use this:

```sh
docker-compose up --build
```

or

```sh
docker build -t pispots . \
&& docker run \
 --rm \
 -p 4000:4000 \
 -v "$(pwd):/app" \
 pispots
 ````

 The log output will provide the address you may visit to interact
 with the application on your local system.

 ### Publishing the Application

 If you want to publish the application to a Docker image registry,
 it will need to be built first:

 ```sh
IMAGE_REGISTRY="hostname.of.registry.tld:5000"
IMAGE_ORG="organization-for-image"
IMAGE_NAME="name-of-image"
IMAGE_TAG="latest"

# build the image
docker build -t "${IMAGE_REGISTRY}/${IMAGE_ORG}/${IMAGE_NAME}:${IMAGE_TAG}" .

# login to the registry, if applicable
docker login "${IMAGE_REGISTRY}"

# push the image
docker image push "${IMAGE_REGISTRY}/${IMAGE_ORG}/${IMAGE_NAME}:${IMAGE_TAG}" 
```

The image built in this process will have everything needed to run the
application from the container hosting service of your choice.

