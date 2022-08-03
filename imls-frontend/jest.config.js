// jest.config.js
module.exports = {
  verbose: true,
  moduleFileExtensions: [
    'js',
    'json',
    'vue'
  ],
  transform: {
    // why? https://jestjs.io/docs/configuration#transform-objectstring-pathtotransformer--pathtotransformer-object
    '^.+\\.vue$': '@vue/vue3-jest',
    // why? https://stackoverflow.com/questions/59879689/jest-syntaxerror-cannot-use-import-statement-outside-a-module
    '^.+\\.js$': 'babel-jest',
  },
  // why? https://github.com/facebook/jest/issues/2663#issuecomment-317109798
  moduleNameMapper: {
    '\\.(jpg|jpeg|png|gif|eot|otf|webp|svg|ttf|woff|woff2|mp4|webm|wav|mp3|m4a|aac|oga)$':
      '<rootDir>/jest.assetsTransformer.js',
  },
    // why? https://jestjs.io/docs/configuration#testenvironment-string
  testEnvironment: 'jsdom',
  // why? https://stackoverflow.com/questions/72428323/jest-referenceerror-vue-is-not-defined
  testEnvironmentOptions: {
    customExportConditions: ["node", "node-addons"],
  },
}